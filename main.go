package main

import (
	"errors"
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

type Editor struct {
	Args
	Options
}

type Args struct {
	Func func()
	Dir string
}

type Options struct {
	N      int
	filter string
}

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
func readImg(path string) (*Info, error) {
	// img情報を保持する構造体を作成
	info := Info{
		Path: path,
	}

	// imgファイルを開く
	f, err := os.Open(info.Path)
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	// おまじない...
	exif.RegisterParsers(mknote.All...)

	// xは *exif.Exif 型
	x, errDecode := exif.Decode(f)
	if errDecode != nil {
		// xがnilの時 exifが存在しないケース
		if x == nil {
			// エラーを返して離脱
			return nil, errors.New("画像情報が存在しません")
		} else {
			fmt.Println(err)
		}
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
	return &info, nil
}

func (e *Editor) reName() {
	iterateFunc(e.Dir, e.filter, func(path string) {
		img, err := readImg(path)

		// Exifが存在しない場合はエラー
		if err != nil {
			fmt.Printf("%v: %v\n", err, path)
			return
		}

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
		newPath := filepath.Join(e.Dir + fName)
		errRename := os.Rename(path, newPath)
		if errRename != nil {
			fmt.Println(errRename)
		}
	})
}

func (e *Editor) del() {
	iterateFunc(e.Dir, e.filter, func(path string) {
		oldName := filepath.Base(path)
		fName := oldName[e.N:]

		newPath := filepath.Join(e.Dir + fName)

		errRename := os.Rename(path, newPath)
		if errRename != nil {
			fmt.Println(errRename)
		}
	})
}

// 対象のディレクトリ内の全ての.jpgに対して関数を実行する
func iterateFunc(dir string, filter string, f func(string)) {
	paths := getPath(dir, filter)

	for _, path := range paths {
		f(path)
	}
}

func (e *Editor) setOptions() {
	// コマンドオプション
	// 削除する文字の長さ デフォルト: 0
	flag.IntVar(&e.N, "n", 0, "set delete length")
	// 絞り込みを行いたい文字列
	filter := flag.String("f", "", "filter file by fileName")
	flag.Parse()
	e.filter = *filter
}

func (e *Editor) setArgs() {
	// コマンド引数: cmd
	// add: 撮影日時を接頭辞としてリネーム
	// del: ファイル名の先頭から指定した文字数分削除する
	cmd := flag.Arg(0)
	if cmd == "" {
		fmt.Println("コマンド引数を設定してください")
		os.Exit(1)
	}else if cmd == "add" {
		e.Func = e.reName
	}else if cmd == "del" {
		e.Func = e.del
	}

	// コマンド引数: dir
	// 対象のディレクトリを指定
	e.Dir = flag.Arg(1)
	if e.Dir == "" {
		fmt.Println("コマンド引数を設定してください")
		os.Exit(1)
	}

	// ディレクトリ名の末尾が"/"でない場合は付与
	length := len(string(e.Dir)) - 1
	if (e.Dir)[length:] != "/" {
		e.Dir += "/"
	}
}

func main() {
	// Editor インスタンスを作成
	e := Editor{}
	// オプションをセットする
	e.setOptions()
	// 引数をセットする
	e.setArgs()
	// 関数を実行する
	e.Func()
}

// ディレクトリにある全ての.jpgのファイルパスを取得する
func getPath(dirname string, filter string) []string {
	var s []string

	// ディレクトリを探索
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		r := regexp.MustCompile(filter)

		threshold := []string{".jpg", ".JPG"}
		//if filepath.Ext(path) == ".jpg" && r.MatchString(filepath.Base(path)){
		if isMatch(filepath.Ext(path), threshold) && r.MatchString(filepath.Base(path)) {
			s = append(s, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return s
}

// 文字列のマッチングの合否を判定
func isMatch(needle string, threshold []string) bool {
	for _, val := range threshold {
		if needle == val {
			return true
		}
	}

	return false
}
