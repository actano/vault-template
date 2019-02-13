FROM golang:1.11-alpine as build

# install go dep (and git which is needed by dep)
RUN apk update && apk add --no-cache git ca-certificates
RUN wget https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 -O /usr/local/bin/dep
RUN chmod +x /usr/local/bin/dep

# install dependencies
WORKDIR $GOPATH/src/github.com/actano/vault-template
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

# build binary
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /vault-template

FROM scratch

COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /vault-template /vault-template

ENTRYPOINT ["/vault-template"]
