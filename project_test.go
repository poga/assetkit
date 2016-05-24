package main

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProject(t *testing.T) {
	// Only pass t into top-level Convey calls
	Convey("Given a project path", t, func() {
		projectPath := "testdata/testproject"
		proj, err := NewProject(projectPath)
		So(err, ShouldBeNil)
		So(proj, ShouldNotBeNil)
		So(proj.Path, ShouldEqual, abs("testdata/testproject"))

		Convey("Can create a project from path with trailing slash", func() {
			proj2, err := NewProject(projectPath + "/")
			So(err, ShouldBeNil)
			So(proj, ShouldNotBeNil)
			So(proj.Path, ShouldEqual, proj2.Path)
		})

		Convey("Can parse path into absolute path", func() {
			So(proj.Name(), ShouldEqual, "Testproject")
			So(proj.Path, ShouldEqual, abs("testdata/testproject"))

			So(proj.LicensePath(), ShouldEqual, abs("testdata/testproject/license.md"))
			So(proj.LogoPath(), ShouldEqual, abs("testdata/testproject/logo.png"))
		})

		Convey("Can return relative path", func() {
			So(proj.GetRelPath(proj.Path), ShouldEqual, ".")

			So(proj.GetRelPath(proj.LicensePath()), ShouldEqual, "license.md")
			So(proj.GetRelPath(proj.LogoPath()), ShouldEqual, "logo.png")
		})
	})
}

func abs(path string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, path)
}
