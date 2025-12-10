package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"password_manager/internal/config"
	"password_manager/internal/storage"
	"password_manager/internal/storage/database"
	"password_manager/internal/storage/models"
	"password_manager/internal/transport"
	"password_manager/internal/transport/cli"
	"password_manager/internal/transport/message"
	"password_manager/internal/utils"
)

const (
	GetList = iota + 1
	GetEntry
	CreateEntry
	UpdateEntry
	DeleteEntry

	//ChangeMasterPassword = 9
	Exit = 0
)

type App struct {
	config    *config.Config
	ticker    *time.Ticker
	storage   storage.Storage
	transport transport.Core
}

func (a *App) Run(stop context.CancelFunc) {
	exists, err := a.storage.CheckStorageExists()
	if err != nil {
		a.transport.SendMessageToUser(message.UnhandledError + err.Error())
		stop()
	}

	var hiddenString string

	if !exists {
		var hiddenStringRepeat string

		for mismatch := true; mismatch; {
			a.transport.SendMessageToUser(message.FirstTimeEnter)

			hiddenString, err = a.transport.GetPasswordHidden()
			if err != nil {
				a.transport.SendMessageToUser(message.UnhandledError + err.Error())
				stop()
			}

			a.transport.SendMessageToUser(message.NewLine + message.RepeatMasterPassword)

			hiddenStringRepeat, err = a.transport.GetPasswordHidden()
			if err != nil {
				a.transport.SendMessageToUser(message.UnhandledError + err.Error())
				stop()
			}

			if hiddenStringRepeat == hiddenString {
				mismatch = false
			} else {
				a.transport.SendMessageToUser(message.NewLine + message.PasswordMismatch + message.NewLine + message.LongSeparator)
			}
		}

		err = a.storage.Init(hiddenString)
		if err != nil {
			a.transport.SendMessageToUser(message.UnhandledError + err.Error())
			stop()
		}

		a.transport.SendMessageToUser(message.NewLine + message.LongSeparator + message.NewLine)
	}

	for auth := false; !auth; {
		a.transport.SendMessageToUser(message.AuthWithPassword)

		hiddenString, err = a.transport.GetPasswordHidden()
		if err != nil {
			a.transport.SendMessageToUser(message.UnhandledError + err.Error())
			stop()
		}

		err = a.storage.Connect(hiddenString)
		if err != nil {
			if err.Error() == "file is not a database" {
				a.transport.SendMessageToUser(message.NewLine + message.WrongPassword + message.NewLine + message.LongSeparator)
			} else {
				a.transport.SendMessageToUser(message.UnhandledError + err.Error())
				stop()
			}
		} else {
			auth = true

			a.transport.SendMessageToUser(message.NewLine + message.LongSeparator + message.NewLine)
		}
	}

	conn := a.storage.GetConnection()
	if conn == nil {
		a.transport.SendMessageToUser(message.UnhandledError + err.Error())
		stop()
	}

	model := models.NewEntryModel(conn)

	go func() {
		err = a.transport.StartInputScanner()
		if err != nil {
			a.transport.SendMessageToUser(message.UnhandledError + err.Error())
			stop()
		}
	}()

	ch := *a.transport.GetChannels()

	for {
		a.transport.SendMessageToUser(message.NewLine + message.AwaitInput)
		select {
		case n := <-ch.NavigationCh:
			switch n {
			case GetList:
				a.resetTicker()

				var list []models.Entry

				list, err = model.List()
				if err != nil {
					a.transport.SendMessageToUser(message.UnhandledError + err.Error())
					stop()
				}

				if len(list) == 0 {
					a.transport.SendMessageToUser(message.NoEntriesFound)
				}

				for _, entry := range list {
					a.transport.SendMessageToUser(message.LongSeparator)
					a.transport.SendMessageToUser(fmt.Sprintf("Id: %#v\nService name:%#v\n", entry.Id, entry.ServiceName))
					a.transport.SendMessageToUser(message.LongSeparator)
				}
			case GetEntry:
				a.resetTicker()

				var inputId int

				a.transport.SwitchFocus(true)
				a.transport.SendMessageToUser(message.RequestInputId)

				input := <-ch.InputCh

				inputId, err = strconv.Atoi(input)
				if err != nil {
					a.transport.SendMessageToUser(message.UnhandledError + err.Error())
					stop()
				}

				var entry *models.Entry

				entry, err = model.Get(inputId)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						a.transport.SendMessageToUser(message.LongSeparator + message.WrongId + message.LongSeparator)
						a.transport.SwitchFocus(false)
					} else {
						a.transport.SendMessageToUser(message.UnhandledError + err.Error())
						stop()
					}
				} else {
					a.transport.SendMessageToUser(message.LongSeparator)
					a.transport.SendMessageToUser(
						fmt.Sprintf(
							"Id: %#v\nService name:%#v\nLogin: %#v\nPassword: %#v\n",
							entry.Id,
							entry.ServiceName,
							entry.Login,
							entry.Password,
						),
					)
					a.transport.SendMessageToUser(message.LongSeparator)
				}

				a.transport.SwitchFocus(false)
			case CreateEntry:
				a.resetTicker()
				a.transport.SwitchFocus(true)
				a.transport.SendMessageToUser(message.RequestServiceName)

				serviceName := <-ch.InputCh

				a.transport.SendMessageToUser(message.NewLine + message.RequestServiceLogin)

				login := <-ch.InputCh

				a.transport.SendMessageToUser(message.NewLine + message.RequestServicePassword)

				password := <-ch.InputCh

				err = model.Create(serviceName, login, password)
				if err != nil {
					a.transport.SendMessageToUser(message.UnhandledError + err.Error())
					stop()
				}

				a.transport.SendMessageToUser(
					message.NewLine + message.CreationSuccess + message.NewLine + message.LongSeparator,
				)
				a.transport.SwitchFocus(false)
			case UpdateEntry:
				a.resetTicker()

				var inputId int

				a.transport.SwitchFocus(true)
				a.transport.SendMessageToUser(message.RequestInputId)

				input := <-ch.InputCh

				inputId, err = strconv.Atoi(input)
				if err != nil {
					a.transport.SendMessageToUser(message.UnhandledError + err.Error())
					stop()
				}

				if model.CheckEntryExists(inputId) {
					a.transport.SendMessageToUser(message.NewLine + message.RequestServiceName)

					serviceName := <-ch.InputCh

					a.transport.SendMessageToUser(message.NewLine + message.RequestServiceLogin)

					login := <-ch.InputCh

					a.transport.SendMessageToUser(message.NewLine + message.RequestServicePassword)

					password := <-ch.InputCh

					err = model.Update(inputId, serviceName, login, password)
					if err != nil {
						a.transport.SendMessageToUser(message.UnhandledError + err.Error())
						stop()
					}
				} else {
					a.transport.SendMessageToUser(message.NewLine + message.WrongId)
				}

				a.transport.SwitchFocus(false)
			case DeleteEntry:
				a.resetTicker()

				var inputId int

				a.transport.SendMessageToUser(message.RequestInputId)

				input := <-ch.InputCh

				inputId, err = strconv.Atoi(input)
				if err != nil {
					a.transport.SendMessageToUser(message.UnhandledError + err.Error())
					stop()
				}

				err = model.Delete(inputId)
				if err != nil {
					a.transport.SendMessageToUser(message.UnhandledError + err.Error())
					stop()
				}
			//case ChangeMasterPassword:
			//	var currentPassword string
			//	currentPassword, err = a.transport.GetPasswordHidden()
			//	if err != nil {
			//		err = a.storage.Connect(currentPassword)
			//		if err != nil {
			//			return
			//		}
			//	}
			//	newPassword, err := a.transport.GetPasswordHidden()
			//	if err != nil {
			//		a.transport.SendMessageToUser(message.UnhandledError + errors.New("could not connect to storage").Error())
			//		stop()
			//	}
			//	repeatNewPassword, err := a.transport.GetPasswordHidden()
			//	if err != nil {
			//		a.transport.SendMessageToUser(message.UnhandledError + errors.New("could not connect to storage").Error())
			//		stop()
			//	}
			//	if newPassword != repeatNewPassword {
			//		a.transport.SendMessageToUser(message.PasswordMismatch)
			//	}
			//	stop()
			case Exit:
				stop()
				return
			default:
				a.transport.SendMessageToUser(message.WrongActionId)
			}
		case <-a.ticker.C:
			a.transport.SendMessageToUser(message.TimeoutMessage)
			stop()

			return
		}
	}
}

func NewApp(config *config.Config) (*App, error) {
	//nolint:staticcheck
	s, err := setupStorageInstance(config)
	t, err := setupTransportInstance(config)
	ticker := setupTickerInstance(config)

	if err != nil || s == nil || t == nil {
		return nil, err
	}

	return &App{
		storage:   s,
		config:    config,
		ticker:    ticker,
		transport: t,
	}, nil
}

func setupTransportInstance(c *config.Config) (transport.Core, error) {
	var t transport.Core

	var err error

	switch c.TransportType {
	case "cli":
		t = cli.NewCli()
	default:
		err = fmt.Errorf("unknown transport type: %s", c.TransportType)
	}

	return t, err
}

func setupTickerInstance(c *config.Config) *time.Ticker {
	configTimeout, err := strconv.Atoi(c.Timeout)
	if err != nil {
		configTimeout = 60
	}

	timeout := time.Duration(configTimeout) * time.Minute

	return time.NewTicker(timeout)
}

func (a *App) resetTicker() {
	configTimeout, err := strconv.Atoi(a.config.Timeout)
	if err != nil {
		configTimeout = 60
	}

	timeout := time.Duration(configTimeout) * time.Minute
	a.ticker.Reset(timeout)
}

func setupStorageInstance(c *config.Config) (storage.Storage, error) {
	var s storage.Storage

	var err error

	switch c.StorageType {
	case "sqlite":
		storageAbsPath := utils.GetAbsolutePath(c.StorageFilePath)
		migrationAbsPath := utils.GetAbsolutePath(c.MigrationPath)
		s = database.NewSQLite(c.StorageUsername, storageAbsPath, migrationAbsPath)
	default:
		err = fmt.Errorf("unknown storage type: %s", c.StorageType)
	}

	return s, err
}

func (a *App) GracefulShutdown() error {
	a.ticker.Stop()

	if a.storage != nil {
		err := a.storage.CloseConnections()
		if err != nil {
			return err
		}
	}

	return nil
}
