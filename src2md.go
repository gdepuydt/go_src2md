/*
comment block 1
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//comment1
type Comment struct {
	comment   string
	codeStart int
	codeEnd   int
}

//comment2
func ExtractComments(content []byte) []Comment {

	lineCount := 1
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
		if content[index] == '/' {
			index++
			if content[index] == '/' {
				index++
				codeEnd += lineCount - 1
				commentStart = index
				for content[index] != '\n' && index < sourceLength-1 {
					index++
				}
				commentEnd = index
				comments = append(comments, Comment{strings.Trim(string(content[commentStart:commentEnd]), "\n"), codeStart, codeEnd})
				if content[index] == '\n' {
					lineCount++
				}
				codeStart += lineCount
			} else if content[index] == '*' {
				index++
				codeEnd += lineCount - 1
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
							codeStart += lineCount + 1
							break
						}
					}
					index++
				}
				comments = append(comments, Comment{string(content[commentStart:commentEnd]), codeStart, codeEnd})
			}
		}
		index++
	}
	return comments
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
					if _, err := f.Write([]byte(comment.comment)); err != nil {
						log.Fatal(err)
					}
				}
				if err := f.Close(); err != nil {
					log.Fatal(err)
				}

			}
		}
	}
}
