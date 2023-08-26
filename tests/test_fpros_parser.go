package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"runtime/debug"

	"github.com/davecgh/go-spew/spew"
	p "github.com/my_projects/ff-draft-dashboard-api/parsers"
	t "github.com/my_projects/ff-draft-dashboard-api/types"
)

func Recover() {
	if r := recover(); r != nil {
		fmt.Println(r)
		debug.PrintStack()
	}
}

func main() {
	defer Recover()

	client := p.NewHttpClient()
	out, err := p.HttpHtmlRequest(client, "GET", p.FProsApiUrl, map[string][]string{}, nil)
	if err != nil {
		log.Fatal(err)
	}

	rgx := regexp.MustCompile(`var ecrData = ({.*})`)
	rs := rgx.FindStringSubmatch(out)
	byt := []byte(rs[1])

	resp := t.FproEcrData{}

	if err := json.Unmarshal(byt, &resp); err != nil {
		panic(err)
	}
	spew.Dump(resp)
}
