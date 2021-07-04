# gakujo-api

## これはなに

某の api を go で書いたやつです。**動　き　ま　せ　ん**(2021/07/04現在)

## メモ

https://gakujo.shizuoka.ac.jp/portal/home/home/initialize の一歩前ぐらいで不正な操作呼ばわりされる。どうして

## test

```console
$ echo -e 'J_USERNAME=学情のID\nJ_PASSWORD=学情のPSWD' > ./.env
$ go test -timeout 30s -run ^TestLogin$ github.com/szpp-dev-team/gakujo-api/api -v
.
.
2021/07/04 20:54:55 は？
exit status 1
```
