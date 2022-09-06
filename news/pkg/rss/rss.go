package rss

import (
	"GoNews/pkg/storage"
	"encoding/json"
	"encoding/xml"
	"github.com/grokify/html-strip-tags-go"
	"log"
	"net/http"
	"os"
	"time"
)

type Channel struct {
	Items []Item `xml:"channel>item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type config struct {
	Rss           []string `json:"rss"`
	RequestPeriod int      `json:"request_period"`
}

//Выгрузка RSS-ленты по заданному URL
func GetRss(url string) ([]storage.Post, error) {
	var c Channel
	//запрос к rss-ленте
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	err = xml.NewDecoder(res.Body).Decode(&c)
	if err != nil {
		return nil, err
	}

	//преобразование данных из rss в список публикаций
	var posts []storage.Post
	for _, i := range c.Items {
		var p storage.Post
		p.Title = i.Title
		p.Content = i.Description
		p.Content = strip.StripTags(p.Content)
		p.Link = i.Link

		t, err := time.Parse(time.RFC1123, i.PubDate)
		if err != nil {
			t, err = time.Parse(time.RFC1123Z, i.PubDate)
		}
		if err != nil {
			t, err = time.Parse("Mon, _2 Jan 2006 15:04:05 -0700", i.PubDate)
		}
		if err == nil {
			p.PubTime = t.Unix()
		}

		posts = append(posts, p)
	}
	return posts, nil
}

//Чтение RSS-лент из конфига с заданным интервалом
func GoNews(configURL string, chP chan<- []storage.Post, chE chan<- error) error {
	//чтение конфига
	file, err := os.Open(configURL)
	if err != nil {
		return err
	}
	var conf config
	err = json.NewDecoder(file).Decode(&conf)
	if err != nil {
		return err
	}

	log.Println("starting to watch rss feeds")
	//запуск горутины для каждой rss-ленты
	for i, r := range conf.Rss {
		go func(r string, i int, chP chan<- []storage.Post, chE chan<- error) {
			for {
				log.Println("launched goroutine", i, "by the link", r)
				p, err := GetRss(r)
				if err != nil {
					chE <- err
					continue
				}
				chP <- p
				log.Println("insert posts from goroutine", i, "by the link", r)
				log.Println("Goroutine ", i, ": waiting for the next iteration")
				time.Sleep(time.Duration(conf.RequestPeriod) * time.Minute)
			}
		}(r, i, chP, chE)
	}
	return nil
}
