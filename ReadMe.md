# [`mllint` — Linter for Machine Learning projects](https://bvobart.github.io/mllint/)

<p align="center">
  <img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/bvobart/mllint/Build%20mllint%20and%20upload%20to%20PyPI">
  <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/bvobart/mllint">
  <a href="https://pkg.go.dev/github.com/bvobart/mllint"><img src="https://pkg.go.dev/badge/github.com/bvobart/mllint.svg" alt="Go Reference"></a>
  <a href="https://codecov.io/gh/bvobart/mllint"><img alt="Code coverage" src="https://codecov.io/gh/bvobart/mllint/branch/main/graph/badge.svg?token=K9PJMGMFVI"/></a>
  <a href="https://goreportcard.com/report/github.com/bvobart/mllint"><img alt="GoReportCard" src="https://goreportcard.com/badge/github.com/bvobart/mllint"/></a>
  <img alt="Platform" src="https://img.shields.io/badge/platform-Linux%20%7C%20MacOS%20%7C%20Windows-informational">
</p>
<p align="center">
  <a href="https://pypi.org/project/mllint/"><img alt="PyPI" src="https://img.shields.io/pypi/v/mllint"></a>
  <a href="https://pypi.org/project/mllint/"><img alt="PyPI - Status" src="https://img.shields.io/pypi/status/mllint"></a>
  <a href="https://pypi.org/project/mllint/"><img alt="PyPI - Downloads - Daily" src="https://img.shields.io/pypi/dd/mllint"></a>
  <a href="https://pypi.org/project/mllint/"><img alt="PyPI - Downloads - Monthly" src="https://img.shields.io/pypi/dm/mllint"></a>
  <a href="https://pypi.org/project/mllint/"><img alt="PyPI - Python Version" src="https://img.shields.io/pypi/pyversions/mllint"></a>
</p>

> ## Attention! This tool is no longer maintained <br/>
> As detailed below, I wrote `mllint` during my MSc thesis in Computer Science between February and October of 2021. I have since graduated and am now no longer developing or actively maintaining this package.
> 
> `mllint` does still work, so feel free to use it! If you find any bugs, feel free to create an issue, I still receive notifications of new issues and there's a good chance that I'll look at them in my free time, but I won't guarantee a timely response or a fix for your issue.
>
> For those interested in the research output produced in my MSc thesis:
> - Full MSc Thesis: http://resolver.tudelft.nl/uuid:b20883f8-a921-487a-8a65-89374a1f3867
> - [The Prevalence of Code Smells in Machine Learning projects](https://arxiv.org/abs/2103.04146) <br />
>   Bart van Oort, Luís Cruz, Maurício Aniche, Arie van Deursen <br />
>   published at WAIN 2021 (1st Workshop on AI Engineering - Software Engineering for AI, co-located with ICSE)
> - ["Project smells" -- Experiences in Analysing the Software Quality of ML Projects with mllint](https://arxiv.org/abs/2201.08246) <br />
>   Bart van Oort, Luís Cruz, Babak Loni, Arie van Deursen <br />
>   published at ICSE SEIP 2022

`mllint` is a command-line utility to evaluate the technical quality of Machine Learning (ML) and Artificial Intelligence (AI) projects written in Python by analysing the project's source code, data and configuration of supporting tools. `mllint` aims to ...

- ... help data scientists and ML engineers in creating and maintaining production-grade ML and AI projects, both on their own personal computers as well as on CI.
- ... help ML practitioners inexperienced with Software Engineering (SE) techniques explore and make effective use of battle-hardended SE for ML tools in the Python ecosystem.
- ... help ML project managers assess the quality of their ML and AI projects and receive recommendations on what aspects of their projects they should focus on improving.

`mllint` does this by measuring the project's adherence to ML best practices, as collected and deduced from [SE4ML](https://se-ml.github.io/) and Google's [Rules for ML](https://developers.google.com/machine-learning/guides/rules-of-ml). Note that these best practices are rather high-level, while `mllint` aims to give practical, down-to-earth advice to its users. `mllint` may therefore be somewhat opinionated, as it tries to advocate specific tools to best fit these best practices. However, `mllint` aims to only recommend open-source tooling and publically verifiable practices. Feedback is of course always welcome!

`mllint` was created during my MSc thesis in Computer Science at the Software Engineering Research Group ([SERG](https://se.ewi.tudelft.nl/)) at [TU Delft](https://tudelft.nl/) and [ING](https://www.ing.com/)'s [AI for FinTech Research Lab](https://se.ewi.tudelft.nl/ai4fintech/) on the topic of [**Code Smells & Software Quality in Machine Learning projects**](http://resolver.tudelft.nl/uuid:b20883f8-a921-487a-8a65-89374a1f3867).

<p align="center"><img src="./docs/gh-pages/static/example-run.svg"></p>

> See [`docs/example-report.md`](docs/example-report.md) for the full report generated for this example project.
>
> See also the [`mllint-example-projects`](https://github.com/bvobart/mllint-example-projects) repository to explore the reports of an example project using `mllint` to measure and improve its project quality over several iterations.
>
> See also [`mllint`'s website](https://bvobart.github.io/mllint/) for online documentation of all of its linting rules and categories.

---

## Installation

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

## Usage

`mllint` is designed to be used both on your personal computer as well as on CI systems. So, open a terminal in your project folder and run one of the following commands, or add it to your project's CI script.

To run `mllint` on the project in the current folder, simply run:
```sh
mllint
```

To run `mllint` on a project in another folder, simply run:
```sh
mllint path/to/my-ml-project
```

`mllint` will analyse your project and create a Markdown-formatted report of its analysis. By default, this will be pretty printed to your terminal. 

If you instead prefer to export the raw Markdown text to a file, which may be particularly useful when running on CI, the `--output` or `-o` flag and provide a filename. `mllint` does not overwrite the destination file if it already exists, unless `--force` or `-f` is used. For example:
```sh
mllint --output report.md
```

Using `-` (a dash) as the filename prints the raw Markdown directly to your terminal:
```sh
mllint -o -
```

In CI scripts, such raw markdown output (whether as a file or printed to the standard output) can be used to e.g. make comments on pull/merge requests or create Wiki pages on your repository.

See [docs/example-report.md](docs/example-report.md) for an example of a report that `mllint` generates, or explore those generated for the [example projects](https://github.com/bvobart/mllint-example-projects).

Of course, feel free to explore `mllint help` for more information about its commands and to discover additional flags that can be used.

### Linters, Categories and Rules

`mllint` analyses your project by evaluating several categories of linting rules. Each category, as well as each rule, has a 'slug', i.e., a lowercased piece of text with dashes or slashes for spaces, e.g., `code-quality/pylint/no-issues`. This slug identifies a rule and is often (if not always) displayed next to the category or rule that it references.

#### Command-line

To list all available (implemented) categories and linting rules, run:
```sh
mllint list all
```

To list all enabled linting rules, run (optionally providing the path to the project's folder):
```sh
mllint list enabled
```

By default, all of `mllint`'s rules are enabled. See [Configuration](#configuration) to learn how to selectively disable certain rules.

To learn more about a certain rule or category, use `mllint describe` along with the slug of the category or rule:
```sh
# Describe the Version Control category. This will also list the rules that it checks.
mllint describe version-control

# Use the exact slug of a rule to describe one rule,
# e.g., the rule on DVC usage in the Version Control category
mllint describe version-control/data/dvc

# Use a partial slug to describe all rules whose slug starts with this snippet, 
# e.g., all rules about version controlling data
mllint describe version-control/data
```

#### Online Documentation

Alternatively, visit the [Categories](https://bvobart.github.io/mllint/docs/categories/) and [Rules](https://bvobart.github.io/mllint/docs/rules/) pages on [`mllint`'s website](https://bvobart.github.io/mllint/) to view the latest online documentation of these rules.


### Custom linting rules

It is also possible to define your own custom linting rules by implementing a script or program that `mllint` will run while performing its analysis.
These custom rules need to be defined in `mllint`'s configuration. For more information on how to do this, see `mllint describe custom` or view the documentation online [here](https://bvobart.github.io/mllint/docs/categories/custom/).

---

## Configuration

`mllint` can be configured either using a `.mllint.yml` file or through the project's `pyproject.toml`. This allows you to:
- selectively disable specific linting rules or categories using their slug
- define custom linting rules
- configure specific settings for various linting rules.

See the code snippets and commands provided below for examples of such configuration files.

#### Commands

To print `mllint`'s current configuration in YAML format, run (optionally providing the path to the project's folder):
```sh
mllint config
```

To print `mllint`'s default configuration in YAML format, run (unless there is a folder called `default` in the current directory):
```sh
mllint config default
```

To create a `.mllint.yml` file from `mllint`'s default configuration, run:
```sh
mllint config default -q > .mllint.yml
```

#### YAML

An example `.mllint.yml` that disables some rules looks as follows:

```yaml
rules:
  disabled:
    - version-control/code/git
    - dependency-management/single
```

Similar to the `describe` command, this also matches partial slugs. So, to disable all rules regarding version controlling data, use `version-control/data`.

#### TOML

If no `.mllint.yml` is found, `mllint` searches the project's `pyproject.toml` for a `[tool.mllint]` section. TOML has a slightly different syntax, but the structure is otherwise the same as the config in the YAML file. 

An example `pyproject.toml` configuration of `mllint` is as follows. Note that it is identical to the YAML example above.

```toml
[tool.mllint.rules]
disabled = ["version-control/code/git", "dependency-management/single"]
```

---

## Getting Started (development)

While `mllint` is a tool for the Python ML ecosystem and distributed through PyPI, it is actually written in Go, compiled to a static binary and published as platform-specific Python wheels. 

To run `mllint` from source, install the latest version of Go for your operating system, then clone this repository and run `go run .` in the root of this repository. Use `go test ./...` or execute `test.sh` to run all of `mllint`'s tests.

To test compiling and packaging `mllint` into a Python wheel for your current platform, run `test.package.sh`.
