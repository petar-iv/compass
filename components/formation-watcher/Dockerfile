FROM golang:1.14.4-alpine3.12 as builder

ENV BASE_APP_DIR /go/src/github.com/kyma-incubator/compass/components/formation-watcher
WORKDIR ${BASE_APP_DIR}

COPY . .

RUN go build -v -o /app/formation-watcher ./cmd/main.go

FROM alpine:edge
LABEL source = git@github.com:kyma-incubator/compass.git
WORKDIR /app

RUN apk --no-cache add bash curl ca-certificates git jq

COPY --from=builder /app /app
COPY --from=builder /go/src/github.com/kyma-incubator/compass/components/formation-watcher/resources /app/

CMD ["/app/formation-watcher"]
