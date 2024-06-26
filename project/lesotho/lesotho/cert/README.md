Put your certificate and private key here.

```sh
openssl genrsa -out server.key 2048

openssl req -x509 -key server.key -sha256 -days 365 -nodes -out server.crt -subj '/CN=lesotho' -addext 'subjectAltName=IP:127.0.0.1'
```

`local_certificates.zip` contains a self signed certificate and key for
development use. The password is `lesothocert1357`. It probably won't work if
hosted outside of `127.0.0.1` because of Python's `requests` library (without
jumping through additional hoops).
