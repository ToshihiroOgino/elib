# ELib

## DB 初期化

- Windows: `get-Content .\sqlite\init.sql | sqlite3.exe .\sqlite\db.sqlite3`
- Linux: `sqlite3 ./sqlite/db.sqlite3 < ./sqlite/init.sql`

## Domain と Repository のコード生成

`go run .\infra\sqlite\generate\main.go'

`generated`以下にコードが生成される。これらのファイルは手動で編集しないこと。

## 実行

`go run ./main.go`
