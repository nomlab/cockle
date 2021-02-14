# Disc
構成変更の実行タイミング制御したりベンチマークとったりする

# Run
+ 対象システムをまず起動してから
```
docker-compose up -d
```
# how it work (simply)
+ 構成の設定ファイル読む(scripts/conf)
  + 1行づつよんで，カンマで3個に区切れる場合はこれに従った構成に変更．そうでない場合はそのままの構成．
+ システムの起動をまつ (scripts/wait_for_runnning.sh実行)
+ 構成変更の指示（対象コンテナのinitにgRPCで要求)
+ ベンチマーク取得 (scripts/testscript.shの実行)
+ 構成変更，ベンチマーク取得の繰り返し

# Note
+ 色々できていないところおおい
