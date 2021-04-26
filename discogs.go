package main

import (
	"fmt"
	"github.com/irlndts/go-discogs"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func init() {
	logLevel := os.Getenv("DJ2DC_LOG")
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	switch logLevel {
	case "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}
}

func loginDiscogs(token string) discogs.Discogs {
	client, err := discogs.New(&discogs.Options{
		UserAgent: "Etienne Journet",
		Currency:  "EUR",
		Token:     token,
		//        Token:     "ujTHWtDyVfxHDwCvfNSvVlYTrzIfkgSVZORsnSmG",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return client
}

func searchDiscogs(artist string, labelNo string, title string, client discogs.Discogs) (discogsID []int) {
	log.Print("Search Based on Catalog Number first : ", labelNo)
	request := discogs.SearchRequest{Catno: labelNo, Type: "release", Format: "Vinyl"}
	search, _ := client.Search(request)

	if len(search.Results) == 0 {
		request = discogs.SearchRequest{Artist: artist, Type: "release", Format: "Vinyl"}
		search, _ = client.Search(request)
	}
	if len(search.Results) == 0 {
		request = discogs.SearchRequest{ReleaseTitle: title, Type: "release", Format: "Vinyl"}
		search, _ = client.Search(request)
	}

	return searchItem(search.Results, artist, title)
}

func searchItem(results []discogs.Result, artist string, title string) (probableRelease []int) {
	probability := make(map[int]int)

	for _, each := range results {
		if strings.Contains(strings.ToLower(strings.Join(each.Format, " ")), "vinyl") {
			for _, word := range strings.Split(artist, " ") {
				if strings.Contains(strings.ToLower(each.Title), strings.ToLower(word)) {
					probability[each.ID]++
				}
			}
			for _, word := range strings.Split(title, " ") {
				if strings.Contains(strings.ToLower(each.Title), strings.ToLower(word)) {
					probability[each.ID]++
				}
			}
			log.Print("This release as a probability number of ", probability[each.ID])
		}
	}
	var maxProba int
	for _, maxProba = range probability {
		break
	}
	for _, value := range probability {
		if value > maxProba {
			maxProba = value
		}
	}
	for key := range probability {
		if probability[key] == maxProba {
			probableRelease = append(probableRelease, key)
		}
	}

	return
}

func searchCollection(discogsID int, client discogs.Discogs) int {
	collection, err := client.CollectionItemsByRelease(client.GetUsername(), discogsID)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return len(collection.Items)
}

func AddRelease(discogsID int, client discogs.Discogs) {
	client.AddToCollection(client.GetUsername(), 1, discogsID)
}
