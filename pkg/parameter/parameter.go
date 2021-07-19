package parameter

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

type Parameter struct {
	Name     string
	Location string
	Repo     *url.URL
	APIonly  bool
}

func Parse(title string, args []string) (*Parameter, error) {
	var repo string
	param := Parameter{}

	fs := flag.NewFlagSet(title, flag.ExitOnError)
	fs.StringVar(&param.Name, "n", "", "name of the app to scafold")
	fs.StringVar(&param.Location, "d", "", "location on disk")
	fs.StringVar(&repo, "r", "", "Github repository")
	fs.BoolVar(&param.APIonly, "s", false, "API backend only or website")

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	if param.Repo, err = url.ParseRequestURI(repo); err != nil {
		return nil, fmt.Errorf("invalid repository url")
	}

	if param.Name == "" {
		return nil, fmt.Errorf("invalid name, parameter -n can't be emtpy")
	}

	if err = invalidLocation(param.Location); err != nil {
		return nil, fmt.Errorf("invalid location parameter -l, %v", err)
	}

	return &param, nil
}

func invalidLocation(dir string) error {
	p, err := os.MkdirTemp(dir, "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(p)
	return nil
}
