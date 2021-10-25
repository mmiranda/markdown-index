package markdown

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TESTDIR = "../test"
)

func TestShouldFindReadme(t *testing.T) {
	files := findFiles(TESTDIR, []string{})
	assert.Equal(t, 6, len(files))
}

func TestShouldReadFirstFileContent(t *testing.T) {
	files := findFiles(TESTDIR, []string{})

	content := files[0].readFile()

	assert.Contains(t, string(content), "# Root Level Markdown")
}

func TestGetFirstParagraph(t *testing.T) {
	files := findFiles(TESTDIR, []string{})

	title := "# " + files[0].FirstParagraph().title
	content := files[0].FirstParagraph().content

	assert.Equal(t, "This is a sample paragraph text for test purpose only. This paragraph will be used as an abstract on the global TOC.", content)
	assert.Equal(t, "# Root Level Markdown", title)
}

func TestGetFirstParagraphInEveryFile(t *testing.T) {
	files := findFiles(TESTDIR, []string{})

	for _, file := range files {
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

	createMDFile(filePath, contentString, false)
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

	createMDFile(filePath, contentString, false)
	assert.FileExists(t, filePath)

	var file RawMarkdown
	file.path = filePath
	content := file.readFile()

	assert.Equal(t, contentString, string(content))
	deleteFile(filePath)
}

func TestCompareFinalFilePlainContent(t *testing.T) {
	var file RawMarkdown
	file.path = TESTDIR + "/final-file-without-metadata.mock"
	mockFile := file.readFile()

	contentNode, contentByte := buildIndexContent(TESTDIR, []string{})

	assert.Equal(t, string(mockFile), contentNode.renderPlainMarkdown(contentByte))
}

func TestCompareFinalFileHTMLContent(t *testing.T) {
	var file RawMarkdown
	file.path = TESTDIR + "/final-file-html.mock"
	mockFile := file.readFile()

	contentNode, contentByte := buildIndexContent(TESTDIR, []string{})

	assert.Equal(t, string(mockFile), contentNode.renderHTMLMarkdown(contentByte))
}

func TestCompareFinalFilePlainContentWithIgnore(t *testing.T) {
	var file RawMarkdown
	file.path = TESTDIR + "/final-file-ignored.mock"
	mockFile := file.readFile()

	contentNode, contentByte := buildIndexContent(TESTDIR, []string{"folder2"})

	assert.Equal(t, string(mockFile), contentNode.renderPlainMarkdown(contentByte))
}

func TestCompareFinalFileHTMLContentWithIgnore(t *testing.T) {
	contentNode, contentByte := buildIndexContent(TESTDIR, []string{"folder2"})

	var file RawMarkdown
	file.path = TESTDIR + "/final-file-html-ignored.mock"
	mockFile := file.readFile()

	assert.Equal(t, string(mockFile), contentNode.renderHTMLMarkdown(contentByte))
}

func TestCompareFinalFilePlainBytes(t *testing.T) {
	var finalFile, mockFile RawMarkdown
	finalFile.path = "/tmp/test-toc.md"

	contentNode, contentString := buildIndexContent(TESTDIR, []string{})
	createMDFile(finalFile.path, contentNode.renderPlainMarkdown(contentString), false)
	fileGenerated := finalFile.readFile()

	mockFile.path = TESTDIR + "/final-file-without-metadata.mock"
	mock := mockFile.readFile()
	assert.True(t, bytes.Equal(fileGenerated, mock))

	deleteFile(finalFile.path)

}

func TestFinalFileIdempotency(t *testing.T) {
	var finalFile, mockFile RawMarkdown
	finalFile.path = TESTDIR + "/" + "test-file.md"
	contentNode, contentString := buildIndexContent(TESTDIR, []string{})
	createMDFile(finalFile.path, contentNode.renderPlainMarkdown(contentString), true)
	fileGenerated := finalFile.readFile()

	mockFile.path = TESTDIR + "/final-file-with-metadata.mock"
	mock := mockFile.readFile()
	assert.True(t, bytes.Equal(fileGenerated, mock))

	contentNode, contentString = buildIndexContent(TESTDIR, []string{})
	createMDFile(finalFile.path, contentNode.renderPlainMarkdown(contentString), true)

	fileGenerated = finalFile.readFile()

	assert.Equal(t, string(mock), string(fileGenerated))
	deleteFile(finalFile.path)
}

func TestCompareFinalFileHTMLBytes(t *testing.T) {
	var finalFile, mockFile RawMarkdown
	finalFile.path = "/tmp/test-toc.md"

	contentNode, contentByte := buildIndexContent(TESTDIR, []string{})
	createMDFile(finalFile.path, contentNode.renderHTMLMarkdown(contentByte), false)
	fileGenerated := finalFile.readFile()

	mockFile.path = TESTDIR + "/final-file-html.mock"
	mockFileContent := mockFile.readFile()
	assert.True(t, bytes.Equal(fileGenerated, mockFileContent))

	deleteFile(finalFile.path)
}

func TestFilterAbstractHeading(t *testing.T) {
	var file RawMarkdown
	file.path = TESTDIR + "/README.md"
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

	file.path = TESTDIR + "/README.md"
	assert.Equal(t, 1, file.calculatePathDepth())

	absPath, _ := filepath.Abs(TESTDIR)
	fmt.Println(absPath)
	file.path = absPath + "/README.md"
	// assert.Equal(t, 1, file.calculatePathDepth())
	file.path = TESTDIR + "/folder1/README.md"
	assert.Equal(t, 2, file.calculatePathDepth())
	file.path = TESTDIR + "/folder1/folder11/README.md"
	assert.Equal(t, 3, file.calculatePathDepth())
	file.path = TESTDIR + "/folder2/README.md"
	assert.Equal(t, 2, file.calculatePathDepth())

}

func TestBuildTableOfContents(t *testing.T) {
	var doc, mock RawMarkdown
	doc.content = []byte("# Title\nContent")

	tocNode, tocByte := doc.buildTableOfContents()

	render := tocNode.renderHTMLMarkdown(tocByte)
	mock.path = TESTDIR + "/table-of-content.mock"
	mockToc := mock.readFile()

	assert.Equal(t, string(mockToc), render)
}

func TestCobraExecutionFlow(t *testing.T) {
	directory := TESTDIR
	output := "my-test-file.md"
	Execute(output, directory)

	var finalFile RawMarkdown
	finalFile.path = directory + "/" + output
	fileGenerated := finalFile.readFile()

	var mockFile RawMarkdown
	mockFile.path = TESTDIR + "/final-file-with-metadata.mock"
	mock := mockFile.readFile()

	assert.True(t, bytes.Equal(fileGenerated, mock))
	assert.FileExists(t, directory+"/"+output)

	deleteFile(directory + "/" + output)
}

func TestFileWithoutMetadata(t *testing.T) {
	var file RawMarkdown
	file.path = TESTDIR + "/README.md"

	assert.False(t, file.isFileAutoGenerated())
}
func TestFileWittMetadata(t *testing.T) {
	var file RawMarkdown
	file.path = TESTDIR + "/final-file-with-metadata.mock"

	assert.True(t, file.isFileAutoGenerated())
}

func TestAddMetadataInFile(t *testing.T) {
	var file, mock RawMarkdown
	file.content = []byte("# Root Level Markdown\n\nThis file should be completely ignored because it has metadata Title == markdown-index\n")

	file.addMetadata()

	mock.path = TESTDIR + "/file-with-metadata.md"

	assert.Equal(t, string(mock.readFile()), string(file.content))
}
