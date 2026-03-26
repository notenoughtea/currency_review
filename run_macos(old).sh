#!/bin/bash
CONFIG_PATH="$(pwd)/config.yaml"

# подключаем зависимости
go mod tidy
go mod download

# генерируем gRPC код
protoc -I proto \
  --go_out=./pkg --go_opt=paths=source_relative \
  --go-grpc_out=./pkg --go-grpc_opt=paths=source_relative \
  proto/currency.proto

#запускаем сервисы с задержкой в 1 сек
osascript <<EOF
tell application "Terminal"
    set wd to do shell script "pwd"

    do script "export CONFIG_PATH=$CONFIG_PATH; cd " & quoted form of wd & " && go run ./currency/cmd/currency/main.go"
    delay 1

    do script "export CONFIG_PATH=$CONFIG_PATH; cd " & quoted form of wd & " && go run ./currency/cmd/cron/main.go"
    delay 1

    do script "export CONFIG_PATH=$CONFIG_PATH; cd " & quoted form of wd & " && go run ./currency/cmd/migrator/main.go"
    delay 1

    do script "export CONFIG_PATH=$CONFIG_PATH; cd " & quoted form of wd & " && go run ./gateway/cmd/gateway/main.go"
end tell
EOF
