package rss

import (
	"encoding/json"
	"encoding/xml"
	strip "github.com/grokify/html-strip-tags-go"
	"log"
	"net/http"
	"news/pkg/storage"
	"os"
	"time"
)

type PostItem struct {
	Title   string
	Content string
	PubDate string
	Link    string
}
type XMLStruct struct {
	Items []PostItem
}
type config struct {
	Rss           []string
	RequestPeriod int
}

func getRssStruct(link string) ([]storage.Post, error) { // get rss
	var postsChan XMLStruct
	res, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	err = xml.NewDecoder(res.Body).Decode(&postsChan)
	if err != nil {
		return nil, err
	}
	var posts []storage.Post
	for _, item := range postsChan.Items {
		var p storage.Post
		p.Title = item.Title
		p.Content = item.Content
		p.Content = strip.StripTags(p.Content)
		p.Link = item.Link

		t, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			t, err = time.Parse(time.RFC1123Z, item.PubDate)
		}
		if err != nil {
			t, err = time.Parse("Mon, _2 Jan 2006 15:04:05 -0700", item.PubDate)
		}
		if err == nil {
			p.PubTime = t.Unix()
		}

		posts = append(posts, p)
	}
	return posts, nil
}

func GetNews(configURL string, chanPosts chan<- []storage.Post, chanErrs chan<- error) error {
	file, err := os.Open(configURL)
	if err != nil {
		return err
	}
	var conf config
	err = json.NewDecoder(file).Decode(&conf)
	if err != nil {
		return err
	}

	log.Println("watching rss channels")
	for i, r := range conf.Rss {
		go func(r string, i int, chanPosts chan<- []storage.Post, chanErrs chan<- error) {
			for {
				log.Println("run  goroutine", i, "on link", r)
				p, err := getRssStruct(r)
				if err != nil {
					chanErrs <- err
					time.Sleep(time.Second * 10)
					continue
				}
				chanPosts <- p
				log.Println("insert posts from goroutine", i, "on link", r)
				log.Println("Goroutine ", i, ": waiting to continue")
				time.Sleep(time.Duration(conf.RequestPeriod) * time.Second * 15)
			}
		}(r, i, chanPosts, chanErrs)
	}
	return nil
}
