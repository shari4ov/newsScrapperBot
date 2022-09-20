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
	href  string
}

var url string = "https://www.oxu.az"

func CreateNewsService() {
	Id, href := getLastNews()
	response := getHtml(url)
	doc, err := goquery.NewDocumentFromReader(response.Body)
	checkError(err)
	if href != "" {
		go scrapePageDataFrom(doc, href, Id)
	} else {
		go scrapePageDataForZero(doc)
	}
	defer response.Body.Close()
}
func WriteToDB(ch chan News) error {
	news := <-ch
	fmt.Println(news)
	db := storage.OpenConnection()
	sqlStatement := `INSERT INTO news (title,news_id,href) VALUES($1,$2,$3);`
	if news.title != "" && news.id != "" {
		_, err := db.Exec(sqlStatement, news.title, news.id, news.href)
		if err != nil {
			fmt.Println(err)
		}
		close(ch)
		defer db.Close()
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
func scrapePageDataForZero(doc *goquery.Document) {
	doc.Find(".news-list").Find(".pagination").PrevAll().Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		fmt.Println(text)
		url, _ := s.Find(".news-i-inner").Attr("href")
		spl := strings.Split(url, "/")
		id := spl[len(spl)-1]
		scrapedData := News{
			title: text,
			id:    id,
			href:  url,
		}
		ch := make(chan News, 10)
		go func() {
			ch <- scrapedData
		}()
		WriteToDB(ch)
	})
}
func scrapePageDataFrom(doc *goquery.Document, href string, Id string) {
	doc.Find(fmt.Sprintf(".news-i-inner[href='%v']", href)).Parent().PrevAll().Filter(".news-i").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		url, _ := s.Find(".news-i-inner").Attr("href")
		spl := strings.Split(url, "/")
		id := spl[len(spl)-1]
		fmt.Println("sdsdsd", text)
		if text == "" {
			fmt.Println("SS")
		}
		if Id != id {
			scrapedData := News{
				title: text,
				id:    id,
				href:  url,
			}
			ch := make(chan News, 10)
			go func() {
				ch <- scrapedData
			}()
			WriteToDB(ch)
		}
	})
}
func getLastNews() (string, string) {
	db := storage.OpenConnection()
	sqlStatement := `SELECT news_id,href FROM news ORDER BY id DESC LIMIT 1;`
	row := db.QueryRow(sqlStatement)
	var Id string
	var href string
	err := row.Scan(&Id, &href)
	if err != nil {
		return "", ""
	}
	db.Close()
	return Id, href
}
