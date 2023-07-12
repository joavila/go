# Creates TLS Server

## Building
Building embeds private key PEM format plus certificate chain in PEM format this way:

```bash
go build -ldflags "-X main.key_data=$(base64 -w 0 certs/tls-key.pem) -X main.chain_data=$(base64 -w 0 certs/tls-chain.cert.pem)" .
```

## Host
Host requires couple of ports open, one for healthcheck and the other for serving requests, e.g.:

```bash
sudo yum update -y
sudo firewall-cmd --zone=public --add-port=8888/tcp --timeout=1h
sudo firewall-cmd --zone=public --add-port=4443/tcp --timeout=1h
nohup ./hello &
```

## Testing

Testing could be performed this way:

```bash
curl -kv https://localhost:8888/healthcheck
curl -kv https://localhost:4443/hello
curl -kv https://localhost:4443/poisonous
```

Last entry stops server.

# References

https://github.com/denji/golang-tls
https://stackoverflow.com/questions/47857573/passing-certificate-and-key-as-string-to-listenandservetls
https://jamielinux.com/docs/openssl-certificate-authority/index.html

