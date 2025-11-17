package subActions

import (
	"fmt"
	"password_manager/cmd/storage"
	"slices"
)

func GetListFromStorage() error {
	filenames := storage.GetFilenamesInArchive()

	return nil
}

//func CreateProfile() error {
//
//}

func CreateEntry() error {

}

func getStorage() {
	storage.GetStorage()
}
