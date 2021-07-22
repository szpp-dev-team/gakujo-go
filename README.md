# gakujo-api

## これはなに

某の api を go で書いたやつです。  
ログイン API がようやく動いたので少しずつ実装していきます。

## test

- login API

```console
$ echo -e 'J_USERNAME=学情のID\nJ_PASSWORD=学情のPSWD' > ./.env
$ go test -timeout 30s -run ^TestLogin$ github.com/szpp-dev-team/gakujo-api/api -v
.
.
--- PASS: TestLogin (8.25s)
PASS
ok      github.com/szpp-dev-team/gakujo-api/api 8.731s
```
