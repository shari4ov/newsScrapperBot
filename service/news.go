package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rabbitmq/amqp091-go"
)

type News struct {
	Title string `json:"title"`
	Id    string `json:"id"`
	Href  string `json:"href"`
}
type LastNews struct {
	Title   string `json:"title"`
	Id      string `json:"id"`
	Href    string `json:"href"`
	News_id string `json:"news_id"`
}

var url string = "https://www.oxu.az"

func CreateNewsService() {
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
	fmt.Println(err)
	amqpChannel, err := conn.Channel()
	fmt.Println(err)
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
	if error != nil {
		fmt.Println(error)
	}
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
	handleError(err, "Line 87")
	defer conn.Close()
	ch, err := conn.Channel()
	handleError(err, "Line 92")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"lastNews.queue",
		true,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Line 104")
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Line 111")
	channel := make(chan LastNews, 1)
	go func() {
		for d := range msgs {
			channel <- JsonToObject(d.Body)
		}
	}()
	return channel
}

func handleError(err error, msg string) {
	if err != nil {
		fmt.Println(err, msg)
	}
}
func JsonToObject(jsonData []byte) LastNews {
	var n LastNews
	_ = json.Unmarshal(jsonData, &n)
	return n
}
