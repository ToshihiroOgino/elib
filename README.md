# ELib

シンプルなWebベースのメモアプリケーション

## 機能

- ユーザー認証（ログイン・登録）
- メモの作成・編集・削除
- リアルタイム統計情報表示（文字数・行数・カーソル位置）
- 自動保存機能
- 共有機能
  - 閲覧のみの共有リンク
  - 編集可能な共有リンク
  - リンク再取得・削除

## 共有機能の使い方

### 1. メモを共有する

1. メモエディター画面で「共有(閲覧のみ)」または「共有(編集可)」ボタンをクリック
2. 生成された共有リンクがサイドバーの「共有リンク」セクションに表示される
3. 📋ボタンでリンクをクリップボードにコピー

### 2. 共有リンクを管理する

- **リンクのコピー**: 📋ボタンをクリック
- **リンクの削除**: 🗑️ボタンをクリック

### 3. 共有メモを閲覧・編集する

- **閲覧のみ**: 共有リンクにアクセスすると読み取り専用でメモが表示される
- **編集可能**: 共有リンクにアクセスするとメモの編集が可能

## DB

### 初期化

- Windows: `get-Content .\sqlite\init.sql | sqlite3.exe .\sqlite\db.sqlite3`
- Linux: `sqlite3 ./sqlite/db.sqlite3 < ./sqlite/init.sql`

### Migration

`sqlite\migration` のSQLファイルをタイムスタンプ順に実行する。

`get-Content .\sqlite\migration\xxx.sql | sqlite3.exe .\sqlite\db.sqlite3`

### Domain, Repositoryのコード生成

`sqlite\db.sqlite3` を参照し、スキーマベースのコード生成を行う。

`go run .\infra\sqlite\generate\main.go`

## 実行

`go run ./main.go` または、 `air`
