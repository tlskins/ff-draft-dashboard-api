package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
)

func main() {
	client := p.NewHttpClient()
	out := map[string]interface{}{}
	if err := p.HttpRequest(client, "GET", p.GetEspnApiUrl(2023), p.EspnQueryHeader(350, 0), nil, &out); err != nil {
		log.Fatal(err)
	}
	spew.Dump(out)
}
