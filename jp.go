package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"

	"github.com/jmespath/jp/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/jmespath/jp/Godeps/_workspace/src/github.com/jmespath/go-jmespath"
)

const version = "0.1.2"

func main() {
	app := cli.NewApp()
	app.Name = "jp"
	app.Version = version
	app.Usage = "jp [<options>] <expression>"
	app.Author = ""
	app.Email = ""
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "filename, f",
			Usage: "Read input JSON from a file instead of stdin.",
		},
		cli.StringFlag{
			Name:  "expr-file, e",
			Usage: "Read JMESPath expression from the specified file.",
		},
		cli.BoolFlag{
			Name:   "unquoted, u",
			Usage:  "If the final result is a string, it will be printed without quotes.",
			EnvVar: "JP_UNQUOTED",
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
	var input interface{}
	var inputStream io.Reader
	if c.String("filename") != "" {
		f, err := os.Open(c.String("filename"))
		if err != nil {
			return errMsg("Error opening input file: %s", err)
		}
		inputStream = f

	} else {
		inputStream = os.Stdin
	}
	newlineNumberReader := NewLineNumberReader(inputStream)
	jsonParser := json.NewDecoder(newlineNumberReader)
	if err := jsonParser.Decode(&input); err != nil {
		syntaxError, ok := err.(*json.SyntaxError)
		if ok && syntaxError.Offset == int64(int(syntaxError.Offset)) {
			line, char := newlineNumberReader.ConvertOffset(int(syntaxError.Offset))
			errMsg("Error parsing input json: %s (line: %d, char: %d)\n",
				syntaxError, line, char)
		} else {
			errMsg("Error parsing input json: %s", err)
		}
		return 2
	}
	result, err := jmespath.Search(expression, input)
	if err != nil {
		if syntaxError, ok := err.(jmespath.SyntaxError); ok {
			return errMsg("%s\n%s\n",
				syntaxError,
				syntaxError.HighlightLocation())
		}
		return errMsg("Error evaluating JMESPath expression: %s", err)
	}
	converted, isString := result.(string)
	if c.Bool("unquoted") && isString {
		os.Stdout.WriteString(converted)
	} else {
		toJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			errMsg("Error marshalling result to JSON: %s\n", err)
			return 3
		}
		os.Stdout.Write(toJSON)
	}
	os.Stdout.WriteString("\n")
	return 0
}

type LineNumberReader struct {
	actualReader     io.Reader
	newlinePositions []int
	bytesRead        int
}

func NewLineNumberReader(actualReader io.Reader) *LineNumberReader {
	return &LineNumberReader{
		actualReader: actualReader,
	}
}

func (lnr *LineNumberReader) Read(p []byte) (n int, err error) {
	n, err = lnr.actualReader.Read(p)

	if err != nil || n == 0 {
		return
	}

	for i, v := range p {
		if i >= n {
			return
		}

		if v == '\n' {
			// add 1 so we record the position of the first character, not the '\n'
			lnr.newlinePositions = append(lnr.newlinePositions, lnr.bytesRead+i+1)
		}
	}

	lnr.bytesRead = lnr.bytesRead + n
	return
}

func (lnr *LineNumberReader) ConvertOffset(offset int) (linePos int, charPos int) {
	index := sort.SearchInts(lnr.newlinePositions, offset)
	// Humans are 1 indexed...
	if index == 0 {
		return 1, offset
	} else {
		return index + 1, offset - lnr.newlinePositions[index-1]
	}
}
