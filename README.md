# 要約
コンテナ内のプロセスの構成を変更する

# ディレクトリ構成
+ seigyo
  コンテナの実行やプロセス移行などの実行タイミングを制御
+ taisyou
  対象となるシステム

# 実行確認環境
+ Linux 4.9.0
+ Docker 19.03
+ Docker-Compose 1.27.4

# 実行
```
# 開始
$ ./up.sh
# 後片付け
$ ./down.sh
```

# 実行結果確認
```
$ docker logs frontend
```
で，いろいろでるログの中にSUCCESS!って文字があれば多分成功．
```
 $ docker exec -it backend /bin/bash
```
で，psコマンドうって，frontendというプロセスが動くことが確認できる．

# その他
+ 5回に2回くらい失敗する．
