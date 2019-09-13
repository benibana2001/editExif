package main

import (
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
	"os"
	"time"
)

type Info struct {
	Path string
	Data Data
}

type Data struct {
	CamModel *tiff.Tag
	DateTime time.Time
}

func main() {
	info := Info{
		Path: "testdata/img02.jpg",
	}

	f, err := os.Open(info.Path)
	if err != nil {
		fmt.Println(err)
	}

	exif.RegisterParsers(mknote.All...)

	// xは *exif.Exif 型
	x, errDecode := exif.Decode(f)
	if errDecode != nil {
		fmt.Println(err)
	}
	//fmt.Printf("type: %T,\nvalue: %v\n", x, x)

	// Camera Model 取得
	camModel, _ := x.Get(exif.Model)

	// DateTime 取得
	tm, _ := x.DateTime()

	// 構造体にセット
	info.Data.CamModel = camModel
	info.Data.DateTime = tm

	fmt.Printf("camModel: %v\n", info.Data.CamModel)
	fmt.Printf("Taken: %v\n", info.Data.DateTime)

	// 正規表現を使用してマッチを確認
	/*
	validator := map[string]string{
		"camModel": "",
		"dateTime": ".{4}-.{2}-.{2}",
	}
	validDateTime := regexp.MustCompile(validator["dateTime"])
	fmt.Printf("Match Date OK: %v\n", validDateTime.MatchString(info.Data.DateTime))
	*/
}
