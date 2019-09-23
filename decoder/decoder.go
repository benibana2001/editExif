package decoder

import (
	"errors"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// Decodeを行う主体
type Decoder struct {}

// img情報を保持する構造体
type Info struct {
	path string
	Data
}

// img情報を保持する構造体
type Data struct {
	camModel *tiff.Tag
	dateTime time.Time
}

// Getter
func (d *Decoder) CamModel(info *Info) *tiff.Tag{
	return info.camModel
}

// Getter
func (d *Decoder) DateTime(info *Info) time.Time{
	return info.dateTime
}


// decoder
// 単一の画像ファイルを読み込み その情報を返す
func (d *Decoder) ReadImg(path string) (*Info, error) {
	// img情報を保持する構造体を作成
	info := Info{
		path: path,
	}

	// imgファイルを開く
	f, err := os.Open(info.path)
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

	// dateTime 取得
	tm, errDateTime := x.DateTime()

	if errDateTime != nil {
		fmt.Printf("Failed to read dateTime: %v\n", errDateTime)
	}

	// 構造体にセット
	info.camModel = m
	info.dateTime = tm
	return &info, nil
}

// todo: 逆フィルターを追加
// ディレクトリにある全ての.jpgのファイルパスを取得する
func GetPath(dirname string, ext string, filter string) []string {
	var s []string

	// ディレクトリを探索
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		r := regexp.MustCompile(filter)

		n, extErr := sanitizeExt(ext)
		if extErr != nil {
			// 未対応の拡張子が設定された場合
			fmt.Println(extErr)
			os.Exit(3)
		}

		if isMatch(filepath.Ext(path), n) && r.MatchString(filepath.Base(path)) {
			s = append(s, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return s
}

// 拡張子を返す
func sanitizeExt(e string) ([]string, error) {
	if e == "jpg" || e == ".jpg" {
		return []string{".jpg", ".JPG"}, nil
	}
	if e == "png" || e == ".png" {
		return []string{".png", ".PNG"}, nil
	}
	return nil, errors.New("指定された拡張子に対応していません")
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
