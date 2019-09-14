package main

import (
	"flag"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

/************************

Windowsで cloneした際に exifパッケージが存在しないのでエラーになる。

todo: go.mod を作成して依存関係を解決する

***********************/

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
	// todo: Finderで確認すると作成日時が存在するにも関わらず、0000-00-00として吐き出される。読み込みのフィールドが別にある？
	if errDateTime != nil {
		fmt.Printf("Failed to read DateTime: %v\n", errDateTime)
	}

	// 構造体にセット
	info.Data.CamModel = m
	info.Data.DateTime = tm
	return &info
}

func reName(dir string) {
	iterateFunc(dir, func(path string) {
		img := readImg(path)

		// Camera Modelを文字列に変換
		m := img.Data.CamModel
		ms := ""
		if m != nil {
			l := `"(.*)"`
			r := regexp.MustCompile(l)
			ms = r.ReplaceAllString(m.String(), "$1") + "-"
		}

		// ファイル名のフォーマット "新ファイル名" = "日時" + "モデル名" + "旧ファイル名"
		var fNames = map[string]string{
			"DateTime": img.Data.DateTime.String()[:10] + "-",
			"Model":    ms,
		}
		fName := fNames["DateTime"] + fNames["Model"] + filepath.Base(path)

		// ファイル名を変更
		newPath := filepath.Join(dir + fName)
		errRename := os.Rename(path, newPath)
		if errRename != nil {
			fmt.Println(errRename)
		}
	})
}

func del(dir string, n int) {
	iterateFunc(dir, func(path string) {
		oldName := filepath.Base(path)
		fName := oldName[n:]

		newPath := filepath.Join(dir + fName)

		errRename := os.Rename(path, newPath)
		if errRename != nil {
			fmt.Println(errRename)
		}
	})
}

// 対象のディレクトリ内の全ての.jpgに対して関数を実行する
func iterateFunc(dir string, f func(string)) {
	paths := getPath(dir)

	for _, path := range paths {
		f(path)
	}
}

func main() {
	dir := "testdata"

	// ディレクトリ名の末尾が"/"でない場合は付与
	length := len(string(dir)) - 1
	if dir[length:] != "/" {
		dir += "/"
	}

	// コマンドオプション
	var N int
	// 削除する文字の長さ デフォルト: 0
	flag.IntVar(&N, "n", 0, "set delete length")
	flag.Parse()

	// コマンド引数
	// add: 撮影日時を接頭辞としてリネーム
	// del: ファイル名の先頭から指定した文字数分削除する
	cmd := flag.Arg(0)

	if cmd == "add" {
		reName(dir)
	} else if cmd == "del" {
		if N != 0 {
			del(dir, N)
		} else {
			fmt.Println("オプション n に 0以外の値を入力してください。")
		}
	}
}

// ディレクトリにある全ての.jpgのファイルパスを取得する
func getPath(dirname string) []string {
	var s []string

	// ディレクトリを探索
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		// todo: ファイル名で絞り込み機能を追加実装 オプション引数を使用
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
