// TOOD: revision control
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yosssi/gohtml"
)

func main() {
	project, err := NewProject(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(gohtml.Format(string(project.Render())))
}
