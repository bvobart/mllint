# `mllint` — Linter for Machine Learning projects

<p align="center">
  <img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/bvobart/mllint/Build%20mllint%20and%20upload%20to%20PyPI">
  <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/bvobart/mllint">
  <img alt="Platform" src="https://img.shields.io/badge/platform-Linux%20%7C%20MacOS%20%7C%20Windows-informational">
</p>
<p align="center">
  <a href="https://pypi.org/project/mllint/"><img alt="PyPI" src="https://img.shields.io/pypi/v/mllint"></a>
  <a href="https://pypi.org/project/mllint/"><img alt="PyPI - Status" src="https://img.shields.io/pypi/status/mllint"></a>
  <a href="https://pypi.org/project/mllint/"><img alt="PyPI - Downloads - Daily" src="https://img.shields.io/pypi/dd/mllint"></a>
  <a href="https://pypi.org/project/mllint/"><img alt="PyPI - Downloads - Monthly" src="https://img.shields.io/pypi/dm/mllint"></a>
  <a href="https://pypi.org/project/mllint/"><img alt="PyPI - Python Version" src="https://img.shields.io/pypi/pyversions/mllint"></a>
</p>

`mllint` is a command-line utility to evaluate the quality of Machine Learning (ML) projects by means of static analysis of the project's repository. It measures the project's adherence to ML best practices, as collected and deduced from se4ml.github.io and Google's [Rules for ML](https://developers.google.com/machine-learning/guides/rules-of-ml).

TODO: write overview of linting rules or implement `mllint describe <rule>`

---

## Getting Started

---

#### Installing `mllint`

`mllint` works on Linux, MacOS and Windows, and is published to [PyPI](https://pypi.org/project/mllint/), so it can be installed using `pip`:
```sh
pip install mllint
```

#### Running `mllint`

To run `mllint` in its default configuration, use one of the following commands:
```sh
# Run `mllint` on the project in the current folder
mllint

# Run `mllint` on the project in projects/my-ml-project
mllint projects/my-ml-project
```

#### Linters and rules

To list all available or all enabled linting rules, use one of the following commands:
```sh
# List all available (implemented) linting rules
mllint list all

# List only the enabled rules for the project in the current folder.
mllint list enabled

# or for a project in projects/my-ml-project
mllint list enabled projects/my-ml-project
```

---

## Configuration

---

#### YAML

`mllint` can be configured using a `.mllint.yml` file that should be placed at the root of the project directory. This is a YAML file in which you can disable specific rules / linters, as well as configure specific settings for various linters.

An example `.mllint.yml` looks as follows:

```yaml
rules:
  disabled:
    - use-git # disables the 'use-git' linter
    - use-dependency-manager/single # disables the 'single' rule of the 'use-dependency-manager' linter.
    # - use-dependency-manager # this would disable the 'use-dependency-manager' linter and all of its rules entirely.
```

#### TOML

Alternatively, if no `.mllint.yml` exists, `mllint` can be configured from the `pyproject.toml` file in the root of the project. This should be done in the `[tool.mllint]` section. TOML has a slightly different syntax, but the structure is otherwise the same as the config in the YAML file. The example below is identical to the YAML example above.

```toml
[tool.mllint]
[tool.mllint.rules]
disabled = ["use-git", "use-dependency-manager/single"]
```

---

Here are some useful commands related to configuration:
```sh
# Print the configuration of the project in the current folder
mllint config

# Print the configuration of the project in projects/my-ml-project
mllint config projects/my-ml-project

# Print the default configuration (unless there's a folder called 'default' in the current dir)
mllint config default

# Create a valid `.mllint.yml` file with the default configuration
mllint config default -q > .mllint.yml
```

---

## Getting Started (development)

Clone this repository and run `go run .` in the root of this repository.
