// TOOD: revision control
package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/termie/go-shutil"
	"github.com/yosssi/gohtml"
)

var themePath = "themes/summit"

func main() {
	projectPath := strings.TrimSuffix(os.Args[1], "/")
	project, err := NewProject(projectPath)
	if err != nil {
		log.Fatal(err)
	}

	err = shutil.CopyTree(projectPath, filepath.Join(os.Args[2], projectPath), nil)
	if err != nil {
		log.Fatal(err)
	}
	err = shutil.CopyTree(themePath, filepath.Join(os.Args[2], "themes"), nil)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(filepath.Join(os.Args[2], "index.html"), []byte(gohtml.Format(string(project.Render()))), 0644)
}
