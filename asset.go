package main

import (
	"bufio"
	"bytes"
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
}

func (a *Asset) Add(path string) {
	ext := filepath.Ext(path)

	// Renderable images
	if ext == ".png" || ext == ".jpg" {
		a.Images = append(a.Images, NewImage(path))
		return
	}

	if ext == ".txt" {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		a.Desc = string(content)
		return
	}

	// other files
	a.Downloadables = append(a.Downloadables, DownloadablePath(path))
}

func (a Asset) RenderPage() template.HTML {
	tmpl, err := template.ParseFiles("asset.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)
	tmpl.Execute(bufWriter, a)
	bufWriter.Flush()

	return template.HTML(buf.String())
}

func NewAsset(path string) *Asset {
	return &Asset{
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
	Path   string
	Width  int
	Height int
}

func NewImage(path string) Image {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	image, _, err := image.DecodeConfig(file)
	if err != nil {
		panic(err)
	}
	return Image{Path: path, Width: image.Width, Height: image.Height}
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
