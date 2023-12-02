package model

import (
	"time"

	"gorm.io/gorm"
)

type APOD struct {
	gorm.Model     `json:"-"`
	Copyright      string    `json:"copyright"       mapstructure:"copyright"	gorm:"primaryKey"`
	Date           time.Time `json:"date"		     mapstructure:"date"		gorm:"index"`
	Explanation    string    `json:"explanation"     mapstructure:"explanation"`
	MediaType      string    `json:"media_type"      mapstructure:"media_type"`
	ServiceVersion string    `json:"service_version" mapstructure:"service_version"`
	Title          string    `json:"title" 		     mapstructure:"title"`
	URL            string    `json:"url" 			 mapstructure:"url"`
	IDData         string    `json:"id_data" 		 mapstructure:"id_data"`
}
