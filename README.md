# nagilla

## これなに

* Nagios で発報されるアラート通知をホスト単位で有効化, 無効化するツールです
    * 現時点では, Nagios 4.4.2 での動作を確認しております
* [direnv](https://github.com/direnv/direnv) が必要です

## インストール

環境に応じたバイナリを[リリースページ](https://github.com/inokappa/nagilla/releases)からダウンロードして, パスが通っているディレクトリに保存する.

* nagilla_darwin_amd64 ... macOS 向け
* nagilla_linux_amd64 ... Linux 向け
* nagilla_windows_amd64.exe ... Windows 向け

以下, macOS の場合.

```
cp nagilla_darwin_amd64 ~/bin/nagilla
chmod +x ~/bin/nagilla
```

動作確認.

```sh
$ nagilla --version
0.0.1
```

## 事前準備

以下のような内容で .envrc を作成する.

```sh
export NAGIOS_USER=Nagios ユーザー名
export NAGIOS_PASS=Nagios パスワード
```

以下を実行して, 設定内容を反映させる.

```sh
direnv allow
```

## 設定ファイル

以下のような内容で設定ファイル作成する. 設定ファイル名は任意の名前で OK.

```json
{
  "nagios": {
      "url": "http://localhost:8080"
  },
  "targets": {
      "host": "localhost"
  }
}
```

* `nagios.url` に対象の Nagios ホスト名を記載する (http://, https:// を含める)
* `targets.host` に操作したい Nagios のホスト名を記載する (コマンドラインオプション `-hosts` でも指定が可能)

## 操作

### ホストの状態を確認

```sh
nagilla -config=${設定ファイルのパス} -check
```

以下のように Nagios の Host Status が JSON 返却される.

```json
{
  "Host": "localhost",
  "Status": {
    "Active Checks": "ENABLED",
    "Check Latency / Duration": "0.001 / 4.152 seconds",
    "Check Type": "ACTIVE",
    "Current Attempt": "1/10  (HARD state)",
    "Event Handler": "ENABLED",
    "Flap Detection": "ENABLED",
    "Host Status": "UP",
    "In Scheduled Downtime": "NO",
    "Is This Host Flapping": "NO",
    "Last Check Time": "12-16-2018 14:36:01",
    "Last Notification": "N/A (notification 0)",
    "Last State Change": "12-16-2018 10:37:29",
    "Last Update": "12-16-2018 14:37:47  ( 0d  0h  0m  8s ago)",
    "Next Scheduled Active Check": "12-16-2018 14:41:01",
    "Notifications": "DISABLED",
    "Obsessing": "ENABLED",
    "Passive Checks": "ENABLED",
    "Performance Data": "rta=0.111000ms;3000.000000;5000.000000;0.000000 pl=0%;80;100;0",
    "Status Information": "PING OK - Packet loss = 0%, RTA = 0.11 ms"
  }
}
```

### 通知を無効化

```sh
nagilla -config=${設定ファイルのパス} -disable
# -hosts オプションに対象ホストを指定して, 通知を無効にすることが出来る
nagilla -config=${設定ファイルのパス} -disable -hosts=xxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### 通知を有効化

```sh
nagilla -config=${設定ファイルのパス} -enable
# -hosts オプションに対象ホストを指定して, 通知を有効にすることが出来る
nagilla -config=${設定ファイルのパス} -enable -hosts=xxxxxxxxxxxxxxxxxxxxxxxxxxxx
```
