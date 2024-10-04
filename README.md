<div align="center">
    <img alt="Litepage logo" height="200" width="200" src="./example/public/litepage.svg">
</div>
<div align="center">
    <br />
    <p><strong>Litepage</strong> - a minimalist, zero dependency static site generator</p>
</div>

## What is Litepage?

Litepage serves as a template to quickly get up and running for your next static site. The philosophy is there is nothing to install, nothing to maintain, nothing to update. Purely some code written in Go to get you started.

Litepage is hyper focused on just delivering a static site with the assets and content you provide. With the barebones approach, you can add complexity and dependencies to suit your needs.

Website coming to read more about this.

## Install

Inline with the zero dependency philosophy, there is nothing to install. This repo serves as a template to start your own static site.

## Start locally

To build and serve the site, run:

```bash
make serve
```

And see the local page at https://localhost:3000

## Local development

For improved DX, you can use a tool like [Air](https://github.com/air-verse/air) to recompile the app on save and refresh the page.

If you don't already have it installed locally, you can install Air by running:

```bash
go install github.com/air-verse/air@latest
```

Now you can run the dev server in live-reload mode:

```bash
make dev
```

## Build the static site

To build site ready to host somewhere, run:

```bash
make build
```

This builds and puts the static site in the `/dist` directory. From here it can be used by any tool to staticly host the site, be it Github Pages, Cloudflare, S3, anywhere that can host a website.
