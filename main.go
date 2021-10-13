package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func main() {
	// fmt.Println("Hello, world.")
	content := buildIndexContent(".")

	createMDFile("toc-index.md", content)
}

// findFiles looks for files recursively
func findFiles(root string) []string {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".md" {
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

// readFile Reads a markdown file and return its content
func readFile(file string) []byte {
	content, err := ioutil.ReadFile(file)

	if err != nil {
		fmt.Println("Error opening file!")
	}

	return content

}

// getFirstParagraph gets the text of first paragraph in a markdown file
func getFirstParagraph(file string) abstractParagraph {

	doc, source := ParseDocument(file)
	// FilterHeadingAbstract("Another title", file)

	title := "# " + string(doc.FirstChild().Text(source))
	content := string(doc.FirstChild().NextSibling().Text(source))

	return abstractParagraph{
		title,
		content,
	}
}

func ParseDocument(filePath string) (ast.Node, []byte) {
	source := readFile(filePath)
	gm := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithASTTransformers(),
		),
	)

	return gm.Parser().Parse(text.NewReader(source)), source
}

func FilterHeadingAbstract(title string, filePath string) abstractParagraph {
	var content abstractParagraph

	doc, source := ParseDocument(filePath)
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)
		var err error

		if n.Kind().String() == "Heading" && string(n.Text(source)) == title {
			content = abstractParagraph{
				string(n.Text(source)),
				string(n.NextSibling().Text(source)),
			}
			s = ast.WalkStatus(ast.WalkStop)
		}
		return s, err
	})

	return content
}

//abstractParagraph represents the paragraph which will be used as abstract.
type abstractParagraph struct {
	title   string
	content string
}

func createMDFile(filePath string, content []string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range content {
		_, _ = datawriter.WriteString(data)
	}

	datawriter.Flush()
	file.Close()
}

//builds the content of the final file
func buildIndexContent(sourcePath string) []string {

	// var files string
	var files, content []string

	files = findFiles(sourcePath)

	for key, _ := range files {
		content = append(content, getFirstParagraph(files[key]).title)
		content = append(content, "\n\n")
		content = append(content, getFirstParagraph(files[key]).content)
		content = append(content, "\n\n")
		content = append(content, "[Read more on the original file...]("+files[key]+")")
		content = append(content, "\n\n")
	}

	return content
}
