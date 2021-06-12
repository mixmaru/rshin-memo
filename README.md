# rshin-memo

ターミナル上で動作するメモアプリです。  

![Animation](https://user-images.githubusercontent.com/2360858/121778993-82a09780-cbd4-11eb-9c3d-3cbd705db189.gif)


# インストール方法
```
$ go get github.com/mixmaru/rshin-memo/cmd/rshinmemo`
```

# 起動方法
```
$ rshinmemo
```


# 操作方法
## 基本

|キー  |動作            |
|------|----------------|
|j     |カーソル上移動  |
|k     |カーソル下移動  |
|enter |決定            |
|esc   |Viewを閉じる    |
|ctrl-c|アプリを終了する|


## DailyListView

作業メモが日毎に一覧されたView

|キー |動作                             |
|-----|---------------------------------|
|o    |カーソル位置の下に作業メモを挿入 |
|O    |カーソル位置に作業メモを挿入     |
|enter|カーソル位置の作業メモをvimで開く|  

## DateSelectView

作業日を指定するためのView  

|キー |動作                                                                                  |
|-----|--------------------------------------------------------------------------------------|
|enter|カーソル位置の日付でメモを挿入する（「手入力する」を選択した場合、yyyy-mm-ddで入力する）|
  

## NoteSelectView

対象の作業メモを選択するView  

|キー |動作                                                                                 |
|-----|-------------------------------------------------------------------------------------|
|enter|カーソル位置の作業メモをvimで開く（「新規追加」を選択した場合、作業メモ名を入力する）|

