package markdown

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path"
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

//AbstractParagraph represents the paragraph which will be used as abstract.
type AbstractParagraph struct {
	title   string
	content string
}

type RawMarkdown struct {
	path    string
	content []byte
}

type AstNode struct {
	ast.Node
}

func Execute(output string, directory string) {

	// Add Cobra CLI later
	contentNode, contentByte := buildIndexContent(directory, ignoreDirectories)

	createMDFile(directory+"/"+output, contentNode.renderPlainMarkdown(contentByte))
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
func (md RawMarkdown) calculatePathDepth() int {
	return strings.Count(path.Dir(path.Clean("/"+md.path)), "/")
}

// readFile Reads a markdown file and return its content
func (file RawMarkdown) readFile() []byte {
	content, err := ioutil.ReadFile(file.path)

	if err != nil {
		log.Fatalf("Error opening file %s!", file.path)
	}

	return content

}

// FirstParagraph gets the text of first paragraph in a markdown file
func (md *RawMarkdown) FirstParagraph() AbstractParagraph {
	doc, source := md.ParseDocument()

	if paragraph := md.FilterHeadingAbstract("Abstract"); (paragraph != AbstractParagraph{}) {
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

func (md *RawMarkdown) ParseDocument() (ast.Node, []byte) {
	source := md.readFile()
	gm := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithASTTransformers(),
		),
	)

	return gm.Parser().Parse(text.NewReader(source)), source
}

func (md *RawMarkdown) FilterHeadingAbstract(title string) AbstractParagraph {
	var content AbstractParagraph

	doc, source := md.ParseDocument()

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

func buildIndexContent(sourcePath string, ignoreDirectories []string) (AstNode, []byte) {

	files := findFiles(sourcePath, ignoreDirectories)

	var finalDoc AstNode
	finalDoc.Node = ast.NewDocument()
	var file RawMarkdown
	for _, filepath := range files {
		file.path = filepath
		heading := ast.NewHeading(file.calculatePathDepth())

		paragraph := ast.NewParagraph()

		heading.AppendChild(heading, ast.NewString([]byte(file.FirstParagraph().title)))

		paragraphContent := file.FirstParagraph().content + "\n\n[Read more on the original file...](" + strings.TrimPrefix(filepath, "../") + ")"
		// paragraphContent := file.FirstParagraph().content + "\n\n[Read more on the original file...](" + filepath + ")"
		paragraph.AppendChild(paragraph, ast.NewString([]byte(paragraphContent)))

		finalDoc.AppendChild(finalDoc, heading)
		finalDoc.AppendChild(finalDoc, paragraph)

	}

	var markdown RawMarkdown
	markdown.content = []byte(finalDoc.renderPlainMarkdown([]byte("")))

	tocNode, source := markdown.buildTableOfContents()

	return tocNode, source
}

func (document AstNode) renderHTMLMarkdown(content []byte) string {
	var buffer bytes.Buffer

	gm := goldmark.New()

	err := gm.Renderer().Render(&buffer, content, document.Node)

	if err != nil {
		log.Fatalf("An error occured: %s", err)
	}

	return buffer.String()
}

func (document AstNode) renderPlainMarkdown(content []byte) string {
	var buffer bytes.Buffer

	mdrender := mdrender.NewRenderer()
	err := mdrender.Render(&buffer, content, document.Node)
	if err != nil {
		log.Fatalf("An error occured: %s", err)
	}

	return buffer.String()
}

func (source RawMarkdown) buildTableOfContents() (AstNode, []byte) {
	gm := goldmark.New(
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithExtensions(
			&toc.Extender{},
		),
	)

	return AstNode{gm.Parser().Parse(text.NewReader(source.content))}, source.content
}
