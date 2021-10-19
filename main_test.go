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

	title := "# " + getFirstParagraph(files[0]).title
	content := getFirstParagraph(files[0]).content

	assert.Equal(t, "This is a sample paragraph text for test purpose only. This paragraph will be used as an abstract on the global TOC.", content)
	assert.Equal(t, "# Root Level Markdown", title)
}

func TestGetFirstParagraphInEveryFile(t *testing.T) {
	files := findFiles("./test", []string{})

	for key, _ := range files {
		content := getFirstParagraph(files[key]).content
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
	contentString := buildIndexContent("./test", []string{})

	mockFile := readFile(mockFilePath)

	assert.Equal(t, string(mockFile), renderPlainMarkdown(contentString))
}

func TestCompareFinalFileHTMLContent(t *testing.T) {
	mockFilePath := "test/mock-toc-toc-html.mock"
	contentString := buildIndexContent("./test", []string{})

	mockFile := readFile(mockFilePath)

	assert.Equal(t, string(mockFile), renderHTMLMarkdown(contentString))
}

func TestCompareFinalFilePlainContentWithIgnore(t *testing.T) {
	mockFilePath := "test/mock-toc-toc-ignored.mock"

	contentString := buildIndexContent("./test", []string{"folder2"})

	mockFile := readFile(mockFilePath)

	assert.Equal(t, string(mockFile), renderPlainMarkdown(contentString))
}

func TestCompareFinalFileHTMLContentWithIgnore(t *testing.T) {
	mockFilePath := "test/mock-toc-toc-html-ignored.mock"

	contentString := buildIndexContent("./test", []string{"folder2"})

	mockFile := readFile(mockFilePath)

	assert.Equal(t, string(mockFile), renderHTMLMarkdown(contentString))
}
func TestCompareFinalFilePlainBytes(t *testing.T) {
	filePath := "/tmp/test-toc.md"
	contentDocument := buildIndexContent("./test", []string{})
	createMDFile(filePath, renderPlainMarkdown(contentDocument))
	fileGenerated := readFile(filePath)

	mockFile := readFile("./test/mock-toc-toc.mock")
	assert.True(t, bytes.Equal(fileGenerated, mockFile))

	deleteFile(filePath)
}

func TestCompareFinalFileHTMLBytes(t *testing.T) {
	filePath := "/tmp/test-toc.md"
	contentDocument := buildIndexContent("./test", []string{})
	createMDFile(filePath, renderHTMLMarkdown(contentDocument))
	fileGenerated := readFile(filePath)

	mockFile := readFile("./test/mock-toc-toc-html.mock")
	assert.True(t, bytes.Equal(fileGenerated, mockFile))

	deleteFile(filePath)
}

func TestFilterAbstractHeading(t *testing.T) {

	content := FilterHeadingAbstract("Another title", "test/README.md")
	assert.NotEmpty(t, content)

	content = FilterHeadingAbstract("Unexistent title heading", "test/README.md")
	assert.Empty(t, content)
}

func TestContainsHelper(t *testing.T) {
	directories := []string{"A", "B"}

	assert.True(t, contains(directories, "A"))
	assert.False(t, contains(directories, "C"))
}

func TestCalculateDepth(t *testing.T) {
	os.MkdirAll("/tmp/markdown-index/docs/test", os.ModePerm)

	assert.Equal(t, 1, calculatePathDepth("/tmp"))
	createMDFile("/tmp/file.md", "sample")
	assert.Equal(t, 1, calculatePathDepth("/tmp/file.md"))
	assert.Equal(t, 2, calculatePathDepth("/tmp/markdown-index"))
	createMDFile("/tmp/markdown-index/test.file", "sample")
	assert.Equal(t, 2, calculatePathDepth("/tmp/markdown-index/test.file"))
	assert.Equal(t, 3, calculatePathDepth("/tmp/markdown-index/docs"))
	createMDFile("/tmp/markdown-index/docs/readme.md", "sample")
	assert.Equal(t, 3, calculatePathDepth("/tmp/markdown-index/docs/readme.md"))
	assert.Equal(t, 4, calculatePathDepth("/tmp/markdown-index/docs/test"))
	assert.Equal(t, 1, calculatePathDepth("./"))
	assert.Equal(t, 1, calculatePathDepth("test/readme.md"))

	os.RemoveAll("/tmp/markdown-index")

}
