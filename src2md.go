/*
F
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

//comment1
type Comment struct {
	text      string
	codeStart int
	codeEnd   int
}

//comment2
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
				//index++
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

//Write comments to Markdown file
func Write(f *os.File, text string) {
	if _, err := f.Write([]byte(text)); err != nil {
		log.Fatal(err)
	}
}

/*
main function
*/

func main() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if pf := filepath.Ext(file.Name()); pf == ".go" {
			fmt.Println("processing file...", file.Name())
			content, err := ioutil.ReadFile(file.Name())
			if err != nil {
				log.Fatal(err)
			}

			if len(content) > 0 && os.MkdirAll("./mdbook/src/", os.ModeDir) == nil {
				err := os.Chdir("./mdbook/src/")
				err = os.Remove(strings.Replace(file.Name(), ".go", ".md", 1))
				if err != nil {
					log.Fatal(err)
				}
				f, err := os.OpenFile(strings.Replace(file.Name(), ".go", ".md", 1), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatal(err)
				}
				for _, comment := range ExtractComments(content) {
					if len(comment.text) > 0 {
						if comment.codeEnd > 0 {
							Write(f, "```go\n")
							text := "{{#include ../../" + file.Name() + ":" + strconv.Itoa(comment.codeStart) + ":" + strconv.Itoa(comment.codeEnd) + "}}\n\n"
							Write(f, text)
							Write(f, "```\n")
						}
						Write(f, comment.text)
						Write(f, "\n\n")
					} else if comment.codeEnd > 0 {
						Write(f, "```go\n")
						text := "{{#include ../../" + file.Name() + ":" + strconv.Itoa(comment.codeStart) + ":}}\n\n"
						Write(f, text)
						Write(f, "```\n")
					}

				}
				if err := f.Close(); err != nil {
					log.Fatal(err)
				}

			}
		}
	}
}
