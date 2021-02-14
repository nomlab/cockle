# 説明
このリポジトリには，コンテナないでライブマイグレーションする際に，コンテナに詰め込むものがふくまれている

# ディレクトリ構成
+ coredns dnsの設定ファイルがはいってる
+ docker-compose.yml 複数のコンテナの構成が書いた設定ファイル
+ Dockerfile Dockerのイメージを作るためのファイル
+ image プロセスのライブマイグレーションする時のダンプファイルをコンテナ間で受け渡す時に使うディレクトリ．ボリュームマウントされる．
+ phaul ライブマイグレーションを制御する
+ start.sh frontendとbackendを実行するシェルスクリプト
+ zipkin-ruby-example 対象システム（frontendとbackend)

# 実行
```
docker-compose up -d
```

