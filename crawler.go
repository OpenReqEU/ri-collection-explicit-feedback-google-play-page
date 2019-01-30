package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"net/http"

	"os"

	"github.com/anaskhan96/soup"
	"github.com/jehiah/go-strftime"
)

const (
	baseURL        = "https://play.google.com"
	baseURLAppPage = baseURL + "/store/apps/details?id="
	lang           = "&hl=en"

	// common html attributes
	div      = "div"
	span     = "span"
	class    = "class"
	meta     = "meta"
	itemprop = "itemprop"
)

// Crawl the information available on a app page
func Crawl(packageName string) AppPage {
	var appPage AppPage

	doc, httpStatus := retrieveDoc(packageName)
	if httpStatus == http.StatusOK {
		appPage = crawlAppPage(doc, packageName)
		if appPage.Description == "" &&
			appPage.Name == "" &&
			appPage.DeveloperName == "" {

			return appPage // probably captcha
		}
	}

	return appPage
}

func retrieveDoc(packageName string) (soup.Root, int) {
	url := baseURLAppPage + packageName + lang
	//fmt.Println("start to crawl:", packageName, "from", url)

	httpStatus := http.StatusOK
	// retrieving the html page
	resp, err := soup.Get(url)
	if err != nil {
		fmt.Println("\tcould not reach", url, "because of the following error:")
		fmt.Println(err)
		httpStatus = http.StatusBadRequest
	}

	// pre-process html
	resp = strings.Replace(resp, "<br>", "\n", -1)
	resp = strings.Replace(resp, "<b>", "", -1)
	resp = strings.Replace(resp, "</b>", "", -1)
	doc := soup.HTMLParse(resp)

	// check if the captcha came up
	captcha := doc.Find("body").Attrs()["onload"]
	if captcha == "e=document.getElementById('captcha');if(e){e.focus();}" {
		fmt.Printf("%s QUIT PROGRAMM: captcha needed\n", packageName)
		os.Exit(0)
	}

	return doc, httpStatus
}

func crawlAppPage(doc soup.Root, packageName string) AppPage {
	fmt.Println(doc)
	appPage := AppPage{}
	appPage.Name = getAppName(doc)
	appPage.PackageName = packageName
	appPage.DateCrawled = getCurrentDate()
	appPage.Category = getCategory(doc)
	appPage.USK = getUsk(doc)
	appPage.Price, appPage.PriceValue, appPage.PriceCurrency = getPrice(doc)
	appPage.Description = getDescription(doc)
	appPage.WhatsNew = getWhatsNew(doc)
	appPage.Rating = getRating(doc)
	appPage.StarsCount = getStarsCount(doc)
	// appPage.CountPerRating = getCountPerRating(doc)
	appPage.EstimatedDownloadNumber = getEstimatedDownloadNumber(doc)
	appPage.DeveloperName = getDeveloperName(doc)
	appPage.TopDeveloper = getTopDeveloper(doc)
	appPage.ContainsAds = getContainsAds(doc)
	appPage.InAppPurchases = getInAppPurchases(doc)
	appPage.LastUpdate = getLastUpdate(doc)
	appPage.Os = getOs()
	appPage.RequiresOsVersion = getRequiresOsVersion(doc)
	appPage.CurrentSoftwareVersion = getCurrentSoftwareVersion(doc)
	appPage.SimilarApps = getSimilarApps(doc)

	return appPage
}

func getAppName(doc soup.Root) string {
	return doc.Find("h1", itemprop, "name").Find(span).Text()
}

func getPackageName(doc soup.Root) string {
	return ""
}

func getCategory(doc soup.Root) string {
	return doc.Find("a", itemprop, "genre").Text()
}

func getUsk(doc soup.Root) string {
	return doc.Find("meta", itemprop, "contentRating").Attrs()["content"]
}

func getPrice(doc soup.Root) (string, float64, string) {
	var price string
	var priceValue float64
	var priceCurrency string

	htmlPrice := doc.Find(meta, itemprop, "price")
	if htmlPrice.Attrs()["content"] == "0" {
		price = "free"
		priceValue = 0
		priceCurrency = ""
	} else {
		price = "paid"
		priceRaw := htmlPrice.Attrs()["content"]
		if priceRaw != "" {
			priceCurrency = string([]rune(priceRaw)[0])

			p := strings.Split(priceRaw, priceCurrency)[1]
			p = strings.Replace(p, ",", ".", -1)
			//fmt.Println(p)
			val, err := strconv.ParseFloat(p, -1)
			if err != nil {
				fmt.Println(err)
				priceValue = 0
			} else {
				priceValue = val
			}
		} else {
			price = "paid"
			priceValue = 0
			priceCurrency = ""
		}
	}

	return price, priceValue, priceCurrency
}

func getDescription(doc soup.Root) string {
	return doc.Find(div, itemprop, "description").Find(div).Text()
}

func getWhatsNew(doc soup.Root) []string {
	items := doc.FindAll(div, class, "recent-change")
	var whatsNew []string
	for key := range items {
		whatsNew = append(whatsNew, items[key].Text())
	}

	return whatsNew
}

func getRating(doc soup.Root) float64 {
	i, err := strconv.ParseFloat(doc.Find("meta", itemprop, "ratingValue").Attrs()["content"], 64)
	if err != nil {
		//fmt.Println("getRating", err)
		return 0
	}

	return i
}

func getStarsCount(doc soup.Root) int64 {
	i, err := strconv.ParseInt(doc.Find("meta", itemprop, "reviewCount").Attrs()["content"], 0, 64)
	if err != nil {
		//fmt.Println("getStarsCount", err)
		return 0
	}

	return int64(i)
}

func getCountPerRating(doc soup.Root) StarCountPerRating {
	s := StarCountPerRating{}
	re := regexp.MustCompile("[^0-9]+")

	fmt.Println(doc.Find(span, class, "mMF0fd").Attrs()["title"])
	fiveStarRaw := doc.Find(span, class, "mMF0fd").Attrs()["title"]
	i, err := strconv.Atoi(re.ReplaceAllString(fiveStarRaw, ""))
	if err != nil {
		//fmt.Println("fiveStarRaw", err)
		s.Five = 0
	} else {
		s.Five = i
	}

	fourStarRaw := doc.Find(div, class, "rating-bar-container four").Find(span, class, "bar-number").Attrs()["aria-label"]
	i, err = strconv.Atoi(re.ReplaceAllString(fourStarRaw, ""))
	if err != nil {
		//fmt.Println("fourStarRaw", err)
		s.Four = 0
	} else {
		s.Four = i
	}

	threeStarRaw := doc.Find(div, class, "rating-bar-container three").Find(span, class, "bar-number").Attrs()["aria-label"]
	i, err = strconv.Atoi(re.ReplaceAllString(threeStarRaw, ""))
	if err != nil {
		//fmt.Println("threeStarRaw", err)
		s.Three = 0
	} else {
		s.Three = i
	}

	twoStarRaw := doc.Find(div, class, "rating-bar-container two").Find(span, class, "bar-number").Attrs()["aria-label"]
	i, err = strconv.Atoi(re.ReplaceAllString(twoStarRaw, ""))
	if err != nil {
		//fmt.Println("twoStarRaw", err)
		s.Two = 0
	} else {
		s.Two = i
	}

	oneStarRaw := doc.Find(div, class, "rating-bar-container one").Find(span, class, "bar-number").Attrs()["aria-label"]
	i, err = strconv.Atoi(re.ReplaceAllString(oneStarRaw, ""))
	if err != nil {
		//fmt.Println("oneStarRaw", err)
		s.One = 0
	} else {
		s.One = i
	}

	return s
}

func getEstimatedDownloadNumber(doc soup.Root) int64 {
	raw := strings.Replace(doc.FindAll(span, class, "htlgb")[5].Text(), ",", "", -1)
	split := strings.Split(raw, "+")
	if len(split) < 2 {
		return int64(0)
	}

	min, err := strconv.Atoi(strings.TrimSpace(split[0]))
	if err != nil {
		fmt.Println("getEstimatedDownloadNumber", err)
		min = 0
	}
	return int64(min)
}

func getDeveloperName(doc soup.Root) string {
	return doc.Find(span, itemprop, "name").Text()
}

func getTopDeveloper(doc soup.Root) bool {
	return "" != doc.Find(span, class, "badge-title").Text()
}

func getContainsAds(doc soup.Root) bool {
	return "" != doc.Find(span, class, "ads-supported-label-msg").Text()
}

func getInAppPurchases(doc soup.Root) bool {
	return "" != doc.Find(div, class, "inapp-msg").Text()
}

func getLastUpdate(doc soup.Root) int64 {
	unFormattedDate := doc.Find(div, itemprop, "datePublished").Text()
	t, err := time.Parse("January 2, 2006", unFormattedDate)
	if err != nil {
		return -1
	}

	s := strftime.Format("%Y%m%d", t)
	val, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("getLastUpdate", err)
		return -1
	}

	return int64(val)
}

func getOs() string {
	return "ANDROID"
}

func getRequiresOsVersion(doc soup.Root) string {
	raw := doc.Find(div, itemprop, "operatingSystems").Text()
	if raw == "" {
		return "unkown"
	}
	split := strings.Fields(raw)
	requiredOs := split[0]
	if len(split) > 1 {
		requiredOs += "+"
	}
	return requiredOs
}

func getCurrentSoftwareVersion(doc soup.Root) string {
	var version string
	version = doc.Find(div, itemprop, "softwareVersion").Text()

	if "" == version {
		return "unknown"
	}

	return strings.Replace(version, " ", "", -1)
}

func getSimilarApps(doc soup.Root) []string {
	var similarApps []string

	suffix := doc.Find("a", class, "title-link id-track-click").Attrs()["href"]
	if suffix == "" {
		return similarApps // could not find similar apps
	}
	similarAppsURL := "https://play.google.com" + suffix

	resp, err := soup.Get(similarAppsURL)
	if err != nil {
		fmt.Println("could not reach", similarAppsURL, " because of the following error:")
		fmt.Println(err)
		return similarApps
	}
	similarAppsDoc := soup.HTMLParse(resp)
	items := similarAppsDoc.FindAll("a", class, "card-click-target")
	//fmt.Println("items", items)
	for key := range items {
		split := strings.Split(items[key].Attrs()["href"], "id=")
		if len(split) > 1 {
			similarApps = append(similarApps, split[1])
		}
	}

	return removeDuplicates(similarApps)
}

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func getCurrentDate() int64 {
	unformattedCurrentDate := time.Now()
	s := strftime.Format("%Y%m%d", unformattedCurrentDate)
	val, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("ERR", err)
		panic(err)
	}
	currentDate := int64(val)

	return currentDate
}
