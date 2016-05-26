package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/russross/blackfriday"
	"github.com/termie/go-shutil"
	"github.com/yosssi/gohtml"
)

type Project struct {
	Path       string
	categories []*Category
	Meta       Meta
}

func (p *Project) LogoPath() string {
	return filepath.Join(p.Path, "logo.png")
}

func (p *Project) LogoDataPath() string {
	return p.DataPath(p.LogoPath())
}

func (p *Project) LicensePath() string {
	return filepath.Join(p.Path, "license.md")
}

func (p *Project) DataPath(path string) string {
	relPath, err := filepath.Rel(p.Path, path)
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(filepath.Base(p.Path), relPath)
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
	path, err := filepath.Abs(strings.TrimRight(path, string(os.PathSeparator)))
	if err != nil {
		return nil, err
	}
	project := &Project{Path: path}
	var categories []*Category
	filesInDirectories, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, f := range filesInDirectories {
		if !f.IsDir() {
			continue
		}
		category, err := NewCategory(project, filepath.Join(path, f.Name()), nil)
		if err != nil {
			return nil, err
		}
		//spew.Dump(category)

		categories = append(categories, category)
	}

	project.categories = categories
	return project, nil
}

func (p *Project) Revisions() map[string]string {
	result := make(map[string]string)

	filepath.Walk(p.Path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		h := sha1.New()
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		h.Write(bytes)
		result[path] = hex.EncodeToString(h.Sum(nil))

		return nil
	})

	return result
}

func (p *Project) SaveMeta() error {
	file, err := os.Create(filepath.Join(p.Path, ".suisui"))
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	return encoder.Encode(p.Meta)
}

func (p *Project) LoadMeta() error {
	metaFilePath := filepath.Join(p.Path, ".suisui")
	meta := Meta{}

	if _, err := os.Stat(metaFilePath); err != nil {
		return err
	}

	file, err := os.Open(metaFilePath)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)
	return decoder.Decode(&meta)
}

type Meta struct {
	Revisions      map[string]string
	LastCompiledAt time.Time
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
		Traverse(x, func(c *Category) {
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

func (p *Project) CompileTo(outputPath string) error {
	err := shutil.CopyTree(p.Path, filepath.Join(outputPath, filepath.Base(p.Path)), nil)
	if err != nil {
		return err
	}
	err = shutil.CopyTree(themePath, filepath.Join(outputPath, "themes"), nil)
	if err != nil {
		return err
	}

	p.Meta.Revisions = p.Revisions()
	p.Meta.LastCompiledAt = time.Now()
	err = p.SaveMeta()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(outputPath, "index.html"), []byte(gohtml.Format(string(p.Render()))), 0644)
}
