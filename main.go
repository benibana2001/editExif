package main

import (
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
	"os"
	"path/filepath"
	"regexp"
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
func readImg(path string) *Info {
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

	//fmt.Printf("camModel: %v\n", info.Data.CamModel)
	//fmt.Printf("Taken: %v\n", info.Data.DateTime)

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

func reName(dir string) {
	paths := getPath(dir)

	for _, path := range paths {
		img := readImg(path)

		// ファイル名を変更する "新ファイル名" = "日時" + "モデル名" + "旧ファイル名"
		var fNames = map[string]string{
			"DateTime": img.Data.DateTime.String()[:10] + "-",
			"Model": (func() string {
				m := img.Data.CamModel
				// Model値がnilでない場合は文字列化して返す
				if m != nil {
					layout := `"(.*)"`
					r := regexp.MustCompile(layout)
					re := r.ReplaceAllString(m.String(), "$1")
					return re + "-"
				} else {
					return ""
				}
			})(),
		}
		fName := fNames["DateTime"] + fNames["Model"] + filepath.Base(path)

		// ファイル名を変更
		newPath := filepath.Join(dir + fName)
		errRename := os.Rename(path, newPath)
		if errRename != nil {
			fmt.Println(errRename)
		}
	}
}

func main() {
	dir := "testdata"

	// ディレクトリ名の末尾が"/"でない場合は付与
	length := len(string(dir)) - 1
	fmt.Println(dir[length:])
	if dir[length:] != "/" {
		dir += "/"
	}

	reName(dir)
}

// ディレクトリにある全ての.jpgのファイルパスを取得する
func getPath(dirname string) []string {
	var s []string

	// ディレクトリを探索
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".jpg" {
			s = append(s, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return s
}
