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

func init() {
	searchHeading = "Abstract"
}

func TestShouldFindReadme(t *testing.T) {
	files := findFiles(TESTDIR, []string{})
	assert.Equal(t, 6, len(files))
}

func TestShouldReadFirstFileContent(t *testing.T) {
	files := findFiles(TESTDIR, []string{})

	content, _ := files[0].readFile()

	assert.Contains(t, string(content), "# Root Level Markdown")
}

func TestGetFirstParagraph2(t *testing.T) {
	files := findFiles(TESTDIR, []string{})

	title := "# " + files[0].getFirstParagraph().title
	content := files[0].getFirstParagraph().content

	assert.Equal(t, "This is a sample paragraph text for test purpose only. This paragraph will be used as an abstract on the global TOC.", content)
	assert.Equal(t, "# Root Level Markdown", title)
}

func TestGetFirstParagraphInEveryFile(t *testing.T) {
	files := findFiles(TESTDIR, []string{})

	for _, file := range files {
		content := file.getFirstParagraph().content
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

	CreateMDFile(filePath, contentString, false)
	assert.FileExists(t, filePath)

	file := newMarkdownFile("", filePath)
	content, _ := file.readFile()

	assert.Equal(t, contentString, string(content))

	DeleteFile(filePath)
}

func TestCreateMultiLineFile(t *testing.T) {
	filePath := "/tmp/test.md"
	contentString := "# Hello, this is a test\nThis is a sample text"

	CreateMDFile(filePath, contentString, false)
	assert.FileExists(t, filePath)

	file := newMarkdownFile("", filePath)
	content, _ := file.readFile()

	assert.Equal(t, contentString, string(content))
	DeleteFile(filePath)
}

func TestCompareFinalFilePlainContent(t *testing.T) {
	file := newMarkdownFile(TESTDIR, TESTDIR+"/final-file-without-metadata.mock")
	mockFile, _ := file.readFile()

	contentNode, contentByte := buildIndexContent(TESTDIR, []string{})
	// fmt.Println(string(mockFile))
	// fmt.Println("------")
	// fmt.Println(contentNode.renderPlainMarkdown(contentByte))

	assert.Equal(t, string(mockFile), contentNode.renderPlainMarkdown(contentByte))
}

func TestCompareFinalFileHTMLContent(t *testing.T) {
	file := newMarkdownFile(TESTDIR, TESTDIR+"/final-file-html.mock")

	mockFile, _ := file.readFile()

	contentNode, contentByte := buildIndexContent(TESTDIR, []string{})

	assert.Equal(t, string(mockFile), contentNode.RenderHTMLMarkdown(contentByte))
}

func TestCompareFinalFilePlainContentWithIgnore(t *testing.T) {
	file := newMarkdownFile(TESTDIR, TESTDIR+"/final-file-ignored.mock")

	mockFile, _ := file.readFile()

	contentNode, contentByte := buildIndexContent(TESTDIR, []string{"folder2"})

	assert.Equal(t, string(mockFile), contentNode.renderPlainMarkdown(contentByte))
}

func TestCompareFinalFileHTMLContentWithIgnore(t *testing.T) {
	contentNode, contentByte := buildIndexContent(TESTDIR, []string{"folder2"})

	file := newMarkdownFile(TESTDIR, TESTDIR+"/final-file-html-ignored.mock")
	mockFile := file.content

	assert.Equal(t, string(mockFile), contentNode.RenderHTMLMarkdown(contentByte))
}

func TestCompareFinalFilePlainBytes(t *testing.T) {
	finalFile := newMarkdownFile("/tmp", "/tmp/test-toc.md")

	contentNode, contentString := buildIndexContent(TESTDIR, []string{})
	CreateMDFile(finalFile.path, contentNode.renderPlainMarkdown(contentString), false)
	fileGenerated, _ := finalFile.readFile()

	mockFile := newMarkdownFile(TESTDIR, TESTDIR+"/final-file-without-metadata.mock")

	mock := mockFile.content
	assert.True(t, bytes.Equal(fileGenerated, mock))

	DeleteFile(finalFile.path)

}

func TestFinalFileIdempotency(t *testing.T) {
	finalFile := newMarkdownFile(TESTDIR, TESTDIR+"test-file.md")
	contentNode, contentString := buildIndexContent(TESTDIR, []string{})
	CreateMDFile(finalFile.path, contentNode.renderPlainMarkdown(contentString), true)
	fileGenerated, _ := finalFile.readFile()

	mockFile := newMarkdownFile(TESTDIR, TESTDIR+"/final-file-with-metadata.mock")
	mock := mockFile.content
	assert.True(t, bytes.Equal(fileGenerated, mock))

	contentNode, contentString = buildIndexContent(TESTDIR, []string{})
	CreateMDFile(finalFile.path, contentNode.renderPlainMarkdown(contentString), true)

	fileGenerated, _ = finalFile.readFile()

	assert.Equal(t, string(mock), string(fileGenerated))
	DeleteFile(finalFile.path)
}

func TestCompareFinalFileHTMLBytes(t *testing.T) {
	finalFile := newMarkdownFile("/tmp", "/tmp/test-toc.md")

	contentNode, contentByte := buildIndexContent(TESTDIR, []string{})
	CreateMDFile(finalFile.path, contentNode.RenderHTMLMarkdown(contentByte), false)
	fileGenerated, _ := finalFile.readFile()

	mockFile := newMarkdownFile(TESTDIR, TESTDIR+"/final-file-html.mock")
	mockFileContent := mockFile.content
	assert.True(t, bytes.Equal(fileGenerated, mockFileContent))

	DeleteFile(finalFile.path)
}

func TestFilterAbstractHeading(t *testing.T) {
	file := newMarkdownFile(TESTDIR, TESTDIR+"/README.md")
	searchHeading = "Another title"
	content := file.getFirstParagraph()
	assert.NotEmpty(t, content)
	assert.Equal(t, "Another title", content.title)

	searchHeading = "Unexistent title heading"
	content = file.getFirstParagraph()
	assert.NotEmpty(t, content)
	assert.Equal(t, "Root Level Markdown", content.title)
}

func TestContainsHelper(t *testing.T) {
	directories := []string{"A", "B"}

	assert.True(t, contains(directories, "A"))
	assert.False(t, contains(directories, "C"))
}

func TestCalculateDepth(t *testing.T) {
	file := newMarkdownFile(TESTDIR, TESTDIR+"/README.md")
	assert.Equal(t, 1, file.calculatePathDepth())

	absPath, _ := filepath.Abs(TESTDIR)

	file = newMarkdownFile(absPath, TESTDIR+"/README.md")
	assert.Equal(t, 1, file.calculatePathDepth())

	file = newMarkdownFile(TESTDIR, TESTDIR+"/folder1/README.md")
	assert.Equal(t, 1, file.calculatePathDepth())

	file = newMarkdownFile(TESTDIR, TESTDIR+"/folder1/folder11/README.md")
	assert.Equal(t, 2, file.calculatePathDepth())

	file = newMarkdownFile(TESTDIR, TESTDIR+"/folder2/README.md")
	assert.Equal(t, 1, file.calculatePathDepth())

}

func TestBuildTableOfContents(t *testing.T) {
	var doc rawMarkdown
	doc.content = []byte("# Title\nContent")

	tocNode, tocByte := doc.buildTableOfContents()

	render := tocNode.RenderHTMLMarkdown(tocByte)
	mock := newMarkdownFile(TESTDIR, TESTDIR+"/table-of-content.mock")
	mockToc, _ := mock.readFile()

	assert.Equal(t, string(mockToc), render)
}

func TestCobraExecutionFlow(t *testing.T) {
	searchHeading = "Abstract"
	fmt.Println(searchHeading)
	directory := TESTDIR
	output := "my-test-file.md"
	Execute(output, directory)

	finalFile := newMarkdownFile(directory, TESTDIR+"/"+output)
	fileGenerated := finalFile.content

	mockFile := newMarkdownFile(TESTDIR, TESTDIR+"/final-file-with-metadata.mock")
	mock := mockFile.content

	assert.True(t, bytes.Equal(fileGenerated, mock))
	assert.Equal(t, string(mock), string(fileGenerated))
	assert.FileExists(t, directory+"/"+output)

	DeleteFile(directory + "/" + output)
}

func TestFileWithoutMetadata(t *testing.T) {
	file := newMarkdownFile(TESTDIR, TESTDIR+"/README.md")

	assert.False(t, file.isFileAutoGenerated())
}
func TestFileWittMetadata(t *testing.T) {
	file := newMarkdownFile(TESTDIR, TESTDIR+"/final-file-with-metadata.mock")

	assert.True(t, file.isFileAutoGenerated())
}

func TestAddMetadataInFile(t *testing.T) {
	var file rawMarkdown
	file.content = []byte("# Root Level Markdown\n\nThis file should be completely ignored because it has metadata Title == markdown-index\n")

	file.content = []byte(AddMetadataPrefix(string(file.content)))

	mock := newMarkdownFile(TESTDIR, TESTDIR+"/file-with-metadata.md")
	mockcontent := mock.content
	assert.Equal(t, string(mockcontent), string(file.content))
}

func TestNewMarkdownFile(t *testing.T) {
	file := newMarkdownFile("/tmp", "/tmp/dir/file.md")
	assert.IsType(t, &rawMarkdown{}, file)
	assert.Equal(t, "dir", file.basedir)
}

func TestFileRealPath(t *testing.T) {

	assert.Equal(t, "./test/file.md", GetFileRealPath(TESTDIR, TESTDIR+"/test/file.md"))
	assert.Equal(t, "./test/test2/file.md", GetFileRealPath(TESTDIR, TESTDIR+"/test/test2/file.md"))
	assert.Equal(t, "./test/file.md", GetFileRealPath(".", "test/file.md"))
}
