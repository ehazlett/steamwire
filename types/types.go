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

// App is the base object for application news
// See https://developer.valvesoftware.com/wiki/Steam_Web_API#GetNewsForApp_.28v0002.29
// for more details
type App struct {
	AppNews *AppNews `json:"appnews"`
}

// List is the list of results from the app list
type List struct {
	AppList struct {
		Apps struct {
			Info []*AppInfo `json:"app"`
		} `json:"apps"`
	} `json:"applist"`
}

// AppInfo is an individual app from the app list
type AppInfo struct {
	AppID int    `json:"appid"`
	Name  string `json:"name"`
}
