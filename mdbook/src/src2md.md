
# src2md.go
[src2md](https://github.com/gdepuydt/go_src2md) reads source files for the [Go programming language](https://golang.org/) and turns it into [Markdown](https://en.wikipedia.org/wiki/Markdown).

The program calls MDBook (written in the [Rust programming language](https://www.rust-lang.org/) to generate the website **_you are looking at now_**.

You can checkout the [Github repository](https://github.com/gdepuydt/go_src2md) for this project

## Here is the source ;-)

I will need to clean it up still. I also intend to document the source code better.

In fact that's the whole point of this project, document your source code in a [Literate](https://en.wikipedia.org/wiki/Literate_programming) way

It's an experiment and still a work in progress. we'll see where it goes!


```go
{{#include ../../src2md.go:17:38}}

```

### ExtractComments

//TODO: Documentation!


```go
{{#include ../../src2md.go:45:92}}

```

### Write text to Markdown file


```go
{{#include ../../src2md.go:97:102}}

```

### Create a new file, removing the file thatalready exists. all text will be appended upon writing.
// TODO: I think I need to look into makingsure the open files are closed properly!

also ... It woud be cool if TODOs are colored red for the webpage...



```go
{{#include ../../src2md.go:111:125}}

```

### src2md function

The actual generation of the Markdown file happens here


```go
{{#include ../../src2md.go:132:185}}

```

### MDBookBuild builds this actual webpage you are viewing right now


```go
{{#include ../../src2md.go:190:195}}

```

# start of the main function

Keep then main simple [KISS](https://nl.wikipedia.org/wiki/KISS-principe)

It would be better to put the main function at the top of the page.

Also would't it be nice if the function calls are hyperlinks to each other... definitely on the TODO!!



```go
{{#include ../../src2md.go:207:}}

```
