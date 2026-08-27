# Archiver

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/asano69/archiver)

<img src="frontend/public/favicon.svg" width="100" align="right" />


## Purpose


## Uses
## Getting Started
1. PocketBaseで、userアカウントの作成。
2. トークンを発行(期間 300000000s=10y)
3. SingleFileのRest from API設定
```
URL: http://archiver.app.internal
認証トークン: 発行したもの
アーカイブ・データ・ファイル名: file
アーカイブURLフィールド名: url
```

## Tech Stack
### backend
- Go
- PocketBase v0.39+

### frontend
- Solid.js v1.9
- Kobalte v0.13+
- Tailwind v4
- ProseKit


