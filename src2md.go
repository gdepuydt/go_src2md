/*
A simple program to extract comments from Go source files
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

/*
this is a comment block
*/

//ExtractComments uses basic lex scanning of a byte array for comments
func ExtractComments(content []byte) {

	fmt.Println("Extract one line comments using basic lex scanning pf a byte arrray...")
	index := 0
	commentStart := 0
	commentEnd := 0
	sourceLength := len(content)
	for index < sourceLength {
		if content[index] == '/' {
			index++
			if content[index] == '/' {
				if index > 1 && content[index-2] != '"' {
					index++
					commentStart = index
					for content[index] != '\n' && index < sourceLength-1 {
						index++
					}
					commentEnd = index
					fmt.Println("single line comment found:", string(content[commentStart:commentEnd]))
				}
			} else if content[index] == '*' {
				if index > 1 && content[index-2] != '"' {
					index++
					commentStart = index
					for index < sourceLength-1 {
						if content[index] == '*' {
							index++
							if content[index] == '/' {
								commentEnd = index - 1
								index++
								break
							}
						}
						index++
					}
					fmt.Println("Found comment block:", strings.Trim(string(content[commentStart:commentEnd]), "\n"))
				}
			}

		}
		index++
	}
}

//A comment
func main() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if pf := filepath.Ext(file.Name()); pf == ".go" {
			fmt.Println("processing file:", file.Name())
			content, err := ioutil.ReadFile(file.Name())
			if err != nil {
				log.Fatal(err)
			}
			ExtractComments(content)
		}
	}
}
