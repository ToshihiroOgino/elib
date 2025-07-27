# ELib アプリケーション技術レポート

## 概要
ELibは、Go言語（Gin フレームワーク）で構築されたWebベースのメモアプリケーションです。ユーザー認証、メモの作成・編集・共有機能を提供し、セキュリティと使いやすさに重点を置いて設計されています。

## HTML5機能の活用

### 1. Clipboard API
- **場所**: `static/js/editor.js`, `static/js/shared-viewer.js`
- **機能**: 共有リンクの自動クリップボードコピー
- **実装詳細**:
  ```javascript
  if (navigator.clipboard && window.isSecureContext) {
    navigator.clipboard.writeText(shareUrl)
      .then(() => { showToast("クリップボードにコピーしました", "success"); })
      .catch((err) => { console.error("Failed to copy: ", err); });
  }
  ```
- **セキュリティ考慮**: `window.isSecureContext`でHTTPS環境での実行を確認

### 2. Fetch API
- **場所**: `static/js/editor.js`, `static/js/shared-editor.js`
- **機能**: 非同期通信によるメモの保存・削除・共有
- **実装詳細**:
  ```javascript
  fetch("/note/save", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ id: noteId, title: title, content: content })
  })
  ```

### 3. DOM操作とイベントリスナー
- **場所**: `static/js/common-editor.js`, `static/js/datetime-utils.js`
- **機能**: 
  - リアルタイム統計情報の更新（文字数、行数、カーソル位置）
  - 自動保存機能
  - UTC時刻のJST変換
- **実装詳細**:
  ```javascript
  textarea.addEventListener("input", function() {
    updateStats();
    if (config.markUnsavedCallback) {
      config.markUnsavedCallback();
    }
  });
  ```

### 4. Document Ready処理
- **場所**: `static/js/datetime-utils.js`
- **機能**: ページ読み込み時の自動初期化
- **実装詳細**:
  ```javascript
  document.addEventListener("DOMContentLoaded", function() {
    convertAllDatesToJST();
  });
  ```

### 5. HTML5セマンティック要素とフォーム
- **場所**: `templates/auth/login.html`, `templates/auth/register.html`
- **機能**: 入力検証、レスポンシブデザイン
- **実装詳細**:
  ```html
  <input type="email" class="form-control" id="email" name="email" required />
  <input type="password" class="form-control" id="password" name="password" required />
  ```

## セキュリティ機能

### 1. Content Security Policy (CSP)
- **場所**: `secure/security.go`
- **機能**: XSS攻撃の防止
- **実装詳細**:
  ```go
  csp := "default-src 'self'; " +
    "script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; " +
    "style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; " +
    "font-src 'self' https://cdn.jsdelivr.net; " +
    "img-src 'self' data: https:; " +
    "connect-src 'self'; " +
    "object-src 'none'; " +
    "base-uri 'self'; " +
    "form-action 'self'"
  ```

### 2. セキュリティヘッダー
- **場所**: `secure/security.go`
- **機能**: 
  - XSS保護: `X-XSS-Protection: 1; mode=block`
  - MIME型スニッフィング防止: `X-Content-Type-Options: nosniff`
  - クリックジャッキング防止: `X-Frame-Options: DENY`
  - リファラーポリシー: `Referrer-Policy: strict-origin-when-cross-origin`

### 3. セキュアクッキー管理
- **場所**: `secure/cookie.go`
- **機能**: 
  - HttpOnly: JavaScriptからのアクセス防止
  - Secure: HTTPS通信時のみ送信
  - SameSite: CSRF攻撃防止
- **実装詳細**:
  ```go
  return CookieConfig{
    MaxAge:   int((24 * 7 * time.Hour).Seconds()), // 7日間
    Path:     "/",
    Secure:   true,                 // HTTPS only
    HttpOnly: true,                 // No JavaScript access
    SameSite: http.SameSiteLaxMode, // SameSite=Lax
  }
  ```

### 4. パスワードハッシュ化
- **場所**: `usecase/user.go`
- **機能**: bcryptによる安全なパスワード保存
- **実装詳細**:
  ```go
  func hashPassword(password string) ([]byte, error) {
    return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  }
  ```

### 5. JWT認証
- **場所**: `secure/jwt.go`, `secure/auth.go`
- **機能**: ステートレスなユーザー認証
- **特徴**: 
  - トークンベース認証
  - セッション管理
  - 自動ログアウト機能

### 6. 入力検証とエスケープ処理
- **場所**: `secure/util.go`, `secure/template.go`
- **機能**: 
  - HTMLエスケープ: XSS防止
  - JSON文字列の安全化
  - 入力文字数制限
- **実装詳細**:
  ```go
  func escapeHTML(input string) template.HTML {
    escaped := html.EscapeString(input)
    return template.HTML(escaped)
  }
  ```

### 7. セッション管理
- **場所**: `secure/session.go`
- **機能**: 
  - セッションデータの暗号化
  - セッション有効期限管理
  - 自動セッション更新
- **特徴**:
  ```go
  type SessionData struct {
    UserID    string    `json:"user_id"`
    Email     string    `json:"email"`
    LoginTime time.Time `json:"login_time"`
    LastSeen  time.Time `json:"last_seen"`
  }
  ```

## フロントエンド技術

### 1. レスポンシブデザイン
- **フレームワーク**: Bootstrap 5.3.0
- **場所**: CDNから読み込み
- **機能**: モバイル対応、グリッドシステム

### 2. CSS3アニメーション
- **場所**: `static/css/editor.css`, `static/css/shared-note.css`
- **機能**: 保存状態の視覚化、ユーザビリティ向上
- **実装詳細**:
  ```css
  @keyframes pulse {
    0% { opacity: 1; }
    50% { opacity: 0.6; }
    100% { opacity: 1; }
  }
  ```

### 3. リアルタイム機能
- **場所**: `static/js/common-editor.js`
- **機能**: 
  - カーソル位置追跡
  - 文字数・行数・段落数のリアルタイム表示
  - 自動保存（5秒間隔）
  - キーボードショートカット（Ctrl+S）

## データベースセキュリティ

### 1. SQLite3の使用
- **場所**: `infra/sqlite/db.go`
- **特徴**: 
  - ローカルファイルベース
  - SQL注入攻撃に対する保護
  - トランザクション管理

### 2. ORMの使用
- **フレームワーク**: GORM
- **機能**: 
  - SQLインジェクション防止
  - 型安全性
  - マイグレーション管理

## アプリケーション構成

### 1. アーキテクチャ
- **パターン**: クリーンアーキテクチャ
- **層構造**: 
  - Controller層: HTTP リクエスト処理
  - Usecase層: ビジネスロジック
  - Repository層: データアクセス
  - Domain層: エンティティ定義

### 2. ミドルウェア
- **セキュリティミドルウェア**: セキュリティヘッダー設定
- **認証ミドルウェア**: ユーザー認証チェック
- **ログミドルウェア**: アクセスログ記録

## 共有機能

### 1. リンクベース共有
- **場所**: `controller/share.go`, `static/js/editor.js`
- **機能**: 
  - 閲覧専用リンク生成
  - 編集可能リンク生成
  - リンク管理（削除機能）

### 2. 権限制御
- **読み取り専用**: コンテンツの表示のみ
- **編集可能**: コンテンツの編集・保存が可能
- **セキュリティ**: 一意のShareIDによるアクセス制御

## まとめ

ELibアプリケーションは、モダンなWeb技術とセキュリティベストプラクティスを組み合わせて構築されています。特に以下の点で優れています：

1. **HTML5 API の効果的活用**: Clipboard API、Fetch API等の最新Web標準を活用
2. **多層防御のセキュリティ**: CSP、セキュアクッキー、入力検証、認証など複数のセキュリティ対策
3. **ユーザビリティ**: リアルタイム統計、自動保存、レスポンシブデザイン
4. **保守性**: クリーンアーキテクチャによる責任の分離
5. **拡張性**: ミドルウェアベースの構成で機能追加が容易

このような技術選択により、セキュアで使いやすく、保守しやすいWebアプリケーションを実現しています。
