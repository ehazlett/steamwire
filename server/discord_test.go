package server

import (
	"testing"
)

func TestGenerateMessage(t *testing.T) {
	item := &NewsItem{
		AppID:    1234,
		Gid:      "testingGID",
		Title:    "Test Title",
		Url:      "http://example.com/test",
		Author:   "ehazlett",
		Contents: "This is a test message",
	}

	msg, err := generateMessage(item)
	if err != nil {
		t.Fatal(err)
	}

	expected := `**Test Title** by ehazlett

This is a test message

[Read more](http://example.com/test)

[Application Page](http://store.steampowered.com/app/1234/)
`

	if msg != expected {
		t.Fatalf("expected message: %s\n\nreceived: %s", expected, msg)
	}
}
