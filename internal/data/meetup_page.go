package data

type MeetupPage struct {
	Title       string   `yaml:"title"`
	CountryCode string   `yaml:"country"`
	City        string   `yaml:"city"`
	Tags        []string `yaml:"tags,omitempty"`
	Date        string   `yaml:"date"`
	Website     string   `yaml:"website,omitempty"`
	Slug        string   `yaml:"-"`
}
