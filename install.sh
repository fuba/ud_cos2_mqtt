#!/bin/bash
echo `pwd`

# 変数設定
SERVICE_NAME="ud_cos2_mqtt"
BINARY_PATH="/usr/local/bin/$SERVICE_NAME"
WORK_DIR=`pwd`
SERVICE_FILE="/etc/systemd/system/$SERVICE_NAME.service"
ENV_FILE="$HOME/.$SERVICE_NAME.env"
CURRENT_USER=`whoami`
PATH=$PATH:/usr/local/go/bin

# 1. ユーザー名とパスワードの入力
read -p "Enter MQTT username: " USER
read -sp "Enter MQTT password: " PASSWORD
echo

# 2. 環境ファイルの作成
echo "Creating environment file..."
sudo tee $ENV_FILE > /dev/null <<EOL
MQTT_USERNAME=$USER
MQTT_PASSWORD=$PASSWORD
EOL

# 環境ファイルのパーミッション設定
sudo chmod 600 $ENV_FILE
sudo chown root:root $ENV_FILE

# 3. Go プログラムのビルド
echo "Building Go program..."
cd $WORK_DIR || exit 1
echo `pwd`
go build -o $SERVICE_NAME

if [ $? -ne 0 ]; then
    echo "Go build failed."
    exit 1
fi

# 4. ビルドされたバイナリを /usr/local/bin にコピー
echo "Copying binary to /usr/local/bin..."
sudo cp $WORK_DIR/$SERVICE_NAME $BINARY_PATH
if [ $? -ne 0 ]; then
    echo "Failed to copy binary to $BINARY_PATH."
    exit 1
fi

# 5. systemd サービスファイルの作成
echo "Creating systemd service file..."
sudo tee $SERVICE_FILE > /dev/null <<EOL
[Unit]
Description=$SERVICE_NAME Service
After=network.target

[Service]
EnvironmentFile=$ENV_FILE
ExecStart=$BINARY_PATH
WorkingDirectory=$WORK_DIR
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
User=$CURRENT_USER

[Install]
WantedBy=multi-user.target
EOL

# 6. サービスの有効化と起動
echo "Reloading systemd daemon..."
sudo systemctl daemon-reload

echo "Enabling $SERVICE_NAME service..."
sudo systemctl enable $SERVICE_NAME.service

echo "Starting $SERVICE_NAME service..."
sudo systemctl start $SERVICE_NAME.service

# 7. サービスの状態確認
sudo systemctl status $SERVICE_NAME.service

echo "Installation complete."
