# [`mllint`'s GitHub Pages website](https://bvobart.github.io/mllint/)

This directory contains the source code for `mllint`'s GitHub Pages website. It's a static website generated from Markdown files using [Hugo](https://gohugo.io/) and uses the [PaperMod](https://git.io/hugopapermod) theme with some custom tweaks and accent colours by me, Bart van Oort (@bvobart).

## Folder structure

```sh
.
├── assets # contains custom CSS and HighlightJS theme
├── content # contains all of the Markdown files that make up the content of the website. Part of it needs to be generated.
├── layouts # contains overridden layout files from the theme for some layout tweaks.
├── public # when the website is built, this is where it ends up
├── scripts # contains the scripts to generate the category and rule documentation files
├── static # contains static files, e.g. images.
├── themes # contains a submodule with the PaperMod theme
├── config.yml # Hugo's configuration file.
└── ReadMe.md
```

## Building & Developing

There's a GitHub Actions workflow (see `.github/workflows/gh-pages.yml` in this repo) that automatically builds and publishes this website on every push to `main` containing changes to this folder and on every pushed tag. The built website can be viewed on the `gh-pages` branch of this repo.

For developing this website, you'll want to be able to run it locally. To do so, make sure you have Go and [Hugo](https://gohugo.io/) installed. Then, from this folder (`docs/gh-pages`), run the following commands:

```sh
git submodule update --init # to download the theme
./scripts/clean-docs.sh && go run ./scripts/generate-docs.go # to generate the rules & categories documentation
hugo server # runs a development server, see http://localhost:1313/mllint/ when it's running
```

> Note: you can re-generate the rules & categories documentation _while_ the Hugo server is running and Hugo will pick up the changes automatically.
