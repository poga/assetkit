package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCategory(t *testing.T) {
	projectPath := "testdata/testproject"
	proj, _ := NewProject(projectPath)

	Convey("Given a relative category path", t, func() {
		categoryPath := "testdata/testproject/CategoryFoo"

		Convey("It should return error", func() {
			c, err := NewCategory(proj, categoryPath, nil)
			So(err, ShouldEqual, ErrRelPath)
			So(c, ShouldBeNil)
		})
	})

	Convey("Given a absolute category path", t, func() {
		categoryPath := abs("testdata/testproject/CategoryBar")
		c, err := NewCategory(proj, categoryPath, nil)
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		Convey("It should automatically convert it into relative path", func() {
			So(c.Path, ShouldEqual, abs("testdata/testproject/CategoryBar"))
		})

		Convey("It should ignore trailing slash", func() {
			c2, err := NewCategory(proj, categoryPath+"/", nil)
			So(err, ShouldBeNil)
			So(c, ShouldNotBeNil)
			So(c.Path, ShouldEqual, c2.Path)
		})

		Convey("It should have correct hierarchy", func() {
			So(c.Project, ShouldPointTo, proj)
			So(c.Parent, ShouldBeNil)
		})

		Convey("Should have name", func() {
			So(c.Name(), ShouldEqual, "CategoryBar")
		})

		Convey("Should have page ID", func() {
			So(c.PageID(), ShouldEqual, "categorybar")
		})

		Convey("Should have page name", func() {
			So(c.PageName(), ShouldEqual, "CategoryBar")
		})

		Convey("It can handle categories without assets", func() {
			So(len(c.Assets), ShouldEqual, 0)
		})

		Convey("It can have nested categories", func() {
			So(len(c.Children), ShouldEqual, 3)

			nestedC1 := c.Children[0]
			nestedC2 := c.Children[1] // name with small word
			nestedEmpty := c.Children[2]

			Convey("Nest categories can have assets", func() {
				So(nestedC1.Name(), ShouldEqual, "CategoryBarC1")
				So(nestedC1.PageID(), ShouldEqual, "categorybar_-_categorybarc1")
				So(nestedC1.PageName(), ShouldEqual, "CategoryBar - CategoryBarC1")

				So(len(nestedC1.Assets), ShouldEqual, 2)
			})

			Convey("categories name with small worlds will be correctly handled", func() {
				So(nestedC2.Name(), ShouldEqual, "CategoryBarC2 the")
				So(nestedC2.PageID(), ShouldEqual, "categorybar_-_categorybarc2_the")
				So(nestedC2.PageName(), ShouldEqual, "CategoryBar - CategoryBarC2 the")

				So(len(nestedC2.Assets), ShouldEqual, 2)
			})

			Convey("Nested category can have no asset", func() {
				So(nestedEmpty.Name(), ShouldEqual, "EmptyCategory")
				So(nestedEmpty.PageID(), ShouldEqual, "categorybar_-_emptycategory")
				So(nestedEmpty.PageName(), ShouldEqual, "CategoryBar - EmptyCategory")

				So(len(nestedEmpty.Assets), ShouldEqual, 0)
			})

			Convey("Traverse category tree should in deep first order", func() {
				var traversed []string

				Traverse(c, func(c *Category) {
					traversed = append(traversed, c.Name())
				})

				So(len(traversed), ShouldEqual, 5)
				So(traversed, ShouldResemble, []string{"CategoryBar", "CategoryBarC1", "CategoryBarC1CC", "CategoryBarC2 the", "EmptyCategory"})
			})
		})

	})

}
