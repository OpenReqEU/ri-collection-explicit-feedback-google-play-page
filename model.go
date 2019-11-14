package main

import (
	"strconv"
	"strings"
)

// AppPage model
type AppPage struct {
	Name                    string             `json:"name" bson:"name"`
	PackageName             string             `json:"package_name" bson:"package_name"`
	DateCrawled             int64              `json:"date_crawled" bson:"date_crawled"`
	Category                string             `json:"category" bson:"category"`
	USK                     string             `json:"usk" bson:"usk"`
	Price                   string             `json:"price" bson:"price"`
	PriceValue              float64            `json:"price_value" bson:"price_value"`
	PriceCurrency           string             `json:"price_currency" bson:"price_currency"`
	Description             string             `json:"description" bson:"description"`
	WhatsNew                []string           `json:"whats_new" bson:"whats_new"`
	Rating                  float64            `json:"rating" bson:"rating"`
	StarsCount              int64              `json:"stars_count" bson:"stars_count"`
	CountPerRating          StarCountPerRating `json:"count_per_rating" bson:"count_per_rating"`
	EstimatedDownloadNumber int64              `json:"estimated_download_number" bson:"estimated_download_number"`
	DeveloperName           string             `json:"developer" bson:"developer"`
	TopDeveloper            bool               `json:"top_developer" bson:"top_developer"`
	ContainsAds             bool               `json:"contains_ads" bson:"contains_ads"`
	InAppPurchases          bool               `json:"in_app_purchase" bson:"in_app_purchase"`
	LastUpdate              int64              `json:"last_update" bson:"last_update"`
	Os                      string             `json:"os" bson:"os"`
	RequiresOsVersion       string             `json:"requires_os_version" bson:"requires_os_version"`
	CurrentSoftwareVersion  string             `json:"current_software_version" bson:"current_software_version"`
	SimilarApps             []string           `json:"similar_apps" bson:"similar_apps"`
	Errors                  []string           `json:"errors" bson:"errors"`
}

// StarCountPerRating model
type StarCountPerRating struct {
	Five  int `json:"5"`
	Four  int `json:"4"`
	Three int `json:"3"`
	Two   int `json:"2"`
	One   int `json:"1"`
}

// style property model
type AttributeStyle struct {
	Name  string
	Value string
	Unit  string
}

// fills an object with given definition
func (style AttributeStyle) fill(definition string) AttributeStyle {
	definition = strings.TrimSpace(definition)
	definitionParts := strings.Split(definition, ":")
	value := definitionParts[1]
	unit := ""
	units := [4]string{"%", "px", "rem", "em"}
	for position := range units {
		if len(value) > len(strings.Replace(value, units[position], "", -1)) {
			value = strings.Replace(value, units[position], "", -1)
			unit = units[position]
			break
		}
	}
	style.Name = strings.Trim(definitionParts[0], " ")
	style.Value = strings.TrimSpace(value)
	style.Unit = unit

	return style
}

// converts the value of the struct into int and returns it
func (style AttributeStyle) getValueAsInt() int {
	valueAsInt, success := strconv.Atoi(style.Value)
	if success != nil {
		valueAsInt = 0
	}
	return valueAsInt
}
