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
	f   func(string) string
	Dir string
}

// コマンドオプションを元に情報を格納する構造体
type Options struct {
	delNum   int
	filter   string
	ext      string
	tag      string
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
	//e.f()
	e.reName(e.Dir, e.ext, e.filter, e.f)
}

func (e *Editor) setOptions() {
	// コマンドオプション
	// 削除する文字の長さ デフォルト: 0
	flag.IntVar(&e.delNum, "n", 0, "set delete length")
	// 絞り込みを行いたい文字列
	flag.StringVar(&e.filter, "f", "", "filter file by fileName")
	// リネーム対象とする拡張子
	flag.StringVar(&e.ext, "e", "jpg", "filter file by extension")
	// todo: 逆フィルターを追加
	// todo: イベントのタイトルをタグとしてファイル名に設定できるように機能を追加
	flag.StringVar(&e.tag, "t", "", "add tag")
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
		e.f = e.addDate
	} else if cmd == "del" {
		e.f = e.delName
	} else if cmd == "addTag" {
		e.f = e.addTag
	} else if cmd == "delTag" {
		e.f = e.delTag
	} else if cmd == "addModel" {
		e.f = e.addModel
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

func (e *Editor) addDate(path string) string {
	// decoderのインスタンスを作成
	// todo: インスタンスの作成は無駄？ 直で呼ぶ
	d := decoder.Decoder{}
	img, err := d.ReadImg(path)
	// Exifが存在しない場合はエラー
	if err != nil {
		fmt.Printf("%v: %v\n", err, path)
		os.Exit(3)
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
	return newPath
}

// カメラ機種名を追加する
func (e *Editor) addModel(path string) string {
	// decoderのインスタンスを作成
	// todo: addDataと処理が被っている。インスタンスの作成、imgのデコード
	d := decoder.Decoder{}
	img, err := d.ReadImg(path)
	// Exifが存在しない場合はエラー
	if err != nil {
		fmt.Printf("%v: %v\n", err, path)
		os.Exit(3)
	}

	// Camera Modelを文字列に変換
	//m := img.camModel
	m := d.CamModel(img)
	if m == nil {
		fmt.Println("カメラ機種名が存在しません")
		os.Exit(3)
	}

	// Model名のサニタイズ
	l := `"(.*)"`
	r := regexp.MustCompile(l)
	ms := r.ReplaceAllString(m.String(), "[${1}]")

	// 書き換え
	oldName := filepath.Base(path)
	var fName string
	// 2001-01-01-#Kobe-[GR2]-img02.jpg
	// 日時-#タグ-[機種名] の順とする
	// タグの差し込み位置を決定
	matchDate := `(.{4}-.{2}-.{2}-)`
	matchTag := `(#.*)-`
	rDate := regexp.MustCompile(matchDate)
	rTag := regexp.MustCompile(matchTag)

	// タグが記載済みの場合は、そのあとに機種名を追記する
	if rTag.MatchString(oldName) {
		fName = rTag.ReplaceAllString(oldName, "${1}" + "-" + ms + "-")
	} else if rDate.MatchString(oldName) {
		// 日時が記載済みの場合は、そのあとにタグを追加する
		fName = rDate.ReplaceAllString(oldName, "${1}" + ms + "-")
	} else {
		// タグも日時も未記載の場合は、先頭にタグを追加する
		fName = ms + "-" + oldName
	}
	newPath := e.newPath(fName, path)
	return newPath
}

// tagを追加する
func (e *Editor) addTag(path string) string {
	// img01.jpg ==>> #Kobe-img01.jpg
	// 2001-01-01-img02.jpg ==> 2001-01-01-#Kobe-img02.jpg
	// 2001-01-01-#Kobe-img02.jpg ==>> 2001-01-01-#Kobe#Family-img02.jpg
	oldName := filepath.Base(path)
	// 新しいタグ
	newTag := "#NewTag"
	// タグの差し込み位置を決定
	matchDate := `(.{4}-.{2}-.{2}-)`
	matchTag := `(#.*)-`
	rDate := regexp.MustCompile(matchDate)
	rTag := regexp.MustCompile(matchTag)

	var fName string

	// タグが記載済みの場合は、そのあとにタグを追記する
	if rTag.MatchString(oldName) {
		fName = rTag.ReplaceAllString(oldName, "${1}"+newTag+"-")
	} else if rDate.MatchString(oldName) {
		// 日時が記載済みの場合は、そのあとにタグを追加する
		fName = rDate.ReplaceAllString(oldName, "${1}"+newTag+"-")
	} else {
		// タグも日時も未記載の場合は、先頭にタグを追加する
		fName = newTag + "-" + oldName
	}
	newPath := e.newPath(fName, path)
	return newPath
}

// tagを削除する
func (e *Editor) delTag(path string) string {
	// todo: not implemented yet
	// #Kyoto#Family-img.01.jpg ==>> img.01.jpg
	// #Kobe-img02.jpg ==>> img02.jpg
	// 2001-01-01-#Kyoto#Family-img03.jpg ==>> 2001-01-01-img03.jpg
	// 2001-01-01-#Kobe-img04.jpg ==>> 2001-01-01-img04.jpg
	return path
}

func (e *Editor) delName(path string) string {
	oldName := filepath.Base(path)
	fName := oldName[e.delNum:]

	newPath := e.newPath(fName, path)
	return newPath
}

// ファイルのWalk, iterate, Rename処理をラップ
func (e *Editor) reName(dir string, ext string, filter string, f func(string) string) {
	// 全てのファイルに共通の処理を実行
	paths := decoder.GetPath(dir, ext, filter)

	for _, path := range paths {
		// oldPathを関数に渡す
		newPath := f(path)
		// reNameを実行
		errRename := os.Rename(path, newPath)
		if errRename != nil {
			fmt.Println(errRename)
		}
	}
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
