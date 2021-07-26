---
title: "Usage"
description: "Shows how to use `mllint` and which commands can be run"
weight: 6
summary: "mllint is compiled for Linux, MacOS and Windows and is published to [PyPI](https://pypi.org/project/mllint/), so it can be installed using `pip install -U mllint` Alternatively, use one of the Docker containers at `bvobart/mllint`"
---

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

See this [`example-report.md`](https://github.com/bvobart/mllint/blob/main/docs/example-report.md) for an example of a report that `mllint` generates, or explore those generated for the [example projects](https://github.com/bvobart/mllint-example-projects).

Of course, feel free to explore `mllint help` for more information about its commands and to discover additional flags that can be used.

### Linters, Categories and Rules

`mllint` analyses your project by evaluating several categories of linting rules. Each category, as well as each rule, has a 'slug', i.e., a lowercased piece of text with dashes or slashes for spaces, e.g., `code-quality/pylint/no-issues`. This slug identifies a rule and is often (if not always) displayed next to the category or rule that it references.

To list all available (implemented) categories and linting rules, run:
```sh
mllint list all
```

To list all enabled linting rules, run (optionally providing the path to the project's folder):
```sh
mllint list enabled
```

By default, all of `mllint`'s rules are enabled. See [Configuration](../configuration) to learn how to selectively disable certain rules.

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

### Custom linting rules

It is also possible to define your own custom linting rules by implementing a script or program that `mllint` will run while performing its analysis.
These custom rules need to be defined in `mllint`'s configuration. For more information on how to do this, see `mllint describe custom`.
