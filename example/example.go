package main

import ( 
    "github.com/gnur/gopiratebay"
	"fmt"
    )

func main() {
    _, torrents := gopiratebay.Search("go lang piratebay")
	for _, tor := range torrents {
		fmt.Println(tor.Magnetlink)
		fmt.Println(tor.Title)
		fmt.Println(tor.Size)
		fmt.Println(tor.Uploaded)
		fmt.Println(tor.User)
		fmt.Println(tor.Vipuser)
		fmt.Println(tor.Category)
	}
}
