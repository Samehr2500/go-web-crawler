package main

import (
	"crawler/data/models"
	"crawler/data/repositories"
	"fmt"
	"runtime"
	"sync"

	crawler "crawler/crawler"
)

var Wg sync.WaitGroup

func main() {
	fmt.Println("Version", runtime.Version())
	fmt.Println("NumCPU", runtime.NumCPU())          //4
	fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0)) //4

	limit_just := -1 // -1 is all

	crawler.GetPaperLinks(limit_just)

	size := crawler.GetLengthPaperLinks()
	fmt.Println("found ", size, " links")

	divisor := 4
	sizeOfList := int(size) / divisor
	for i := 0; i < sizeOfList; i++ {
		Wg.Add(1)
		if (divisor - i) == 1 { //if is the last iteration, take care of summing the remainder
			go crawler.GetInfoFromURL(crawler.PopPaperLinks(sizeOfList+(int(size)%divisor)), &Wg)
		} else {
			go crawler.GetInfoFromURL(crawler.PopPaperLinks(sizeOfList), &Wg)
		}
	}
	Wg.Wait()
	paperSize := crawler.GetLengthPapersInfo()
	fmt.Printf("No of info returned: %d ", paperSize)

	fmt.Println("\n - WRITING TO DB -")

	for i := 0; i < int(paperSize); i++ {
		info := crawler.PopPaperInfo()
		var stock models.Stock
		stock.MarketValue = float32(info.MarketValue)
		stock.CompanyName = crawler.ParseNonUtf8(info.CompanyName)
		stock.PaperName = info.CompanyName
		stock.DailyRate = info.DailyRate

		repositories.Insert(stock)
	}
	fmt.Println("\n - READING FROM DB -")
	result := repositories.GetAllIds()
	fmt.Println("\n - PRINTING FROM DATABASE ")
	for i := 0; i < len(result); i++ {
		item := repositories.Get(result[i])
		fmt.Printf("#%d - \t Company: %s \n \t Market Value: %.2f \n", i, item.CompanyName, item.MarketValue)
	}

	fmt.Println("\n.. FINISHED .. ")
}
