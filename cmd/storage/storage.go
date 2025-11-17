package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/alexmullins/zip"
)

var MasterPassword string

type Processor interface {
	Init(string) error
	CheckStorageInitiated() (bool, error)
	CheckAccess(string) (bool, error)
	ReadEntries() (*Entries, error)
	ReadEntry(uint) (*Entry, error)
	WriteEntry(*Entry) error
	ModifyEntry(*Entry) error
	DeleteEntry(uint) error
	ChangeMasterPassword(string) error
}
type Storage struct {
	MasterPassword *string
	DirPath        string
	ZipFile        string
	DataFile       string
}

type Entries struct {
	Entries []Entry `json:"entries"`
}

type Entry struct {
	Id          uint   `json:"id"`
	ServiceName string `json:"service_name"`
	Login       string `json:"login"`
	Password    string `json:"password"`
}

func GetStorage() Storage {
	return Storage{
		MasterPassword: nil,
		DirPath:        "./storage",
		ZipFile:        "data.zip",
		DataFile:       "data.json",
	}
}

func (s Storage) Init(password string) error {
	MasterPassword = password
	if _, err := os.Stat(s.DirPath); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(Path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	fZip, err := os.Create(filepath.Join(s.DirPath, s.ZipFile))
	if err != nil {
		log.Fatalln(err)
	}
	zipW := zip.NewWriter(fZip)
	defer func(zipW *zip.Writer) {
		err = zipW.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(zipW)

	w, err := zipW.Encrypt(s.DataFile, password)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(w, bytes.NewReader([]byte(`{"entries":[]}`)))
	if err != nil {
		log.Fatal(err)
	}
	err = zipW.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) CheckStorageInitiated() (bool, error) {
	return false, err
}

func (s Storage) CheckAccess(password string) (bool, error) {
	r, err := zip.OpenReader(filepath.Join(s.DirPath, s.ZipFile))
	if err != nil {
		log.Fatal(err)
	}
	defer func(r *zip.ReadCloser) {
		err = r.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(r)

	var f *zip.File
	for _, v := range r.File {
		if v.Name == s.DataFile {
			//ensure filename is correct
			f = v
		}
	}
	if f != nil {
		f.SetPassword(password)
		_, err = f.Open()
		if err == nil {
			s.MasterPassword = &password
			return true, nil
		}
	}
	return false, err
}

func (s Storage) ReadEntries() (*Entries, error) {
	b, err := s.extractStorageData()
	if err == nil && b != nil {
		var data Entries
		err = json.Unmarshal(b.Bytes(), &data)
		if err == nil {
			return &data, nil
		}
	}
	return nil, err
}

func (s Storage) ReadEntry(id uint) (*Entry, error) {
	r, err := zip.OpenReader(Path + DirectorySeparator + ZipFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func(r *zip.ReadCloser) {
		err = r.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(r)

	for _, f := range r.File {
		f.SetPassword(MasterPassword)
		var rc io.ReadCloser
		rc, err = f.Open()
		if err != nil {
			log.Fatal(err)
		}
		err = rc.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s Storage) CheckStorageEntryExists() bool {
	_, err := os.Stat(Path + DirectorySeparator + ZipFile)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Println(err)
		}
		return false
	}
	return true
}

func (s Storage) WriteEntry(data *Entry) error {
	return nil
}

func (s Storage) ModifyEntry(data *Entry) error {
	return nil
}

func (s Storage) DeleteEntry(id uint) error {
	return nil
}

func (s Storage) ChangeMasterPassword(password string) error {
	return nil
}

func (s Storage) extractStorageData() (*bytes.Buffer, error) {
	if _, err := os.Stat(filepath.Join(s.DirPath, s.ZipFile)); errors.Is(err, os.ErrNotExist) {
		return nil, errors.New("storage file doesnt exist")
	}
	r, err := zip.OpenReader(filepath.Join(s.DirPath, s.ZipFile))
	if err == nil {
		defer func(r *zip.ReadCloser) {
			err = r.Close()
			if err != nil {

				log.Fatal(err)
			}
		}(r)

		f := r.File[0]
		f.SetPassword(*s.MasterPassword)
		var rc io.ReadCloser
		rc, err = f.Open()
		if err == nil {
			defer func(rc io.ReadCloser) {
				err = rc.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(rc)

			var b bytes.Buffer
			_, err = io.Copy(&b, rc)

			return &b, err
		}
	}

	return nil, err
}
