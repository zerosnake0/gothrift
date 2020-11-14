package app

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/zerosnake0/gothrift/pkg/format"
	"github.com/zerosnake0/gothrift/pkg/parser"
)

var (
	write bool
	debug bool
)

func process(r io.Reader, w io.Writer) error {
	doc, err := parser.ParseReader(r)
	if err != nil {
		return err
	}
	formatter := format.Formatter{
		Doc:    doc,
		Writer: w,
	}
	formatter.Encode()
	return nil
}

func output(r io.Reader, w io.Writer) {
	if false {
		if _, err := io.Copy(w, r); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		format := "%s\n"
		if debug {
			format = "%s|\n"
		}
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			txt := scanner.Text()
			if debug {
				txt = strings.ReplaceAll(txt, "\t", "--->")
			}
			fmt.Fprintf(w, format, txt)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if _, err := w.Write([]byte{'\n'}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Main() {
	flag.BoolVar(&write, "w", false, "write directly to file")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	buf := bytes.NewBuffer(nil)
	if flag.NArg() == 0 {
		if err := process(os.Stdin, buf); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		output(buf, os.Stdout)
	} else {
		for _, filename := range flag.Args() {
			// parse
			func() {
				fp, err := os.Open(filename)
				if err != nil {
					fmt.Println(filename, err)
					os.Exit(1)
				}
				defer fp.Close()
				if err := process(fp, buf); err != nil {
					fmt.Println(filename, err)
					os.Exit(1)
				}
			}()
			if write {
				// output
				var name string
				func() {
					fp, err := ioutil.TempFile("", "")
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					defer fp.Close()
					name = fp.Name()
					output(buf, fp)
				}()
				if err := os.Rename(name, filename); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				output(buf, os.Stdout)
			}
			buf.Reset()
		}
	}
}
