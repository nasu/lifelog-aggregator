package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/nasu/lifelog-aggregator/domain/google/maps/locationhistory"
)

func main() {
	filepath := flag.String("file", "", "semantic location filepath")
	flag.Parse()
	f, err := os.Open(*filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	hist := locationhistory.SemanticLocationHistory{}
	if err := json.Unmarshal(buf, &hist); err != nil {
		log.Fatal(err)
	}
	seg := hist.GetActivitySegment()
	vis := hist.GetPlaceVisits()
	fmt.Println(len(seg), len(vis), len(hist.TimelineObjects))
	fmt.Printf("%#+v", vis[0])
}
