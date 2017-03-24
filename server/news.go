package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ehazlett/steamwire/types"
	"github.com/sirupsen/logrus"
)

const (
	baseURL = "http://api.steampowered.com/ISteamNews/GetNewsForApp/v0002/?appid=%s&count=%d&maxlength=%d&format=json"
)

func buildURL(appID string, count int, maxLength int) string {
	return fmt.Sprintf(baseURL, appID, count, maxLength)
}

// getNews gets the latest news for the specified application
// This is limited to a single item as well as 1024 characters in the content
func (s *Server) getNews(appID string) (*types.AppNews, error) {
	u := buildURL(appID, 1, 1024)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	app := &types.App{}
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		logrus.Errorf("error decoding: %s", err)
		return nil, err
	}

	return app.AppNews, nil
}
