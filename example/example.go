package main

import ( 
    "github.com/gnur/gopiratebay"
	"fmt"
    )

func main() {
    err, torrents := gopiratebay.Search()
	for _, tor := range torrents {
		fmt.Println(tor.Magnetlink)
	}
}
