package service

import (
	"encoding/json"
	"fmt"
	dto "home/DTO"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

type News dto.News
type LastNews dto.LastNews

func CreateNewsService() {
	err := godotenv.Load("local.env")
	checkError(err)
	url := os.Getenv("URL")
	news := <-getLastNews()
	fmt.Println("Id: ", news.News_id)
	fmt.Println("Href: ", news.Href)
	response := getHtml(url)
	doc, err := goquery.NewDocumentFromReader(response.Body)
	checkError(err)
	if news.Href != "" {
		go scrapePageDataFrom(doc, news.Href, news.News_id)
	} else {
		go scrapePageDataForZero(doc)
	}
	defer response.Body.Close()
}
func WriteToRM(ch chan News) error {
	news := <-ch
	fmt.Println(news)
	conn, err := amqp091.Dial("amqp://admin:admin@localhost:5672")
	checkError(err)
	amqpChannel, err := conn.Channel()
	checkError(err)
	defer conn.Close()
	defer amqpChannel.Close()
	q, error := amqpChannel.QueueDeclare(
		"news.queue",
		true,
		false,
		false,
		false,
		nil,
	)
	checkError(error)
	amqpChannel.QueueBind(
		q.Name,
		"",
		"amq.fanout",
		false,
		nil,
	)
	news_Marshalled, _ := json.Marshal(news)
	msg := amqp091.Publishing{
		Body: []byte(news_Marshalled),
	}
	amqpChannel.Publish(
		"",
		q.Name,
		false,
		false,
		msg,
	)
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
			Title: text,
			Id:    id,
			Href:  url,
		}
		ch := make(chan News, 10)
		go func() {
			ch <- scrapedData
		}()
		WriteToRM(ch)
	})
}
func scrapePageDataFrom(doc *goquery.Document, href string, Id string) {
	doc.Find(fmt.Sprintf(".news-i-inner[href='%v']", href)).Parent().PrevAll().Filter(".news-i").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		url, _ := s.Find(".news-i-inner").Attr("href")
		spl := strings.Split(url, "/")
		id := spl[len(spl)-1]
		if text == "" {
			fmt.Println("SS")
		}
		if Id != id {
			scrapedData := News{
				Title: text,
				Id:    id,
				Href:  url,
			}
			ch := make(chan News, 10)
			go func() {
				ch <- scrapedData
			}()
			WriteToRM(ch)
		}
	})
}
func getLastNews() (newsChannel chan LastNews) {
	conn, err := amqp091.Dial("amqp://admin:admin@localhost:5672")
	checkError(err)
	defer conn.Close()
	ch, err := conn.Channel()
	checkError(err)
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"lastNews.queue",
		true,
		false,
		false,
		false,
		nil,
	)
	checkError(err)
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	checkError(err)
	channel := make(chan LastNews, 1)
	go func() {
		for d := range msgs {
			channel <- JsonToObject(d.Body)
		}
	}()
	return channel
}

func JsonToObject(jsonData []byte) LastNews {
	var n LastNews
	_ = json.Unmarshal(jsonData, &n)
	return n
}
