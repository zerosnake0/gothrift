package test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zerosnake0/gothrift/pkg/format"
	"github.com/zerosnake0/gothrift/pkg/parser"
)

const (
	inputDir  = "input"
	outputDir = "output"
)

func proc(must *require.Assertions, input []byte) (output string) {
	doc, err := parser.Parse(input)
	must.NoError(err)
	buf := bytes.NewBuffer(nil)
	f := format.Formatter{
		Doc:    doc,
		Writer: buf,
	}
	f.Encode()
	return buf.String()
}

func testSingleFile(t *testing.T, relpath string) {
	t.Run(relpath, func(t *testing.T) {
		must := require.New(t)

		// output file check
		output, err := ioutil.ReadFile(filepath.Join(outputDir, relpath))
		must.NoError(err)

		formatted := proc(must, output)
		if !assert.Equal(t, string(output), formatted, "exp:\n%s\ngot:\n%s", output, formatted) {
			t.FailNow()
		}

		// input file check
		input, err := ioutil.ReadFile(filepath.Join(inputDir, relpath))
		must.NoError(err)

		formatted = proc(must, input)
		assert.Equal(t, string(output), formatted, "exp:\n%s\ngot:\n%s", output, formatted)
	})
}

func TestFormat(t *testing.T) {
	filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(inputDir, path)
		require.NoError(t, err)
		testSingleFile(t, rel)
		return nil
	})
}

func TestFormatOne(t *testing.T) {
	testSingleFile(t, "definition/service/1.thrift")
}
