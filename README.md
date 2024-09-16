<div align="center">
    <img alt="Litepage logo" height="200" width="200" src="./public/litepage.svg">
</div>
<div align="center">
    <br />
    <p><strong>Litepage</strong> - a minimalist, zero dependency static site generator</p>
</div>

## What is Litepage?

Litepage serves as a template to quickly get up and running for your next static site. The philosophy is there is nothing to install, nothing to maintain, nothing to update. Purely some simple code written in Go to get you started.

It is hyper focused on just delivering a static site with the assets and content you provide. What it lacks in features is made up in simplicity, which should pay off overtime.

Website coming to read more about this.

## Install

Inline with the zero dependency philosophy, there is nothing to install. This repo serves as a template to start your own static site.

Simply clone the repo and get started. Once cloned, `cd` into the repo and start serving the static site with:

```bash
make serve
# go to https://localhost:3001
```

And see the local page at https://localhost:3001

## Local development

With the above command, we start a dev server to host the static web pages. You can work with this and refresh everytime you make a change and re-build the app, however we can take advantage of a live-reloading tool called [Air](https://github.com/air-verse/air) to reload the app for us. The template already includes the config for you to get started.

Air is a command line utility you install on your machine. While technically it is not a dependency of the project, I think the benifits of live-reloading justifies the use, as it is a big boost to the developer experience.

If you don't already have it, you can install Air by running:

```bash
go install github.com/air-verse/air@latest
```

Now you can run the dev server in live-reload mode:

```bash
make dev
```

And go to https://localhost:3000

Note the port 3000, as now we are connecting to the proxy server via Air to automatically refresh the site after saving changes.

## Build the static site

To build the static file assest, run:

```bash
make build
```

This builds and puts the static site in the `/dist` directory. From here it can be used by any tool to staticly host the site, be it Github Pages, Cloudflare, etc.
