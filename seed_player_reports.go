package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ayush6624/go-chatgpt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/my_projects/ff-draft-dashboard-api/api"
	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
	"github.com/my_projects/ff-draft-dashboard-api/store"
)

func main() {
	cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	openAiApiKey := os.Getenv("OPEN_AI_API_KEY")
	draftDashDbName := os.Getenv("DRAFT_DASHBOARD_DB_NAME")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPwd := os.Getenv("MONGO_PWD")
	mongoHost := os.Getenv("MONGO_HOST")

	year := 2023

	// init

	client, err := chatgpt.NewClient(openAiApiKey)
	if err != nil {
		log.Fatal(err)
	}
	httpClient := p.NewHttpClient()
	store, err := store.NewStore(draftDashDbName, mongoHost, mongoUser, mongoPwd)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	playerReportsCol := store.PlayerReportsCol()
	upsertTrue := true
	upsertOpts := &options.UpdateOptions{Upsert: &upsertTrue}

	// work

	espnPlayers, err := p.GetEspnPlayersForYear(httpClient, year)
	if err != nil {
		log.Fatal(err)
	}
	harrisPlayers := p.ParseHarrisRanksV2(year)
	players, _, err := p.MatchHarrisAndEspnPlayers(harrisPlayers, espnPlayers)
	if err != nil {
		panic(err)
	}

	for _, player := range players {
		fmt.Printf("Processing: %s\n", player.Name)
		report, err := p.CalcPlayerReport(player, client)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(report.Pros)
		fmt.Println(report.Cons)

		saveCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err = playerReportsCol.UpdateOne(saveCtx, api.M{"_id": report.Id}, api.M{"$set": report}, upsertOpts)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Saved: %s\n", player.Name)
	}
}
