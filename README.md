# web2ppm

>A tiny tool to take a screenshot of an entire web page as [Plain Portable Pixel Map](http://netpbm.sourceforge.net/doc/ppm.html).

## Why

See [How To Build A Blog That Lasts 1000 Years](https://lambdar.me/archives/how-to-build-a-blog-that-lasts-1000-years/) (only Japanese version is currently available.)

## Installation

```bash
$ go get -u github.com/builtinnya/web2ppm/cmd/web2ppm
```

## Usage

```bash
$ web2ppm https://google.com > google.ppm
```

## Development

It is recommended to use [Remote - Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) on [Visual Studio Code](https://code.visualstudio.com/).
See [Developing inside a Container using Visual Studio Code Remote Development](https://code.visualstudio.com/docs/remote/containers).

Otherwise, you can launch a development environment directly using [Docker Compose](https://docs.docker.com/compose/).

```bash
$ docker-compose run --rm web2ppm bash
nya@<container-id>:/go/src/github.com/builtinnya/web2ppm$
```

## License

See [LICENSE](./LICENSE).
