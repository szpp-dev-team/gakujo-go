# gakujo-api

## これはなに

某の api を go で書いたやつです。  
ログイン API がようやく動いたので少しずつ実装していきます。

## Installation

```console
$ go get -u github.com/szpp-dev-team/gakujo-api
```

## Documentation

### import

以下のようにインポートしてください。

```go
import "github.com/szpp-dev-team/gakujo-api/gakujo"
```

### Declartion & Login

基本的にクライアントを作成するときは `NewClient` と `Login` を同時に行ってください。特に `Login` をしないと何もできないので注意してください。

```go
client := gakujo.NewClient() // api -> gakujo でもいいかも・・？
if err := client.Login(username, password); err != nil {
    // ログイン失敗
}
```

### Home

ホーム画面から **お知らせ** と **未提出課題一覧** を取得します。  
これで事足りるとは思っていますが、追加で欲しいものがある場合は issue に投げてください。

```go
homeInfo, _ := client.Home()
for _, task := range homeInfo.TaskRows {
    fmt.Println(task)
}
for _, notice := range homeInfo.NoticeRows {
    fmt.Println(notice)
}
```

#### model

- 課題

```go
type TaskRow struct {
    Type     TaskType  // 課題のタイプ
    Deadline time.Time // 締め切り日時
    Name     string    // 課題名
    Index    int       // index
}
```

- お知らせ

```go
type NoticeRow struct {
    Type        NoticeType    // お知らせのタイプ
    SubType     SubNoticeType // お知らせのサブタイプ
    Important   bool          // 重要ラベルの有無
    Date        time.Time     // お知らせが届いた日
    Title       string        // タイトル
    Affiliation string        // 所属
    Index       int           // index
}
```

## test

- login API

```console
$ echo -e 'J_USERNAME=学情のID\nJ_PASSWORD=学情のPSWD' > ./.env
$ go test -timeout 30s -run ^TestLogin$ github.com/szpp-dev-team/gakujo-api/gakujo -v
.
.
--- PASS: TestLogin (8.25s)
PASS
ok      github.com/szpp-dev-team/gakujo-api/gakujo 8.731s
```
