DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS users;

-- ユーザー情報を格納するテーブル
CREATE TABLE users (
    -- ユーザーID (UUIDをテキストとして保存)
    id TEXT PRIMARY KEY NOT NULL,
    -- メールアドレス (ユニーク制約)
    email TEXT UNIQUE NOT NULL,
    -- ハッシュ化されたパスワード
    password_hash BLOB NOT NULL
);

-- 書籍情報を格納するテーブル
CREATE TABLE books (
    -- 書籍ID
    id TEXT PRIMARY KEY NOT NULL,
    -- 外部キー (usersテーブルのidを参照)
    user_id TEXT NOT NULL,
    -- タイトル
    title TEXT NOT NULL,
    -- 著者
    author TEXT NOT NULL,
    /*
     * 提供元
     * 0: PDF
     * 1: BookLive
     * 2: Kindle
     */
    provider INTEGER NOT NULL,
    /*
     * カテゴリ
     * 0: Unknown
     * 1: Novel
     * 2: Comic
     * 3: Other
     */
    category INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- booksテーブルのuser_idカラムにインデックスを作成し、検索を高速化
CREATE INDEX idx_books_user_id ON books(user_id);
