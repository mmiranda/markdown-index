[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![codecov](https://codecov.io/gh/mmiranda/markdown-index/branch/main/graph/badge.svg?token=3B0LZEZ6XN)](https://codecov.io/gh/mmiranda/markdown-index)
[![https://goreportcard.com/report/github.com/mmiranda/markdown-index](https://goreportcard.com/badge/github.com/mmiranda/markdown-index)](https://goreportcard.com/report/github.com/mmiranda/markdown-index)
![[Test](https://github.com/mmiranda/markdown-index/actions/workflows/test-coverage.yml)](https://github.com/mmiranda/markdown-index/actions/workflows/test-coverage.yml/badge.svg)


# markdown index
Markdown-index is a library to help you generate a global index for multiple markdown files recursively in a directory, containing a summary of every file found.

![](show.gif)

## Installation

The easiest way to install it is using Homebrew:

```bash
brew tap mmiranda/mdindex
brew install markdown-index
```

If you prefer, you also can download the latest binary on the [release section](https://github.com/mmiranda/markdown-index/releases), or simply use the pre-built [dockerfile image](#dockerfile)

## Usage
You can use this tool using multiple ways:

### Running Local
```bash
cd some-directory
markdown-index
```

### Dockerfile
```bash
docker pull ghcr.io/mmiranda/markdown-index:latest
docker run --rm -it -v /path/to/root/md/files:/data ghcr.io/mmiranda/markdown-index:latest
```


After running the command, a new markdown file will be created containing a summary of every other file found.

## Contributing
Contributions, issues, and feature requests are welcome!

Give a ⭐️ if you like this project!

## License
[MIT](https://choosealicense.com/licenses/mit/)
