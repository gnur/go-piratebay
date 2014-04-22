package main

import ( 
    "github.com/gnur/gopiratebay"
    )

func main() {
    err, torrents := gopiratebay.Search()
	for tor := range torrents {
		fmt.Println(tor.magnetlink)
	}
}
