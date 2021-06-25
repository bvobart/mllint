package testing

import "github.com/bvobart/mllint/api"

var RuleHasTests = api.Rule{
	Name: "Project has automated tests",
	Slug: "testing/has-tests",
	Details: `Every ML project should have a set of automated tests to assess the quality, consistency and correctness of their application in a repeatable and reproducible manner.

This rule checks how many test files your project contains. ` + "In accordance with `pytest`'s [conventions](https://docs.pytest.org/en/6.2.x/goodpractices.html#conventions-for-python-test-discovery) for Python tests, test files are Python files starting with `test_` or ending with `_test.py`." + `
Per default, ` + "`mllint`" + ` expects **at least one test file** to be implemented in your project ` + "(i.e. a Python file starting with `test_` or ending with `_test.py`)" + `
and recommends that you have **at least 1 test file** for **every 4 non-test files**, though both these targets are configurable.

In order to configure different targets for how many tests a project should have, use the following ` + "`mllint`" + ` configuration snippet:
` + "```yaml" + `
testing:
  report: tests-report.xml # JUnit report for rule testing/pass

  # Specify targets for testing/has-tests.
  # Both the minimum required amount of tests as well as the desired ratio of tests to other Python files will be checked.
  targets:
    # Minimum number of tests expected to be in your project. Default: 1
    minimum: 1

    # Define a target ratio of Python test files to other Python files in your project.
    # By default, mllint expects that 20% of all Python files in your project are tests,
    # i.e., 1 test file is implemented for every 4 other Python files.
    ratio:
      tests: 1
      other: 4
` + "```" + `

or equivalent TOML (without the explaining comments):
` + "```toml" + `
[tool.mllint.testing]
report = "tests-report.xml"
targets = { minimum = 1, ratio = { tests = 1, other = 4 }}
` + "```",
	Weight: 1,
}

var RuleTestsPass = api.Rule{
	Name: "Project passes all of its automated tests",
	Slug: "testing/pass",
	Details: `Of course, the point of having automated tests is to ensure that they pass.
While ` + "`mllint`" + ` will **not run** your tests as part of its static analysis, ` + "`mllint`" + ` expects you to run these on your own terms
and provide a the filenames to a JUnit-compatible XML test report and a Cobertura-compatible XML coverage
report in your project's ` + "`mllint`" + ` configuration. Specifically for this rule, the JUnit test report is analysed.

` + howToMakeJUnitXML + `

You can then configure mllint to pick up your test report as follows:

` + "```yaml" + `
testing:
  report: tests-report.xml # JUnit report for rule testing/pass
` + "```" + `

or equivalent TOML:
` + "```toml" + `
[tool.mllint.testing]
report = "tests-report.xml"
` + "```" + `
`,
	Weight: 1,
}

var RuleTestCoverage = api.Rule{
	Name: "Project provides a test coverage report",
	Slug: "testing/coverage",
	Details: `One way of measuring the effectiveness of automated tests, is by measuring how many lines of code are touched while the tests are being executed.
This is called **test coverage**. The idea is that the more lines are being executed by your tests, the more of your code's behaviour is being exercised,
thus yielding a greater probability of bugs surfacing and being detected or prevented.

Note, however, that line test coverage only measures whether a line of code is executed.
This does **not** mean that the result or side-effect of that line's execution is being assessed for correctness.

Additionally, one line may cause two different pathways through your application, e.g., an *if*-statement.
When testing one such path, line test coverage will show that the *if*-statement was covered, 
yet it does not always show that only one of the possible paths through your application has been exercised.
This can especially occur in complex one-line operations. For this use-case, there is also the concept of **branch coverage**,
though ` + "`mllint`" + ` currently does not assess this though.

Furthermore, for testing ML systems, there is also academic discussion as to whether line coverage or branch coverage makes sense,
or whether different forms of coverage are required. While ` + "`mllint`" + ` currently does not check or support any of these novel forms of test coverage for ML,
we are looking for suggestions on what novel forms of ML code coverage should be assessed and how these can be measured.

---

While ` + "`mllint`" + ` will **not run** your tests as part of its static analysis, ` + "`mllint`" + ` expects you to run these on your own terms
and provide a the filenames to a JUnit-compatible XML test report and a Cobertura-compatible XML coverage
report in your project's ` + "`mllint`" + ` configuration. Specifically for this rule, the Cobertura-compatible coverage report is analysed.

` + howToMakeCoverageXML + `

You can then configure mllint to pick up your test report as follows:

` + "```yaml" + `
testing:
  coverage:
    report: coverage.xml
    targets:
      line: 80 # percent line coverage. Default is 80%
` + "```" + `

or equivalent TOML:
` + "```toml" + `
[tool.mllint.testing.coverage]
report = "coverage.xml"
targets = { line = 80.0 } 

# Note: unlike YAML, TOML distinguishes between floats and integers, so be sure to use 80.0 instead of 80
` + "```" + `
`,
	Weight: 1,
}

var RuleTestsFolder = api.Rule{
	Name: "Tests should be placed in the tests folder",
	Slug: "testing/tests-folder",
	Details: "In accordance with `pytest`'s [conventions](https://docs.pytest.org/en/6.2.x/goodpractices.html#conventions-for-python-test-discovery) for Python tests and [recommendations on test layout](https://docs.pytest.org/en/6.2.x/goodpractices.html#tests-outside-application-code), test files are Python files starting with `test_` or ending with `_test.py`" + `
and should be placed in a folder called ` + "`tests`" + ` at the root of your project.

This rule therefore simply checks whether all test files in your projects are indeed in this ` + "`tests`" + ` folder at the root of your project.`,
	Weight: 1,
}
