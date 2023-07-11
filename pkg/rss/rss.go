package rss

import (
	"encoding/xml"
	"fmt"
	"io"
	"mod36a41/pkg/storage"
	"net/http"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
	"golang.org/x/text/encoding/charmap"
)

// RSS
type Feed struct {
	Name    xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// channel info and RSS-feed data
type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

// read and decode RSS, returns news array
func ParseFeed(url string) ([]storage.Post, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	var f Feed

	d := xml.NewDecoder(resp.Body)
	d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}

	err = d.Decode(&f)
	if err != nil {
		return nil, err
	}

	var news []storage.Post
	for _, item := range f.Channel.Items {
		var p storage.Post
		p.Link = item.Link
		p.Title = item.Title
		p.Content = strip.StripTags(item.Description)
		t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			t, err = time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", item.PubDate)
		}
		if err == nil {
			p.PubTime = t.Unix()
		}
		news = append(news, p)
	}
	return news, nil
}
