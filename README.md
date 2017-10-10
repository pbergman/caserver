## CA Server

caserver is a small ca (certificate authority) server that can be used for signing, creating and storing certificates.
This was build to easily get/create certificates and no security is implemented and should not be used for production.

See API.md for examples and endpoint of the server.

## Installing

```
make build
sudo make install
```

after that you should edit the `/etc/caserver.cnf` file.

## Chrome

to install the ca in chrome you should get the ca cert first:

```
curl http://127.0.0.1:8080/api/v1/ca  > ca.pem
```

then go to `chrome://settings/certificates` and in the tab Authorities you can import the
download certificate.

## Nginx

create a certificate:

```
curl -X POST -d 'cn=dev&host=*.dev' http://127.0.0.1:8080/api/v1/cert --output /etc/nginx/ssl/dev.pem
```

a simple ssl config could be:

```
server {
    listen              443 ssl;
    server_name		    *.dev;

    ssl_certificate     /etc/nginx/ssl/dev.pem;
    ssl_certificate_key /etc/nginx/ssl/dev.pem;

    ssl                         on;
    ssl_session_cache           builtin:1000  shared:SSL:10m;
    ssl_protocols               TLSv1 TLSv1.1 TLSv1.2;
    ssl_ciphers                 HIGH:!aNULL:!eNULL:!EXPORT:!CAMELLIA:!DES:!MD5:!PSK:!RC4;
    ssl_prefer_server_ciphers   on;

    location / {
      proxy_set_header        Host $host;
      proxy_set_header        X-Real-IP $remote_addr;
      proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header        X-Forwarded-Proto $scheme;
      proxy_pass              http://localhost:80;
      proxy_read_timeout      90;

    }
}
```

## Debug

us the --debug flag to output all debug messages. This will also setup the debug routes for the server see [pprof](https://golang.org/pkg/net/http/pprof/)