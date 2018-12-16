# nagilla

## これなに

* Nagios で発報されるアラート通知をホスト単位で有効化, 無効化するツールです
    * 現時点では, Nagios 4.4.2 での動作を確認しております
* [direnv](https://github.com/direnv/direnv) が必要です

## インストール

環境に応じたバイナリを[リリースページ]()からダウンロードして, パスが通っているディレクトリに保存する.

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
