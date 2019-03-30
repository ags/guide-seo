package main

import (
	"context"
	"errors"
	"flag"
	"html/template"
	"log"
	"os"

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
		collectionID  = fs.Int("collection", 0, "collection id")
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
	if *collectionID == 0 {
		return errors.New("missing collection ID")
	}

	tmpl := template.Must(template.ParseFiles("template.html"))

	gc := guide.NewClient(*guideAPIKey)

	c, err := gc.FindCollection(context.Background(), guide.FindCollectionInput{
		RegionID:      *regionID,
		CollectionID:  *collectionID,
		CompanyAPIKey: *companyAPIKey,
	})
	if err != nil {
		return err
	}

	p := page{
		Rows: make([]row, len(c.Destinations)),
	}
	for i, d := range c.Destinations {
		p.Rows[i] = row{
			Name:        d.Name,
			Description: template.HTML(d.Description),
			ImageURL:    "https://guide.app" + d.BannerImages[0] + "?w=240&h=160",
		}
	}

	tmpl.Execute(os.Stdout, p)
	return nil
}

type page struct {
	Rows []row
}

type row struct {
	Name        string
	Description template.HTML
	ImageURL    string
}
