package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ayush6624/go-chatgpt"
	"github.com/joho/godotenv"

	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
	s "github.com/my_projects/ff-draft-dashboard-api/store"
	t "github.com/my_projects/ff-draft-dashboard-api/types"
)

func main() {
	cfgPath := flag.String("config", "../config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	openAiApiKey := os.Getenv("OPEN_AI_API_KEY")
	draftDashDbName := os.Getenv("DRAFT_DASHBOARD_DB_NAME")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPwd := os.Getenv("MONGO_PWD")
	mongoHost := os.Getenv("MONGO_HOST")

	year := 2023
	replaceAll := true

	// init

	client, err := chatgpt.NewClient(openAiApiKey)
	if err != nil {
		log.Fatal(err)
	}
	httpClient := p.NewHttpClient()
	store, err := s.NewStore(draftDashDbName, mongoHost, mongoUser, mongoPwd)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	reportsCol := store.PlayerReportsCol()

	// work

	players, err := p.GetEspnPlayersForYear(httpClient, year, 400)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %v players\n", len(players))

	for _, player := range players {
		fmt.Printf("Processing: %s\n", player.Name)

		report := &t.PlayerReport{}
		err = store.FindOne(reportsCol, s.M{"_id": player.Id}, report)
		// report not found or force replace all reports
		if report == nil || replaceAll {
			// calc and persist report
			report, err = p.CalcPlayerReport(player, client)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(report.Pros)
			fmt.Println(report.Cons)

			if err = store.Upsert(reportsCol, s.M{"_id": report.Id}, s.M{"$set": report}, nil); err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Printf("Skipping %s\n", player.Name)
		}
	}
}
