package dependencymgmt

import (
	"fmt"

	"github.com/bvobart/mllint/api"
)

// RuleUse is a linting rule to check whether the project is using a proper dependency management solution.
var RuleUse = api.Rule{
	Slug: "use",
	Name: "Project properly keeps track of its dependencies",
	Details: fmt.Sprintf(`Most, if not all, ML projects either implicitly or explicitly depend on external packages
such as %s.

While you can install these manually and individually onto your local environment with %s,
it is very easy to lose track of which exact packages and versions thereof you've installed.
In turn, that makes it very difficult for your colleagues (or even yourself) to replicate the
set of packages that you had installed while developing your project. This could result in your code
simply not working due to missing packages or displaying unexpected bugs because of an updated dependency.

Proper dependency management is thus important for the maintainability and reproducibility of your project,
yet [research on open-source ML projects](https://arxiv.org/abs/2103.04146) has shown that very few ML applications
actually manage their dependencies correctly. Many use basic %s files, often generated using %s, but these have a 
high tendency to include unrelated packages or packages that cannot be resolved from [PyPI](https://pypi.org/) (Pip's standard package index),
are hard to maintain as they have no distinction between run-time dependencies and development-time dependencies, nor direct and indirect dependencies, 
and may hamper the reproducibility of your ML project by underspecifying their exact versions and checksums.
Managing your project's packages with a %s file is similarly flawed and thus also not recommended, 
except if there is a direct need to build your project into a platform-specific Pip package.

The [Python Packaging User Guide](https://packaging.python.org/tutorials/managing-dependencies/#managing-dependencies) 
recommends using either [Poetry](https://python-poetry.org/) or [Pipenv](https://pipenv.pypa.io/en/latest/) as dependency managers.
The recommendation is to use Pipenv if your project is an application and to use Poetry if it is a library or otherwise needs to be built into a Python package.

If you're seeing this in a report, it means your project is currently not using a dependency manager,
or one that is not recommended. 

Learn more about Poetry and Pipenv using the links below, pick the one that most suits you, your project and your team, 
then start managing your dependencies with it.
- [Poetry](https://python-poetry.org/)
- [Pipenv](https://pipenv.pypa.io/en/latest/)
`,
		"`numpy`, `scikit-learn`, `pandas`, `matplotlib`, `tensorflow`, `pytorch`, etc.",
		"`pip install`",
		"`requirements.txt`",
		"`pip freeze`",
		"`setup.py`"),
	Weight: 1,
}

// Details to be added when RuleUse detects a requirements.txt
var DetailsNoRequirementsTxt = fmt.Sprintf(`Your project seems to be managing its dependencies using a %s file.
Such %s files have a high tendency to include unrelated packages or packages that cannot be resolved from [PyPI](https://pypi.org/) (Pip's standard package index),
are hard to maintain as they have no distinction between run-time dependencies and development-time dependencies, nor direct and indirect dependencies,
and may hamper the reproducibility of your ML project by underspecifying their exact versions and checksums.

We therefore recommend switching to Poetry or Pipenv and keeping track of all your dependencies there.`,
	"`requirements.txt`", "`requirements.txt`")

// Details to be added when RuleUse detects a setup.py
var DetailsNoSetupPy = fmt.Sprintf(`Your project seems to be managing its dependencies using a %s file.
Perhaps you are also using your setup.py to be able to bundle your project into a Python package.
However, %s files have similar flaws to using %s files to manage your dependencies:
it can be hard to maintain and may hamper the reproducibility of your ML project by underspecifying 
your dependencies' exact versions and checksums.

We therefore recommend switching to Poetry or Pipenv and keeping track of all your dependencies there.
Specifically, we recommend Poetry as it also supports building Python packages, as opposed to Pipenv which does not.`,
	"`setup.py`", "`setup.py`", "`requirements.txt`")

// RuleSingle is a linting rule to check whether the project is only using a single dependency manager instead of multiple.
var RuleSingle = api.Rule{
	Slug: "single",
	Name: "Project should only use one dependency manager",
	Details: fmt.Sprintf(`In most cases, using multiple different dependency managers only creates confusion in your team regarding 
which manager to install, which to use for installing the project's dependencies, and in what order.
It can also be confusing for your team to figure out where a new dependency should be added,
or where an existing dependency should be updated (just in one dependency manager (but which one?), or in both?).

We therefore recommend using only one dependency manager, preferably either Poetry or Pipenv.
Please see the description of rule %s for more information.`, "`dependency-management/use`"),
	Weight: 1,
}

// Details to be appended when RuleSingle detects certain combinations of dependency managers
var DetailsRequirementsTxtSetupPy = fmt.Sprintf("Consider using Poetry to replace both your %s and %s", "`requirements.txt`", "`setup.py`")
var DetailsRequirementsTxtPipenv = fmt.Sprintf("Since you are using Pipenv, the %s file in your project is redundant. Migrate any dependencies left in there to Pipenv and remove it.", "`requirements.txt`")
var DetailsRequirementsTxtPoetry = fmt.Sprintf("Since you are using Poetry, the %s file in your project is redundant. Migrate any dependencies left in there to Pipenv and remove it.", "`requirements.txt`")
var DetailsPipenvSetupPy = fmt.Sprintf("Consider using Poetry instead of Pipenv. Poetry is very similar to Pipenv, but also supports building and publishing Python packages, which is what I presume you're using %s for now.", "`setup.py`")
var DetailsPoetrySetupPy = fmt.Sprintf("The %s in your project is redundant and should be removed, as you can also use Poetry to build your project into a Python package using %s, see the [Poetry Docs](https://python-poetry.org/docs/libraries/#packaging) to learn more.", "`setup.py`", "`poetry build`")
