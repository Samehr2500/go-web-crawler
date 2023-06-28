package crawler

import (
	driver "crawler/data/driver"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

var baseURL = "https://www.fundamentus.com.br/"
var papersLinks = "papers_links"
var papersData = "papers_data"

// GetPaperLinks find all the links in the main page
func GetPaperLinks(limit int) {
	// Request the HTML page.
	res, err := http.Get(baseURL + "detalhes.php")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	rdb, ctx := driver.CreateRedisConnection()
	defer rdb.Close()
	// delete any existing list with the same name
	err = rdb.Del(ctx, papersLinks).Err()
	if err != nil {
		panic(err)
	}
	leng := 0
	doc.Find("tbody td a").Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href") //get the link itself
		// append an item to the end of a list
		if leng < limit || limit == -1 {
			err := rdb.RPush(ctx, papersLinks, link).Err()
			if err != nil {
				panic(err)
			}
			leng++
		}
	})
}

func GetInfoFromURL(urls []string, Wg *sync.WaitGroup) {
	defer Wg.Done()
	var paperName string
	var companyName string
	var marketValue float64
	var dailyRate string

	for i := 0; i < len(urls); i++ {
		paperInfo := PapersInfo{} //declaring a struct of paper with its information
		res, err := http.Get(baseURL + urls[i])

		checkError(err, "at Get: "+baseURL+urls[i])
		defer res.Body.Close()

		doc, err := goquery.NewDocumentFromReader(res.Body)
		checkError(err, "at NewDocFromReader: "+baseURL+urls[i])

		pNameSelector := "body > div.center > div.conteudo.clearfix > table:nth-child(2) > tbody > tr:nth-child(1) > td.data.w35"
		doc.Find(pNameSelector).Each(func(index int, item *goquery.Selection) { //paper's Name
			paperName = item.Find("span").Text()
			paperInfo.PaperName = paperName
		})

		cNameSelector := "body > div.center > div.conteudo.clearfix > table:nth-child(2) > tbody > tr:nth-child(3) > td:nth-child(2)"
		doc.Find(cNameSelector).EachWithBreak(func(index int, item *goquery.Selection) bool { //company's name
			companyName = item.Find("span").Text()
			paperInfo.CompanyName = companyName
			if index == 0 {
				return false
			}
			return true
		})

		mvalueSelector := "body > div.center > div.conteudo.clearfix > table:nth-child(3) > tbody > tr:nth-child(1) > td.data.w3"
		doc.Find(mvalueSelector).EachWithBreak(func(index int, item *goquery.Selection) bool { //market's value
			if index == 0 {
				marketV := item.Find("span").Text()             //text as string
				noDots := strings.Replace(marketV, ".", "", -1) //-1 means all occurrencies (taking out the dots in the string to convert it to float later)
				marketValue, _ = strconv.ParseFloat(noDots, 64) //converting to float
				paperInfo.MarketValue = marketValue
				return false
			}
			return true
		})

		doc.Find("span.oscil").EachWithBreak(func(index int, item *goquery.Selection) bool { //daily's rate
			if index == 0 {
				dailyRate = item.Text()
				paperInfo.DailyRate = dailyRate
				return false
			}
			return true
		})

		jsonPaperInfo, err := json.Marshal(paperInfo)
		if err != nil {
			panic(err)
		}
		rdb, ctx := driver.CreateRedisConnection()
		defer rdb.Close()
		// append the JSON object to the end of a list
		err = rdb.RPush(ctx, papersData, string(jsonPaperInfo)).Err()
		if err != nil {
			panic(err)
		}
	}
}

func checkError(err error, url string) {
	if err != nil {
		fmt.Println("URL: " + url)
		panic(err)
	}
}

func GetLengthPaperLinks() int64 {
	rdb, ctx := driver.CreateRedisConnection()
	defer rdb.Close()
	len, err := rdb.LLen(ctx, papersLinks).Result()
	if err != nil {
		panic(err)
	}

	return len
}

func GetLengthPapersInfo() int64 {
	rdb, ctx := driver.CreateRedisConnection()
	defer rdb.Close()
	len, err := rdb.LLen(ctx, papersData).Result()
	if err != nil {
		panic(err)
	}

	return len
}

func PopPaperLinks(size int) []string {
	rdb, ctx := driver.CreateRedisConnection()
	defer rdb.Close()
	var urls []string

	for j := 0; j < size; j++ {
		item, err := rdb.LPop(ctx, papersLinks).Result()
		if err != nil {
			panic(err)
		}
		urls = append(urls, item)
	}

	return urls
}

func PopPaperInfo() PapersInfo {
	rdb, ctx := driver.CreateRedisConnection()
	defer rdb.Close()
	var info PapersInfo

	// pop a JSON string from the Redis list
	jsonStr, err := rdb.LPop(ctx, papersData).Result()
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(jsonStr), &info)
	if err != nil {
		panic(err)
	}

	return info
}

func ParseNonUtf8(s string) string {
	if !utf8.ValidString(s) {
		v := make([]rune, 0, len(s))
		for i, r := range s {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(s[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		s = string(v)
	}
	return s
}
