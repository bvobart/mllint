import platform
import setuptools

# TODO: based on platform, include the right platform-specific mllint binary as package_data
# TODO: ensure that building this creates platform-specific wheels

with open("../ReadMe.md", "r", encoding="utf-8") as fh:
  long_description = fh.read()

setuptools.setup(
  name="mllint",
  version="0.1.0",
  author="Bart van Oort",
  author_email="bart@vanoort.is",
  description="Linter for Machine Learning projects",
  license='MIT',
  long_description=long_description,
  long_description_content_type="text/markdown",
  url="https://gitlab.com/bvobart/mllint",
  project_urls={
      "Bug Tracker": "https://gitlab.com/bvobart/mllint/-/issues",
  },
  classifiers=[
      "Development Status :: 2 - Pre-Alpha",
      "Environment :: Console",
      "Intended Audience :: Developers",
      "Intended Audience :: Information Technology",
      "Intended Audience :: Science/Research",
      "License :: OSI Approved :: MIT License",
      "Natural Language :: English",
      "Operating System :: MacOS",
      "Operating System :: Microsoft :: Windows",
      "Operating System :: POSIX :: BSD :: FreeBSD",
      "Operating System :: POSIX :: Linux",
      "Programming Language :: Python :: 3",
      "Topic :: Scientific/Engineering :: Artificial Intelligence",
      "Topic :: Scientific/Engineering :: Information Analysis",
      "Topic :: Software Development :: Build Tools",
      "Topic :: Software Development :: Pre-processors",
      "Topic :: Software Development :: Quality Assurance",
      "Topic :: Software Development :: Version Control :: Git",
  ],
  packages=['mllint'],
  package_data={'mllint': ['mllint-exe']},
  python_requires=">=3.6",
  scripts=['mllint/mllint'],
)