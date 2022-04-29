
# Icons

Inline icons plugin, based on external packages (like heroicons).

## Installing

First, install icons package.

```bash
go get github.com/kyoto-framework/icons
```

Then, use FuncMap during template creation or attach FuncMap to your existing one.

```go
...

tmpl := template.Must(template.New("name").Funcs(icons.FuncMap).ParseGlob("templates/*.html"))


// or

func FuncMap() template.FuncMap {
	return render.ComposeFuncMap(
		...
		icons.FuncMap(),
	)
}
```

## Usage

In Go files:

```go
...

icon1 := icons.Icon("heroicons-outline", "archive")
icon2 := icons.Icon("heroicons-outline", "archive", "class1 class2")
icon3 := icons.Icon("heroicons-outline", "archive", "class1", "class2")

```

In templates:

```html
{{ icon `heroicons-outline` `archive` }}
{{ icon `heroicons-outline` `archive` `h-6 w-6` }}
```
