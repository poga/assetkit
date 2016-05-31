package main

import (
	"bufio"
	"bytes"
	"errors"
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
	Children []*Category
	Assets   []*Asset
	Parent   *Category
	Project  *Project
}

var ErrRelPath = errors.New("Path can't be relative")

func (c Category) RenderMenu() template.HTML {
	tmpl, err := template.ParseFiles("category_menu.tmpl")
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
		return c.Name()
	}
	return c.Parent.PageName() + " - " + c.Name()
}

func (c Category) PageID() string {
	return strings.ToLower(strings.Replace(c.PageName(), " ", "_", -1))
}

func NewCategory(project *Project, path string, parentCategory *Category) (*Category, error) {
	if !filepath.IsAbs(path) {
		return nil, ErrRelPath
	}

	path = strings.TrimSuffix(path, string(os.PathSeparator))
	category := Category{
		Path:     path,
		Children: make([]*Category, 0),
		Assets:   make([]*Asset, 0),
		Parent:   parentCategory,
		Project:  project,
	}
	fileInCategory, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// grouping asset by asset's filename
	assetGroup := make(map[string]*Asset)

	for _, f := range fileInCategory {
		absPath := filepath.Join(path, f.Name())
		if f.IsDir() {
			// recursive category creation
			c, err := NewCategory(project, absPath, &category)
			if err != nil {
				return nil, err
			}
			category.Children = append(category.Children, c)
		} else {
			// Ignore Hidden file
			if strings.HasPrefix(f.Name(), ".") {
				continue
			}
			name := AssetName(absPath)
			if _, exists := assetGroup[name]; !exists {
				assetGroup[name] = NewAsset(project, absPath)
			}

			err := assetGroup[name].Add(absPath)
			if err != nil {
				return nil, err
			}
		}
	}

	for _, asset := range assetGroup {
		category.Assets = append(category.Assets, asset)
	}

	return &category, nil
}

// DFS
func Traverse(c *Category, f func(*Category)) {
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
