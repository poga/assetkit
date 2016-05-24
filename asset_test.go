package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAsset(t *testing.T) {
	Convey("An asset is a group of files with related filename", t, func() {
		projectPath := "testdata/testproject"
		project, _ := NewProject(projectPath)
		assetTXTPath := "testdata/testproject/CategoryBar/CategoryBarC1/asset.txt"
		assetPNGPath := "testdata/testproject/CategoryBar/CategoryBarC1/asset.png"
		assetPNGPath2 := "testdata/testproject/CategoryBar/CategoryBarC1/asset-2.png"
		assetDATPath := "testdata/testproject/CategoryBar/CategoryBarC1/asset.dat"

		Convey("files with related filename will be grouped into one asset", func() {
			assetTXT := NewAsset(project, assetTXTPath)
			So(assetTXT.Name, ShouldEqual, "asset")
			assetPNG := NewAsset(project, assetPNGPath)
			So(assetPNG.Name, ShouldEqual, "asset")
			assetPNG2 := NewAsset(project, assetPNGPath2)
			So(assetPNG2.Name, ShouldEqual, "asset")
			assetDAT := NewAsset(project, assetDATPath)
			So(assetDAT.Name, ShouldEqual, "asset")

			Convey("Assets with different name can not be grouped", func() {
				err := assetTXT.Add("testdata/testproject/CategoryBar/CategoryBarC1/other asset.png")
				So(err, ShouldEqual, ErrAssetNameIncorrect)
			})

			Convey("Assets with the same name can be grouped", func() {
				Convey("PNG will be saved as Images", func() {
					err := assetTXT.Add(assetPNGPath)
					So(err, ShouldBeNil)

					So(len(assetTXT.Downloadables), ShouldEqual, 0)
					So(len(assetTXT.Images), ShouldEqual, 1)
					So(assetTXT.Images[0].Name(), ShouldEqual, "asset")
					So(assetTXT.Images[0].Path, ShouldEqual, assetPNGPath)
					So(assetTXT.Images[0].Width, ShouldEqual, 16)
					So(assetTXT.Images[0].Height, ShouldEqual, 16)
				})
				Convey(".txt will be saved as desc", func() {
					err := assetTXT.Add(assetTXTPath)
					So(err, ShouldBeNil)

					So(len(assetTXT.Downloadables), ShouldEqual, 0)
					So(len(assetTXT.Images), ShouldEqual, 0)
					So(assetTXT.Desc, ShouldEqual, "asset description\n")
				})
				Convey("other extension will be saved as Downloadables", func() {
					err := assetTXT.Add(assetDATPath)
					So(err, ShouldBeNil)

					So(len(assetTXT.Downloadables), ShouldEqual, 1)
					So(len(assetTXT.Images), ShouldEqual, 0)
					So(assetTXT.Downloadables[0], ShouldEqual, DownloadablePath(assetDATPath))
					So(assetTXT.Downloadables[0].Ext(), ShouldEqual, ".dat")
					So(assetTXT.Downloadables[0].Name(), ShouldEqual, "asset")
				})
			})
		})
	})
}
