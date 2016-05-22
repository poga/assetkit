// TOOD: revision control
package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/termie/go-shutil"
	"github.com/yosssi/gohtml"
)

var themePath = "themes/summit"

func main() {
	project, err := NewProject(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	err = shutil.CopyTree(os.Args[1], filepath.Join(os.Args[2], os.Args[1]), nil)
	if err != nil {
		log.Fatal(err)
	}
	err = shutil.CopyTree(themePath, filepath.Join(os.Args[2], "themes"), nil)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(filepath.Join(os.Args[2], "index.html"), []byte(gohtml.Format(string(project.Render()))), 0644)
}
