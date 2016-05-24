package main

import (
	"bufio"
	"bytes"
	"errors"
	"html/template"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// a group of asset files
type Asset struct {
	Name          string
	Desc          string
	Downloadables []DownloadablePath
	Images        []Image
	Project       *Project
}

var ErrAssetNameIncorrect = errors.New("Assets with different name can not be grouped")

func (a *Asset) Add(path string) error {
	if AssetName(path) != a.Name {
		return ErrAssetNameIncorrect
	}
	ext := filepath.Ext(path)

	// Renderable images
	if ext == ".png" || ext == ".jpg" {
		img, err := NewImage(a.Project, path)
		if err != nil {
			return err
		}
		a.Images = append(a.Images, *img)
		return nil
	}

	if ext == ".txt" {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		a.Desc = string(content)
		return nil
	}

	// other files
	a.Downloadables = append(a.Downloadables, DownloadablePath(path))
	return nil
}

func (a Asset) RenderPage() template.HTML {
	tmpl, err := template.ParseFiles(filepath.Join(themePath, "asset.tmpl"))
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)
	tmpl.Execute(bufWriter, a)
	bufWriter.Flush()

	return template.HTML(buf.String())
}

func NewAsset(project *Project, path string) *Asset {
	return &Asset{
		Project:       project,
		Images:        make([]Image, 0),
		Downloadables: make([]DownloadablePath, 0),
		Name:          AssetName(path),
	}
}

func AssetName(path string) string {
	ext := filepath.Ext(path)
	fn := strings.TrimSuffix(filepath.Base(path), ext)

	// remove trailing "-01"
	re := regexp.MustCompile(`-\d+$`)
	res := re.ReplaceAllString(fn, "")
	return res
}

type Image struct {
	Path    string
	Width   int
	Height  int
	Project *Project
}

func NewImage(project *Project, path string) (*Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	image, _, err := image.DecodeConfig(file)
	if err != nil {
		return nil, err
	}

	return &Image{Path: path, Width: image.Width, Height: image.Height}, nil
}

func (i Image) Name() string {
	return AssetName(i.Path)
}

type DownloadablePath string

func (dp DownloadablePath) Ext() string {
	return filepath.Ext(string(dp))
}

func (dp DownloadablePath) Name() string {
	return AssetName(string(dp))
}
