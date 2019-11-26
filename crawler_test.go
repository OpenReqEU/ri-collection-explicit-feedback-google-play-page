package main

import (
	"testing"

	"github.com/OlegSchmidt/soup"
)

var mailformedHTML = `
<html>
  <head>
    <title>Sample "Hello, World" Application</title>
  </head>
  <body>
 	<div class="LXrl4c"> 
 	</div> 
  </body>
  
<html>
`

func TestCrawlAppPage(t *testing.T) {
	for _, document := range []string{
		mailformedHTML,
	} {
		appPage := crawlAppPage(soup.HTMLParse(document), "com.test")
		if appPage.Name != "" {
			t.Errorf("name should be empty")
		}
		if appPage.Category != "" {
			t.Errorf("category should be empty")
		}
		if appPage.USK != "" {
			t.Errorf("usk should be empty")
		}
		if appPage.Price != "" {
			t.Errorf("price should be empty")
		}
		if appPage.Description != "" {
			t.Errorf("description should be empty")
		}
		if len(appPage.WhatsNew) > 0 {
			t.Errorf("whatsNew should be empty")
		}
		if appPage.Rating > 0 {
			t.Errorf("Rating should be empty")
		}
		if appPage.StarsCount > 0 {
			t.Errorf("StarsCount should be empty")
		}
		if appPage.CountPerRating.Five > 0 {
			t.Errorf("CountPerRating.Five should be empty")
		}
		if appPage.EstimatedDownloadNumber > 0 {
			t.Errorf("EstimatedDownloadNumber should be empty")
		}
		if appPage.DeveloperName != "" {
			t.Errorf("DeveloperName should be empty")
		}
		if appPage.TopDeveloper {
			t.Errorf("TopDeveloper should be false")
		}
		if appPage.ContainsAds {
			t.Errorf("ContainsAds should be false")
		}
		if appPage.InAppPurchases {
			t.Errorf("InAppPurchases should be false")
		}
		if appPage.LastUpdate > 0 {
			t.Errorf("LastUpdate should be empty")
		}
		if appPage.RequiresOsVersion != "" {
			t.Errorf("RequiresOsVersion should be empty")
		}
		if appPage.CurrentSoftwareVersion != "" {
			t.Errorf("CurrentSoftwareVersion should be empty")
		}
		if len(appPage.SimilarApps) > 0 {
			t.Errorf("SimilarApps should be empty")
		}
		if len(appPage.Errors) == 0 {
			t.Errorf("there should be errors")
		}
	}
}
