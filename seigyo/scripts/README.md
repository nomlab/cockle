# Disc
コンテナ環境で動くプロセスの分散配置を設定ファイルに基づいて動的に変更して，ベンチマークを取っていく

# Setting
+ scripts/wait_for_running


# Run
+ 対象システムをまず起動．(所属するネットワークを

+ run.shとかあるけど，汎用的なものではありません．
`$ docker run -it --rm -v <scriptがあるディレクトリ>:/root/scripts --net="<対象の属するネットワーク>" --dns="<対象のシステムにたてたDNSのIP>" <このコンテナのイメージ名> /bin/bash -c 'python /root/main.py --composefile /root/scripts/docker-compose.yml --conffile /root/scripts/conf --check /root/scripts/wait_for_running.sh --test /root/scripts/testscript.sh'`



# Note
+ 力技で解決していること，そもそも解決していないことがおおい．
+ よく失敗する．
