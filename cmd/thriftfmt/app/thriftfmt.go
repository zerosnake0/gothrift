package app

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/zerosnake0/gothrift/pkg/format"
	"github.com/zerosnake0/gothrift/pkg/parser"
)

var (
	write       bool
	debug       bool
	showVersion bool
)

const versionStr = "0.0.2"

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
	if debug {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			txt := scanner.Text()
			txt = strings.ReplaceAll(txt, "\t", "--->")
			fmt.Fprintf(w, "%s|\n", txt)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		if _, err := io.Copy(w, r); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func Main() {
	flag.BoolVar(&showVersion, "V", false, "print version")
	flag.BoolVar(&write, "w", false, "write directly to file")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	if showVersion {
		fmt.Println(versionStr)
		return
	}

	if write {
		debug = false
	}

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
				name := filename
				for {
					name += "~"
					if _, err := os.Stat(name); err != nil {
						if os.IsNotExist(err) {
							break
						}
						fmt.Println(name, err)
						os.Exit(1)
					}
				}
				func() {
					fp, err := os.Create(name)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					defer fp.Close()
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
