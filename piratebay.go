package gopiratebay

import (
	"code.google.com/p/go.net/html"
	"net/http"
	"net/url"
	"strings"
	"strconv"
)

type Torrent struct {
	Title      string
	Magnetlink string
	Size       int
	Uploaded   string
	User       string
	Vipuser    bool
	Seeders    int
	Leechers   int
	Category   string
	CategoryId int
}

func Search(q string) (error, []Torrent) {
	resp, err := http.Get("http://thepiratebay.se/search/" + url.QueryEscape(q) + "/0/99/0")
	if err != nil {
		return err, nil
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err, nil
	}
	torrentsChannel := make(chan Torrent)
	result := make(chan []Torrent)
	go torrentReceiver(torrentsChannel, result)
	loopdom(doc, torrentsChannel)
	close(torrentsChannel)
	torrents := <-result
	return nil, torrents
}

func torrentReceiver(torCh chan Torrent, result chan []Torrent) {
	torrents := make([]Torrent, 0)
	for tor := range torCh {
		torrents = append(torrents, tor)
	}
	result <- torrents
}

func loopdom(n *html.Node, tc chan Torrent) {
	if n.Type == html.ElementNode && n.Data == "tbody" {
		// Do something with n...
		extractTorrents(n, tc)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		loopdom(c, tc)
	}
}

func extractTorrents(n *html.Node, tc chan Torrent) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "tr" {
			var torrent Torrent
			torrent.Vipuser = false
			err := getTorrent(c, &torrent)
			if err == nil {
				tc <- torrent
			}
		}
	}
}

func getTorrent(n *html.Node, t *Torrent) error {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if n.Data == "a" && a.Key == "href" && a.Val[:6] == "magnet" {
				t.Magnetlink = a.Val
			} else if n.Data == "a" &&  a.Key == "href" && a.Val[:9] == "/torrent/" {
				if t.Title == "" {
					t.Title = getNodeText(n)
				}
			} else if n.Data == "font" && a.Key == "class" && a.Val == "detDesc" {
				parts := strings.Split(getNodeText(n), ", ")
				if len(parts) > 1 {
					t.Uploaded = strings.Split(parts[0], " ")[1]
					t.Size = sizeStrToInt(strings.Split(parts[1], " ")[1])
				}
			} else if n.Data == "img" && a.Key == "alt" && a.Val == "VIP" {
				t.Vipuser = true
			} else if n.Data == "a" && a.Key == "class" && a.Val == "detDesc" {
				t.User = getNodeText(n)
			} else if n.Data == "a" && a.Key == "href" && a.Val[:8] == "/browse/" && t.Category == "" {
				t.Category = getNodeText(n)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getTorrent(c, t)
	}
	return nil
}

func sizeStrToInt(s string) int {
	var multiply int
	if len(s) < 5 {
		return 0
	}
	multiply = 1
	ext := s[len(s)-3:]
	if ext == "MiB" {
		multiply = 1024*1024
	} else if ext == "KiB" {
		multiply = 1024
	} else if ext == "GiB" {
		multiply = 1024*1024*1024
	}
	size, err := strconv.ParseFloat(s[:len(s)-5], 64)
	if err != nil {
		return 0
	}
	return int(size) * multiply
}

func getNodeText(n *html.Node) string {
	for a := n.FirstChild; a != nil; a = a.NextSibling {
		if a.Type == html.TextNode {
			return strings.TrimSpace(a.Data)
		}
	}
	return ""
}
