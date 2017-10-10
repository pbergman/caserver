## Accept header

You can set the accept header to switch the response output.

Most endpoint will support:

|        |                      |
|--------|----------------------|
|json    |application/json      |
|tar     |application/tar       |
|tar.gz  |application/tar+gzip  |
|pem     |application/pkix-cert |
|text    |text/plain            |


## Get Certificate Authority
##### \[GET\]   /api/v1/ca

this will return the CA certificate that can be add to the
browser to verify all signed certificated.

```
curl http://127.0.0.1:8080/api/v1/ca > ca.pem
```

## Sign an Request Certificate
##### \[PUT\]   /api/v1/ca

```
openssl req -new -newkey rsa:2048 -keyout test.key -out test.csr -nodes -subj "/C=NL/O=Exmaple Company/OU=Org/CN=www.example.com"

curl -X PUT -F "csr=@test.csr" http://127.0.0.1:8080/api/v1/cert
```

## Create an Certificate
##### \[POST\] /api/v1/ca

###### Certificate subject post paramters:

| name                  |required              |
|-----------------------|----------------------|
|country                |false                 |
|organization           |false                 |
|organizational_unit    |false                 |
|locality               |false                 |
|province               |false                 |
|street_address         |false                 |
|postalcode             |false                 |
|common_name            |true                  |

###### Extra paramters:


| name                  |description                                           |
|-----------------------|----------------------------------------------------- |
|host                   |the host to bind the certificate to (can be multiple) |
|bits                   |the bit for creating the private key (default to 2048)|


```
curl -X POST -d 'cn=example&host=*.example.com&host=example.com' http://127.0.0.1:8080/api/v1/cert
```

If a certificate exist for the given host a 400 response will be returned
with a link header fot the matching record.

```
> curl -i -X POST -d 'cn=example&host=*.example.com&host=example.com' http://127.0.0.1:8080/api/v1/cert


< HTTP/1.1 400 Bad Request
< Content-Type: text/plain; charset=utf-8
< Link: href="/api/v1/cert/bf7ff32915a37e2b20230def4d1405a09eeada11", rel="record"
< X-Content-Type-Options: nosniff
< Date: Tue, 10 Oct 2017 21:07:15 GMT
< Content-Length: 25
<
< a csr exists for example
```

## Remove an Certificate
##### \[DELETE\] /api/v1/ca/\<id\>

The id needs to be a full hash of 40 characters.

```
curl -i -X DELETE http://127.0.0.1:8080/api/v1/cert/bf7ff32915a37e2b20230def4d1405a09eeada11
```

## Get an Certificate
##### \[GET\] /api/v1/ca/\<id\>

The id can be a short hash (of a minimal of 4 character).

```
curl -H 'Accept: application/tar+gzip' http://127.0.0.1:8080/api/v1/cert/bf7ff32915a37e2b20230def4d1405a09eeada11 --output file.tar.gz
```

## List All Certificates info
##### \[GET\] /api/v1/list

```
curl -i http://127.0.0.1:8080/api/v1/list
```

## List All Certificate Requests info
##### \[GET\] /api/v1/list/csr

```
curl -i http://127.0.0.1:8080/api/v1/list/csr
```

## List All Certificate info (CA excluded)
##### \[GET\] /api/v1/list/cert

```
curl -i http://127.0.0.1:8080/api/v1/list/cert
```

The host query parameter can be used to search for certificates the match that given hostname.

```
curl -i http://127.0.0.1:8080/api/v1/list/cert?host=example.com
```

## List CA Certificate info
##### \[GET\] /api/v1/list/ca

```
curl -i http://127.0.0.1:8080/api/v1/list/ca
```