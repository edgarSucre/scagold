package main

import (
	"fmt"
	"os"

	"github.com/edgarSucre/scagold/pkg/parameter"
)

func main() {
	sc, err := parameter.Parse(os.Stdout, os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	errors := parameter.Validate(sc)
	if len(errors) > 0 {
		for _, e := range errors {
			fmt.Println(e)
		}
		os.Exit(1)
	}
	fmt.Printf("Generating scaffold for project %s in %s\n", sc.Name, sc.Location)
}
