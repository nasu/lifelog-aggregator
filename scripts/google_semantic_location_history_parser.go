package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/nasu/lifelog-aggregator/constant"
	"github.com/nasu/lifelog-aggregator/domain/google/maps/locationhistory"
	"github.com/nasu/lifelog-aggregator/infrastructure/dynamodb"
)

var db *dynamodb.DB

func init() {
	ctx := context.Background()
	dynamodbUrl := os.Getenv("DYNAMODB_URL")
	if dynamodbUrl == "" {
		panic("DYNAMODB_URL is required")
	}
	dynamodbRegion := os.Getenv("DYNAMODB_REGION")
	if dynamodbRegion == "" {
		panic("DYNAMODB_REGION is required")
	}

	var err error
	db, err = dynamodb.NewDB(ctx, dynamodbUrl, dynamodbRegion)
	if err != nil {
		panic(err)
	}
}

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
	registerActivities(hist.GetActivitySegment())
	registerVisits(hist.GetPlaceVisits())
}

func registerActivities(act []*locationhistory.ActivitySegment) {
	ctx := context.Background()
	repo := locationhistory.NewMoveRepository(db)
	for _, a := range act {
		if err := repo.Save(ctx, constant.USER_ID, a); err != nil {
			log.Fatal(err)
		}
	}
}

func registerVisits(vis []*locationhistory.PlaceVisit) {
	ctx := context.Background()
	repo := locationhistory.NewVisitRepository(db)
	for _, v := range vis {
		if v.Location.SemanticType == locationhistory.SemanticType_HOME {
			continue
		}
		if v.Location.SemanticType == locationhistory.SemanticType_WORK {
			continue
		}
		if err := repo.Save(ctx, constant.USER_ID, v); err != nil {
			log.Fatal(err)
		}
	}
}
