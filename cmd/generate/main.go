package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

type Icon struct {
	Name    string
	Type    string
	Content string
}

var iconsetfiletmpl = "package icons\n\n" +
	"var IconSet%s = IconSet{\n" +
	"%s\n" +
	"}\n"

func main() {
	heroicons()
}

func heroicons() {
	// Fetch repository
	cmd := exec.Command("git", "clone", "https://github.com/tailwindlabs/heroicons.git", "/tmp/heroicons")
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	// Find icon paths
	matches, err := filepath.Glob("/tmp/heroicons/src/**/*.svg")
	if err != nil {
		panic(err)
	}
	// Generate icons set
	iconset := []Icon{}
	for _, match := range matches {
		// Path tokens
		tokens := strings.Split(match, "/")
		// Extract path vars
		icontype := tokens[4]
		iconname := strings.ReplaceAll(tokens[5], ".svg", "")
		// Extract content
		iconbts, err := ioutil.ReadFile(match)
		if err != nil {
			panic(err)
		}
		// Save icon to set
		iconset = append(iconset, Icon{
			Name:    iconname,
			Type:    icontype,
			Content: strings.ReplaceAll(string(iconbts), "\n", ""),
		})
	}
	// Generate iconsets content
	iconsetoutline := ""
	iconsetsolid := ""
	for _, icon := range iconset {
		if icon.Type == "outline" {
			iconsetoutline += "\"" + icon.Name + "\": `" + icon.Content + "`,\n"
		} else if icon.Type == "solid" {
			iconsetsolid += "\"" + icon.Name + "\": `" + icon.Content + "`,\n"
		}
	}
	// Generate iconset files
	ioutil.WriteFile(
		"heroicons.outline.iconset.go",
		[]byte(fmt.Sprintf(iconsetfiletmpl, "HeroiconsOutline", iconsetoutline)),
		0644,
	)
	ioutil.WriteFile(
		"heroicons.solid.iconset.go",
		[]byte(fmt.Sprintf(iconsetfiletmpl, "HeroiconsSolid", iconsetsolid)),
		0644,
	)
	// Generate snippetsfunc file
	snippetsfunc := []string{}
	snippetsinline := []string{}
	for _, icon := range iconset {
		snippetsfunc = append(snippetsfunc, fmt.Sprintf(
			`"Icon: Heroicons %s %s": {"scope": "tmpl,html", "prefix": "i:heroicons-%s-%s", "body": ["%s"]}`,
			icon.Name, icon.Type, icon.Name, icon.Type,
			fmt.Sprintf("{{ icon `heroicons-%s` `%s` }}", icon.Type, icon.Name),
		))
		snippetsinline = append(snippetsinline, fmt.Sprintf(
			`"Inline Icon: Heroicons %s %s": {"scope": "tmpl,html", "prefix": "ii:heroicons-%s-%s", "body": ["%s"]}`,
			icon.Name, icon.Type, icon.Name, icon.Type, strings.ReplaceAll(icon.Content, "\"", "\\\""),
		))
	}
	ioutil.WriteFile("heroicons.code-snippets", []byte(fmt.Sprintf(`{%s}`, strings.Join(snippetsfunc, ",\n"))), 0644)
	ioutil.WriteFile("heroicons-inline.code-snippets", []byte(fmt.Sprintf(`{%s}`, strings.Join(snippetsinline, ",\n"))), 0644)
	// Cleanup
	exec.Command("rm", "-rf", "/tmp/heroicons").Run()
}
