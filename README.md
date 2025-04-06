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

### Connect to the services to start playing

Using websocat for the example but you can use any other tool of your preference
```
websocat ws://localhost:3000/publish -H "Authorization: Bearer {JTW_TOKEN}"

websocat ws://localhost:3001/subscribe -H "Authorization: Bearer {JTW_TOKEN}"
```

### To access metrics dashboard click in the following link: [Metrics Dashboard](http://localhost:3002/d/realtime-metrics-dashboard)

### Demo

Build and run the project:
```
docker compose up --build
```

##### Obtain a JWT token and run demo container

Linux:
```
JWT_TOKEN=$(docker compose run --rm auth-cli /app/auth-cli-bin -action create -username demo-user | grep "generated JWT" | awk -F': ' '{print $2}') && docker compose run --rm -e JWT_TOKEN="$JWT_TOKEN" demo /app/demo-bin
```

Windows Powershell:
```
$JWT_TOKEN = (docker compose run --rm auth-cli /app/auth-cli-bin -action create -username demo-user | Select-String "generated JWT" | ForEach-Object { $_.Line -replace ".*: ", "" }); docker compose run --rm -e JWT_TOKEN=$JWT_TOKEN demo /app/demo-bin
```
