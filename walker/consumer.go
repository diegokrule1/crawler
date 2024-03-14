package walker

import (
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"time"
)

type Consumer struct {
	PageChan chan string
	KillChan chan bool
	Repo     UrlRepository
	Logger   *zap.Logger
	Walker   Walker
	Producer *Producer
	Ticker   *time.Ticker
}

func (c *Consumer) Consume() {
	for {
		select {
		case id := <-c.PageChan:
			go c.processNewUrl(id)
		case <-c.Ticker.C:
			c.Logger.Info("Got	 Tick")
			total, err := c.Repo.countNonterminated()
			if err != nil {
				c.Logger.Error("Could not check termination condition", zap.Error(err))
				return
			}
			if total == 0 {
				c.Logger.Info("No more url to process")
				c.KillChan <- true
				break
			}
		}
	}
}

func (c *Consumer) processNewUrl(id string) {
	urlToProcess, err := c.Repo.getUrlToProcess(id)
	if err != nil {
		c.Logger.Error("Could not get message to process")
		return
	}

	_, err = c.Repo.changeState(urlToProcess.Id, Processing)

	if err != nil {
		c.Logger.Error("Could not update url state", zap.Error(err))
		return
	}
	c.Logger.Info(fmt.Sprintf("url %s on processing", urlToProcess.Id))

	fullUrl := fmt.Sprintf("%s://%s%s", urlToProcess.Scheme, urlToProcess.Host, urlToProcess.Path)
	newUrl, err := c.Walker.Walk(fullUrl, fmt.Sprintf("%s://%s", urlToProcess.Scheme, urlToProcess.Host))
	if err != nil {
		c.Logger.Error("could not analyze url", zap.Stringp("url", &fullUrl))
	}
	for _, urlToSave := range newUrl {
		parsedUrl, err := url.Parse(urlToSave)
		if err != nil {
			c.Logger.Error("Could not parse url", zap.String("url", urlToSave))
			continue
		}
		c.Producer.Produce(parsedUrl.Scheme, parsedUrl.Host, parsedUrl.Path, &urlToProcess.Id)

	}

	total, err := c.Repo.changeState(urlToProcess.Id, "processed")
	if err != nil {
		c.Logger.Error("could not change state to processed to url with id", zap.String("id", urlToProcess.Id))
	} else {
		c.Logger.Info(fmt.Sprintf("Message processed with id %s and total %d ", urlToProcess.Id, total))
	}
}
