# URL Shortener Service

シンプルなURL短縮サービスです。長いURLを短いURLに変換し、短いURLにアクセスすると元のURLにリダイレクトします。

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
404 Not Foundが返ります。

2. 不正なリクエストの場合：
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": ""}'
```
400 Bad Requestが返ります。

## テストの実行

以下のコマンドでユニットテストを実行できます：

```bash
go test -v
```

データ競合のチェックを含むテストを実行する場合：

```bash
go test -race -v
```
