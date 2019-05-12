# 未読記事の管理

- 記事の未読/既読の切り替え
- 優先度の設定
- 記事の追加

## 使用方法
### 管理画面
1. `make local` によりバイナリを作成する
2. `bin/XXX` を実行する
---
1. `make` する
2. `fabfile` を使ってリモートにupload
3. `nginx`, `supervisord` をrestart

### 記事の追加
`script/article.py` によりDBに追加される  
追加対象の記事は`script/article.py`に記載

cronにより自動実行している
