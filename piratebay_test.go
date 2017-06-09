package piratebay

import (
	"testing"
)

func TestSearchCount(t *testing.T) {
	torrents, _ := Search("go lang piratebay")
	if len(torrents) != 1 {
		t.Error("Expected 1 result, got ", len(torrents))
	}
}

func TestSearchFields(t *testing.T) {
	torrents, _ := Search("go lang piratebay")
	testTor := torrents[0]
	if testTor.Title != "Go lang piratebay api" {
		t.Error("Title should be: Go lang piratebay api, but it is:", testTor.Title)
	}
	if testTor.User != "go-piratebay" {
		t.Error("user should be: go-piratebay, but it is:", testTor.User)
	}
	if testTor.Category != "Applications" {
		t.Error("Category should be: Applications, but it is:", testTor.Category)
	}
	if testTor.VIP {
		t.Error("VIP should be: false, but it is:", testTor.VIP)
	}
}
