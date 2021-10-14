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
	files := findFiles("./test", []string{})
	assert.Equal(t, 5, len(files))
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

	assert.Equal(t, "[0] This is a sample paragraph text for test purpose only. This paragraph will be used as an abstract on the global TOC.", content)
	assert.Equal(t, "# Root Level Markdown", title)
}

func TestGetFirstParagraphInEveryFile(t *testing.T) {
	files := findFiles("./test", []string{})

	for key, _ := range files {
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
	contentString := buildIndexContent("./test", []string{})

	mockFile := readFile(mockFilePath)

	assert.Equal(t, string(mockFile), strings.Join(contentString[:], ""))
}

func TestCompareFinalFileContentWithIgnore(t *testing.T) {
	mockFilePath := "test/mock-toc-toc-ignored.mock"

	contentString := buildIndexContent("./test", []string{"folder2"})

	mockFile := readFile(mockFilePath)

	assert.Equal(t, string(mockFile), strings.Join(contentString[:], ""))
}

func TestCompareFinalFileBytes(t *testing.T) {
	filePath := "/tmp/test-toc.md"
	contentString := buildIndexContent("./test", []string{})
	createMDFile(filePath, contentString)
	fileGenerated := readFile(filePath)

	mockFile := readFile("./test/mock-toc-toc.mock")
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
