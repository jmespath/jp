package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/jmespath/go-jmespath"
)

func main() {
	app := cli.NewApp()
	app.Name = "jp"
	app.Version = "0.0.1"
	app.Usage = "jp expression"
	app.Author = ""
	app.Email = ""
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
	jsonParser := json.NewDecoder(os.Stdin)
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
	toJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Error marshalling result to JSON: %s\n", err)
		return 3
	}
	os.Stdout.Write(toJSON)
	os.Stdout.WriteString("\n")
	return 0
}
