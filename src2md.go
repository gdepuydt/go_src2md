/*
# src2md.go
[src2md](https://github.com/gdepuydt/go_src2md) reads source files for the [Go programming language](https://golang.org/) and turns it into [Markdown](https://en.wikipedia.org/wiki/Markdown).

The program calls MDBook (written in the [Rust programming language](https://www.rust-lang.org/) to generate the website **_you are looking at now_**.

You can checkout the [Github repository](https://github.com/gdepuydt/go_src2md) for this project

## Here is the source ;-)

I will need to clean it up still. I also intend to document the source code better.

In fact that's the whole point of this project, document your source code in a [Literate](https://en.wikipedia.org/wiki/Literate_programming) way

It's an experiment and still a work in progress. we'll see where it goes!
*/
package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Comment struct {
	text          string
	codeStartLine int
	codeEndLine   int
	offset        int
}

/*
### ExtractComments

//TODO: Documentation!
*/

func ExtractComments(filename string, content []byte) []Comment {

	var comments []Comment
	var s scanner.Scanner
	// positions are relative to fset
	// register input "file"
	fset := token.NewFileSet()
	file := fset.AddFile(filename, fset.Base(), len(content))
	s.Init(file, content, nil, scanner.ScanComments)

	newCodeBlock := true
	codeStartLine := 0
	codeEndLine := 0
	// Repeated calls to Scan yield the token sequence found in the input.
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		if tok != token.COMMENT && newCodeBlock == true {
			codeStartLine = fset.Position(pos).Line
			newCodeBlock = false
		}
		if tok == token.COMMENT {

			if !strings.HasPrefix(lit, "//") {
				codeEndLine = fset.Position(pos).Line - 1
				lit = strings.TrimSuffix(strings.TrimPrefix(lit, "/*"), "*/")
				comments = append(comments, Comment{lit, codeStartLine, codeEndLine, fset.Position(pos).Offset})
				newCodeBlock = true
			}

			pos, tok, lit = s.Scan()
			for tok == token.COMMENT {
				if !strings.HasPrefix(lit, "//") {
					lit = strings.TrimSuffix(strings.TrimPrefix(lit, "/*"), "*/")
					comments = append(comments, Comment{lit, 0, 0, fset.Position(pos).Offset})

				}
				pos, tok, lit = s.Scan()
			}

		}
	}
	comments = append(comments, Comment{"", codeStartLine, 0, 0})
	return comments
}

/*
### Write text to Markdown file
*/

func Write(f *os.File, text string) {
	if _, err := f.Write([]byte(text)); err != nil {
		log.Fatal(err)
	}
}

/*
### Create a new file, removing the file thatalready exists. all text will be appended upon writing.
// TODO: I think I need to look into makingsure the open files are closed properly!

also ... It woud be cool if TODOs are colored red for the webpage...

*/

func createFile(fileName string) *os.File {
	err := os.Remove(fileName)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: do the defer destructor function here to correctly close the opened file

	return f
}

/*
### src2md function

The actual generation of the Markdown file happens here
*/

func src2md() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	if len(files) > 0 {

		for _, file := range files {
			if pf := filepath.Ext(file.Name()); pf == ".go" {
				fmt.Println("processing file...", file.Name())
				content, err := ioutil.ReadFile(file.Name())
				if err != nil {
					log.Fatal(err)
				}
				err = os.MkdirAll("./mdbook/src/", os.ModeDir)
				if err != nil {
					log.Fatal(err)
				}
				err = os.Chdir("./mdbook/src/")
				summary := createFile("SUMMARY.md")
				Write(summary, "# "+file.Name()+"\n\n")
				Write(summary, "- ["+strings.Replace(file.Name(), ".go", "", 1)+"](./"+strings.Replace(file.Name(), ".go", ".md", 1)+")")
				if len(content) > 0 {
					mdfile := createFile(strings.Replace(file.Name(), ".go", ".md", 1))
					//loop over all found comments
					for _, comment := range ExtractComments(mdfile.Name(), content) {
						if comment.codeEndLine != 0 {
							Write(mdfile, "```go\n")
							text := "{{#include ../../" + file.Name() + ":" + strconv.Itoa(comment.codeStartLine) + ":" + strconv.Itoa(comment.codeEndLine) + "}}\n\n"
							Write(mdfile, text)
							Write(mdfile, "```\n")
						}
						if len(comment.text) > 1 {
							Write(mdfile, comment.text)
							Write(mdfile, "\n\n")
						}
						if comment.codeStartLine > comment.codeEndLine {
							Write(mdfile, "```go\n")
							text := "{{#include ../../" + file.Name() + ":" + strconv.Itoa(comment.codeStartLine) + ":}}\n\n"
							Write(mdfile, text)
							Write(mdfile, "```\n")
						}

					}
					if err := mdfile.Close(); err != nil {
						log.Fatal(err)
					}

				}
			}
		}
	}
}

/*
### MDBookBuild builds this actual webpage you are viewing right now
*/

func MDBookBuild() {
	os.Chdir("../")
	cmd := exec.Command("mdbook", "build", "-o")
	cmd.Run()
}

/*
# start of the main function

Keep then main simple [KISS](https://nl.wikipedia.org/wiki/KISS-principe)

It would be better to put the main function at the top of the page.

Also would't it be nice if the function calls are hyperlinks to each other... definitely on the TODO!!

*/

func main() {
	src2md()
	MDBookBuild()

}
