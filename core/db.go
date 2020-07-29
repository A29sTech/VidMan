package core

import (
	"errors"
	"strconv"

	"github.com/jinzhu/gorm"
)

// Some Importent Const ;
const (
	TableName = "videos"
)

// ViData : is a representaion of SQL video table's row ;
type ViData struct { // ViData means VideoData.
	ID       int    `gorm:"AUTOINCREMENT;PRIMARY_KEY" json:"id"`
	Filename string `gorm:"Type:TEXT;UNIQUE;NOT NULL" json:"filename"`
	Title    string `gorm:"Type:TEXT;NOT NULL" json:"title"`
	Subject  string `gorm:"Type:TEXT;NOT NULL" json:"subject"`
	Author   string `gorm:"Type:TEXT" json:"author"`
	Tags     string `gorm:"Type:TEXT" json:"tags"`
	Desc     string `gorm:"Type:TEXT" json:"desc"`
	Indx     int    `gorm:"NOT NULL" json:"indx"`
}

// OpenViDB Database from given path ;
func OpenViDB(path string) (*gorm.DB, error) {
	vdb, err := gorm.Open("sqlite3", path)
	// Check Errors.
	if err != nil {
		return vdb, err
	}
	// As Not Error.
	if !vdb.HasTable(TableName) {
		err = vdb.Table(TableName).CreateTable(&ViData{}).Error
	}
	return vdb, err
}

// Validate : ViData Model ;
func (vidata ViData) Validate(ignoreFileName bool) bool {
	// Impliment ViData Validation ...
	return false
}

// ViDataHasColumName : to check if ViData Tabel Contain Colum Name ? ;
func ViDataHasColumName(columName string) bool {
	switch columName {
	case "filename":
		return true
	case "title":
		return true
	case "subject":
		return true
	case "author":
		return true
	case "tags":
		return true
	case "desc":
		return true
	case "indx":
		return true
	default:
		return false
	}
}

// Form2ViData : FormData to ViData with Validation.
func Form2ViData(formData map[string][]string) (ViData, error) {

	formFieldCount := 6

	vidata := ViData{}

	for key, values := range formData {
		switch key {
		case "title":
			formFieldCount--
			vidata.Title = values[0]
		case "subject":
			formFieldCount--
			vidata.Subject = values[0]
		case "author":
			formFieldCount--
			vidata.Author = values[0]
		case "tags":
			formFieldCount--
			vidata.Tags = values[0]
		case "desc":
			formFieldCount--
			vidata.Desc = values[0]
		case "indx":
			// Convert String To Int ;
			intIndx, err := strconv.Atoi(values[0])
			if err != nil {
				return vidata, errors.New("Can't Parse field 'indx' to int")
			}
			formFieldCount--
			vidata.Indx = intIndx
		}
	}

	// Check If All Fields Are Found or Not ;
	if formFieldCount == 0 {
		return vidata, nil
	}

	// Return error ;
	return vidata, errors.New("All required field not found")

}
