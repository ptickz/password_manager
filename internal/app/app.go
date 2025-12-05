package app

import (
	"database/sql"
	"errors"
	"fmt"
	"password_manager/internal/config"
	"password_manager/internal/storage"
	"password_manager/internal/storage/database"
	"password_manager/internal/storage/models"
	"password_manager/internal/transport"
	"password_manager/internal/transport/cli"
	"password_manager/internal/transport/message"
	"password_manager/internal/utils"
	"strconv"
	"time"
)

const (
	GetList = iota + 1
	GetEntry
	CreateEntry
	UpdateEntry
	DeleteEntry

	Exit = 0
)

type App struct {
	config    *config.Config
	ticker    *time.Ticker
	storage   storage.Storage
	transport transport.Core
}

func (a *App) Run() error {
	exists, err := a.storage.CheckStorageExists()
	if err != nil {
		return err
	}
	if !exists {
		a.transport.SendMessageToUser(message.FirstTimeEnter)
		hiddenString, err := a.transport.GetPasswordHidden()
		if err != nil {
			return err
		}
		//todo
		panic(hiddenString)
	}
	a.transport.SendMessageToUser(message.AuthWithPassword)
	hiddenString, err := a.transport.GetPasswordHidden()
	if err != nil {
		return err
	}
	err = a.storage.Connect(hiddenString)
	if err != nil {
		return err
	}
	conn := a.storage.GetConnection()
	if conn == nil {
		return errors.New("could not connect to storage")
	}
	model := models.NewEntryModel(conn)
	go a.transport.StartInputScanner()
	ch := *a.transport.GetChannels()
	for {
		a.transport.SendMessageToUser(message.NewLine)
		a.transport.SendMessageToUser(message.AwaitInput)
		select {
		case n := <-ch.NavigationCh:
			switch n {
			case GetList:
				var list []models.Entry
				list, err := model.List()
				if err != nil {
					return err
				}
				for _, entry := range list {
					a.transport.SendMessageToUser(message.LongSeparator)
					a.transport.SendMessageToUser(fmt.Sprintf("Id: %#v\nService name:%#v\n", entry.Id, entry.ServiceName))
					a.transport.SendMessageToUser(message.LongSeparator)
				}
			case GetEntry:
				a.transport.SwitchFocus(true)
				a.transport.SendMessageToUser(message.RequestInputId)
				input := <-ch.InputCh
				inputId, err := strconv.Atoi(input)
				if err != nil {
					return err
				}
				entry, err := model.Get(inputId)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						a.transport.SendMessageToUser(message.LongSeparator)
						a.transport.SendMessageToUser(message.WrongId)
						a.transport.SendMessageToUser(message.LongSeparator)

						a.transport.SwitchFocus(false)
					} else {
						return err
					}
				} else {
					a.transport.SendMessageToUser(message.LongSeparator)
					a.transport.SendMessageToUser(fmt.Sprintf("Id: %#v\nService name:%#v\n", entry.Id, entry.ServiceName))
					a.transport.SendMessageToUser(message.LongSeparator)
				}
				a.transport.SwitchFocus(false)
			case CreateEntry:
				a.transport.SwitchFocus(true)
				//todo
				a.transport.SwitchFocus(false)
			case UpdateEntry:
				a.transport.SwitchFocus(true)
				//todo
				a.transport.SwitchFocus(false)
			case DeleteEntry:
				//todo
			case Exit:
				//todo
				return nil
			default:
				a.transport.SendMessageToUser(message.WrongActionId)
			}
		case <-a.ticker.C:
			a.transport.SendMessageToUser(message.TimeoutMessage)
			return nil
		}
	}
	return err
}

func NewApp(config *config.Config) (*App, error) {
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
		err = errors.New(fmt.Sprintf("unknown transport type: %s", c.TransportType))
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

func setupStorageInstance(c *config.Config) (storage.Storage, error) {
	var s storage.Storage
	var err error
	switch c.StorageType {
	case "sqlite":
		storageAbsPath := utils.GetAbsolutePath(c.StorageFilePath)
		migrationAbsPath := utils.GetAbsolutePath(c.MigrationPath)
		s = database.NewSQLite(c.StorageUsername, storageAbsPath, migrationAbsPath)
		break
	default:
		err = errors.New(fmt.Sprintf("unknown storage type: %s", c.StorageType))
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
