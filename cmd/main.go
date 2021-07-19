package main

import (
	"fmt"
	"os"

	"github.com/edgarSucre/scagold/pkg/parameter"
)

func main() {
	sc, err := parameter.Parse(os.Args[0], os.Args[1:])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sc)
}
