package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldFindReadme(t *testing.T) {
	files := findFiles("./test", []string{})
	assert.Equal(t, 6, len(files))
}

func TestShouldReadFileContent(t *testing.T) {
	files := findFiles("./test", []string{})

	content := readFile(files[0])

	assert.Contains(t, string(content), "# Root Level Markdown")
}

func TestGetFirstParagraph(t *testing.T) {
	files := findFiles("./test", []string{})

	file := RawMarkdown{files[0]}
	title := "# " + file.FirstParagraph().title
	content := file.FirstParagraph().content

	assert.Equal(t, "This is a sample paragraph text for test purpose only. This paragraph will be used as an abstract on the global TOC.", content)
	assert.Equal(t, "# Root Level Markdown", title)
}

func TestGetFirstParagraphInEveryFile(t *testing.T) {
	files := findFiles("./test", []string{})

	for key, _ := range files {
		file := RawMarkdown{files[key]}
		content := file.FirstParagraph().content
		if len(content) > 0 {
			assert.Equal(t, "This is a sample paragraph text for test purpose only. This paragraph will be used as an abstract on the global TOC.", content)
		} else {
			assert.Equal(t, "", content)
		}

	}

	assert.True(t, true)
}

func TestCreateSingleLineFile(t *testing.T) {
	filePath := "/tmp/test.md"
	contentString := "# Hello, this is a test"
	// contentString = append(contentString, "\nMultiline test")
	createMDFile(filePath, contentString)
	assert.FileExists(t, filePath)

	content := readFile(filePath)

	assert.Equal(t, contentString, string(content))

	deleteFile(filePath)
}

func TestCreateMultiLineFile(t *testing.T) {
	filePath := "/tmp/test.md"
	contentString := "# Hello, this is a test\nThis is a sample text"

	createMDFile(filePath, contentString)
	assert.FileExists(t, filePath)

	content := readFile(filePath)

	assert.Equal(t, contentString, string(content))
	deleteFile(filePath)
}

func deleteFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		panic("Test failed deleting file")
	}
}

func TestCompareFinalFilePlainContent(t *testing.T) {
	mockFilePath := "test/mock-toc-toc.mock"
	mockFile := readFile(mockFilePath)

	contentNode, contentByte := buildIndexContent("./test", []string{})

	assert.Equal(t, string(mockFile), renderPlainMarkdown(contentNode, contentByte))
}

func TestCompareFinalFileHTMLContent(t *testing.T) {
	mockFilePath := "test/mock-toc-toc-html.mock"
	mockFile := readFile(mockFilePath)

	contentNode, contentByte := buildIndexContent("./test", []string{})

	assert.Equal(t, string(mockFile), renderHTMLMarkdown(contentNode, contentByte))
}

func TestCompareFinalFilePlainContentWithIgnore(t *testing.T) {
	mockFilePath := "test/mock-toc-toc-ignored.mock"
	mockFile := readFile(mockFilePath)

	contentNode, contentByte := buildIndexContent("./test", []string{"folder2"})

	assert.Equal(t, string(mockFile), renderPlainMarkdown(contentNode, contentByte))
}

func TestCompareFinalFileHTMLContentWithIgnore(t *testing.T) {
	contentNode, contentByte := buildIndexContent("./test", []string{"folder2"})

	mockFile := readFile("test/mock-toc-toc-html-ignored.mock")

	assert.Equal(t, string(mockFile), renderHTMLMarkdown(contentNode, contentByte))
}

func TestCompareFinalFilePlainBytes(t *testing.T) {
	filePath := "/tmp/test-toc.md"
	contentDocument, contentString := buildIndexContent("./test", []string{})
	createMDFile(filePath, renderPlainMarkdown(contentDocument, contentString))
	fileGenerated := readFile(filePath)

	mockFile := readFile("./test/mock-toc-toc.mock")
	assert.True(t, bytes.Equal(fileGenerated, mockFile))

	deleteFile(filePath)
}

func TestCompareFinalFileHTMLBytes(t *testing.T) {
	filePath := "/tmp/test-toc.md"
	contentNode, contentByte := buildIndexContent("./test", []string{})
	createMDFile(filePath, renderHTMLMarkdown(contentNode, contentByte))
	fileGenerated := readFile(filePath)

	mockFile := readFile("./test/mock-toc-toc-html.mock")
	assert.True(t, bytes.Equal(fileGenerated, mockFile))

	deleteFile(filePath)
}

func TestFilterAbstractHeading(t *testing.T) {
	file := RawMarkdown{"test/README.md"}
	content := file.FilterHeadingAbstract("Another title")
	assert.NotEmpty(t, content)

	content = file.FilterHeadingAbstract("Unexistent title heading")
	assert.Empty(t, content)
}

func TestContainsHelper(t *testing.T) {
	directories := []string{"A", "B"}

	assert.True(t, contains(directories, "A"))
	assert.False(t, contains(directories, "C"))
}

func TestCalculateDepth(t *testing.T) {

	assert.Equal(t, 1, calculatePathDepth("test/README.md"))
	assert.Equal(t, 2, calculatePathDepth("test/folder1/README.md"))
	assert.Equal(t, 3, calculatePathDepth("test/folder1/folder11/README.md"))
	assert.Equal(t, 2, calculatePathDepth("test/folder2/README.md"))

}

func TestBuildTableOfContents(t *testing.T) {
	doc := []byte("# Title\nContent")

	tocNode, tocByte := buildTableOfContents(doc)

	render := renderHTMLMarkdown(tocNode, tocByte)
	mockToc := readFile("test/mock-table-of-content.mock")

	assert.Equal(t, string(mockToc), render)
}
