# Go Basic Web API

This repository contains the very basic web api using golang, it uses http package instead of popular package like Gin.

## Docker

### Dockerfile

This web api is containerized using multi-staged Dockerfile below.
```Dockerfile
# build the golang binary file
FROM golang:alpine3.24 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# run golang server
FROM alpine
RUN apk update && apk add ca-certificates --no-cache
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
ENTRYPOINT ["./server"]
```

The `Dockerfile` is divided into two sections, the first one is to compile the golang app with `golang:alpine3.24` image.

1. We use alpine based golang image and named it `builder`.
2. We change the working directory in the container to `/app`
3. Files in this project are moved into `/app` directory inside the container
4. `CGO_ENABLED=0` is used to disable the C bindings during compilation, so the binary file is not relying to external C libraries. Where `GOOS=linux` tells the system OS that used to compile this go project, in this case is linux because we want to run it on alpine, which is based on linux. The output of this compilation will be named `server`.

The second section runs the compiled binary using a minimal `alpine` image.
1. We use `alpine` as base image to run the compiled binary. This is lightweight image that the size is relatively small
2. `apk update && apk add ca-certificates --no-cache` installs root CA certificates, which are needed if the app makes outbound HTTPS requests (e.g. calling external APIs) — the base Alpine image doesn't include these by default.
3. Change the working directory in the container to `/app`.
4. Copies the `server` file from `builder` to `/app` directory in `alpine` image.
5. `EXPOSE 8080` documents that the app listens on port 8080 inside the container. Note that this only serves as documentation — the port must still be published (e.g. via docker-compose or -p) to be reachable from outside the container
6. `ENTRYPOINT ["./server"]` runs the binary when the container starts.

### Docker Compose

We also use docker compose to integrate every image used in this project, which is image that we built and nginx for reverse proxy.

```yaml
version: '3'

services:
  app:
    build: .
    expose:
      - 8080

  nginx:
    image: nginx:alpine
    ports:
      - 80:80
    volumes:
      - ./nginx.conf:/etc/nginx/conf.d/default.conf
    depends_on:
      - app
```

1. `app` service is built from `Dockerfile` in the current directory, and only exposes port `8080` internally, which is not reachable from outside Docker directly, and it only communicates with other service, in this case `nginx` through same docker network.
2. `nginx` service uses `nginx:alpine` image and maps port `80` from host to port `80` in container
3. `volumes` mounts `nginx.conf` file in this project to `/etc/nginx/conf.d/deafult.conf` in nginx container, so nginx uses the configuration from this project instead of its default config
4. `depends_on` ensures the `app` container starts before `nginx` service, so nginx doesn't fail trying to reach a service that isn't up yet.

## NGINX Configuration

```nginx
server {
    listen 80;

    location / {
        proxy_pass http://app:8080;
        
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

- `proxy_pass http://app:8080` forwards all incoming requests on listened port, in this case it's `80`, to the `app` service on the port `8080`. The hostname `app` is resolved using Docker's built-in DNS, which runs on the n  etwork that created by Docker Compose.
- The `proxy_set_header` directives forward the original client information (host, IP) to the backend, since by default the Go app would otherwise only see requests coming from nginx itself.

  > NGINX forwards `Host`, `X-Real-IP`, and `X-Forwarded-For` headers to the backend, but the current Go application does not yet read or use them.