# URL Shortener Service

シンプルなURL短縮サービスです。長いURLを短いURLに変換し、短いURLにアクセスすると元のURLにリダイレクトします。

## 特徴
- コンテキストサポート（キャンセル処理対応）
- 構造化ログ（INFO, WARN, ERROR）
- スレッドセーフなin-memoryストア
- 詳細なエラーハンドリング

## サービス起動方法

1. プロジェクトのルートディレクトリで以下のコマンドを実行します：

```bash
go run main.go
```

サーバーは8080ポートで起動します。

## 動作確認手順

### 1. URLの短縮

以下のcurlコマンドで長いURLを短縮URLに変換できます：

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/very/long/url/that/needs/shortening"}'
```

成功すると以下のようなレスポンスが返ります：

```json
{
  "short_url": "http://localhost:8080/Ab3Cd9"
}
```

### 2. リダイレクトの確認

生成された短縮URLにアクセスすると、元のURLにリダイレクトされます：

```bash
curl -i http://localhost:8080/Ab3Cd9
```

レスポンスには302ステータスコードと元のURLへのLocationヘッダーが含まれます。

### エラーケース

1. 存在しないコードの場合：
```bash
curl -i http://localhost:8080/nonexistent
```
404 Not Foundが返り、WARNレベルのログが出力されます。

2. 不正なリクエストの場合：
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": ""}'
```
400 Bad Requestが返り、WARNレベルのログが出力されます。

3. コンテキストのキャンセル：
リクエストがキャンセルされた場合（タイムアウトなど）、ERRORレベルのログが出力され、適切なエラーレスポンスが返されます。

### ログレベル
- INFO: 正常な操作（URL登録、リダイレクトなど）
- WARN: 不正なリクエスト、存在しないURLなど
- ERROR: 内部エラー、コンテキストのキャンセルなど

## テストの実行

以下のコマンドでユニットテストを実行できます：

```bash
go test -v
```

データ競合のチェックを含むテストを実行する場合：

```bash
go test -race -v
```

### テストケース
1. URL短縮機能
   - 正常なURL登録と短縮URLの生成
   - 不正なリクエストのハンドリング

2. リダイレクト機能
   - 正常なリダイレクト
   - 存在しないコードの処理

3. コンテキストとキャンセル処理
   - コンテキストのキャンセル時の動作
   - タイムアウト処理

4. ログ出力
   - 各ログレベル（INFO/WARN/ERROR）の出力確認
   - ログメッセージの内容検証
