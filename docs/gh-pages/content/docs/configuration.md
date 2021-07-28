---
title: "Configuration"
description: "Shows `mllint` can be configured, including YAML and TOML examples"
weight: 7
# summary: ""
---

`mllint` can be configured either using a `.mllint.yml` file or through the project's `pyproject.toml`. This allows you to:
- selectively disable specific linting rules or categories using their slug
- define custom linting rules
- configure specific settings for various linting rules.

See the code snippets and commands provided below for examples of such configuration files.

### Commands

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

### YAML

An example `.mllint.yml` that disables some rules looks as follows:

```yaml
rules:
  disabled:
    - version-control/code/git
    - dependency-management/single
```

Similar to the `describe` command, this also matches partial slugs. So, to disable all rules regarding version controlling data, use `version-control/data`.

### TOML

If no `.mllint.yml` is found, `mllint` searches the project's `pyproject.toml` for a `[tool.mllint]` section. TOML has a slightly different syntax, but the structure is otherwise the same as the config in the YAML file. 

An example `pyproject.toml` configuration of `mllint` is as follows. Note that it is identical to the YAML example above.

```toml
[tool.mllint.rules]
disabled = ["version-control/code/git", "dependency-management/single"]
```