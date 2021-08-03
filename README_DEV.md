# 開発者向けドキュメント

gakujo-api 開発者向けのドキュメントです。  
初学者向けにスクレイピングなどの説明もしておきます。

## チュートリアル

ホーム画面から時間割を取得してスクレイピングの基本的な手順を学びましょう。

### 0. 準備

#### (1) firefox のインストール

https://www.mozilla.org/ja/firefox/new/

このチュートリアルでは firefox を使っていることを前提として進めていきます。firefox は開発者ツールが結構良い感じなので使います。

#### (2) gakujo-api の clone とブランチ切り替え

ssh の設定だけしておいてください。  
[szpp の git 勉強会の資料](https://docs.google.com/presentation/d/1VIsmzR-w8zFsqjyv5hX6H_LXooF3SleCmXIYIMwuVbs/edit?usp=sharing) の17ページを参考にしてください。

```console
$ git clone git@github.com:szpp-dev-team/gakujo-api.git
$ cd gakujo-api
$ git switch tutorial
```

### 1. ホーム画面の取得(ブラウザ)

#### (1) 画面を取得するとは

ある画面(html)を取得する時に確認すべきポイントは以下の通りです。

1. URL
2. メソッド(GET, POST など)
3. (POST だったら)form data
4. (GET だったら)パラメータ

これらの情報は「開発者ツール」を使うことで見ることができます。  
早速開発者ツールを使ってホーム画面をどのように取得しているかを見てみましょう。

#### (2) 開発者ツールの使い方

開発者ツールは F12 キーで開くことができます。  
開くと色々と出てきますが、とりあえず以下の機能は留意しておいてください。

- 矢印のボタン

<img width="32" alt="Screen Shot 2021-08-03 at 17 06 32" src="https://user-images.githubusercontent.com/43411965/127980816-5ff11df8-6b39-4741-99b4-1db8ecb1afaa.png">

左上にあるやつです。  
これをクリックしてから画面内の要素にカーソルを持ってくると、その要素に対応する html のタグを表示してくれます。これかなり便利なので是非積極的に使いましょう。

- インスペクター

現在行事しているページの html を表示します。

- ネットワーク

ブラウザが行っている http リクエスト/レスポンスを表示します。  
ページを取得するときに使用している情報(URL, メソッド, フォームデータなど)を確認することができるので、ネットワークの使い方も抑えておきましょう。

#### (3) ホーム画面取得時のhttpリクエスト/レスポンスを覗いてみる

開発者ツールの使い方を知ったところで早速ホーム画面取得時の通信を見てみましょう。  

まずは開発者ツールを開き、ネットワークタブをクリックしてください。  
次に画面左上にある「ゴミ箱アイコン」をクリックし、表示されている通信全てをクリアしてください。  
そしたら、学情の画面に戻り、「Live Campus」の右隣にある「ホームボタン」をクリックしてください。

何やらたくさん通信が表示されたと思いますが、見るべき情報は基本的に一番上だけでokです。
一番上の通信をクリックしてください。そしたら以下の情報について一つ一つみていきましょう。

- 「ヘッダー」 - URL とメソッド他

<img width="300" alt="Screen Shot 2021-08-03 at 17 16 32" src="https://user-images.githubusercontent.com/43411965/127982482-4ff9061d-1b92-48d3-9ced-797ac45c482f.png">

ここでは「URL」、「メソッド」、「ステータス」などを確認することができます。  
今回の場合は URL は「https://gakujo.shizuoka.ac.jp/...」、メソッドは POST、ステータスは「200 OK」です。

- 「要求」 - フォームデータ

<img width="300" alt="Screen Shot 2021-08-03 at 17 22 34" src="https://user-images.githubusercontent.com/43411965/127983079-da2695ac-0537-410a-9c37-39ac8110d465.png">

ここでは「フォームデータ」などを確認することができます。

```go
"org.apache.struts.taglib.html.TOKEN" ... 
"headTitle"
"menuCode"
"nextPath"
"_screenidentifier"
"_screeninfoDisp"
"_scrollTop"
```

### 2. ホーム画面の取得(go)

ここまで、ホーム画面を取得するために必要な情報を調べてきました。  
では、これらの情報を元に早速 go で実装をしてみましょう。

#### (1) ソースファイルを作成する

`gakujo` ディレクトリ内に `timetable.go` ファイルを作成してください。  
次に、以下のコードを記述してください。

```go
package gakujo

import (
    "io"
)

func (c *Client) fetchHomeHtml() (io.ReadCloser, error) {

}
```

`fetchHomeHtml()` は `Client` 構造体のメソッドで、`io.ReadCloser` と `error` を返します。ここで、`Client` は学情の api クライアント構造体で、`io.ReadCloser` は html、`error` はエラーのことを言います。  
これから `fetchHomeHtml()` 内にホーム画面の html を取得するコードを記述していきます。

#### (2) form data を作成する

先ほど確認したフォームデータをまずは作成してみましょう。  
フォームデータは `url.Values` を使うことで作成することができます。先ほど確認した情報を確認しながらコーディングしていきましょう。コードは以下の通りになるはずです。

```go
func (c *Client) fetchHomeHtml() (io.ReadCloser, error) {
    data := url.Values{} // url.Values 構造体の変数を作成
    datas.Set("headTitle", "headTitle の文字列")
	datas.Set("menuCode", "menuCode の文字列")
	datas.Set("nextPath", "nextPath の文字列")
}
```

これでフォームデータの準備はできました。

#### (3) リクエストを飛ばして html を拾う

フォームデータさえ準備できればリクエストは可能となるので、早速飛ばしてみましょう。  
学情内でのリクエストは `c.getPage()` ですることができます。

```go
func (c *Client) fetchHomeHtml() (io.ReadCloser, error) {
    data := url.Values{} // url.Values 構造体の変数を作成
    data.Set("headTitle", "headTitle の文字列")
    data.Set("menuCode", "menuCode の文字列")
    data.Set("nextPath", "nextPath の文字列")

    homeHtmlRc, err := c.getPage(GeneralPurposeUrl, data)
    if err != nil {
        return nil, err // エラーハンドリングはお忘れなく
    }

    return homeHtmlRc, nil
}
```

これでホーム画面の html を取得することができました。

> Topic  
"org.apache.struts.taglib.html.TOKEN" は `c.getPage` でセットしているので、data にわざわざセットしなくても大丈夫です。

#### (4) 本当に取得できたのか？

これでホーム画面を取得できた **はず** ですが、本当にできたかどうかは実際にみてみないとわかりません。

「本当に取得できたのか確認する」ということはプログラマっぽく言い換えれば 「`fetchHomeHtml()` をテストする」ということです。これをテストするためにテスト関数を作成してみましょう。

#### (5) テスト関数を作成する

まずは `gakujo/api_test.go` を開いてください。このファイルは api 関連のコードをテストするためのファイルとなっています。  
次に、そのファイルの一番下に以下のコードを記述してください。

```go
func TestTimetable(t *testing.T) {
    c := NewClient()
    if err := c.Login(username, password); err != nil {
        t.Fatal(err)
    }
    homeHtmlRc, err := c.fetchHomeHtml()
    if err != nil {
        t.Fatal(err)
    }
    b, err := io.ReadAll(homeHtmlRc) // homeHtmlRc は io.ReadCloser なので バイト列に変換
    if err != nil {
        t.Fatal(err)
    }
    fmt.Println(string(b)) // ここで html を表示。バイト列を string にキャスト
}
```

流れとしては「Client 構造体を作成」->「ログイン」->「fetchHomeHtml()」、「html を表示」となっています。これで html を正常に表示できたら勝ちです。

#### (6) テストを実行する

ではテストを実行してみましょう。  
以下のコマンドで実行できます。(ボタンクリックでもできますが説明は省きます)

```console
$ go test -timeout 120s -run ^TestTimetable$ github.com/szpp-dev-team/gakujo-api/gakujo -v -count=1
```

html が表示されましたか？  
されたら成功です。
