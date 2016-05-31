package main

import (
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
)

type Project struct {
	Path       string
	Categories []*Category
	Meta       Meta
}

func (p *Project) LogoPath() string {
	return filepath.Join(p.Path, "logo.png")
}

func (p *Project) LicensePath() string {
	return filepath.Join(p.Path, "license.md")
}

func (p *Project) Rel(path string) string {
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

func (p *Project) LicenseText() template.HTML {
	renderer := blackfriday.HtmlRenderer(0, "", "")
	md, err := ioutil.ReadFile(p.LicensePath())
	if err != nil {
		return ""
	}

	return template.HTML(blackfriday.Markdown(md, renderer, blackfriday.EXTENSION_HARD_LINE_BREAK))
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

		categories = append(categories, category)
	}

	project.Categories = categories
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

func (p *Project) LoadMeta() (*Meta, error) {
	metaFilePath := filepath.Join(p.Path, ".suisui")
	meta := Meta{}

	if _, err := os.Stat(metaFilePath); err != nil {
		return nil, err
	}

	file, err := os.Open(metaFilePath)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&meta)
	if err != nil {
		return nil, err
	}

	return &meta, nil
}

type Meta struct {
	Revisions      map[string]string
	LastCompiledAt time.Time
}

func (p *Project) Pages() []*Category {
	var result []*Category

	for _, x := range p.Categories {
		Traverse(x, func(c *Category) {
			result = append(result, c)
		})
	}

	return result
}

type Status struct {
	New    []string
	Remove []string
	Change []string
}

func (p *Project) Status() (Status, error) {
	status := Status{New: make([]string, 0), Remove: make([]string, 0), Change: make([]string, 0)}
	savedMeta, err := p.LoadMeta()
	if err != nil {
		return Status{}, err
	}
	currentRevision := p.Revisions()

	for path, hash := range savedMeta.Revisions {
		currentHash, exists := currentRevision[path]
		if !exists {
			status.Remove = append(status.Remove, path)
			continue
		}

		if currentHash != hash {
			status.Change = append(status.Change, path)
		}
	}

	for path, _ := range currentRevision {
		_, exists := savedMeta.Revisions[path]
		if !exists {
			status.New = append(status.New, path)
		}
	}

	return status, nil
}
