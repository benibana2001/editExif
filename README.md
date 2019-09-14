# editExif
jpgイメージのexif情報をを元に,対象のjpgファイルの名前を編集します。

## 使い方
インストール
Mac
```bash
go build main.go
```
Windows(64bit)
```bash
GOOS=windows GOARCH=amd64 go build 
```

### 日時を接頭辞として追加
- Mac
```bash
./editExif add
```
- Windows(64bit)
```bash
./editExif.exe add
```

### ファイル名の先頭から指定した文字数分削除
- Mac
```bash
./editExif -n=4 del
```
- Windows(64bit)
```bash
./editExif.exe -n=4 del
