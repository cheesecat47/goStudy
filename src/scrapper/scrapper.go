package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	limit   int    = 50
	baseURL string = fmt.Sprintf("https://kr.indeed.com/jobs?q=python&limit=%d", limit)
)

type extractedJob struct {
	id       string
	title    string
	location string
	summary  string
}

func Scrapper() {
	startTime := time.Now()
	totalPages := getPages(baseURL, 0)

	var jobs []extractedJob
	c := make(chan []extractedJob)
	for i := 0; i < totalPages; i++ {
		go getPage(i, c)
	}

	for i := 0; i < totalPages; i++ {
		jobs = append(jobs, <-c...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs), "in", time.Since(startTime))
}

func getPage(page int, mainC chan<- []extractedJob) {
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("pageUrl:", pageURL)

	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	var jobs []extractedJob
	c := make(chan extractedJob)
	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(
		func(i int, card *goquery.Selection) {
			go extractJob(card, c)
		})

	for i := 0; i < searchCards.Length(); i++ {
		jobs = append(jobs, <-c)
	}

	mainC <- jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-jk")
	c <- extractedJob{
		id:       id,
		title:    cleanString(card.Find(".title>a").Text()),
		location: cleanString(card.Find(".sjcl").Text()),
		summary:  cleanString(card.Find(".summary").Text()),
	}
}

func cleanString(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func getPages(url string, prevLast int) (lastPage int) {
	fmt.Println("Finding the last page...")
	if prevLast != 0 {
		url = baseURL + "&start=" + strconv.Itoa((prevLast-1)*limit)
	}

	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages := s.Find("a")
		pageLen := pages.Length()

		if pageLen == 3 {
			lastPage = prevLast
		} else {
			nextLast := 0
			pages.Each(func(i int, s *goquery.Selection) {
				if i == pageLen-2 {
					nextLast, _ = strconv.Atoi(s.Text())
				}
			})
			lastPage = getPages(baseURL, nextLast)
		}
	})
	return
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)
	utf8bom := []byte{0xEF, 0xBB, 0xBF}
	file.Write(utf8bom)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Link", "Title", "Location", "Summary"}
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{
			"https://kr.indeed.com/viewjob?jk=" + job.id,
			job.title,
			job.location,
			job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalf("Request failed with status code: %d\n", res.StatusCode)
	}
}
