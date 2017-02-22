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

const baseUrl string = "http://maps.wireless.utoronto.ca/stg"

type building struct {
	name                  string
	connectionCount       int
	accessPointsUsed      int
	accessPointsAvailable int
}

func scrapeBuildingIds() []string {
	var buildingIds []string

	doc, err := goquery.NewDocument(fmt.Sprintf("%s/index.php", baseUrl))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("area[shape=\"circle\"]").Each(func(_ int, s *goquery.Selection) {
		buildingId, exists := s.Attr("alt")
		if exists {
			buildingIds = append(buildingIds, buildingId)
		}
	})

	return buildingIds
}

func scrapeWifiUsage(id string) building {
	doc, err := goquery.NewDocument(fmt.Sprintf("%s/popUp.php?name=%s", baseUrl, id))
	if err != nil {
		log.Fatal(err)
	}

	name := doc.Find("font[size=\"5\"]").Text()
	b := building{name: name}

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
	buildingPtr := flag.String("buildings", "", "Building IDs")
	verbosePtr := flag.Bool("verbose", false, "Show detailed output")
	flag.Parse()

	var buildings []building

	buildingIds := scrapeBuildingIds()

	if *buildingPtr != "" {
		// Filter inputted building IDs (remove nonexistent IDs)
		var filteredIds []string

		for _, requestedId := range strings.Split(*buildingPtr, ",") {
			if contains(buildingIds, requestedId) {
				filteredIds = append(filteredIds, requestedId)
			} else if *verbosePtr {
				log.Printf("Building with ID `%s` doesn't exist.", requestedId)
			}
		}

		buildingIds = filteredIds
	}

	var wg sync.WaitGroup
	wg.Add(len(buildingIds))

	for _, buildingId := range buildingIds {
		go func(id string) {
			defer wg.Done()
			buildings = append(buildings, scrapeWifiUsage(id))
		}(buildingId)
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
