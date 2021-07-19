package parameter_test

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/edgarSucre/scagold/pkg/parameter"
)

type response struct {
	p *parameter.Parameter
	e string
}

func TestParse(t *testing.T) {
	dir, _ := os.UserHomeDir()

	tests := []struct {
		name string
		args []string
		want response
	}{
		{
			"simple",
			[]string{"-n", "simple app", "-d", dir, "-r", "https://github.com/test", "-s", "true"},
			response{
				p: &parameter.Parameter{
					"simple app",
					dir,
					&url.URL{Scheme: "https", Host: "github.com", Path: "test"},
					true,
				},
				e: "",
			},
		},
		{
			"invalid repo",
			[]string{"-n", "simple app", "-d", dir, "-r", "should fail", "-s", "true"},
			response{
				p: nil,
				e: "invalid repository",
			},
		},
		{
			"empty name",
			[]string{"-n", "", "-d", dir, "-r", "https://github.com/test", "-s", "true"},
			response{
				p: nil,
				e: "invalid name",
			},
		},
		{
			"empty location",
			[]string{"-n", "app", "-d", "not a dir", "-r", "https://github.com/test", "-s", "true"},
			response{
				p: nil,
				e: "invalid location",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parameter.Parse("", test.args)

			if test.want.e != "" {
				if !strings.Contains(err.Error(), test.want.e) {
					t.Fatalf("\nexpected: %#v \ngot: %#v", test.want.e, err.Error())
				}
			} else {
				if equals(*test.want.p, *got) {
					t.Errorf("\nexpected: %#v \ngot: %#v", test.want.p, got)
				}
			}
		})
	}
}

func equals(p1, p2 parameter.Parameter) bool {
	return p1.Name == p2.Name &&
		p1.Location == p2.Location &&
		fmt.Sprintf(
			"%s:%s/%s",
			p1.Repo.Scheme,
			p1.Repo.Host,
			p1.Repo.Path,
		) == fmt.Sprintf(
			"%s:%s/%s",
			p2.Repo.Scheme,
			p2.Repo.Host,
			p2.Repo.Path,
		) &&
		p1.APIonly == p2.APIonly
}
