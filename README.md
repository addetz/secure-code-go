# secure-code-go
Demo repository for the talk ["The shimmy to the left: why security is coming for engineers"](https://www.gophercon.co.uk/schedule/)

## Environment variables
Certificates & keys can be easily generated for local testing with [`mkcert`](https://github.com/FiloSottile/mkcert). 

```bash
export SERVER_CERT_FILE="localhost.pem"
export SERVER_KEY_FILE="localhost-key.pem"
export SERVER_PORT="1232"
```

## Execute demos
Run the demo servers one by one: 
```bash
go run demo1/server.go
```
