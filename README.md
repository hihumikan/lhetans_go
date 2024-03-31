![image](https://github.com/hihumikan/lhetans_go/assets/26848713/c4ef8ced-64d0-4aed-b479-7e9a682455e9)

# Location-based Home ETA Notification System

帰宅報告(帰宅予想時間と現在位置)をwebhooksに通知する奴のバックエンドサーバー

## Setup

### 1. .envファイルを作成&編集

```bash
cp .env.example .env
vim .env
```

### 2. Dockerコンテナを起動

```bash
docker compose up -d
```

## Usage

```bash
curl -X POST -H "Content-Type: application/json" -d '{"home_location": "35.112133,136.912307", "current_location": "35.688484,139.693222", "travel_mode": "driving", "webhook_url": "https://discord.com/api/webhooks/..."}' http://localhost:3000/notification
```

## License

MIT
