package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"net/http"

	"github.com/OlegSchmidt/soup"
	"github.com/jehiah/go-strftime"
)

const (
	baseURL        = "https://play.google.com"
	baseURLAppPage = baseURL + "/store/apps/details?id="
	lang           = "&hl=en"

	// common html nodes
	a    = "a"
	h1   = "h1"
	h2   = "h2"
	div  = "div"
	span = "span"
	meta = "meta"
	img  = "img"

	// common html attributes
	class    = "class"
	style    = "style"
	href     = "href"
	itemprop = "itemprop"
	alt      = "alt"
	content  = "content"

	// common style attribute values
	styleWidth = "width"

	// CSS classes for finding the right elements
	classMainInformationAppContainer        = "oQ6oV"
	classMainInformationSimilarContainer    = "Ktdaqe"
	classMainInformationApp                 = "rlnrKc"
	classMainInformationSimilar             = "ZmHEEd"
	classMainInformationHeadline            = "Rm6Gwb"
	classMainInformationAdditionalContainer = "IxB2fe"
	classAppPage                            = "LXrl4c"
	classAppCategoryUsk                     = "ZVWMWc"
	classAppRating                          = "BHMmbe"
	classAppStarsCount                      = "EymY4b"
	classAppCountPerRating                  = "VEF2C"
	classAppContainsAds                     = "bSIuKf"
	classAppInAppPurchases                  = "bSIuKf"

	// itemprop values
	itempropAppName         = "name"
	itempropAppCategory     = "genre"
	itempropAppPrice        = "price"
	itempropAppDescription  = "description"
	itempropAppTopDeveloper = "editorsChoiceBadgeUrl"

	// element values as strings
	valueContainsAds                   = "Contains Ads"
	valueInAppPurchases                = "Offers in-app purchases"
	valueRequiresOsVersion             = "Varies with device"
	valueCurrentSoftwareVersionDefault = "unknown"

	// block types
	blockTypeAppName    = "appName"
	blockTypeReview     = "review"
	blockTypeWhatsNew   = "whats new"
	blockTypeAdditional = "additional"

	// errors
	errorPageNotFound = "Page content not found, please update the CSS class in the constant \"classAppPage\""
)

// Crawl the information available on a app page
func Crawl(packageName string) AppPage {
	var appPage AppPage

	document, httpStatus := retrieveDoc(packageName)
	if httpStatus == http.StatusOK {
		appPage = crawlAppPage(document, packageName)
		if appPage.Description == "" && appPage.Name == "" && appPage.DeveloperName == "" {
			// probably captcha
		}
	}

	return appPage
}

// parses the website and returns the DOM struct
func retrieveDoc(packageName string) (soup.Root, int) {
	url := baseURLAppPage + packageName + lang
	var document soup.Root
	httpStatus := http.StatusOK
	// retrieving the html page
	response, soupError := soup.Get(url)
	if soupError != nil {
		fmt.Println("\tcould not reach", url, "because of the following error:")
		fmt.Println(soupError)
		httpStatus = http.StatusBadRequest
	} else {
		// pre-process html
		response = strings.Replace(response, "<br>", "\n", -1)
		response = strings.Replace(response, "<b>", "", -1)
		response = strings.Replace(response, "</b>", "", -1)
		document = soup.HTMLParse(response)
	}

	return document, httpStatus
}

// crawls the page and fills the struct with values
func crawlAppPage(document soup.Root, packageName string) AppPage {
	var lastError error
	appPage := AppPage{}
	appPage.DateCrawled = getCurrentDate()
	appPage.PackageName = packageName
	appPage.Os = getOs()
	if document.Error != nil {
		appPage.Errors = append(appPage.Errors, document.Error.Error())
	} else {
		appPageDocument, appPageDocumentError := getPageDocument(document)
		if appPageDocumentError == nil {
			appPage.Name, lastError = getAppName(appPageDocument)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}

			appPage.Category, lastError = getCategory(appPageDocument)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}

			appPage.USK, lastError = getUsk(appPageDocument)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}

			appPage.Price, appPage.PriceValue, appPage.PriceCurrency, lastError = getPrice(appPageDocument)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}

			appPage.Description, lastError = getDescription(appPageDocument)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}

			appPage.WhatsNew, lastError = getWhatsNew(appPageDocument)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}

			appPage.Rating, lastError = getRating(appPageDocument)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}

			appPage.StarsCount, lastError = getStarsCount(appPageDocument)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}

			appPage.CountPerRating, lastError = getCountPerRating(appPageDocument)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}

			appPage.EstimatedDownloadNumber, lastError = getEstimatedDownloadNumber(appPageDocument)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}

			appPage.DeveloperName, lastError = getDeveloperName(appPageDocument)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}

			appPage.TopDeveloper, lastError = getTopDeveloper(appPageDocument)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}

			appPage.ContainsAds, lastError = getContainsAds(appPageDocument)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}

			appPage.InAppPurchases, lastError = getInAppPurchases(appPageDocument)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}

			appPage.LastUpdate, lastError = getLastUpdate(appPageDocument)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}

			appPage.RequiresOsVersion, lastError = getRequiresOsVersion(appPageDocument)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}

			appPage.CurrentSoftwareVersion, lastError = getCurrentSoftwareVersion(appPageDocument)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}
			// here the whole page is needed, not the app block
			appPage.SimilarApps, lastError = getSimilarApps(document)
			if lastError != nil {
				appPage.Errors = append(appPage.Errors, lastError.Error())
			}
		} else {
			appPage.Errors = append(appPage.Errors, appPageDocumentError.Error())
		}
	}

	return appPage
}

// returns the object of the content area of app information
func getPageDocument(document soup.Root) (soup.Root, error) {
	pageDom := document.Find(div, class, classAppPage)
	var pageDomError error = nil
	if pageDom.Error != nil {
		pageDomError = errors.New(errorPageNotFound)
	}

	return pageDom, pageDomError
}

// returns the 3 main app information blocks : reviews, new functions, additional information
func getMainInformationBlocks(document soup.Root) []soup.Root {
	var informationBlocks []soup.Root

	informationBlockHeadlines := document.FindAll(h2, class, classMainInformationHeadline)
	for position := range informationBlockHeadlines {
		if informationBlockHeadlines[position].Error == nil {
			informationBlocks = append(informationBlocks, informationBlockHeadlines[position].FindParent().FindParent())
		}
	}
	return informationBlocks
}

// returns the requested information block child while using placeholder
func getMainInformationBlockValidated(document soup.Root, position int, property string, blockType string) (soup.Root, error) {
	var informationBlock soup.Root
	var informationBlockError error = nil

	informationBlocks := getMainInformationBlocks(document)
	if len(informationBlocks) >= 3 {
		informationBlockContainer := informationBlocks[position]
		informationBlockChildren := informationBlockContainer.Children()
		if len(informationBlockChildren) >= 2 {
			informationBlock = informationBlockChildren[1]
		} else {
			informationBlockError = errors.New(property + " : main information block \"" + blockType + "\" should contain at least 2 children")
		}
	} else {
		informationBlockError = errors.New(property + " : main information blocks couldn't be found, looking for 2 levels above <h2 class=\"" + classMainInformationHeadline + "\"></h2>")
	}
	return informationBlock, informationBlockError
}

// returns the app information block (basic information at the top of the site)
func getMainInformationBlockApp(document soup.Root, property string) (soup.Root, error) {
	var informationBlockApp soup.Root
	var informationBlockAppError error = nil

	informationBlock := document.Find(div, class, classMainInformationAppContainer)
	if informationBlock.Error == nil {
		informationBlockAppContainer := informationBlock.Find(div, class, classMainInformationApp)
		if informationBlockAppContainer.Error == nil {
			informationBlockApp = informationBlockAppContainer
		} else {
			informationBlockAppError = errors.New(property + " : main information block \"app\" should contain <div class=\"" + classMainInformationApp + "\"></div>")
		}
	} else {
		informationBlockAppError = errors.New(property + " : main information block \"app\" couldn't be found, looking for <div class=\"" + classMainInformationAppContainer + "\"></div>")
	}
	return informationBlockApp, informationBlockAppError
}

// returns the similar information block (list of similar apps or apps from same developer)
func getMainInformationBlockSimilar(document soup.Root, property string) (soup.Root, error) {
	var informationBlockSimilar soup.Root
	var informationBlockSimilarError error = nil

	informationBlock := document.Find(div, class, classMainInformationSimilarContainer)
	if informationBlock.Error == nil {
		informationBlockSimilar = informationBlock
	} else {
		informationBlockSimilarError = errors.New(property + " : main information block \"similar apps\" couldn't be found, looking for <div class=\"" + classMainInformationSimilarContainer + "\"></div>")
	}
	return informationBlockSimilar, informationBlockSimilarError
}

// returns the similar information block (list of similar apps or apps from same developer)
func getMainInformationBlockSimilarChildren(document soup.Root, property string) ([]soup.Root, error) {
	var informationBlockSimilarChildren []soup.Root
	var informationBlockSimilarChildrenError error = nil

	informationBlockSimilar, informationBlockSimilarError := getMainInformationBlockSimilar(document, property)
	if informationBlockSimilarError == nil {
		informationBlockSimilarLink := informationBlockSimilar.Find(a)
		if informationBlockSimilarLink.Error == nil && informationBlockSimilarLink.HasAttribute(href) && informationBlockSimilarLink.GetAttribute(href) != "" {
			similarAppsPageHTML, similarAppsPageHTMLError := soup.Get(baseURL + informationBlockSimilarLink.GetAttribute(href))
			if similarAppsPageHTMLError == nil {
				similarAppsDocument := soup.HTMLParse(similarAppsPageHTML)
				similarAppsAreas := similarAppsDocument.Find(div, class, classMainInformationSimilar)
				if similarAppsAreas.Error == nil {
					informationBlockSimilarChildren = similarAppsAreas.Children()
				}
			}
		}
		if len(informationBlockSimilarChildren) == 0 {
			similarAppsAreas := informationBlockSimilar.Find(div, class, classMainInformationSimilar)
			if similarAppsAreas.Error == nil {
				informationBlockSimilarChildren = similarAppsAreas.Children()
			}
		}
	} else {
		informationBlockSimilarChildrenError = informationBlockSimilarError
	}
	return informationBlockSimilarChildren, informationBlockSimilarChildrenError
}

// returns the review information block
func getMainInformationBlockReview(document soup.Root, property string) (soup.Root, error) {
	return getMainInformationBlockValidated(document, 0, property, blockTypeReview)
}

// returns the whats new information block
func getMainInformationBlockWhatsNew(document soup.Root, property string) (soup.Root, error) {
	return getMainInformationBlockValidated(document, 1, property, blockTypeWhatsNew)
}

// returns the additional information block
func getMainInformationBlockAdditional(document soup.Root, property string) (soup.Root, error) {
	return getMainInformationBlockValidated(document, 2, property, blockTypeAdditional)
}

// returns the additional information block children (updated, size, installs etc.)
func getMainInformationBlockAdditionalChildren(document soup.Root, property string) ([]soup.Root, error) {
	var informationAdditionalChildren []soup.Root
	var informationAdditionalChildrenError error = nil

	informationBlockAdditional, informationBlockAdditionalError := getMainInformationBlockAdditional(document, property)
	if informationBlockAdditionalError == nil {
		informationBlockAdditionalContainer := informationBlockAdditional.Find(div, class, classMainInformationAdditionalContainer)
		if informationBlockAdditionalContainer.Error == nil {
			informationAdditionalChildren = informationBlockAdditionalContainer.Children()
			if len(informationAdditionalChildren) < 11 {
				informationAdditionalChildrenError = errors.New(property + " : <div class=\"" + classMainInformationAdditionalContainer + "\"></div> in main information block \"additional information\" should contain at least 11 children")
			}
		} else {
			informationAdditionalChildrenError = errors.New(property + " : there is no <div class=\"" + classMainInformationAdditionalContainer + "\"></div> in main information block \"additional information\"")
		}
	} else {
		informationAdditionalChildrenError = informationBlockAdditionalError
	}
	return informationAdditionalChildren, informationAdditionalChildrenError
}

// returns specific additional information block child (for example "updated")
func getMainInformationBlockAdditionalChild(document soup.Root, property string, position int) (soup.Root, error) {
	var informationAdditionalChild soup.Root
	var informationAdditionalChildError error = nil

	informationBlockAdditionalChildren, informationBlockAdditionalChildrenError := getMainInformationBlockAdditionalChildren(document, property)
	if informationBlockAdditionalChildrenError == nil {
		informationAdditionalChild = informationBlockAdditionalChildren[position]
	} else {
		informationAdditionalChildError = informationBlockAdditionalChildrenError
	}
	return informationAdditionalChild, informationAdditionalChildError
}

// returns the name of the app
func getAppName(document soup.Root) (string, error) {
	property := "appName"
	appName := ""
	var appNameError error = nil

	informationBlockApp, informationBlockAppError := getMainInformationBlockApp(document, property)
	if informationBlockAppError == nil {
		headline := informationBlockApp.Find(h1, itemprop, itempropAppName)
		if headline.Error == nil {
			headlineSpan := headline.Find(span)
			if headlineSpan.Error == nil {
				appName = headlineSpan.Text()
				if appName == "" {
					appNameError = errors.New(property + " : span inside of <h1 itemprop=\"" + itempropAppName + "\"></h1> is empty")
				}
			} else {
				appNameError = errors.New(property + " : there is no span inside of <h1 itemprop=\"" + itempropAppName + "\"></h1>")
			}
		} else {
			appNameError = errors.New(property + " : there is no <h1 itemprop=\"" + itempropAppName + "\"></h1>")
		}
	} else {
		appNameError = informationBlockAppError
	}

	return appName, appNameError
}

// returns the category of the app
func getCategory(document soup.Root) (string, error) {
	property := "category"
	category := ""
	var categoryError error = nil

	informationBlockApp, informationBlockAppError := getMainInformationBlockApp(document, property)
	if informationBlockAppError == nil {
		categoryElement := informationBlockApp.Find(a, itemprop, itempropAppCategory)
		if categoryElement.Error == nil {
			category = categoryElement.Text()
			if category == "" {
				categoryError = errors.New(property + " : <a itemprop=\"" + itempropAppCategory + "\"></a> is empty")
			}
		} else {
			categoryError = errors.New(property + " : there is no <a itemprop=\"" + itempropAppCategory + "\"></a>")
		}
	} else {
		categoryError = informationBlockAppError
	}

	return category, categoryError
}

// returns the USK of the app
func getUsk(document soup.Root) (string, error) {
	property := "usk"
	usk := ""
	var uskError error = nil

	informationBlockApp, informationBlockAppError := getMainInformationBlockApp(document, property)
	if informationBlockAppError == nil {
		elementCategoryUsk := informationBlockApp.Find(div, class, classAppCategoryUsk)
		if elementCategoryUsk.Error == nil {
			elementCategoryUskChildren := elementCategoryUsk.Children()
			if len(elementCategoryUskChildren) >= 2 {
				blockUsk := elementCategoryUskChildren[1]
				uskImage := blockUsk.Find(img)
				if uskImage.Error == nil {
					usk = uskImage.GetAttribute(alt)
					if usk == "" {
						uskError = errors.New(property + " : the alt of the image of second child of <div class=\"" + classAppCategoryUsk + "\"></div> is empty")
					}
				} else {
					uskError = errors.New(property + " : the second child of <div class=\"" + classAppCategoryUsk + "\"></div> should contain an image some levels lower")
				}
			} else {
				uskError = errors.New(property + " : there should be at least 2 children in <div class=\"" + classAppCategoryUsk + "\"></div>")
			}
		} else {
			uskError = errors.New(property + " : there is no <div class=\"" + classAppCategoryUsk + "\"></div> in main information block \"app\"")
		}
	} else {
		uskError = informationBlockAppError
	}

	return usk, uskError
}

// returns the marker (free or paid), the price of the app and the currency
func getPrice(document soup.Root) (string, float64, string, error) {
	property := "price"
	var price string
	var priceValue float64
	var priceCurrency string
	var priceError error = nil

	informationBlockApp, informationBlockAppError := getMainInformationBlockApp(document, property)
	if informationBlockAppError == nil {
		blockPrice := informationBlockApp.Find(meta, itemprop, itempropAppPrice)
		if blockPrice.Error == nil {
			if blockPrice.HasAttribute(content) {
				attributeContent := blockPrice.GetAttribute(content)
				if attributeContent == "0" {
					price = "free"
					priceValue = 0
					priceCurrency = ""
				} else if attributeContent == "" {
					price = "paid"
					priceValue = 0
					priceCurrency = ""
				} else {
					priceCurrency = string([]rune(attributeContent)[0])

					p := strings.Split(attributeContent, priceCurrency)[1]
					p = strings.Replace(p, ",", ".", -1)
					//fmt.Println(p)
					priceParsed, parseError := strconv.ParseFloat(p, -1)
					if parseError == nil {
						priceValue = priceParsed
					} else {
						priceValue = 0
					}
				}
			} else {
				priceError = errors.New(property + " : <meta itemprop=\"" + itempropAppPrice + "\"></meta> should contain attribute \"" + content + "\"")
			}
		} else {
			priceError = errors.New(property + " : there is no <meta itemprop=\"" + itempropAppPrice + "\"></meta> in main information block \"app\"")
		}
	} else {
		priceError = informationBlockAppError
	}

	return price, priceValue, priceCurrency, priceError
}

// returns the description of the app
func getDescription(doc soup.Root) (string, error) {
	description := ""
	var descriptionError error = nil

	blockDescription := doc.Find(div, itemprop, itempropAppDescription)
	if blockDescription.Error == nil {
		descriptionElement := blockDescription.Find(div)
		if descriptionElement.Error == nil {
			description = descriptionElement.Text()
			if description == "" {
				descriptionError = errors.New("description : the first div below <div itemprop=\"" + itempropAppDescription + "\"></div> is empty")
			}
		} else {
			descriptionError = errors.New("description : <div itemprop=\"" + itempropAppDescription + "\"></div> should contain an div some levels lower")
		}
	} else {
		descriptionError = errors.New("description : there is no <div itemprop=\"" + itempropAppDescription + "\"></meta>")
	}

	return description, descriptionError
}

// returns a list of entries what is new in the app
func getWhatsNew(document soup.Root) ([]string, error) {
	var whatsNew []string
	var whatsNewError error = nil

	informationBlock, informationBlockError := getMainInformationBlockWhatsNew(document, "whatsNew")
	if informationBlockError == nil {
		whatsNewContainer := informationBlock.Find(span)
		if whatsNewContainer.Error == nil {
			whatsNewElements := whatsNewContainer.Children(true)
			for position := range whatsNewElements {
				if whatsNewElements[position].NodeValue != "br" {
					whatsNew = append(whatsNew, whatsNewElements[position].NodeValue)
				}
			}
		} else {
			whatsNewError = errors.New("whatsNew : second child of main information block should contain a span at some level below")
		}
	} else {
		whatsNewError = informationBlockError
	}

	return whatsNew, whatsNewError
}

// returns the star rating of the app
func getRating(document soup.Root) (float64, error) {
	var rating float64 = 0
	var ratingError error = nil

	informationBlockReview, informationBlockReviewError := getMainInformationBlockReview(document, "rating")
	if informationBlockReviewError == nil {
		ratingContainer := informationBlockReview.Find(div, class, classAppRating)
		if ratingContainer.Error == nil {
			ratingString := ratingContainer.Text()
			if ratingString != "" {
				ratingFloat, parseError := strconv.ParseFloat(ratingString, 64)
				if parseError == nil {
					rating = ratingFloat
				} else {
					ratingError = errors.New("rating : <div class=\"" + classAppRating + "\"></div> is not a float and contains \"" + ratingString + "\"")
				}
			} else {
				ratingError = errors.New("rating : <div class=\"" + classAppRating + "\"></div> is empty")
			}
		} else {
			ratingError = errors.New("rating : there is no <div class=\"" + classAppRating + "\"></div> inside of main information block \"reviews\"")
		}
	} else {
		ratingError = informationBlockReviewError
	}

	return rating, ratingError
}

// returns the amount of the stars for the app
func getStarsCount(document soup.Root) (int64, error) {
	var starsCount int64 = 0
	var starsCountError error = nil

	informationBlockReview, informationBlockReviewError := getMainInformationBlockReview(document, "starsCount")
	if informationBlockReviewError == nil {
		starsCountContainer := informationBlockReview.Find(span, class, classAppStarsCount)
		if starsCountContainer.Error == nil {
			starsCountContainerChildren := starsCountContainer.Children()
			if len(starsCountContainerChildren) >= 2 {
				starsCountString := starsCountContainerChildren[1].Text()
				if starsCountString != "" {
					starsCountStringReplaced := strings.Replace(starsCountString, ",", "", -1)
					starsCountStringReplaced = strings.Replace(starsCountStringReplaced, ".", "", -1)
					starsCountNumber, parseError := strconv.ParseInt(starsCountStringReplaced, 0, 64)
					if parseError == nil {
						starsCount = starsCountNumber
					} else {
						starsCountError = errors.New("starsCount : <span class=\"" + classAppStarsCount + "\"></span> is not an integer and contains \"" + starsCountString + "\"")
					}
				} else {
					starsCountError = errors.New("starsCount : <span class=\"" + classAppStarsCount + "\"></span> is empty")
				}
			} else {
				starsCountError = errors.New("starsCount : <span class=\"" + classAppStarsCount + "\"></span> should contain at least 2 children")
			}
		} else {
			starsCountError = errors.New("starsCount : there is no <span class=\"" + classAppStarsCount + "\"></span> inside of main information block \"reviews\"")
		}
	} else {
		starsCountError = informationBlockReviewError
	}

	return starsCount, starsCountError
}

// returns the distribution in percentage between the amount of stars you can give based on the current rating
func getCountPerRating(document soup.Root) (StarCountPerRating, error) {
	countPerRating := StarCountPerRating{}
	var countPerRatingError error = nil

	informationBlockReview, informationBlockReviewError := getMainInformationBlockReview(document, "countPerRating")
	if informationBlockReviewError == nil {
		countPerRatingContainer := informationBlockReview.Find(div, class, classAppCountPerRating)
		if countPerRatingContainer.Error == nil {
			countPerRatingElements := countPerRatingContainer.Children()
			if len(countPerRatingElements) >= 5 {
				for position := range countPerRatingElements {
					countPerRatingElementChildren := countPerRatingElements[position].Children()
					if len(countPerRatingElementChildren) >= 2 {
						rating := countPerRatingElementChildren[0].Text()
						width := 0
						if rating != "" {
							if countPerRatingElementChildren[1].HasAttribute(style) {
								style := countPerRatingElementChildren[1].GetAttribute("style")
								styleParts := strings.Split(style, ";")
								for positionStyle := range styleParts {
									styleDefinition := AttributeStyle{}.fill(styleParts[positionStyle])
									if styleDefinition.Name == styleWidth {
										width = styleDefinition.getValueAsInt()
										switch rating {
										case "1":
											countPerRating.One = width
											break
										case "2":
											countPerRating.Two = width
											break
										case "3":
											countPerRating.Three = width
											break
										case "4":
											countPerRating.Four = width
											break
										case "5":
											countPerRating.Five = width
											break
										}
										break
									}
								}
							}
						} else {
							countPerRatingError = errors.New("countPerRating : final element doesn't contain a rating")
						}
					} else {
						countPerRatingError = errors.New("countPerRating : child of <div class=\"" + classAppCountPerRating + "\"></div> in main information block \"reviews\" should have at least 2 children")
					}
				}
			} else {
				countPerRatingError = errors.New("countPerRating : <div class=\"" + classAppCountPerRating + "\"></div> in main information block \"reviews\" should have at least 5 children")
			}
		} else {
			countPerRatingError = errors.New("countPerRating : there is no <div class=\"" + classAppCountPerRating + "\"></div> in main information block \"reviews\"")
		}
	} else {
		countPerRatingError = informationBlockReviewError
	}

	ratingWidthSum := float64(countPerRating.One + countPerRating.Two + countPerRating.Three + countPerRating.Four + countPerRating.Five)
	countPerRating.One = int(math.Round(float64(countPerRating.One) / ratingWidthSum * 100))
	countPerRating.Two = int(math.Round(float64(countPerRating.Two) / ratingWidthSum * 100))
	countPerRating.Three = int(math.Round(float64(countPerRating.Three) / ratingWidthSum * 100))
	countPerRating.Four = int(math.Round(float64(countPerRating.Four) / ratingWidthSum * 100))
	countPerRating.Five = int(math.Round(float64(countPerRating.Five) / ratingWidthSum * 100))

	return countPerRating, countPerRatingError
}

// returns the estimated number of downloads of the app
func getEstimatedDownloadNumber(document soup.Root) (int64, error) {
	var estimatedDownloadNumber int64 = 0
	var estimatedDownloadNumberError error = nil
	childPosition := 2

	informationBlockAdditionalChild, informationBlockAdditionalChildError := getMainInformationBlockAdditionalChild(document, "estimatedDownloadNumber", childPosition)
	if informationBlockAdditionalChildError == nil {
		estimatedDownloadNumberElement := informationBlockAdditionalChild.FindAll(span)
		if len(estimatedDownloadNumberElement) > 0 {
			estimatedDownloadNumberString := estimatedDownloadNumberElement[len(estimatedDownloadNumberElement)-1].Text()
			if estimatedDownloadNumberString != "" {
				estimatedDownloadNumberStringReplaced := strings.Replace(estimatedDownloadNumberString, ",", "", -1)
				estimatedDownloadNumberStringReplaced = strings.Replace(estimatedDownloadNumberStringReplaced, "+", "", -1)
				estimatedDownloadNumberInt, parseError := strconv.ParseInt(strings.TrimSpace(estimatedDownloadNumberStringReplaced), 0, 64)
				if parseError == nil {
					estimatedDownloadNumber = estimatedDownloadNumberInt
				} else {
					estimatedDownloadNumberError = errors.New("estimatedDownloadNumber : final element doesn't contain a number of downloads, it contains : \"" + estimatedDownloadNumberString + "\"")
				}
			} else {
				estimatedDownloadNumberError = errors.New("estimatedDownloadNumber : final element doesn't contain a number of downloads")
			}
		} else {
			estimatedDownloadNumberError = errors.New("estimatedDownloadNumber : " + string(childPosition+1) + ". child of <div class=\"" + classMainInformationAdditionalContainer + "\"></div> in main information block \"additional information\" should contain at least one span at lower levels")
		}
	} else {
		estimatedDownloadNumberError = informationBlockAdditionalChildError
	}

	return estimatedDownloadNumber, estimatedDownloadNumberError
}

// returns the link to the developer website
func getDeveloperName(document soup.Root) (string, error) {
	developerName := ""
	var developerNameError error = nil

	informationBlockAdditionalChildren, informationBlockAdditionalChildrenError := getMainInformationBlockAdditionalChildren(document, "developerName")
	if informationBlockAdditionalChildrenError == nil {
		developerNameLink := informationBlockAdditionalChildren[len(informationBlockAdditionalChildren)-1].Find(a)
		if developerNameLink.Error == nil {
			if developerNameLink.HasAttribute(href) == true && developerNameLink.GetAttribute(href) != "" {
				developerName = developerNameLink.GetAttribute("href")
			} else {
				developerNameError = errors.New("developerName : the link in <div class=\"" + classMainInformationAdditionalContainer + "\"></div> in main information block \"additional information\" doesn't have \"href\" Attribute or its empty")
			}
		} else {
			developerNameError = errors.New("developerName : <div class=\"" + classMainInformationAdditionalContainer + "\"></div> in main information block \"additional information\" should contain a link at some lower levels")
		}
	} else {
		developerNameError = informationBlockAdditionalChildrenError
	}
	return developerName, developerNameError
}

// returns the badge if the app was marked as "redaction suggestion"
func getTopDeveloper(document soup.Root) (bool, error) {
	topDeveloper := false
	var topDeveloperError error = nil

	informationBlockApp, informationBlockAppError := getMainInformationBlockApp(document, "topDeveloper")
	if informationBlockAppError == nil {
		topDeveloper = informationBlockApp.Find(meta, itemprop, itempropAppTopDeveloper).Error == nil
	} else {
		topDeveloperError = informationBlockAppError
	}
	return topDeveloper, topDeveloperError
}

// returns if the app has advertisements or not
func getContainsAds(document soup.Root) (bool, error) {
	containsAds := false
	var containsAdsError error = nil

	informationBlockApp, informationBlockAppError := getMainInformationBlockApp(document, "containsAds")
	if informationBlockAppError == nil {
		containsAdsBlock := informationBlockApp.Find(div, class, classAppContainsAds)
		if containsAdsBlock.Error == nil {
			containsAdsBlockChildren := containsAdsBlock.Children(true)
			if len(containsAdsBlockChildren) == 0 && containsAdsBlock.Text() == valueContainsAds {
				containsAds = true
			} else {
				for position := range containsAdsBlockChildren {
					if containsAdsBlockChildren[position].NodeValue == valueContainsAds {
						containsAds = true
						break
					}
				}
			}
		} else {
			containsAdsError = errors.New("containsAds : there is no <div class=\"" + classAppContainsAds + "\"></div> in main information block \"app\"")
		}
	} else {
		containsAdsError = informationBlockAppError
	}
	return containsAds, containsAdsError
}

// returns if the app offers purchases
func getInAppPurchases(document soup.Root) (bool, error) {
	inAppPurchases := false
	var inAppPurchasesError error = nil

	informationBlockApp, informationBlockAppError := getMainInformationBlockApp(document, "inAppPurchases")
	if informationBlockAppError == nil {
		inAppPurchasesBlock := informationBlockApp.Find(div, class, classAppInAppPurchases)
		if inAppPurchasesBlock.Error == nil {
			inAppPurchasesBlockChildren := inAppPurchasesBlock.Children(true)
			if len(inAppPurchasesBlockChildren) == 0 && inAppPurchasesBlock.Text() == valueInAppPurchases {
				inAppPurchases = true
			} else {
				for position := range inAppPurchasesBlockChildren {
					if inAppPurchasesBlockChildren[position].NodeValue == valueInAppPurchases {
						inAppPurchases = true
						break
					}
				}
			}
		} else {
			inAppPurchasesError = errors.New("inAppPurchases : there is no <div class=\"" + classAppInAppPurchases + "\"></div> in main information block \"app\"")
		}
	} else {
		inAppPurchasesError = informationBlockAppError
	}
	return inAppPurchases, inAppPurchasesError
}

// return the date of last update
func getLastUpdate(document soup.Root) (int64, error) {
	var lastUpdate int64 = 0
	var lastUpdateError error = nil
	childPosition := 0

	informationBlockAdditionalChild, informationBlockAdditionalChildError := getMainInformationBlockAdditionalChild(document, "lastUpdate", childPosition)
	if informationBlockAdditionalChildError == nil {
		lastUpdateElements := informationBlockAdditionalChild.FindAll(span)
		if len(lastUpdateElements) > 0 {
			lastUpdateString := lastUpdateElements[len(lastUpdateElements)-1].Text()
			lastUpdateString = strings.TrimSpace(lastUpdateString)
			if lastUpdateString != "" {
				lastUpdateObject, lastUpdateObjectError := time.Parse("January 2, 2006", lastUpdateString)
				if lastUpdateObjectError == nil {
					lastUpdateFormatted := strftime.Format("%Y%m%d", lastUpdateObject)
					lastUpdateNumber, lastUpdateNumberError := strconv.ParseInt(lastUpdateFormatted, 0, 64)
					if lastUpdateNumberError == nil {
						lastUpdate = lastUpdateNumber
					} else {
						lastUpdateError = errors.New("lastUpdate : content of last span of " + string(childPosition+1) + ". child of <div class=\"" + classMainInformationAdditionalContainer + "\"></div> in main information block \"additional information\" couldn't be converted into a number")
					}
				} else {
					lastUpdateError = errors.New("lastUpdate : content of last span of " + string(childPosition+1) + ". child of <div class=\"" + classMainInformationAdditionalContainer + "\"></div> in main information block \"additional information\" doesn't contain a date")
				}
			} else {
				lastUpdateError = errors.New("lastUpdate : last span of " + string(childPosition+1) + ". child of <div class=\"" + classMainInformationAdditionalContainer + "\"></div> in main information block \"additional information\" should contain a date but is empty")
			}
		} else {
			lastUpdateError = errors.New("lastUpdate : " + string(childPosition+1) + ". child of <div class=\"" + classMainInformationAdditionalContainer + "\"></div> in main information block \"additional information\" should contain at least one span at lower levels")
		}
	} else {
		lastUpdateError = informationBlockAdditionalChildError
	}

	return lastUpdate, lastUpdateError
}

// returns the needed operation system for the app
func getOs() string {
	return "ANDROID"
}

// returns the required version of operating system
func getRequiresOsVersion(document soup.Root) (string, error) {
	requiresOsVersion := ""
	var requiresOsVersionError error = nil
	childPosition := 4

	informationBlockAdditionalChild, informationBlockAdditionalChildError := getMainInformationBlockAdditionalChild(document, "requiresOsVersion", childPosition)
	if informationBlockAdditionalChildError == nil {
		requiresOsVersionElements := informationBlockAdditionalChild.FindAll(span)
		if len(requiresOsVersionElements) > 0 {
			requiresOsVersionString := requiresOsVersionElements[len(requiresOsVersionElements)-1].Text()
			requiresOsVersionString = strings.TrimSpace(requiresOsVersionString)
			if requiresOsVersionString != "" {
				if requiresOsVersionString == valueRequiresOsVersion {
					requiresOsVersion = requiresOsVersionString
				} else {
					osVersionParts := strings.Fields(requiresOsVersionString)
					osVersion := osVersionParts[0]
					if len(osVersionParts) > 1 {
						osVersion += "+"
					}
					requiresOsVersion = osVersion
				}
			} else {
				requiresOsVersionError = errors.New("requiresOsVersion : last span of " + string(childPosition+1) + ". child of <div class=\"" + classMainInformationAdditionalContainer + "\"></div> in main information block \"additional information\" should contain a string but is empty")
			}
		} else {
			requiresOsVersionError = errors.New("requiresOsVersion : " + string(childPosition+1) + ". child of <div class=\"" + classMainInformationAdditionalContainer + "\"></div> in main information block \"additional information\" should contain at least one span at lower levels")
		}
	} else {
		requiresOsVersionError = informationBlockAdditionalChildError
	}

	return requiresOsVersion, requiresOsVersionError
}

// returns the current version of the app
func getCurrentSoftwareVersion(document soup.Root) (string, error) {
	currentSoftwareVersion := ""
	var currentSoftwareVersionError error = nil
	childPosition := 3

	informationBlockAdditionalChild, informationBlockAdditionalChildError := getMainInformationBlockAdditionalChild(document, "currentSoftwareVersion", childPosition)
	if informationBlockAdditionalChildError == nil {
		currentSoftwareVersionElements := informationBlockAdditionalChild.FindAll(span)
		if len(currentSoftwareVersionElements) > 0 {
			requiresOsVersionString := currentSoftwareVersionElements[len(currentSoftwareVersionElements)-1].Text()
			currentSoftwareVersion = strings.TrimSpace(requiresOsVersionString)
			if requiresOsVersionString == "" {
				currentSoftwareVersion = valueCurrentSoftwareVersionDefault
			}
		} else {
			currentSoftwareVersionError = errors.New("requiresOsVersion : " + string(childPosition+1) + ". child of <div class=\"" + classMainInformationAdditionalContainer + "\"></div> in main information block \"additional information\" should contain at least one span at lower levels")
		}
	} else {
		currentSoftwareVersionError = informationBlockAdditionalChildError
	}

	return currentSoftwareVersion, currentSoftwareVersionError
}

// returns
func getSimilarApps(document soup.Root) ([]string, error) {
	var similarApps []string
	var similarAppsError error = nil

	similarAppElements, similarAppElementsError := getMainInformationBlockSimilarChildren(document, "similarApps")
	if similarAppElementsError == nil {
		for position := range similarAppElements {
			similarAppLink := similarAppElements[position].Find("a")
			if similarAppLink.Error == nil {
				if similarAppLink.HasAttribute(href) == true {
					similarAppLinkParts := strings.Split(similarAppLink.GetAttribute(href), "?")
					if len(similarAppLinkParts) == 2 {
						getParameter := similarAppLinkParts[1]
						parameters := strings.Split(getParameter, "&")
						for parameterPosition := range parameters {
							parameterParts := strings.Split(parameters[parameterPosition], "=")
							if len(parameterParts) == 2 {
								if parameterParts[0] == "id" {
									similarApps = append(similarApps, parameterParts[1])
								}
							}
						}
					} else {
						similarAppsError = errors.New("similarApps : \"href\" attribute of the link to the app suggestion doesn't contain GET-parameters")
					}
				} else {
					similarAppsError = errors.New("similarApps : link to the app suggestion doesn't contain a \"href\" attribute")
				}
			} else {
				similarAppsError = errors.New("similarApps : app suggestion doesn't contain a link to the app")
			}
		}
	} else {
		similarAppsError = similarAppElementsError
	}
	return similarApps, similarAppsError
}

// returns the current date as integer
func getCurrentDate() int64 {
	var currentDate int64 = 0
	dateNow := time.Now()
	currentDateFormatted := strftime.Format("%Y%m%d", dateNow)
	formattedValue, parseError := strconv.Atoi(currentDateFormatted)
	if parseError != nil {
		fmt.Println("ERR", parseError)
		panic(parseError)
	}
	currentDate = int64(formattedValue)

	return currentDate
}
