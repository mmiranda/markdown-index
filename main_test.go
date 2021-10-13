package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldFindReadme(t *testing.T) {
	files := findFiles("./test")
	assert.Equal(t, 4, len(files))
}

func TestShouldReadFileContent(t *testing.T) {
	files := findFiles("./test")

	content := readFile(files[0])

	assert.Contains(t, string(content), "# Root Level Markdown")
}

func TestGetFirstParagraph(t *testing.T) {
	files := findFiles("./test")

	title := getFirstParagraph(files[0]).title
	content := getFirstParagraph(files[0]).content

	assert.Equal(t, "[0] This is a sample paragraph text for test purpose only. This paragraph will be used as an abstract on the global TOC.", content)
	assert.Equal(t, "# Root Level Markdown", title)
}

func TestGetFirstParagraphInEveryFile(t *testing.T) {
	files := findFiles("./test")

	for key, _ := range files {
		// fmt.Println(key)
		// fmt.Println(file)
		content := getFirstParagraph(files[key]).content
		assert.Equal(t, fmt.Sprintf("[%s] This is a sample paragraph text for test purpose only. This paragraph will be used as an abstract on the global TOC.", strconv.Itoa(key)), content)
	}

	assert.True(t, true)
}

func TestCreateSingleLineFile(t *testing.T) {
	filePath := "/tmp/test.md"
	contentString := []string{"# Hello, this is a test"}
	// contentString = append(contentString, "\nMultiline test")
	createMDFile(filePath, contentString)
	assert.FileExists(t, filePath)

	content := readFile(filePath)

	assert.Equal(t, strings.Join(contentString[:], ""), string(content))

	deleteFile(filePath)
}

func TestCreateMultiLineFile(t *testing.T) {
	filePath := "/tmp/test.md"
	contentString := []string{
		"# Hello, this is a test",
		"\n",
		"This is a sample text",
	}
	createMDFile(filePath, contentString)
	assert.FileExists(t, filePath)

	content := readFile(filePath)

	assert.Equal(t, strings.Join(contentString[:], ""), string(content))
	deleteFile(filePath)
}

func deleteFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		panic("Test failed deleting file")
	}
}

func TestCompareFinalFileContent(t *testing.T) {
	mockFilePath := "test/mock-toc-toc.mock"
	contentString := buildIndexContent("test")

	mockFile := readFile(mockFilePath)

	assert.Equal(t, string(mockFile), strings.Join(contentString[:], ""))

}

func TestCompareFinalFileBytes(t *testing.T) {
	filePath := "/tmp/test-toc.md"
	contentString := buildIndexContent("test")
	createMDFile(filePath, contentString)
	file1 := readFile(filePath)

	mockFile := readFile("test/mock-toc-toc.mock")

	assert.True(t, bytes.Equal(file1, mockFile))

	deleteFile(filePath)
}

func TestFilterAbstract(t *testing.T) {

	content := FilterHeadingAbstract("Another title", "test/README.md")
	assert.NotEmpty(t, content)

	content = FilterHeadingAbstract("Unexistent title heading", "test/README.md")
	assert.Empty(t, content)
}
