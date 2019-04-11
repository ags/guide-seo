package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ags/guide-seo/pkg/guide"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("guide-seo", flag.ExitOnError)
	var (
		guideAPIKey   = fs.String("guide-api-key", "", "")
		companyAPIKey = fs.String("company-api-key", "", "")
		regionID      = fs.Int("region", 0, "region id")
		collections   = fs.String("collections", "", "collection ids")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *guideAPIKey == "" {
		return errors.New("missing guide api key")
	}
	if *companyAPIKey == "" {
		return errors.New("missing company api key")
	}
	if *regionID == 0 {
		return errors.New("missing region ID")
	}
	if *collections == "" {
		return errors.New("missing collection IDs")
	}
	collectionIDs, err := parseCollections(*collections)
	if err != nil {
		return fmt.Errorf("invalid collection IDs: %v", err)
	}

	tmpl := template.Must(template.ParseFiles("template.html"))

	gc := guide.NewClient(*guideAPIKey)

	destinationsByID := map[int]guide.Destination{}

	for _, collectionID := range collectionIDs {
		c, err := gc.FindCollection(context.Background(), guide.FindCollectionInput{
			RegionID:      *regionID,
			CollectionID:  collectionID,
			CompanyAPIKey: *companyAPIKey,
		})
		if err != nil {
			return err
		}
		for _, d := range c.Destinations {
			destinationsByID[d.ID] = d
		}
	}

	p := page{
		Rows: make([]row, len(destinationsByID)),
	}
	i := 0
	for _, d := range destinationsByID {
		p.Rows[i] = row{
			Name:        d.Name,
			Description: template.HTML(d.Description),
			ImageURL:    "https://guide.app" + d.BannerImages[0] + "?w=240&h=160",
		}
		i++
	}

	tmpl.Execute(os.Stdout, p)
	return nil
}

func parseCollections(s string) ([]int, error) {
	var ids []int
	strs := strings.Split(s, ",")
	for _, idstr := range strs {
		id, err := strconv.Atoi(idstr)
		if err != nil {
			return []int{}, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

type page struct {
	Rows []row
}

type row struct {
	Name        string
	Description template.HTML
	ImageURL    string
}
