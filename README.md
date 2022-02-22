# AKJG - API Key JWT Generator

A Simple service to generate JWT by API Key (with configurable claims)

## Sample Config

```toml
port = 3000

jwt-secret = "abc" # generate a strong secret (https://cloud.google.com/network-connectivity/docs/vpn/how-to/generating-pre-shared-key)
jwt-expiry-mins = 10

[api-key-jwt-map]
    
	# Setup API Keys and associated JWT claims
    [api-key-jwt-map."f30d3c1a-4144-4321-ba6c-bed3cf4a7308"]
    sub = "1"
    name = "test-user"
    admin = false
       [api-key-jwt-map."f30d3c1a-4144-4321-ba6c-bed3cf4a7308"."https://hasura.io/jwt/claims"]
       x-hasura-default-role = "user"
       x-hasura-allowed-roles = ["editor","user", "mod"]
       x-hasura-org-id = "123"

[log]
format = "text" # valid values are text or json
output = "term" # valid values are term or file
#log-directory = "./logs/" # needed if writing logs to a file
level = "debug"
```

## Deploy using Docker

```
# docker run --rm -p 3000:3000 -v ./config.toml:/app/configs/config.toml pvik/akjg:latest
```

## Usage

```bash
$ curl 'http://localhost:3000/akjg/v1/jwt?apikey=f30d3c1a-4144-4321-ba6c-bed3cf4a7308'
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImh0dHBzOi8vaGFzdXJhLmlvL2p3dC9jbGFpbXMiOnsieC1oYXN1cmEtYWxsb3dlZC1yb2xlcyI6WyJlZGl0b3IiLCJ1c2VyIiwibW9kIl0sIngtaGFzdXJhLWRlZmF1bHQtcm9sZSI6InVzZXIiLCJ4LWhhc3VyYS1vcmctaWQiOiIxMjMifSwibmFtZSI6InRlc3QtdXNlciIsInN1YiI6IjEifQ.jSGbSWE7BWngjXsUyOohw_W7Kay3RdQuHK1kEqnwnW0"}
```

The decoded JWT above:

```json
{
  "admin": false,
  "https://hasura.io/jwt/claims": {
    "x-hasura-allowed-roles": [
      "editor",
      "user",
      "mod"
    ],
    "x-hasura-default-role": "user",
    "x-hasura-org-id": "123"
  },
  "name": "test-user",
  "sub": "1"
}
```

#### Unauthorized Access

```bash
$ curl -vvv 'http://localhost:3000/akjg/v1/jwt?apikey=invalid'
*   Trying 127.0.0.1:3000...
* Connected to localhost (127.0.0.1) port 3000 (#0)
> GET /akjg/v1/jwt?apikey=invalid HTTP/1.1
> Host: localhost:3000
> User-Agent: curl/7.81.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 401 Unauthorized
< Date: Tue, 22 Feb 2022 18:22:42 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
```
