package markdown

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

func TestShouldReadFirstFileContent(t *testing.T) {
	files := findFiles("./test", []string{})

	var file RawMarkdown
	file.path = files[0]
	content := file.readFile()

	assert.Contains(t, string(content), "# Root Level Markdown")
}

func TestGetFirstParagraph(t *testing.T) {
	files := findFiles("./test", []string{})
	var file RawMarkdown
	file.path = files[0]
	title := "# " + file.FirstParagraph().title
	content := file.FirstParagraph().content

	assert.Equal(t, "This is a sample paragraph text for test purpose only. This paragraph will be used as an abstract on the global TOC.", content)
	assert.Equal(t, "# Root Level Markdown", title)
}

func TestGetFirstParagraphInEveryFile(t *testing.T) {
	files := findFiles("./test", []string{})
	var file RawMarkdown
	for key, _ := range files {
		file.path = files[key]
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

	createMDFile(filePath, contentString)
	assert.FileExists(t, filePath)

	var file RawMarkdown
	file.path = filePath
	content := file.readFile()

	assert.Equal(t, contentString, string(content))

	deleteFile(filePath)
}

func TestCreateMultiLineFile(t *testing.T) {
	filePath := "/tmp/test.md"
	contentString := "# Hello, this is a test\nThis is a sample text"

	createMDFile(filePath, contentString)
	assert.FileExists(t, filePath)

	var file RawMarkdown
	file.path = filePath
	content := file.readFile()

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
	var file RawMarkdown
	file.path = "test/mock-toc-toc.mock"
	mockFile := file.readFile()

	contentNode, contentByte := buildIndexContent("./test", []string{})

	assert.Equal(t, string(mockFile), contentNode.renderPlainMarkdown(contentByte))
}

func TestCompareFinalFileHTMLContent(t *testing.T) {
	var file RawMarkdown
	file.path = "test/mock-toc-toc-html.mock"
	mockFile := file.readFile()

	contentNode, contentByte := buildIndexContent("./test", []string{})

	assert.Equal(t, string(mockFile), contentNode.renderHTMLMarkdown(contentByte))
}

func TestCompareFinalFilePlainContentWithIgnore(t *testing.T) {
	var file RawMarkdown
	file.path = "test/mock-toc-toc-ignored.mock"
	mockFile := file.readFile()

	contentNode, contentByte := buildIndexContent("./test", []string{"folder2"})

	assert.Equal(t, string(mockFile), contentNode.renderPlainMarkdown(contentByte))
}

func TestCompareFinalFileHTMLContentWithIgnore(t *testing.T) {
	contentNode, contentByte := buildIndexContent("./test", []string{"folder2"})

	var file RawMarkdown
	file.path = "test/mock-toc-toc-html-ignored.mock"
	mockFile := file.readFile()

	assert.Equal(t, string(mockFile), contentNode.renderHTMLMarkdown(contentByte))
}

func TestCompareFinalFilePlainBytes(t *testing.T) {
	var finalFile, mockFile RawMarkdown
	finalFile.path = "/tmp/test-toc.md"

	contentNode, contentString := buildIndexContent("./test", []string{})
	createMDFile(finalFile.path, contentNode.renderPlainMarkdown(contentString))
	fileGenerated := finalFile.readFile()

	mockFile.path = "./test/mock-toc-toc.mock"
	mock := mockFile.readFile()
	assert.True(t, bytes.Equal(fileGenerated, mock))

	deleteFile(finalFile.path)
}

func TestCompareFinalFileHTMLBytes(t *testing.T) {
	var finalFile, mockFile RawMarkdown
	finalFile.path = "/tmp/test-toc.md"

	contentNode, contentByte := buildIndexContent("./test", []string{})
	createMDFile(finalFile.path, contentNode.renderHTMLMarkdown(contentByte))
	fileGenerated := finalFile.readFile()

	mockFile.path = "./test/mock-toc-toc-html.mock"
	mockFileContent := mockFile.readFile()
	assert.True(t, bytes.Equal(fileGenerated, mockFileContent))

	deleteFile(finalFile.path)
}

func TestFilterAbstractHeading(t *testing.T) {
	var file RawMarkdown
	file.path = "test/README.md"
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
	var file RawMarkdown

	file.path = "test/README.md"
	assert.Equal(t, 1, file.calculatePathDepth())
	file.path = "test/folder1/README.md"
	assert.Equal(t, 2, file.calculatePathDepth())
	file.path = "test/folder1/folder11/README.md"
	assert.Equal(t, 3, file.calculatePathDepth())
	file.path = "test/folder2/README.md"
	assert.Equal(t, 2, file.calculatePathDepth())

}

func TestBuildTableOfContents(t *testing.T) {
	var doc, mock RawMarkdown
	doc.content = []byte("# Title\nContent")

	tocNode, tocByte := doc.buildTableOfContents()

	render := tocNode.renderHTMLMarkdown(tocByte)
	mock.path = "test/mock-table-of-content.mock"
	mockToc := mock.readFile()

	assert.Equal(t, string(mockToc), render)
}
