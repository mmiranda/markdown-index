package markdown

import (
	"bufio"
	"bytes"
	"io/ioutil"

	// "log"
	"os"
	"path/filepath"
	"strings"

	mdrender "github.com/Kunde21/markdownfmt/v2/markdown"
	toc "github.com/abhinav/goldmark-toc"
	log "github.com/sirupsen/logrus"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var (
	ignoreDirectories []string
	searchHeading     string
)
var (
	LogLevel = log.WarnLevel
)

//abstractParagraph represents the paragraph which will be used as abstract.
type abstractParagraph struct {
	title   string
	content string
}

type rawMarkdown struct {
	path     string
	realPath string
	basedir  string
	content  []byte
}

type astNode struct {
	ast.Node
}

// Execute orchestrates the execution of the library
func Execute(output string, rootDirectory string) {

	log.SetLevel(LogLevel)

	// Add Cobra CLI later
	contentNode, contentByte := buildIndexContent(rootDirectory, ignoreDirectories)

	CreateMDFile(rootDirectory+"/"+output, contentNode.renderPlainMarkdown(contentByte), true)
}

// findFiles looks for files recursively
func findFiles(root string, ignoreDirectories []string) []*rawMarkdown {
	var files []*rawMarkdown

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// Ignoring files that are not .md
		if filepath.Ext(path) != ".md" {
			return nil
		}

		file := newMarkdownFile(root, path)

		// It ignores files previously generated by this lib or explicitly ignored directories
		if file.isFileAutoGenerated() || contains(ignoreDirectories, file.basedir) {
			return nil
		}

		files = append(files, file)
		return nil
	})

	if err != nil {
		panic(err)
	}
	return files
}

// newMarkdownFile Creates a new Markdown object
func newMarkdownFile(rootDirectory, path string) *rawMarkdown {
	var file rawMarkdown
	file.path = path
	file.basedir = filepath.Base(filepath.Dir(file.path))

	file.realPath = GetFileRealPath(rootDirectory, path)

	file.content, _ = file.readFile()

	return &file
}

// GetFileRealPath Prepares the relative path of a file based on the root directory
func GetFileRealPath(rootDirectory, path string) string {
	if rootDirectory == "." {
		return rootDirectory + "/" + path
	}
	return "." + strings.ReplaceAll(path, rootDirectory, "")
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
func (md rawMarkdown) calculatePathDepth() int {
	depth := strings.Count(md.realPath, "/")
	if depth == 1 {
		return depth
	}
	return depth - 1
}

// readFile Reads a markdown file and return its content
func (md rawMarkdown) readFile() ([]byte, error) {
	content, err := ioutil.ReadFile(md.path)

	return content, err

}

// getMetadata gets the YAML metadata in a Markdown file
func (md rawMarkdown) getMetadata() map[string]interface{} {
	var buf bytes.Buffer
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	context := parser.NewContext()
	if err := markdown.Convert([]byte(md.content), &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	return meta.Get(context)
}

// AddMetadataPrefix Adds a prefix in some content
func AddMetadataPrefix(content string) string {
	metadata := "---\ngenerated-by: markdown-index\n---\n"

	return metadata + content
}

// isFileAutoGenerated Checks if a Markdown file has a specific metadata or not
func (md rawMarkdown) isFileAutoGenerated() bool {
	metadata := md.getMetadata()
	if metadata != nil && metadata["generated-by"] == "markdown-index" {
		return true
	}

	return false
}

// FirstParagraph gets the text of first heading in a markdown file OR
// the heading accordingly to the custom filter searchHeading
func (md *rawMarkdown) getFirstParagraph() abstractParagraph {
	var abstract abstractParagraph
	doc, source := md.parseDocument()

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)
		var err error
		if n.Kind().String() != "Heading" {
			return s, err
		} else if string(n.Text(source)) != "Table of Contents" {
			rawContent := ""
			// Check prevents single heading document
			if doc.FirstChild() != doc.LastChild() {
				rawContent = string(n.NextSibling().Text(source))
			}

			// If it is already populated, don't replace it
			if (abstract == abstractParagraph{}) {
				abstract = abstractParagraph{
					string(n.Text(source)),
					rawContent,
				}
			}
			// Look if there is a specific Heading and if found overwrite the previous (first) paragraph
			if searchHeading != "" && string(n.Text(source)) == searchHeading {
				abstract = abstractParagraph{
					string(n.Text(source)),
					string(n.NextSibling().Text(source)),
				}

				s = ast.WalkStatus(ast.WalkStop)
			}

		}

		return s, err
	})

	if err != nil {
		log.Fatalf("An error occurred: %s", err)
	}

	return abstract
}

// parseDocument Uses Goldmark library to parse a document
func (md *rawMarkdown) parseDocument() (ast.Node, []byte) {
	gm := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithASTTransformers(),
		),
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	return gm.Parser().Parse(text.NewReader(md.content)), md.content
}

// DeleteFile Deletes a file from the disk
func DeleteFile(file string) {
	if fi, err := os.Stat(file); err == nil {
		if fi.Mode().IsRegular() {
			err := os.Remove(file)
			if err != nil {
				log.Fatalf("An error occured deleting the file %s: %s", file, err.Error())
			}
		}
	}

}

// CreateMDFile persists a file on disk
func CreateMDFile(filePath, content string, metadata bool) {
	DeleteFile(filePath)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)
	if metadata {
		content = AddMetadataPrefix(content)
	}

	datawriter.WriteString(content)

	datawriter.Flush()
	file.Close()
}

// buildIndexContent Builds an ast.Node (Markdown) based on all the files it found recursively
func buildIndexContent(rootDirectory string, ignoreDirectories []string) (astNode, []byte) {

	files := findFiles(rootDirectory, ignoreDirectories)

	var finalDoc astNode
	finalDoc.Node = ast.NewDocument()

	for _, file := range files {

		log.Debugf("Reading file: %s", file.path)
		heading := ast.NewHeading(file.calculatePathDepth())

		paragraph := ast.NewParagraph()

		heading.AppendChild(heading, ast.NewString([]byte(file.getFirstParagraph().title)))

		paragraphContent := file.getFirstParagraph().content + "\n\n[Read more on the original file...](" + file.realPath + ")"
		paragraph.AppendChild(paragraph, ast.NewString([]byte(paragraphContent)))

		finalDoc.AppendChild(finalDoc, heading)
		finalDoc.AppendChild(finalDoc, paragraph)

	}

	var markdown rawMarkdown
	markdown.content = []byte(finalDoc.renderPlainMarkdown([]byte("")))

	tocNode, source := markdown.buildTableOfContents()

	return tocNode, source
}

// RenderHTMLMarkdown returns a Markdown content rendered as HTML
func (document astNode) RenderHTMLMarkdown(content []byte) string {
	var buffer bytes.Buffer

	gm := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	err := gm.Renderer().Render(&buffer, content, document.Node)

	if err != nil {
		log.Fatalf("An error occurred: %s", err)
	}

	return buffer.String()
}

// renderPlainMarkdown returns a Markdown content rendered as plain Markdown
func (document astNode) renderPlainMarkdown(content []byte) string {
	var buffer bytes.Buffer

	mdrender := mdrender.NewRenderer()
	err := mdrender.Render(&buffer, content, document.Node)
	if err != nil {
		log.Fatalf("An error occurred: %s", err)
	}

	return buffer.String()
}

// buildTableOfContents builds a TOC based on markdown content
func (md rawMarkdown) buildTableOfContents() (astNode, []byte) {
	gm := goldmark.New(
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithExtensions(
			&toc.Extender{},
			meta.Meta,
		),
	)

	return astNode{gm.Parser().Parse(text.NewReader(md.content))}, md.content
}
