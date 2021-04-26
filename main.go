package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {
	filePath := flag.String("file", "csv", "path of your csv file")
	token := flag.String("token", "", "token for discogs api")
	flag.Parse()

	csvFile, err := os.Open(*filePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ';'
	csvReader.LazyQuotes = true

	csvData, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	csvHeader := csvData[0]

	var labelNoFld int
	var titleFld int
	var artistFld int
	for key, value := range csvHeader {
		if value == "Label No" {
			labelNoFld = key
		}
		if value == "Title" {
			titleFld = key
		}
		if value == "Artist" {
			artistFld = key
		}
	}

	csvData = csvData[1 : len(csvData)-1]
	client := loginDiscogs(*token)

	for _, each := range csvData {
		idList := searchDiscogs(each[artistFld], each[labelNoFld], each[titleFld], client)
		if len(idList) == 1 {
			if searchCollection(idList[0], client) > 0 {
				fmt.Println("This Item is already in your collection", each[artistFld], each[labelNoFld], each[titleFld])
			} else {
				fmt.Println("Importing this item to your collection", each[artistFld], each[labelNoFld], each[titleFld])
				AddRelease(idList[0], client)
			}
		} else if len(idList) > 1 {
			fmt.Println("Found ", len(idList), "releases, unable to choose one, please do it manually : ", each[artistFld], each[labelNoFld], each[titleFld])
		} else {
			fmt.Println("Couldn't find any release for ", each[artistFld], each[labelNoFld], each[titleFld])
		}
	}
}
