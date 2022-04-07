package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/**
Website: Hacker News

This package is designed to provide the function
of pulling data with the help of the API provided
by the website: <https://news.ycombinator.com/>.
And the API doc can be leran form the folloing
website: <https://github.com/HackerNews/API>.
*/

const (
	apiBase = "https://hacker-news.firebaseio.com/v0"
)

//Exposed to provide methods to pull data
type Client struct {
	apiBase string
}

func (c *Client) defaultify() {
	if c.apiBase == "" {
		c.apiBase = apiBase
	}
}

//Return a int slice of top stories.
func (c *Client) GetTopStories() ([]int, error) {
	c.defaultify()
	//This API will returns a slice of IDs of stories.
	url := fmt.Sprintf("%s/topstories.json", c.apiBase)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//decode
	var ids []int
	err = json.NewDecoder(resp.Body).Decode(&ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (c *Client) GetItem(id int) (Item, error) {
	c.defaultify()
	var item Item
	url := fmt.Sprintf("%s/item/%d.json", c.apiBase, id)
	resp, err := http.Get(url)
	if err != nil {
		return item, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return item, err
	}
	return item, nil
}

type Item struct {
	Id int `json:"id"`
	//One of ["job", "story", "comment", "poll", "pollopt"]
	Type string `json:"type"`
	//Author
	By   string `json:"by"`
	Time int    `json:"time"`
	//Comments of the item
	Kids  []int  `json:"kids"`
	URL   string `json:"url"`
	Title string `json:"title"`
	//hot score? jusk like ranking
	Score int `json:"score"`
	//In the case of stories or polls, the total comment count
	Descendants int `json:"descendants"`

	//have no this attr, so it will be nil, add omitempty tag
	Host string `json:"host,omitempty"`
}
