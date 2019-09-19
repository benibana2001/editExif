package main

import (
	"flag"
	"fmt"
	"github.com/benibana2001/editExif/decoder"
	"os"
	"path/filepath"
	"regexp"
)

// ファイルを編集する構造体
type Editor struct {
	Args
	Options
}

// コマンド引数を元に情報を格納する構造体
type Args struct {
	f   func()
	Dir string
}

// コマンドオプションを元に情報を格納する構造体
type Options struct {
	delNum   int
	filter   string
	dirBreak bool
}

func main() {
	// Editor インスタンスを作成
	e := Editor{}
	// オプションをセットする
	e.setOptions()
	// 引数をセットする
	e.setArgs()
	// 関数を実行する
	e.f()
}

func (e *Editor) setOptions() {
	// コマンドオプション
	// 削除する文字の長さ デフォルト: 0
	flag.IntVar(&e.delNum, "n", 0, "set delete length")
	// 絞り込みを行いたい文字列
	flag.StringVar(&e.filter, "f", "", "filter file by fileName")
	// ディレクトリ階層内の全てのファイルを単一のディレクトリ直下に配置したい時
	flag.BoolVar(&e.dirBreak, "b", false, "ignoring directory layer")
	flag.Parse()
}

func (e *Editor) setArgs() {
	// コマンド引数: cmd
	// add: 撮影日時を接頭辞としてリネーム
	// del: ファイル名の先頭から指定した文字数分削除する
	cmd := flag.Arg(0)
	if cmd == "" {
		fmt.Println("コマンド引数を設定してください")
		os.Exit(1)
	} else if cmd == "add" {
		e.f = e.add
	} else if cmd == "del" {
		e.f = e.del
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

func (e *Editor) add() {
	d := decoder.Decoder{}
	// 全てのファイルに共通の処理を実行
	d.IterateFunc(e.Dir, e.filter, func(path string) {
		img, err := d.ReadImg(path)

		// Exifが存在しない場合はエラー
		if err != nil {
			fmt.Printf("%v: %v\n", err, path)
			return
		}

		// Camera Modelを文字列に変換
		//m := img.camModel
		m := d.CamModel(img)
		ms := ""
		if m != nil {
			l := `"(.*)"`
			r := regexp.MustCompile(l)
			ms = r.ReplaceAllString(m.String(), "$1") + "-"
		}

		// ファイル名のフォーマット "新ファイル名" = "日時" + "モデル名" + "旧ファイル名"
		var fNames = map[string]string{
			"dateTime": d.DateTime(img).String()[:10] + "-",
			"Model":    ms,
		}
		fName := fNames["dateTime"] + fNames["Model"] + filepath.Base(path)

		// ファイル名を変更
		// フラグに応じてディレクトリ階層を無視してトップディレクトリ直下に全ファイルを展開
		newPath := e.newPath(fName, path)
		// エラー
		errRename := os.Rename(path, newPath)
		if errRename != nil {
			fmt.Println(errRename)
		}
	})
}

func (e *Editor) del() {
	d := decoder.Decoder{}
	// 全てのファイルに共通の処理を実行
	d.IterateFunc(e.Dir, e.filter, func(path string) {
		oldName := filepath.Base(path)
		fName := oldName[e.delNum:]

		newPath := e.newPath(fName, path)

		errRename := os.Rename(path, newPath)
		if errRename != nil {
			fmt.Println(errRename)
		}
	})
}

func (e *Editor) newPath(name string, path string) string {
	// 階層構造を維持したママ名前を変更
	if e.dirBreak == true {
		// 階層構造を無視 フラグ有りの時
		return filepath.Join(e.Dir + name)
	} else {
		// 階層構造を保持 デフォルト（フラグ無し）
		newDir := filepath.Dir(path) + "/"
		return filepath.Join(newDir + name)
	}
}
