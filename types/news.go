package types

// NewsItem is a steam application news item
type NewsItem struct {
	AppID     int    `json:"appid"`
	Gid       string `json:"gid"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Author    string `json:"author"`
	Contents  string `json:"contents"`
	FeedLabel string `json:"feedlabel"`
	FeedName  string `json:"feedname"`
	Date      int    `json:"date"`
}

// AppNews is the object that contains the application ID, count and news items
type AppNews struct {
	AppID     int         `json:"appid"`
	NewsItems []*NewsItem `json:"newsitems"`
	Count     int         `json:"count"`
}
