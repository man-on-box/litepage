<div align="center">
    <img alt="Litepage logo" height="200" width="200" src="./example/public/litepage.svg">
</div>
<div align="center">
    <br />
    <p><strong>Litepage</strong> - a minimalist, zero dependency static site generator</p>
</div>

## What is Litepage?

Litepage is a lightweight library written in Go that simplifies building static sites. It enables you to create simple, efficient websites without unnecessary features, churn, or bloat. Prioritizing simplicity over convenience, it offers a small API surface and maintains zero additional dependencies.

**Features:**

- 🎁 Builds your static site ready to be hosted on GitHub Pages, Cloudflare Pages, etc.
- ⚡ Serves your static site locally during development
- 🧹 No plugins or extra packages to install
- 📖 Common recipes to help with Markdown, Tailwind CSS, live reloading and more
- 🐢 Stable API with minimal changes over time

## Motivation

There are many static site generators out there, but many try to do too much, such as:

- Requiring project config files
- Requiring specific folder structures
- Introducing a new DSL or templating syntax
- Including features you don’t use
- Introducing dependencies in your project

Litepage stays away from this and avoids too much abstraction. You work with three methods and the standard `io.Writer` interface - that’s all! For this reason, it is easy to introduce Litepage to your existing project (or strip it out later if you don’t need it) and can be used without any library specific knowledge required.

If you are looking for more opinionated and feature rich static site generators in Go, you could look at:

- [Hugo](https://gohugo.io/)
- [Zas](https://github.com/darccio/zas)

## Installation

```
go get github.com/man-on-box/litepage
```

## Project structure

The only folders Litepage interacts with by default are:

- `public/*` - to place your static assets (js, css, icons, images, etc.)
- `dist/*` - contains the outputted site when built, ready to be hosted

## Example

See an example project in [./example/main.go](./example/main.go)

**As simple as it gets:**

```go
func main() {
	lp, _ := litepage.New("hello-world.com")

	lp.Page("/index.html", func(w io.Writer) {
		t := template.Must(template.New("helloWorld").Parse("<h1>Hello, World!</h1>"))
		t.Execute(w, nil)
	})

	lp.Build()
}
```

When built, `index.html` will be created in your `/dist` folder containing your parsed template, along with any static assets contained within your `/public` folder.

## Installation

```
go get github.com/man-on-box/litepage
```

## Usage

### Initialization

Create a new Litepage instance:

```go
lp, err := litepage.New("hello-world.com")
```

Optionally you can pass configuration options when creating the instance:

```go
lp, err = litepage.New("hello-world.com",
    litepage.WithDistDir("custom_dist"),
    litepage.WithPublicDir("custom_public"),
    litepage.WithoutSitemap(),
)
```

#### Options

- `WithDistDir` - Specify a custom dist directory to be used, that is created/written to when building the static site. Default value is `dist`.
- `WithPublicDir` - Specify a custom public directory to be used, that is read to retrieve static assets when building or serving the static site. Default value is `public`.
- `WithoutSitemap` - Do not create a sitemap of your site. By default a `sitemap.xml` is created mapping all pages of the static site. Disable this if you do not want this, or if you want to create your own sitemap.

### Creating pages

Create a new page by passing in the relative filename that will be used when building the site, such as `/index.html` or nested pages like `/articles/new-recipes.html`. **Note:** Paths must start with a `/`, include a file extension and be a valid filepath.

Here you also pass a function that receives the standard `io.Writer` interface. Write your templates to this interface to generate your pages with content. This means you can use any html templating library that supports this interface, such as the Go standard [html/template](https://pkg.go.dev/html/template) package or custom packages like [templ](https://templ.guide/).

```go
lp, _ := litepage.New("hello-world.com")

err := lp.Page("/index.html", func (w io.Writer) {
	    t := template.Must(template.New("helloWorld").Parse("<h1>Hello, World!</h1>"))
	    t.Execute(w, nil)
})
```

An error is returned if the file path is not valid, or if you have already created a page with the same file path previously.

### Building your site

Once you have created all your pages, you can build your site. The outputted contents will be placed in the `/dist` directory, including all your assets in the `/public` directory.

```go
lp, _ := litepage.New("hello-world.com")

// ... add all your pages

err := lp.Build()
```

The result in `/dist` directory can then be used with your preferred static site hosting service like GitHub Pages or CloudFlare pages. It should be an easy process to automate your build and host the outputted files during continuous integration.

### Previewing your site

So we know how to build, but what about when developing your site? Ideally instead of writing files to `/dist`, you could see your site hosted locally on your machine for a better developer experience. For this you can call `Serve` which serves your static site locally, instead of outputting it to your dist folder.

```go
lp, _ := litepage.New("hello-world.com")

// ... add all your pages

err := lp.Serve("3000")
```

This will start a web server at http://localhost:3000 to preview your site.

### Build or Serve

Both methods explicitly build or serve your site, however if you want to be able to serve your site locally, while building it during CI, you can take advantage of the `BuildOrServe` method.

```go
lp, _ := litepage.New("hello-world.com")

// ... add all your pages

err := lp.BuildOrServe()
```

By default, this will build your site the same as if you had called `Build`, however you can optionally `Serve` your site with the use of environment variables.

The following environment variables are checked when calling this method:

- `LP_MODE` - set this to `serve` to serve your site
- `LP_PORT` - set this to customise the port to serve your site on (default '3000')

See the [example Makefile](./example/Makefile) on how you could build or serve by specifying environment variables when building your application.

## Contributing

If you like this project, consider giving it a star ⭐.

If you have any feedback please raise it as an issue 🎁.
