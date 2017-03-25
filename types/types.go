package types

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
