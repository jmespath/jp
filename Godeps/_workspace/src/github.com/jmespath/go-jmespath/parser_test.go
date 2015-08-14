package jmespath

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var parsingErrorTests = []struct {
	expression string
	msg        string
}{
	{"foo.", "Incopmlete expression"},
	{"[foo", "Incopmlete expression"},
	{"]", "Invalid"},
	{")", "Invalid"},
	{"}", "Invalid"},
	{"foo..bar", "Invalid"},
	{`foo."bar`, "Forwards lexer errors"},
	{`{foo: bar`, "Incomplete expression"},
	{`{foo bar}`, "Invalid"},
	{`[foo bar]`, "Invalid"},
	{`foo@`, "Invalid"},
	{`&&&&&&&&&&&&t(`, "Invalid"},
	{`[*][`, "Invalid"},
}

func TestParsingErrors(t *testing.T) {
	assert := assert.New(t)
	parser := NewParser()
	for _, tt := range parsingErrorTests {
		_, err := parser.Parse(tt.expression)
		assert.NotNil(err, fmt.Sprintf("Expected parsing error: %s, for expression: %s", tt.msg, tt.expression))
	}
}

func BenchmarkParseIdentifier(b *testing.B) {
	runParseBenchmark(b, exprIdentifier)
}

func BenchmarkParseSubexpression(b *testing.B) {
	runParseBenchmark(b, exprSubexpr)
}

func BenchmarkParseDeeplyNested50(b *testing.B) {
	runParseBenchmark(b, deeplyNested50)
}

func BenchmarkParseDeepNested50Pipe(b *testing.B) {
	runParseBenchmark(b, deeplyNested50Pipe)
}

func BenchmarkParseDeepNested50Index(b *testing.B) {
	runParseBenchmark(b, deeplyNested50Index)
}

func BenchmarkParseQuotedIdentifier(b *testing.B) {
	runParseBenchmark(b, exprQuotedIdentifier)
}

func BenchmarkParseQuotedIdentifierEscapes(b *testing.B) {
	runParseBenchmark(b, quotedIdentifierEscapes)
}

func BenchmarkParseRawStringLiteral(b *testing.B) {
	runParseBenchmark(b, rawStringLiteral)
}

func BenchmarkParseDeepProjection104(b *testing.B) {
	runParseBenchmark(b, deepProjection104)
}

func runParseBenchmark(b *testing.B, expression string) {
	parser := NewParser()
	for i := 0; i < b.N; i++ {
		parser.Parse(expression)
	}
}
