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
)

type Category struct {
	Path     string
	Desc     string
	Children []Category
	Assets   []*Asset
	Parent   *Category
}

func (c Category) RenderMenu() template.HTML {
	tmpl, err := template.ParseFiles(filepath.Join(themePath, "category.tmpl"))
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)
	tmpl.Execute(bufWriter, c)
	bufWriter.Flush()

	return template.HTML(buf.String())
}

func (c Category) RenderPage() template.HTML {
	tmpl, err := template.ParseFiles(filepath.Join(themePath, "category_page.tmpl"))
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)
	tmpl.Execute(bufWriter, c)
	bufWriter.Flush()

	return template.HTML(buf.String())
}

func (c Category) Name() string {
	comps := strings.Split(c.Path, string(os.PathSeparator))
	return NormalizeName(comps[len(comps)-1])
}

func (c Category) PageName() string {
	if c.Parent == nil {
		return ""
	}
	return c.Parent.Name() + " - " + c.Name()
}

func (c Category) PageID() string {
	return strings.ToLower(strings.Replace(c.PageName(), " ", "_", -1))
}

func NewCategory(path string, parentCategory *Category) Category {
	path = strings.TrimSuffix(path, string(os.PathSeparator))
	category := Category{
		Path:     path,
		Children: make([]Category, 0),
		Assets:   make([]*Asset, 0),
		Parent:   parentCategory,
	}
	fileInCategory, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// grouping asset by asset's filename
	assetGroup := make(map[string]*Asset)

	for _, f := range fileInCategory {
		absPath := filepath.Join(path, f.Name())
		if f.IsDir() {
			// recursive category creation
			category.Children = append(category.Children, NewCategory(absPath, &category))
		} else {
			// Ignore Hidden file
			if strings.HasPrefix(f.Name(), ".") {
				continue
			}
			name := AssetName(absPath)
			if _, exists := assetGroup[name]; !exists {
				assetGroup[name] = NewAsset(absPath)
			}

			assetGroup[name].Add(absPath)
		}
	}

	for _, asset := range assetGroup {
		category.Assets = append(category.Assets, asset)
	}

	return category
}

// BFS
func Traverse(c Category, f func(Category)) {
	f(c)

	for _, child := range c.Children {
		Traverse(child, f)
	}
}

func NormalizeName(s string) string {
	words := strings.Fields(s)
	smallwords := " a an on the to "

	for index, word := range words {
		if strings.Contains(smallwords, " "+word+" ") {
			words[index] = word
		} else {
			words[index] = strings.Title(word)

		}

	}
	return strings.Join(words, " ")
}
