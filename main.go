package main

import (
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
	"os"
	"time"
)

// img情報を保持する構造体
type Info struct {
	Path string
	Data Data
}

type Data struct {
	CamModel *tiff.Tag
	DateTime time.Time
}

// 単一の画像ファイルを読み込み その情報を返す
func readImg(path string) *Info{
	// img情報を保持する構造体を作成
	info := Info{
		Path: path,
	}

	// imgファイルを開く
	f, err := os.Open(info.Path)
	if err != nil {
		fmt.Println(err)
	}

	// おまじない...
	exif.RegisterParsers(mknote.All...)

	// xは *exif.Exif 型
	x, errDecode := exif.Decode(f)
	if errDecode != nil {
		fmt.Println(err)
	}

	// Camera Model 取得
	m, errModel := x.Get(exif.Model)
	if errModel != nil {
		fmt.Printf("Failed to read Camera Model: %v\n", errModel)
	}

	// DateTime 取得
	tm, errDateTime := x.DateTime()
	if errDateTime != nil {
		fmt.Printf("Failed to read DateTime: %v\n", errDateTime)
	}

	// 構造体にセット
	info.Data.CamModel = m
	info.Data.DateTime = tm

	fmt.Printf("camModel: %v\n", info.Data.CamModel)
	fmt.Printf("Taken: %v\n", info.Data.DateTime)

	// 正規表現を使用してマッチを確認
	/*
	validator := map[string]string{
		"camModel": ".{1,}",
		"dateTime": ".{4}-.{2}-.{2}",
	}
	validDateTime := regexp.MustCompile(validator["dateTime"])
	validCamModel := regexp.MustCompile(validator["camModel"])
	*/
	// フォーマットに合っているか判定結果を出力
	/*
	fmt.Printf("Match CamModel OK: %v\n", validCamModel.MatchString(info.Data.CamModel))
	fmt.Printf("Match Date OK: %v\n", validDateTime.MatchString(info.Data.DateTime.String()))
	*/

	return &info
}

func main() {
	i := readImg("testdata/img01.jpg")
	fmt.Printf("img: %v\n", i)

	// ディレクトリを探索

}
