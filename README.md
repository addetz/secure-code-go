# secure-code-go
Demo repository for the talk ["The shimmy to the left: why security is coming for engineers"](https://www.gophercon.co.uk/schedule/)

## Environment variables
Certificates & keys can be easily generated for local testing with [`mkcert`](https://github.com/FiloSottile/mkcert). 

```bash
export SERVER_CERT_FILE="localhost.pem"
export SERVER_KEY_FILE="localhost-key.pem"
export SERVER_PORT="1232"
export SIGNING_KEY="SUPER-DUPER-SECRET"
```

## Execute demos
Run the demo servers one by one. Each demo builds upon the previous one.

### Demo 1: Server with HTTPS
```bash
go run demo1/server.go
```

### Demo 2: Server with JWT
```bash
go run demo2/server.go
```

### Demo 3: Server with access control checks
```bash
go run demo3/server.go
```


