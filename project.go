package main

import (
	"bufio"
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday"
)

type Project struct {
	Path       string
	categories []Category
}

func (p *Project) LogoPath() string {
	return filepath.Join(p.Path, "logo.png")
}

func (p *Project) LicensePath() string {
	return filepath.Join(p.Path, "license.md")
}

func (p *Project) Name() string {
	comps := strings.Split(p.Path, string(os.PathSeparator))
	return NormalizeName(comps[len(comps)-1])

}

func (p *Project) License() (template.HTML, error) {
	renderer := blackfriday.HtmlRenderer(0, "", "")
	md, err := ioutil.ReadFile(p.LicensePath())
	if err != nil {
		return "", err
	}

	return template.HTML(blackfriday.Markdown(md, renderer, blackfriday.EXTENSION_HARD_LINE_BREAK)), nil
}

func NewProject(path string) (*Project, error) {
	path = strings.TrimRight(path, string(os.PathSeparator))
	var categories []Category
	filesInDirectories, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, f := range filesInDirectories {
		if !f.IsDir() {
			continue
		}
		category := NewCategory(filepath.Join(path, f.Name()), nil)
		//spew.Dump(category)

		categories = append(categories, category)
	}

	return &Project{Path: path, categories: categories}, nil
}

func (p *Project) RenderMenu() template.HTML {
	result := ""

	for _, category := range p.categories {
		result += string(category.RenderMenu())
	}

	return template.HTML(result)
}

func (p *Project) RenderContent() template.HTML {
	result := ""
	for _, x := range p.categories {
		Traverse(x, func(c Category) {
			result += string(c.RenderPage())
		})
	}

	return template.HTML(result)
}

func (p *Project) Render() template.HTML {
	tmpl, err := template.ParseFiles(filepath.Join(themePath, "project.tmpl"))
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)
	tmpl.Execute(bufWriter, p)
	bufWriter.Flush()

	return template.HTML(buf.String())
}
