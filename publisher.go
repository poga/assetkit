package main

import (
	"bufio"
	"bytes"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/termie/go-shutil"
	"github.com/yosssi/gohtml"
)

type Publisher struct {
	themePath string
	copyAsset bool
	Project   *Project
}

func (publisher Publisher) renderProject() (template.HTML, error) {
	tmpl, err := template.ParseFiles("layout.tmpl", "category_page.tmpl")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)
	err = tmpl.Execute(bufWriter, publisher)
	if err != nil {
		return "", err
	}
	err = bufWriter.Flush()
	if err != nil {
		return "", err
	}

	return template.HTML(buf.String()), nil
}

func (publisher Publisher) ProjectLogoPath() string {
	if publisher.copyAsset {
		return publisher.Project.Rel(publisher.Project.LogoPath())
	}

	return publisher.Project.LogoPath()
}

func (publisher Publisher) Publish(outputPath string) error {
	var err error
	project := publisher.Project
	if publisher.copyAsset {
		err = shutil.CopyTree(project.Path, filepath.Join(outputPath, filepath.Base(project.Path)), nil)
		if err != nil {
			return err
		}
	}

	err = shutil.CopyTree(themePath, filepath.Join(outputPath, "themes"), nil)
	if err != nil {
		return err
	}

	project.Meta.Revisions = project.Revisions()
	project.Meta.LastCompiledAt = time.Now()
	err = project.SaveMeta()
	if err != nil {
		return err
	}

	output, err := publisher.renderProject()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(outputPath, "index.html"), []byte(gohtml.Format(string(output))), 0644)
}

func NewPublisher(project *Project, themePath string, copyAsset bool) Publisher {
	return Publisher{Project: project, themePath: themePath, copyAsset: copyAsset}
}
