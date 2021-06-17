package main

import (
	"bytes"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/jmespath/go-jmespath"
	"github.com/spf13/cobra"
)

const version = "0.1.3.7"

type JppConfig struct {
	Accumulate bool `json:"accumulate"`
	Ast bool `json:"ast"`
	Compact bool `json:"compact"`
	ExprFile string `json:"expr-file"`
	Filename string `json:"filename"`
	RawInput bool `json:"raw-input"`
	RawOutput bool `json:"raw-output"`
	Slurp bool `json:"slurp"`
	Unbox bool `json:"unbox"`
}

func JppCommand() *cobra.Command {
	var jppCmd = &cobra.Command{
		Use:   "jpp [<options>] [expression]",
		Short: "An extended superset of the jp CLI for JMESPath",
		Args:  cobra.MaximumNArgs(1),
		RunE: JppCobraLaunchMain,
	}

	flags := jppCmd.PersistentFlags()

	flags.BoolP(
		"accumulate",
		"a",
		false,
		"Accumulate all output objects into a single recursively merged output object.",
	)
	flags.Bool(
		"ast",
		false,
		"Only print the AST of the parsed expression.  Do not rely on this output, only useful for debugging purposes.",
	)
	flags.BoolP(
		"compact",
		"c",
		false,
		"Produce compact JSON output that omits nonessential whitespace.",
	)
	flags.StringP(
		"filename",
		"f",
		"",
		"Read input JSON from a file instead of stdin.",
	)
	flags.StringP(
		"expr-file",
		"e",
		"",
		"Read JMESPath expression from the specified file.",
	)
	flags.BoolP(
		"raw-output",
		"r",
		false,
		"If the final result is a string, it will be printed without quotes (an alias for --unquoted).",
	)
	flags.BoolP(
		"raw-input",
		"R",
		false,
		"Read raw string input and box it as JSON strings.",
	)
	flags.BoolP(
		"slurp",
		"s",
		false,
		"Read one or more input JSON objects into an array and apply the JMESPath expression to the resulting array.",
	)
	flags.BoolP(
		"unbox",
		"u",
		false,
		"If the final result is a list, unbox it into a stream of output objects that is suitable for consumption by --slurp mode.",
	)
	flags.Bool(
		"unquoted",
		false,
		"If the final result is a string, it will be printed without quotes.",
	)
	flags.BoolP(
		"help",
		"h",
		false,
		"show usage and exit",
	)

	return jppCmd
}

func main() {
	jppCmd := JppCommand()
	jppCmd.SetArgs(os.Args[1:])
	if err := jppCmd.Execute(); err != nil {
		os.Exit(errMsg(err.Error()))
	}
}

func errMsg(msg string, a ...interface{}) int {
	fmt.Fprintf(os.Stderr, msg, a...)
	fmt.Fprintln(os.Stderr)
	return 1
}

func MustGetString(cmd *cobra.Command, name string) string {
	value, err := cmd.Flags().GetString(name)
	if err != nil {
		panic(err)
	}
	return value
}

func MustGetBool(cmd *cobra.Command, name string) bool {
	value, err := cmd.Flags().GetBool(name)
	if err != nil {
		panic(err)
	}
	return value
}

func JppCobraLaunchMain(cmd *cobra.Command, args []string) error {
	config := &JppConfig{
		MustGetBool(cmd, "accumulate"),
		MustGetBool(cmd, "ast"),
		MustGetBool(cmd, "compact"),
		MustGetString(cmd, "expr-file"),
		MustGetString(cmd, "filename"),
		MustGetBool(cmd, "raw-input"),
		MustGetBool(cmd, "raw-output") || MustGetBool(cmd, "unquoted"),
		MustGetBool(cmd, "slurp"),
		MustGetBool(cmd, "unbox"),
	}
	return JppMain(config, args)
}

func JppMain(config *JppConfig, args []string) error {
	var expression string
	if config.ExprFile != "" {
		byteExpr, err := ioutil.ReadFile(config.ExprFile)
		expression = string(byteExpr)
		if err != nil {
			return fmt.Errorf("Error opening expression file: %w", err)
		}
	} else {
		if len(args) == 0 {
			expression = "@"
		} else if len(args) > 1 {
			return fmt.Errorf("Must not provide more than one argument.")
		} else {
			expression = args[0]
		}
	}
	if config.Ast {
		parser := jmespath.NewParser()
		parsed, err := parser.Parse(expression)
		if err != nil {
			if syntaxError, ok := err.(jmespath.SyntaxError); ok {
				fmt.Errorf("%s\n%s\n",
					syntaxError,
					syntaxError.HighlightLocation())
			}
			return err
		}
		fmt.Println("")
		fmt.Printf("%s\n", parsed)
		return nil
	}
	var jsonParser *json.Decoder
	var f *os.File
	var  rawInput *bufio.Scanner
	var rawInputBuffer []byte
	if config.Filename != "" {
		var err error
		f, err = os.Open(config.Filename)
		if err != nil {
			fmt.Errorf("Error opening input file: %w", err)
		}
	} else {
		f = os.Stdin
	}

	if config.RawInput && config.Slurp {
		var err error
		rawInputBuffer, err = ioutil.ReadAll(f)
		if err != nil {
			fmt.Errorf("Error reading input file: %w", err)
		}
	} else if config.RawInput {
		 rawInput = bufio.NewScanner(f)
	} else {
		jsonParser = json.NewDecoder(f)
	}

	var slurpInput []interface{}
	if config.Slurp && !config.RawInput {
		for {
			var element interface{}
			if  rawInput != nil {
				if !rawInput.Scan() {
					break
				}
				element =  rawInput.Text()
			} else if err := jsonParser.Decode(&element); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("Error parsing input json: %w", err)
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
			if config.Slurp {
				eof = true
				if config.RawInput {
					input = string(rawInputBuffer)
				} else {
					input = slurpInput
				}
			} else if  rawInput != nil {
				if !rawInput.Scan() {
					eof = true
					break
				}
				input =  rawInput.Text()
			} else if err = jsonParser.Decode(&input); err == io.EOF {
				eof = true
				break
			} else if err != nil {
				return fmt.Errorf("Error parsing input json: %w", err)
			}
			result, err = jmespath.Search(expression, input)
			if err != nil {
				if syntaxError, ok := err.(jmespath.SyntaxError); ok {
					return fmt.Errorf("%s\n%s\n",
						syntaxError,
						syntaxError.HighlightLocation())
				}
				return fmt.Errorf("Error evaluating JMESPath expression: %w", err)
			}

			if config.Accumulate {
				if accumulator == nil {
					accumulator = result
				} else {
					accumulator, err = merge(result, accumulator); if err != nil {
						return fmt.Errorf("Error merging output json: %w", err)
					}
				}
				if config.Slurp {
					break
				}
			} else {
				break
			}
		}

		if config.Accumulate {
			result = accumulator
		} else if eof && !config.Slurp {
			break
		}

		var unboxed bool
		if config.Unbox {
			switch result := result.(type) {
				case []interface{}:
					unboxed = true
					for _, element := range result {
						if err := OutputResult(element, config); err != nil {
							return fmt.Errorf("Error marshalling result to JSON: %w", err)
						}
					}
			}
		}

		if !unboxed {
			if err := OutputResult(result, config); err != nil {
				return fmt.Errorf("Error marshalling result to JSON: %w", err)
			}
		}

		if eof || config.Accumulate || config.Slurp {
			break
		}
	}
	return nil
}

func OutputResult(result interface{}, config *JppConfig) error {
	converted, isString := result.(string)
	quoted := ! (config.RawOutput && isString)
	if quoted {
		var toJSON []byte
		var err error
		if config.Compact {
			toJSON, err = json.Marshal(result)
		} else {
			toJSON, err = json.MarshalIndent(result, "", "  ")
		}
		if err != nil {
			return err
		}
		os.Stdout.Write(toJSON)
	} else {
		os.Stdout.WriteString(converted)
	}
	os.Stdout.WriteString("\n")
	return nil
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
			if !contains(result, element) {
				result = append(result, element)
			}
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

func equal(lhs interface{}, rhs interface{}) bool {
	switch lhs := lhs.(type) {
		case nil:
			switch rhs.(type) {
				case nil:
					return true
			}
			return false

		case string:
			switch rhs := rhs.(type) {
				case string:
					if lhs == rhs {
						return true
					}
			}
			return false

		case int:
			switch rhs := rhs.(type) {
				case int:
					if lhs == rhs {
						return true
					}
			}
			return false

		case float32:
			switch rhs := rhs.(type) {
				case float32:
					if lhs == rhs {
						return true
					}
			}
			return false

		case float64:
			switch rhs := rhs.(type) {
				case float64:
					if lhs == rhs {
						return true
					}
			}
			return false
		default:
			panic(fmt.Sprintf("unhandled type comparison: %s vs %s", reflect.TypeOf(lhs), reflect.TypeOf(rhs)))
	}
}

func contains(values []interface{}, value interface{}) bool {

	valueData, err := json.Marshal(value); if err != nil {
		panic(err)
	}

	for _, v := range values {
		switch v := v.(type) {
			case map[string]interface{}:
				switch value.(type) {
					case map[string]interface{}:
						data, err := json.Marshal(v); if err != nil {
							panic(err)
						}
						if bytes.Compare(valueData, data) == 0 {
							return true
						}
				}

			case []interface{}:
				switch value.(type) {
					case []interface{}:
						data, err := json.Marshal(v); if err != nil {
							panic(err)
						}
						if bytes.Compare(valueData, data) == 0 {
							return true
						}
				}

			default:
				if equal(v, value) {
					return true
				}
		}
	}

	return false
}
