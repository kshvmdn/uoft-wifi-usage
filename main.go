package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const baseURL string = "http://maps.wireless.utoronto.ca/stg"

type building struct {
	id                    string
	name                  string
	connectionCount       int
	accessPointsUsed      int
	accessPointsAvailable int
}

func scrapeBuildingIDs() []string {
	var buildingIDs []string

	doc, err := goquery.NewDocument(fmt.Sprintf("%s/index.php", baseURL))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("area[shape=\"circle\"]").Each(func(i int, s *goquery.Selection) {
		buildingID, exists := s.Attr("alt")
		if exists {
			buildingIDs = append(buildingIDs, buildingID)
		}
	})

	return buildingIDs
}

func scrapeWifiUsage(id string) building {
	doc, err := goquery.NewDocument(fmt.Sprintf("%s/popUp.php?name=%s", baseURL, id))
	if err != nil {
		log.Fatal(err)
	}

	name := doc.Find("font[size=\"5\"]").Text()
	b := building{id: id, name: name}

	doc.Find("bq").Each(func(i int, s *goquery.Selection) {
		intVal, err := strconv.Atoi(s.Text())
		if err != nil {
			log.Fatal(err)
		}

		switch i {
		case 0:
			b.connectionCount = intVal
		case 1:
			b.accessPointsUsed = intVal
		case 2:
			b.accessPointsAvailable = intVal
		}
	})

	return b
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {
	buildingPtr := flag.String("buildings", "", "Building IDs (leave empty to view all)")
	verbosePtr := flag.Bool("verbose", false, "Show detailed output")
	flag.Parse()

	var wg sync.WaitGroup
	var buildings []building

	buildingIDs := scrapeBuildingIDs()

	if *buildingPtr != "" {
		// Filter inputted building IDs (remove nonexistent IDs)
		var filteredIds []string

		for _, requestedID := range strings.Split(*buildingPtr, ",") {
			requestedID = strings.Trim(requestedID, " ")
			if contains(buildingIDs, requestedID) {
				filteredIds = append(filteredIds, requestedID)
			} else if *verbosePtr {
				log.Printf("Building with ID `%s` doesn't exist.", requestedID)
			}
		}

		buildingIDs = filteredIds
	}

	wg.Add(len(buildingIDs))

	for _, buildingID := range buildingIDs {
		go func(id string) {
			defer wg.Done()
			buildings = append(buildings, scrapeWifiUsage(id))
		}(buildingID)
	}

	wg.Wait()

	for _, building := range buildings {
		if *verbosePtr {
			fmt.Printf("%s - %d active connections, %d active access points (of %d)\n",
				building.name, building.connectionCount, building.accessPointsUsed,
				building.accessPointsAvailable)
		} else {
			fmt.Printf("%s, %d connections\n", building.name, building.connectionCount)
		}
	}
}
