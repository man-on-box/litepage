<div align="center">
    <img alt="Litepage logo" height="200" width="200" src="./example/public/litepage.svg">
</div>
<div align="center">
    <br />
    <p><strong>Litepage</strong> - a minimalist, zero dependency static site generator</p>
</div>

## What is Litepage?

Litepage is a lightweight library to build your static sites, written in Go. You bring the content, and Litepage can get you up and running in seconds with a site ready to deploy with any deployment services that support hosting static websites, such as GitHub Pages and Cloudflare Pages.

**As simple as it gets:**

```go
func main() {
	lp, _ := litepage.New("hello-world.com")

	lp.Page("/index.html", func(w io.Writer) {
		t := template.Must(template.New("helloWorld").Parse("<h1>Hello, World!</h1>"))
		t.Execute(w, nil)
	})

	lp.BuildOrServe()
}
```

Litepage is hyper focused on just delivering a static site with the assets and content you provide. With the barebones approach, you can add complexity and dependencies to suit your needs.

Website coming to read more about this.

## Features

- Automatically copy your assets from `./public` to your `./dist`
- Build static site to your `./dist` folder
- Create a basic 'out of the box' sitemap of your static site
- Provide a dev server to be able to support hot reloading during development
- Comes with 'recipes' on how to handle markdown files, hot reloading, creating pages pragmatically, tailwind and more

## Installation

```
go get github.com/man-on-box/litepage
```

## Usage

WIP
