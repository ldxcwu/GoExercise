package hn

import (
	"fmt"
	"testing"
)

func TestGetTopStories(t *testing.T) {
	c := Client{}
	ids, _ := c.GetTopStories()
	item, _ := c.GetItem(ids[0])
	fmt.Printf("%+v\n", item)
}

// func TestX(t *testing.T) {
// 	start := time.Now()
// 	time.Sleep(time.Second * 2)
// 	fmt.Printf("time.Since(start): %.2f\n", time.Since(start).Seconds())
// }
