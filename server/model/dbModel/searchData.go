package dbModel

import "github.com/jinzhu/gorm"

type TopWords struct {

	gorm.Model  `json:"-"`
	Name string `gorm:"unique" json:"name"`
	Type string `json:"type"`


}
