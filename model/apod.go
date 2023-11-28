package model

import (
	"gorm.io/gorm"
)

type APOD struct {
	gorm.Model
	Copyright      string `json:"copyright"       mapstructure:"copyright"`
	// mapstucture безопасно не может распарсить time.Time из http запроса
	// поэтому здесь будет string
	// нужно отдельно создавать для запроса структуру и для базы данных
	Date           string `json:"date"		      mapstructure:"date"`
	Explanation    string `json:"explanation"     mapstructure:"explanation"`
	MediaType      string `json:"media_type"      mapstructure:"media_type"`
	ServiceVersion string `json:"service_version" mapstructure:"service_version"`
	Title          string `json:"title" 		  mapstructure:"title"`
	URL            string `json:"url" 			  mapstructure:"url"`
	IDData         string `json:"id_data" 		  mapstructure:"id_data"`
}
