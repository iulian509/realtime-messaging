# realtime-messaging
Real time pubsub messaging with Golang and Websockets

### Build and run the project
```
docker compose up --build
```

### Obtain a JWT token
```
docker compose run auth-cli /app/auth-cli-bin -action create -username demo-user
```

### Access Grafana dashboard: localhost:3002/d/realtime-metrics-dashboard
