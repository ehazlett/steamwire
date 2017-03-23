package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type NewsItem struct {
	AppID     int    `json:"appid"`
	Gid       string `json:"gid"`
	Title     string `json:"title"`
	Url       string `json:"url"`
	Author    string `json:"author"`
	Contents  string `json:"contents"`
	FeedLabel string `json:"feedlabel"`
	FeedName  string `json:"feedname"`
	Date      int    `json:"date"`
}
type AppNews struct {
	AppID     int         `json:"appid"`
	NewsItems []*NewsItem `json:"newsitems"`
	Count     int         `json:"count"`
}

type App struct {
	AppNews *AppNews `json:"appnews"`
}

const (
	baseURL = "http://api.steampowered.com/ISteamNews/GetNewsForApp/v0002/?appid=%s&count=%d&maxlength=%d&format=json"
)

func buildURL(appID string, count int, maxLength int) string {
	return fmt.Sprintf(baseURL, appID, count, maxLength)
}

func (s *Server) GetNews(appID string) (*AppNews, error) {
	u := buildURL(appID, 1, 1024)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	app := &App{}
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		logrus.Errorf("error decoding: %s", err)
		return nil, err
	}

	return app.AppNews, nil
}
