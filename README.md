# Malgo Command & Control

## Development Guide

### Required tools

- [Go](https://go.dev/)
- [gRPC for Go](https://grpc.io/docs/languages/go/quickstart/)
- [Bun](https://bun.sh/)
- [Task](https://taskfile.dev/)
- [Docker](https://www.docker.com/)

### Documentation

Documentation may be found under `/documentation` and is developed using Material for MkDocs.

To launch it, run:

```bash
task docs
```

It will launch a Docker container and bind it [here](http://localhost:8888/)

### Working on Documentation

To launch an interactive, hot-reload version, run the following command:

```bash
task docs-dev
```

This will also launch a Docker container and bind it [here](http://localhost:8889)
