package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

var iconmapfiletmpl = "package icons\n\n" +
	"type IconSet map[string]string\n\n" +
	"var IconMap = map[string]IconSet{\n" +
	"%s" +
	"}\n"

var iconmaptmpl = ""

func main() {
	heroicons()
	iconifyicons()
	// Generate iconmap file
	ioutil.WriteFile(
		"iconmap.go",
		[]byte(fmt.Sprintf(iconmapfiletmpl, iconmaptmpl)),
		0644,
	)
}

func heroicons() {
	// Fetch repository
	cmd := exec.Command("git", "clone", "https://github.com/tailwindlabs/heroicons.git", "/tmp/heroicons")
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	// Find icon paths
	matches, err := filepath.Glob("/tmp/heroicons/optimized/**/*.svg")
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

	// Update iconmap
	iconmaptmpl += "\t\"heroicons-outline\": IconSetHeroiconsOutline,\n"
	iconmaptmpl += "\t\"heroicons-solid\": IconSetHeroiconsSolid,\n"
	// Cleanup
	exec.Command("rm", "-rf", "/tmp/heroicons").Run()
}

func iconifyicons() {
	// Fetch repository
	cmd := exec.Command("git", "clone", "https://github.com/iconify/icon-sets.git", "/tmp/iconifyicons")
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	// Find icon paths
	matches, err := filepath.Glob("/tmp/iconifyicons/json/*.json")
	if err != nil {
		panic(err)
	}
	// Generate icons set
	for _, match := range matches {
		// Open icon set file
		iconset, err := os.Open(match)
		if err != nil {
			panic(err)
		}
		defer iconset.Close()

		byteValue, _ := ioutil.ReadAll(iconset)

		var result map[string]interface{}
		json.Unmarshal([]byte(byteValue), &result)

		// Extract type
		icontype := result["prefix"].(string)

		// Prepare data
		iconsettmpl := ""
		snippetsfunc := []string{}
		snippetsinline := []string{}

		for iconname, icon := range result["icons"].(map[string]interface{}) {
			// Define icon
			icon := Icon{
				Name:    iconname,
				Type:    icontype,
				Content: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" aria-hidden="true">` + strings.ReplaceAll(icon.(map[string]interface{})["body"].(string), "\n", "") + `</svg>`,
			}
			// Update icon set template
			iconsettmpl += "\"" + icon.Name + "\": `" + icon.Content + "`,\n"

			// Add snippets
			snippetsfunc = append(snippetsfunc, fmt.Sprintf(
				`"Icon: Iconify %s %s": {"scope": "tmpl,html", "prefix": "i:iconify-%s-%s", "body": ["%s"]}`,
				icon.Name, icon.Type, icon.Name, icon.Type,
				fmt.Sprintf("{{ icon `iconify-%s` `%s` }}", icon.Type, icon.Name),
			))
			snippetsinline = append(snippetsinline, fmt.Sprintf(
				`"Inline Icon: Iconify %s %s": {"scope": "tmpl,html", "prefix": "ii:iconify-%s-%s", "body": ["%s"]}`,
				icon.Name, icon.Type, icon.Name, icon.Type, strings.ReplaceAll(icon.Content, "\"", "\\\""),
			))
		}
		// Generate iconset files
		ioutil.WriteFile(
			"iconify."+icontype+".iconset.go",
			[]byte(fmt.Sprintf(iconsetfiletmpl, "Iconify"+strings.ReplaceAll(strings.Title(strings.ReplaceAll(icontype, "-", " ")), " ", ""), iconsettmpl)),
			0644,
		)
		// Generate snippets files
		ioutil.WriteFile("snippets/iconify."+icontype+".code-snippets", []byte(fmt.Sprintf(`{%s}`, strings.Join(snippetsfunc, ",\n"))), 0644)
		ioutil.WriteFile("snippets/iconify-inline."+icontype+".code-snippets", []byte(fmt.Sprintf(`{%s}`, strings.Join(snippetsinline, ",\n"))), 0644)

		// Update iconmap
		iconmaptmpl += fmt.Sprintf("\t\"iconify-%s\": IconSet%s,\n", icontype, "Iconify"+strings.ReplaceAll(strings.Title(strings.ReplaceAll(icontype, "-", " ")), " ", ""))
	}
	// Cleanup
	exec.Command("rm", "-rf", "/tmp/iconifyicons").Run()
}
