FROM golang:1.16.3-alpine3.13 as builder

ENV BASE_APP_DIR /go/src/github.com/kyma-incubator/compass/components/system-broker
WORKDIR ${BASE_APP_DIR}

COPY go.mod go.sum ${BASE_APP_DIR}/
RUN go mod download -x

COPY . ${BASE_APP_DIR}

RUN CGO_ENABLED="0" go test -c -cover -covermode=set -coverpkg ./...  -tags testrunmain -v -o /app/system-broker ./cmd/

FROM alpine:3.13.5
LABEL source = git@github.com:kyma-incubator/compass.git
WORKDIR /app

RUN apk --no-cache add curl ca-certificates

COPY --from=builder /app /app

CMD ["/app/system-broker"]
