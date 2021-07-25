---
title: "Installation"
description: "Options for installing `mllint`"
weight: 5
summary: "mllint is compiled for Linux, MacOS and Windows and is published to [PyPI](https://pypi.org/project/mllint/), so it can be installed using `pip install -U mllint` Alternatively, use one of the Docker containers at `bvobart/mllint`"
---

`mllint` is compiled for Linux, MacOS and Windows, both 64 and 32 bit x86 (MacOS 64-bit only), as well as 64-bit ARM on Linux and MacOS (Apple M1).

`mllint` is published to [PyPI](https://pypi.org/project/mllint/), so it can be installed globally or in your current environment using `pip`:
```sh
pip install --upgrade mllint
```

Alternatively, to add `mllint` to an existing project, if your project uses Poetry for its dependencies:
```sh
poetry add --dev mllint
```

Or if your project uses Pipenv:
```sh
pipenv install --dev mllint
```

### Tools

`mllint` has a soft dependency on several Python tools that it uses for its analysis. While `mllint` will recommend that you place these tools in your project's development dependencies, these tools are listed as optional dependencies of `mllint` and can be installed along with `mllint` using:

```sh
pip install --upgrade mllint[tools]
```

### Docker

There are also `mllint` Docker containers available on [Docker Hub](https://hub.docker.com/r/bvobart/mllint) at `bvobart/mllint` for Python 3.6, 3.7, 3.8 and 3.9. These may particularly be helpful when running `mllint` in CI environments, such as Gitlab CI or Github Actions. See the Docker Hub for a full list of available tags that can be used.

The Docker containers require that you mount the folder with your project onto the container as a volume on `/app`. Here is an example of how to use this Docker container, assuming that your project is in the current folder. Replace `$(pwd)` with the full path to your project folder if it is somewhere else.

```sh
docker run -it --rm -v $(pwd):/app bvobart/mllint:latest
```
