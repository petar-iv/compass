FROM golang:1.16.3-alpine3.13 as builder

ENV BASE_APP_DIR /go/src/github.com/kyma-incubator/compass/components/external-services-mock
WORKDIR ${BASE_APP_DIR}

#
# Download dependencies
#

COPY go.mod go.sum ${BASE_APP_DIR}/
RUN go mod download -x

#
# Copy files
#

COPY . .

#
# Build app
#
ARG WITH_COVERAGE

RUN if [[ -z "$WITH_COVERAGE" ]]; then go build -v -o main ./cmd/main.go; else CGO_ENABLED="0" go test -c -cover -covermode=set -coverpkg ./...  -tags testrunmain -v -o main ./cmd/; fi
RUN mkdir /app && mv ./main /app/main

FROM alpine:3.13.5
LABEL source = git@github.com:kyma-incubator/compass.git
WORKDIR /app

#
# Copy binary
#

COPY --from=builder /app /app

#
# Run app
#

CMD ["/app/main"]
