/*
# src2md.go
Uses the Go stl "text/scanner" package
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Comment struct {
	text      string
	codeStart int
	codeEnd   int
}

func ExtractComments(content []byte) []Comment {

	lineCount := 0
	index := 0
	commentStart := 0
	commentEnd := 0
	codeStart := 0
	codeEnd := 0
	sourceLength := len(content)
	var comments []Comment

	for index < sourceLength-1 {

		if content[index] == '\n' {
			index++
			lineCount++
		}
		if content[index] == '"' {
			index++
			for content[index] != '"' {
				if content[index] == '\n' {
					lineCount++
				}
				index++
			}
			index++
		}
		if content[index] == '\n' {
			index++
			lineCount++
		}
		if content[index] == '/' {
			index++
			if content[index] == '/' {
				index++
				codeEnd = lineCount
				commentStart = index
				for content[index] != '\n' && index < sourceLength-1 {
					index++
				}
				commentEnd = index
				comments = append(comments, Comment{strings.Trim(string(content[commentStart:commentEnd]), "\n"), codeStart, codeEnd})

				if content[index] == '\n' {
					lineCount++
				}
				codeStart = lineCount + 1
			} else if content[index] == '*' {
				index++
				codeEnd = lineCount
				if content[index] == '\n' {
					lineCount++
					index++
				}
				commentStart = index
				for index < sourceLength-1 {
					if content[index] == '\n' {
						lineCount++
						index++
					}

					if content[index] == '*' {
						index++
						if content[index] == '/' {
							commentEnd = index - 1

							break
						}
					}
					index++
				}
				comments = append(comments, Comment{strings.Trim(string(content[commentStart:commentEnd]), "\n"), codeStart, codeEnd})
				if content[index] == '\n' {
					lineCount++
					index++
				}
				codeStart = lineCount + 3
			}

		}
		index++
	}
	comments = append(comments, Comment{"", codeStart, codeEnd})
	return comments
}

//Write text to Markdown file
func Write(f *os.File, text string) {
	if _, err := f.Write([]byte(text)); err != nil {
		log.Fatal(err)
	}
}

//Create a new file, removing the file thatalready exists. all text will be appended upon writing.
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

/*src2md function
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
					for _, comment := range ExtractComments(content) {
						if len(comment.text) > 0 {
							if comment.codeEnd > 0 {
								Write(mdfile, "```go\n")
								text := "{{#include ../../" + file.Name() + ":" + strconv.Itoa(comment.codeStart) + ":" + strconv.Itoa(comment.codeEnd) + "}}\n\n"
								Write(mdfile, text)
								Write(mdfile, "```\n")
							}
							Write(mdfile, comment.text)
							Write(mdfile, "\n\n")
						} else if comment.codeEnd > 0 {
							Write(mdfile, "```go\n")
							text := "{{#include ../../" + file.Name() + ":" + strconv.Itoa(comment.codeStart) + ":}}\n\n"
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

//MDBookBuild builds this actual webpage you are vieuwing right now
func MDBookBuild() {
	// TODO!!
}

/*
## start of the main function
*/

func main() {
	src2md()
	MDBookBuild()

}
