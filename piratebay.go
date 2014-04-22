package gopiratebay

import (
    "code.google.com/p/go.net/html"
    "fmt"
    "net/http"
    )

type Torrent struct {
    Title string
    Magnetlink string
    //created date
    Size int
    User string
    Seeders int
    Leechers int
}

func Search() (error, []Torrent) {
    resp, err := http.Get("http://pong.erwin.io/tpb.html")
    if err != nil {
        fmt.Println("helaas, mislukt")
        return err, nil
    }
    doc, err := html.Parse(resp.Body)
    if err != nil {
        fmt.Println("parse mislukt")
        return err, nil
    }
    torrentsChannel := make(chan Torrent)
	torrents := make([]Torrent, 0)
    go torrentReceiver(torrentsChannel, torrents)
    loopdom(doc, torrentsChannel)
	close(torrentsChannel)
	return nil, torrents
}

func torrentReceiver(torCh chan Torrent, torrents []Torrent) {
    for tor := range torCh {
		torrents = append(torrents, tor)
    }
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
            err := getTorrent(c, &torrent)
            if err == nil {
                tc <- torrent
            }
        }
    }
}

func getTorrent(n *html.Node, t *Torrent) error {
    if n.Type == html.ElementNode && n.Data == "a" {
        for _, a := range n.Attr {
            if a.Key == "href" && a.Val[:6] == "magnet" {
                t.Magnetlink = a.Val
            }
            if a.Key == "href" && a.Val[:9] == "/torrent/" {
                for a := n.FirstChild; a != nil; a = a.NextSibling {
                    if a.Type == html.TextNode && t.Title == "" {
                        t.Title = a.Data
                    }
                }
            }
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        getTorrent(c, t)
    }
    return nil
}
