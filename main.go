package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	mdrender "github.com/Kunde21/markdownfmt/v2/markdown"
	toc "github.com/abhinav/goldmark-toc"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var ignoreDirectories []string

func main() {

	contentDocument := buildIndexContent(".", ignoreDirectories)

	createMDFile("toc-index.md", renderPlainMarkdown(contentDocument))
}

// findFiles looks for files recursively
func findFiles(root string, ignoreDirectories []string) []string {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// Ignoring pre-defined directories and files that are not .md
		if filepath.Ext(path) != ".md" || contains(ignoreDirectories, filepath.Base(filepath.Dir(path))) {
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

// contains verify if a string is inside a set of strings
func contains(slice []string, searchterm string) bool {
	for _, value := range slice {
		if value == searchterm {
			return true
		}
	}
	return false
}

// calculatePathDepth returns the depth of a folder structure
func calculatePathDepth(path string) int {

	return strings.Count(path, "/")
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
func getFirstParagraph(file string) AbstractParagraph {

	doc, source := ParseDocument(file)

	if paragraph := FilterHeadingAbstract("Abstract", file); (paragraph != AbstractParagraph{}) {
		return AbstractParagraph{
			paragraph.title,
			paragraph.content,
		}
	}

	abstractContent := ""

	// identify single node documents aka markdown with a single heading/paragraph
	if doc.FirstChild() != doc.LastChild() {
		abstractContent = string(doc.FirstChild().NextSibling().Text(source))
	}

	return AbstractParagraph{
		string(doc.FirstChild().Text(source)),
		abstractContent,
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

func FilterHeadingAbstract(title string, filePath string) AbstractParagraph {
	var content AbstractParagraph

	doc, source := ParseDocument(filePath)

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)
		var err error

		if n.Kind().String() == "Heading" && string(n.Text(source)) == title {
			content = AbstractParagraph{
				string(n.Text(source)),
				string(n.NextSibling().Text(source)),
			}
			s = ast.WalkStatus(ast.WalkStop)
		}
		return s, err
	})

	if err != nil {
		log.Fatalf("An error occured: %s", err)
	}

	return content
}

//AbstractParagraph represents the paragraph which will be used as abstract.
type AbstractParagraph struct {
	title   string
	content string
}

func createMDFile(filePath string, content string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	datawriter.WriteString(content)

	datawriter.Flush()
	file.Close()
}

func buildIndexContent(sourcePath string, ignoreDirectories []string) ast.Node {

	files := findFiles(sourcePath, ignoreDirectories)
	finalDoc := ast.NewDocument()
	for key, _ := range files {
		heading := ast.NewHeading(calculatePathDepth(files[key]))

		paragraph := ast.NewParagraph()

		heading.AppendChild(heading, ast.NewString([]byte(getFirstParagraph(files[key]).title)))

		paragraphContent := getFirstParagraph(files[key]).content + "\n\n[Read more on the original file...](" + files[key] + ")"
		paragraph.AppendChild(paragraph, ast.NewString([]byte(paragraphContent)))

		finalDoc.AppendChild(finalDoc, heading)
		finalDoc.AppendChild(finalDoc, paragraph)

	}

	return finalDoc
}

func renderHTMLMarkdown(document ast.Node) string {
	var (
		content []byte
		buffer  bytes.Buffer
	)
	gm := goldmark.New()

	err := gm.Renderer().Render(&buffer, content, document)

	if err != nil {
		log.Fatalf("An error occured: %s", err)
	}
	buffer = prepareTableOfContents(document)
	return buffer.String()
}

func prepareTableOfContents(document ast.Node) bytes.Buffer {
	tree, err := toc.Inspect(document, []byte(""))
	if err != nil {
		log.Fatalf("An error occured: %s", err)
	}

	list := toc.RenderList(tree)
	var output bytes.Buffer

	document.InsertBefore(document, document.FirstChild(), list)

	// Render the document as output.
	gm := goldmark.New()
	gm.Renderer().Render(&output, []byte(""), document)

	return output
}

func renderPlainMarkdown(document ast.Node) string {
	var buffer bytes.Buffer

	bufferHtml := prepareTableOfContents(document)
	mdrender := mdrender.NewRenderer()
	err := mdrender.Render(&buffer, bufferHtml.Bytes(), document)
	if err != nil {
		log.Fatalf("An error occured: %s", err)
	}

	return buffer.String()
}
