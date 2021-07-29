# sns

## Deploy

### Build

`go build -o /bin/mises cmd/main.go`

### Config

```
mkdir upload # mount a storage device

edit .env file
```

### Start

`APP_ENV=production JWT_SECRET="jwt secret" /bin/mises`
