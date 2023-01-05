package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/adrianosela/ff"
)

var (
	value  float64
	userID string
)

func main() {
	flag.Float64Var(&value, "value", 0.5, "flag value")
	flag.StringVar(&userID, "user-id", "adrianosela", "user id to test against flag")
	flag.Parse()

	f, err := ff.NewFeatureFlag("test flag", value)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Flag with value %f for user id \"%s\" is %s\n", value, userID, strconv.FormatBool(f.IsEnabledForUser(userID)))
}
