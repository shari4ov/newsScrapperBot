package dto

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
