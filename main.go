package main

import (
    "code.google.com/p/go.net/html"
    "fmt"
    "net/http"
    )

type Torrent struct {
    title string
    magnetlink string
    //created date
    size int
    user string
    seeders int
    leechers int
}

func main() {
    resp, err := http.Get("http://pong.erwin.io/tpb.html")
    if err != nil {
        fmt.Println("helaas, mislukt")
        return
    }
    doc, err := html.Parse(resp.Body)
    if err != nil {
        fmt.Println("parse mislukt")
        return
    }
    torrentsChannel := make(chan Torrent)
    go torrentReceiver(torrentsChannel)
    loopdom(doc, torrentsChannel)
}

func torrentReceiver(torCh chan Torrent) {
    for tor := range torCh {
        fmt.Println("-------------------------")
        fmt.Println(tor.magnetlink)
        fmt.Println(tor.title)
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
                t.magnetlink = a.Val
            }
            if a.Key == "href" && a.Val[:9] == "/torrent/" {
                for a := n.FirstChild; a != nil; a = a.NextSibling {
                    if a.Type == html.TextNode && t.title == "" {
                        t.title = a.Data
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
