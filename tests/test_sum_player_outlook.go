package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ayush6624/go-chatgpt"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"

	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
	s "github.com/my_projects/ff-draft-dashboard-api/store"
	t "github.com/my_projects/ff-draft-dashboard-api/types"
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

	player := &t.Player{
		Id:                "123",
		Position:          t.WR,
		Name:              "AJ Brown",
		EspnPlayerOutlook: "A.J. Brown was everything the Eagles could have hoped for after making the trade in the offseason. Brown finished as the WR8 in fantasy with career highs across the board. The concerns about Jalen Hurts supporting an elite wide receiver in fantasy quickly dissipated as Hurts emerged as a possible MVP candidate. Brown was eighth in raw target volume (146), seventh in deep targets, and 12th in red zone looks. Brown is entering his prime (age 26 season) with an ascending elite quarterback in one of the best offenses in football. Brown is a locked-in WR1 in 2023.",
	}

	client, err := chatgpt.NewClient(openAiApiKey)
	if err != nil {
		log.Fatal(err)
	}

	report, err := p.CalcPlayerReport(player, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(report.Pros)
	fmt.Println(report.Cons)

	// DB

	store, err := s.NewStore(draftDashDbName, mongoHost, mongoUser, mongoPwd)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	reportsCol := store.PlayerReportsCol()

	update := &t.PlayerReport{
		Id:       "123",
		Name:     "AJ Brown",
		Position: t.WR,
		Pros:     report.Pros,
		Cons:     report.Cons,
	}
	out := &t.PlayerReport{}
	if err = store.Upsert(reportsCol, s.M{"_id": "123"}, s.M{"$set": update}, out); err != nil {
		log.Fatal(err)
	}

	spew.Dump(out)
}
