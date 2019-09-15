# editExif
jpgイメージのexif情報をを元に,対象のjpgファイルの名前を編集します。

## 使い方
### インストール

- Mac
```bash
go build main.go
```
もしくは
```bash
go get github.com/benibana2001/editExif
```

- Windows(64bit)
```bash
GOOS=windows GOARCH=amd64 go build 
```
もしくは
```bash
go get github.com/benibana2001/editExif
```

### 日時を接頭辞として追加
- Mac
```bash
./editExif -d=testdata add
```
- Windows(64bit)
```bash
./editExif.exe -d=testdata add
```

### ファイル名の先頭から指定した文字数分削除
- Mac
```bash
./editExif -d=testdata -n=4 del
```
- Windows(64bit)
```bash
./editExif.exe -d=testdata -n=4 del
```

## オプション 
- d 
    - 対象のディレクトリを指定
- f
    - 対象とするファイルを指定した文字列のマッチングで絞り込み
- n
    - 削除する文字数を指定

