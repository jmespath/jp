jsoncolor
=========

[![GoDoc](https://godoc.org/github.com/nwidger/jsoncolor?status.svg)](https://godoc.org/github.com/nwidger/jsoncolor)

`jsoncolor` is a drop-in replacement for `encoding/json`'s `Marshal`
and `MarshalIndent` functions which produce colorized output using
fatih's [color](https://github.com/fatih/color) package.

## Installation

```
go get -u github.com/nwidger/jsoncolor
```

## Usage

To use as a replacement for `encoding/json`, exchange

`import "encoding/json"` with `import json "github.com/nwidger/jsoncolor"`.

`json.Marshal` and `json.MarshalIndent` will now produce colorized
output.

## Custom Colors

The colors used for each type of token can be customized by creating a
custom `Formatter` and changing its `XXXColor` fields.  See
[color.New](https://godoc.org/github.com/fatih/color#New) for creating
custom color values and the
[GoDocs](https://godoc.org/github.com/nwidger/jsoncolor#pkg-variables)
for the default colors.

``` go
import (
        "bytes"
		"encoding/json"
		"fmt"
		"log"

        "github.com/fatih/color"
        "github.com/nwidger/jsoncolor"
)

// marshal v into src using encoding/json
src, err := json.Marshal(v)
if err != nil {
        log.Fatal(err)
}

// create custom formatter,
f := jsoncolor.NewFormatter()

// set custom colors
f.SpaceColor = color.New(color.FgRed, color.Bold)
f.CommaColor = color.New(color.FgWhite, color.Bold)
f.ColonColor = color.New(color.FgBlue)
f.ObjectColor = color.New(color.FgBlue, color.Bold)
f.ArrayColor = color.New(color.FgWhite)
f.FieldColor = color.New(color.FgGreen)
f.StringColor = color.New(color.FgBlack, color.Bold)
f.TrueColor = color.New(color.FgWhite, color.Bold)
f.FalseColor = color.New(color.FgRed)
f.NumberColor = color.New(color.FgWhite)
f.NullColor = color.New(color.FgWhite, color.Bold)

// colorized output is written to dst
dst := &bytes.Buffer{}
err := f.Format(dst, src)
if err != nil {
        log.Fatal(err)
}

// print colorized output to stdout
fmt.Println(dst.String())
```
