package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/jmespath/jp/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/jmespath/jp/Godeps/_workspace/src/github.com/jmespath/go-jmespath"
)

const version = "0.1.3"

func main() {
	app := cli.NewApp()
	app.Name = "jpp"
	app.Version = version
	app.Usage = "jpp [<options>] <expression>"
	app.Author = ""
	app.Email = ""
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "accumulate, a",
			Usage:  "Accumulate all output objects into a single recursively merged output object.",
		},
		cli.BoolFlag{
			Name:   "compact, c",
			Usage:  "Produce compact JSON output that omits nonessential whitespace.",
		},
		cli.StringFlag{
			Name:  "filename, f",
			Usage: "Read input JSON from a file instead of stdin.",
		},
		cli.StringFlag{
			Name:  "expr-file, e",
			Usage: "Read JMESPath expression from the specified file.",
		},
		cli.BoolFlag{
			Name:   "slurp, s",
			Usage:  "Read one or more input JSON objects into an array and apply the JMESPath expression to the resulting array.",
		},
		cli.BoolFlag{
			Name:   "unquoted, u",
			Usage:  "If the final result is a string, it will be printed without quotes.",
		},
		cli.BoolFlag{
			Name:  "ast",
			Usage: "Only print the AST of the parsed expression.  Do not rely on this output, only useful for debugging purposes.",
		},
	}
	app.Action = runMainAndExit

	app.Run(os.Args)
}

func runMainAndExit(c *cli.Context) {
	os.Exit(runMain(c))
}

func errMsg(msg string, a ...interface{}) int {
	fmt.Fprintf(os.Stderr, msg, a...)
	fmt.Fprintln(os.Stderr)
	return 1
}

func runMain(c *cli.Context) int {
	var expression string
	if c.String("expr-file") != "" {
		byteExpr, err := ioutil.ReadFile(c.String("expr-file"))
		expression = string(byteExpr)
		if err != nil {
			return errMsg("Error opening expression file: %s", err)
		}
	} else {
		if len(c.Args()) == 0 {
			return errMsg("Must provide at least one argument.")
		}
		expression = c.Args()[0]
	}
	if c.Bool("ast") {
		parser := jmespath.NewParser()
		parsed, err := parser.Parse(expression)
		if err != nil {
			if syntaxError, ok := err.(jmespath.SyntaxError); ok {
				return errMsg("%s\n%s\n",
					syntaxError,
					syntaxError.HighlightLocation())
			}
			return errMsg("%s", err)
		}
		fmt.Println("")
		fmt.Printf("%s\n", parsed)
		return 0
	}
	var jsonParser *json.Decoder
	if c.String("filename") != "" {
		f, err := os.Open(c.String("filename"))
		if err != nil {
			return errMsg("Error opening input file: %s", err)
		}
		jsonParser = json.NewDecoder(f)

	} else {
		jsonParser = json.NewDecoder(os.Stdin)
	}

	var slurpInput []interface{}
	if c.Bool("slurp") {
		for {
			var element interface{}
			if err := jsonParser.Decode(&element); err == io.EOF {
				break
			} else if err != nil {
				errMsg("Error parsing input json: %s\n", err)
				return 2
			}
			slurpInput = append(slurpInput, element)
		}
	}

	var accumulator interface{}
	eof := false

	for {
		var result interface{}
		for {
			var input interface{}
			var err error
			if c.Bool("slurp") {
				input = slurpInput
			} else if err = jsonParser.Decode(&input); err == io.EOF {
				eof = true
				break
			} else if err != nil {
				errMsg("Error parsing input json: %s\n", err)
				return 2
			}
			result, err = jmespath.Search(expression, input)
			if err != nil {
				if syntaxError, ok := err.(jmespath.SyntaxError); ok {
					return errMsg("%s\n%s\n",
						syntaxError,
						syntaxError.HighlightLocation())
				}
				return errMsg("Error evaluating JMESPath expression: %s", err)
			}

			if c.Bool("accumulate") {
				if accumulator == nil {
					accumulator = result
				} else {
					accumulator, err = merge(result, accumulator); if err != nil {
						errMsg("Error merging output json: %s\n", err)
						return 2
					}
				}
			} else {
				break
			}
		}

		if c.Bool("accumulate") {
			result = accumulator
		} else if eof {
			break
		}

		converted, isString := result.(string)
		if c.Bool("unquoted") && isString {
			os.Stdout.WriteString(converted)
		} else {
			var toJSON []byte
			var err error
			if c.Bool("compact") {
				toJSON, err = json.Marshal(result)
			} else {
				toJSON, err = json.MarshalIndent(result, "", "  ")
			}
			if err != nil {
				errMsg("Error marshalling result to JSON: %s\n", err)
				return 3
			}
			os.Stdout.Write(toJSON)
		}
		os.Stdout.WriteString("\n")
		if eof || c.Bool("accumulate") || c.Bool("slurp") {
			break
		}
	}
	return 0
}

// The following merge and merge1 functions come from the
// golang playground link posted by Roger Peppe in this
// "Recursively merge JSON structures" thread:
//
// https://groups.google.com/g/golang-nuts/c/nLCy75zMlS8/m/O9ZMubnKCQAJ
// https://play.golang.org/p/8jlJUbEJKf

// merge merges the two JSON-marshalable values x1 and x2,
// preferring x1 over x2 except where x1 and x2 are
// JSON objects, in which case the keys from both objects
// are included and their values merged recursively.
//
// It returns an error if x1 or x2 cannot be JSON-marshaled.
func merge(x1, x2 interface{}) (interface{}, error) {
	data1, err := json.Marshal(x1)
	if err != nil {
		return nil, err
	}
	data2, err := json.Marshal(x2)
	if err != nil {
		return nil, err
	}
	var j1 interface{}
	err = json.Unmarshal(data1, &j1)
	if err != nil {
		return nil, err
	}
	var j2 interface{}
	err = json.Unmarshal(data2, &j2)
	if err != nil {
		return nil, err
	}
	return merge1(j1, j2), nil
}

func merge1(x1, x2 interface{}) interface{} {
	switch x1 := x1.(type) {
	case map[string]interface{}:
		x2, ok := x2.(map[string]interface{})
		if !ok {
			return x1
		}
		for k, v2 := range x2 {
			if v1, ok := x1[k]; ok {
				x1[k] = merge1(v1, v2)
			} else {
				x1[k] = v2
			}
		}
	case []interface{}:
		x2, ok := x2.([]interface{})
		if !ok {
			return x1
		}
		var result []interface{}
		for _, element := range x2 {
			result = append(result, element)
		}
		for _, element := range x1 {
			result = append(result, element)
		}
		return result
	case nil:
		// merge(nil, map[string]interface{...}) -> map[string]interface{...}
		x2, ok := x2.(map[string]interface{})
		if ok {
			return x2
		}
	}
	return x1
}

