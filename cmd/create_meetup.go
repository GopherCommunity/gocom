package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/gophercommunity/gocom/internal/data"
	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var name string
var website string
var city string
var countryCode string
var countryName string

func init() {
	rootCmd.AddCommand(createMeetupCmd)

	createMeetupCmd.Flags().StringVar(&name, "name", "", "Name of the meetup")
	createMeetupCmd.Flags().StringVar(&website, "website", "", "URL of the meetup's website")
	createMeetupCmd.Flags().StringVar(&city, "city", "", "City of the meetup")
	createMeetupCmd.Flags().StringVar(&countryCode, "country-code", "", "Country code of the meetup")
	createMeetupCmd.Flags().StringVar(&countryName, "country", "", "Country of the meetup")
}

var createMeetupCmd = &cobra.Command{
	Use:   "create-meetup",
	Short: "Creates a new meetup page.",
	Run: func(cmd *cobra.Command, args []string) {
		if name == "" {
			log.Fatalf("No name specified")
		}
		if city == "" {
			log.Fatalf("No city specified")
		}
		if countryCode == "" {
			log.Fatalf("No country code specified")
		}
		if countryName == "" {
			log.Fatalf("No country specified")
		}

		countriesFilePath := filepath.Join(rootFolder, "data", "countries.yaml")
		countries, err := data.LoadCountries(countriesFilePath)
		if err != nil {
			log.WithError(err).Fatalf("Failed to load countries")
		}

		slugified := slug.Make(name)
		fpath := filepath.Join(rootFolder, "content", "meetup", countryCode, fmt.Sprintf("%s.md", slugified))

		_, err = os.Stat(fpath)
		if err == nil {
			log.Fatalf("%s already exists", fpath)
		}
		if !os.IsNotExist(err) {
			log.WithError(err).Fatalf("Failed to check if %s already exists", fpath)
		}

		now := time.Now()

		if !countries.ContainsCode(countryCode) {
			countries.Add(data.Country{
				Code: countryCode,
				Name: countryName,
			})
			if err := countries.WriteToFile(countriesFilePath); err != nil {
				log.WithError(err).Fatalf("Failed to update countries file")
			}
		}

		page := data.MeetupPage{
			Title:       name,
			Website:     website,
			City:        city,
			CountryCode: countryCode,
			Date:        now.Format(time.RFC3339),
		}
		out, err := yaml.Marshal(page)
		if err != nil {
			log.WithError(err).Fatal("Failed to encode page data as YAML")
		}
		if err := ioutil.WriteFile(fpath, []byte(fmt.Sprintf("---\n%s---", out)), 0644); err != nil {
			log.WithError(err).Fatalf("Failed to write %s", fpath)
		}
		log.Infof("%s created", fpath)
	},
}
