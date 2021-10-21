[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org) [![codecov](https://codecov.io/gh/mmiranda/markdown-index/branch/main/graph/badge.svg?token=3B0LZEZ6XN)](https://codecov.io/gh/mmiranda/markdown-index)

# markdown index
Markdown-index is a library to help you generate a global index for multiple markdown files recursively in a directory, containing a summary of every file found.


## Installation

The easiest way to install it is using Homebrew:

```bash
brew tap foobar
brew install mmiranda/markdown-index
```

If you prefer, you also can download the latest binary on the release section, or simply use the pre-built dockerfile image

## Usage

### Using Local
```bash
cd some-directory
markdown-index
```

### Dockerfile
```bash
docker run --rm -it --entrypoint=/bin/sh mmiranda/markdown-index:latest --directory /path/to/directory
```

### Github Actions

You can also use this in your Github Actions workflows - TBD


After running the command, a new markdown file will be created containing a summary of every other file found.

## Contributing
Contributions, issues, and feature requests are welcome!

Give a ⭐️ if you like this project!

## License
[MIT](https://choosealicense.com/licenses/mit/)
