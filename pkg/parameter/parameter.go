package parameter

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

type Parameter struct {
	Name     string
	Location string
	Repo     string
	APIonly  bool
}

//Parse computes arguments into flags parameters
func Parse(writer io.Writer, args []string) (*Parameter, error) {
	param := Parameter{}

	fs := flag.NewFlagSet("scaffold-gen", flag.ExitOnError)
	fs.SetOutput(writer)
	fs.StringVar(&param.Name, "n", "", "name of the app to scafold")
	fs.StringVar(&param.Location, "d", "", "location on disk")
	fs.StringVar(&param.Repo, "r", "", "Github repository")
	fs.BoolVar(&param.APIonly, "s", false, "API backend only or website")

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	i := fs.NArg()

	if i != 0 {
		missing := fs.Args()
		_ = missing
		return nil, fmt.Errorf("no positional parameters expected")
	}

	return &param, nil
}

// Validate verifies if all parameters were set. Returns a list of errors found.
func Validate(params *Parameter) []error {
	errors := make([]error, 0)
	if len(strings.TrimSpace(params.Repo)) == 0 {
		errors = append(errors, fmt.Errorf("invalid repository url -r can't be emtpy"))
	}

	if _, err := url.ParseRequestURI(params.Repo); err != nil {
		errors = append(errors, fmt.Errorf("invalid repository url"))
	}

	if len(strings.TrimSpace(params.Name)) == 0 {
		errors = append(errors, fmt.Errorf("invalid name, parameter -n can't be emtpy"))
	}

	if len(strings.TrimSpace(params.Location)) == 0 {
		errors = append(errors, fmt.Errorf("invalid location, parameter -l can't be empty"))
	}

	if err := invalidLocation(params.Location); err != nil {
		errors = append(errors, fmt.Errorf("invalid location parameter -l, %v", err))
	}

	return errors
}

func invalidLocation(dir string) error {
	p, err := os.MkdirTemp(dir, "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(p)
	return nil
}
