
# ud_cos2_mqtt

## 使い方

`ud_cos2_mqtt` は、シリアルポートからセンサーデータを取得し、MQTT ブローカーに送信するプログラムです。以下の方法で利用できます。

### プログラムの実行

`ud_cos2_mqtt` プログラムは `systemd` サービスとして実行されますが、手動で実行する場合は以下のコマンドを使用します。

```bash
/usr/local/bin/ud_cos2_mqtt
```

### コマンドラインオプション

`ud_cos2_mqtt` は以下のコマンドラインオプションをサポートしています:

- `-h`: MQTTブローカーのホスト名またはIPアドレスを指定します。デフォルトは `localhost` です。
- `-p`: MQTTブローカーのポートを指定します。デフォルトは `1883` です。
- `-u`: MQTTブローカーのユーザー名を指定します。
- `-P`: MQTTブローカーのパスワードを指定します。

例:

```bash
/usr/local/bin/ud_cos2_mqtt -h example.com -p 1883 -u myuser -P mypassword
```

### MQTT ブローカーの設定

MQTT ブローカーの設定は環境ファイルを通じて行われます。`ud_cos2_mqtt` は、センサーデータを指定された MQTT ブローカーに送信します。

`install.sh` 実行時に生成された `~/.ud_cos2_mqtt.env` には、以下のような内容が含まれています。

```bash
MQTT_USERNAME=your_username
MQTT_PASSWORD=your_password
```

必要に応じて、これらの値を直接編集することもできます。

## インストール

1. **リポジトリのクローン**

   まず、このリポジトリをクローンします。

   ```bash
   git clone <リポジトリのURL>
   cd <リポジトリのディレクトリ>
   ```

2. **インストールスクリプトの実行**

   `install.sh` スクリプトを実行して、プログラムをインストールします。

   ```bash
   ./install.sh
   ```

   スクリプトを実行すると、以下のプロセスが自動的に行われます。

   - `ud_cos2_mqtt` プログラムのビルド。
   - MQTT のユーザー名とパスワードの入力プロンプト。
   - 環境ファイルの作成と適切なパーミッションの設定。
   - `systemd` サービスファイルの作成とサービスの有効化および起動。

3. **サービスの確認**

   インストールが完了したら、`systemctl` を使ってサービスの状態を確認できます。

   ```bash
   sudo systemctl status ud_cos2_mqtt.service
   ```

   また、ログを確認するには以下を使用します。

   ```bash
   journalctl -u ud_cos2_mqtt.service -n 30
   ```

## 注意点

- パスワードは暗号化されずに環境ファイルに保存されるため、専用のものを発行するとよいです。

## ライセンス

このプロジェクトは、[MITライセンス](LICENSE)のもとで公開されています。
