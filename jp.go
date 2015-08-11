package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jmespath/jp/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/jmespath/jp/Godeps/_workspace/src/github.com/jmespath/go-jmespath"
)

func main() {
	app := cli.NewApp()
	app.Name = "jp"
	app.Version = "0.0.1"
	app.Usage = "jp [<options>] <expression>"
	app.Author = ""
	app.Email = ""
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "filename, f",
			Usage: "Read input JSON from a file instead of stdin.",
		},
		cli.BoolFlag{
			Name:   "unquoted, u",
			Usage:  "If the final result is a string, it will be printed without quotes.",
			EnvVar: "JP_UNQUOTED",
		},
	}
	app.Action = runMainAndExit

	app.Run(os.Args)
}

func runMainAndExit(c *cli.Context) {
	os.Exit(runMain(c))
}

func runMain(c *cli.Context) int {
	if len(c.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Must provide at least one argument.\n")
		return 255
	}
	expression := c.Args()[0]
	var input interface{}
	var jsonParser *json.Decoder
	if c.String("filename") != "" {
		f, err := os.Open(c.String("filename"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening input file: %s\n", err)
			return 1
		}
		jsonParser = json.NewDecoder(f)

	} else {
		jsonParser = json.NewDecoder(os.Stdin)
	}
	if err := jsonParser.Decode(&input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing input json: %s\n", err)
		return 2
	}
	result, err := jmespath.Search(expression, input)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Error evaluating JMESPath expression: %s\n", err)
		return 1
	}
	converted, isString := result.(string)
	if c.Bool("unquoted") && isString {
		os.Stdout.WriteString(converted)
	} else {
		toJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"Error marshalling result to JSON: %s\n", err)
			return 3
		}
		os.Stdout.Write(toJSON)
	}
	os.Stdout.WriteString("\n")
	return 0
}
