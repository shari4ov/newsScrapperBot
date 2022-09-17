package service

import (
	"fmt"
	"home/storage"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type News struct {
	title string
	id    string
}

var url string = "https://www.oxu.az"

func CreateNewsService() {
	for {
		response := getHtml(url)
		defer response.Body.Close()
		doc, err := goquery.NewDocumentFromReader(response.Body)
		checkError(err)
		scrapePageData(doc)
		href, _ := doc.Find(".pagination a.more").Attr("href")
		if href == "" {
			break
		} else {
			url = "https://www.oxu.az"
			url = fmt.Sprintf("%v%v", url, href)
		}
	}
}
func WriteToDB(n News) error {
	db := storage.OpenConnection()
	sqlStatement := `INSERT INTO news (title,news_id) VALUES($1,$3);`
	_, err := db.Exec(sqlStatement, n.title, n.id)
	if err != nil {
		panic(err)
	}
	return nil

}
func getHtml(url string) *http.Response {
	response, err := http.Get(url)
	checkError(err)
	if response.StatusCode > 400 {
		fmt.Println("Status code", response.StatusCode)
	}
	return response
}
func checkError(error error) {
	if error != nil {
		fmt.Println(error)
	}
}
func scrapePageData(doc *goquery.Document) {
	doc.Find(".news-i").Each(func(i int, s *goquery.Selection) {
		text := s.Find(".title").Text()
		url, _ := s.Find(".news-i-inner").Attr("href")

		spl := strings.Split(url, "/")
		id := spl[len(spl)-1]

		scrapedData := News{
			title: text,
			id:    id,
		}
		WriteToDB(scrapedData)
	})
}
