package parameter_test

import (
	"bytes"
	"fmt"
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
			[]string{"-n", "simple app", "-d", dir, "-r", "https://github.com/test", "-s"},
			response{
				p: &parameter.Parameter{
					"simple app",
					dir,
					"https://github.com/test",
					true,
				},
				e: "",
			},
		},
		{
			"error",
			[]string{"fail"},
			response{
				p: nil,
				e: "positional parameters",
			},
		},
		// {
		// 	"empty repo",
		// 	[]string{"-n", "simple app", "-d", dir, "-r", "", "-s", "true"},
		// 	response{
		// 		p: nil,
		// 		e: "invalid repository",
		// 	},
		// },
		// {
		// 	"invalid repo",
		// 	[]string{"-n", "simple app", "-d", dir, "-r", "should fail", "-s", "true"},
		// 	response{
		// 		p: nil,
		// 		e: "invalid repository",
		// 	},
		// },
		// {
		// 	"empty name",
		// 	[]string{"-n", "", "-d", dir, "-r", "https://github.com/test", "-s", "true"},
		// 	response{
		// 		p: nil,
		// 		e: "invalid name",
		// 	},
		// },
		// {
		// 	"empty location",
		// 	[]string{"-n", "app", "-d", "not a dir", "-r", "https://github.com/test", "-s", "true"},
		// 	response{
		// 		p: nil,
		// 		e: "invalid location",
		// 	},
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := new(bytes.Buffer)
			got, err := parameter.Parse(output, test.args)

			if test.want.e != "" {
				if !strings.Contains(err.Error(), test.want.e) {
					t.Fatalf("\nexpected: %#v \ngot: %#v", test.want.e, err.Error())
				}
			} else {
				if !equals(*test.want.p, *got) {
					t.Errorf("\nexpected: %#v \ngot: %#v", test.want.p, got)
				}
			}
		})
	}
}

func TestValidate(t *testing.T) {
	dir, _ := os.UserHomeDir()

	tests := []struct {
		name string
		args *parameter.Parameter
		errs []error
	}{
		{
			"no errors",
			&parameter.Parameter{
				"simple app",
				dir,
				"https://github.com/test",
				false,
			},
			[]error{},
		},
		{
			"empty fields",
			&parameter.Parameter{
				"",
				"",
				"",
				false,
			},
			[]error{
				fmt.Errorf("invalid repository url -r can't be emtpy"),
				fmt.Errorf("invalid repository url"),
				fmt.Errorf("invalid name, parameter -n can't be emtpy"),
				fmt.Errorf("invalid location, parameter -l can't be empty"),
			},
		},
		{
			"dummy fields",
			&parameter.Parameter{
				"name",
				"/home/none",
				"not an url",
				false,
			},
			[]error{
				fmt.Errorf("invalid repository url"),
				fmt.Errorf("invalid location parameter -l, %v", "stat /home/none: no such file or directory"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := parameter.Validate(test.args)

			if len(test.errs) == 0 && len(errs) != 0 {
				t.Errorf("Expected no errors, got: %v", errs)
			}

			if len(test.errs) != 0 {
				for i, e := range test.errs {
					if errs[i] == nil || errs[i].Error() != e.Error() {
						t.Errorf("Expected error: %v, Got: %v", e, errs[i])
					}
				}
			}
		})
	}
}

func equals(p1, p2 parameter.Parameter) bool {
	return p1.Name == p2.Name &&
		p1.Location == p2.Location &&
		p1.Repo == p2.Repo &&
		p1.APIonly == p2.APIonly
}
