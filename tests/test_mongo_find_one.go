package main

import (
	"flag"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"

	s "github.com/my_projects/ff-draft-dashboard-api/store"
	t "github.com/my_projects/ff-draft-dashboard-api/types"
)

func main() {
	cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	draftDashDbName := os.Getenv("DRAFT_DASHBOARD_DB_NAME")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPwd := os.Getenv("MONGO_PWD")
	mongoHost := os.Getenv("MONGO_HOST")

	// init

	store, err := s.NewStore(draftDashDbName, mongoHost, mongoUser, mongoPwd)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	// work

	player := &t.PlayerReport{}
	err = store.FindOne(store.PlayerReportsCol(), s.M{"_id": "4362628"}, player)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(player)

	players := []*t.PlayerReport{}
	err = store.Find(store.PlayerReportsCol(), s.M{}, &players)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(len(players))
}
