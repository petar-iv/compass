FROM golang:1.18.0-alpine3.15 as builder

ARG GH_ACCESS_TOKEN

RUN apk update && apk add --no-cache git curl

RUN go env -w GOPRIVATE=github.tools.sap
RUN git config --global url."https://${GH_ACCESS_TOKEN}@github.tools.sap".insteadOf "https://github.tools.sap"

ENV BASE_APP_DIR /go/src/github.com/kyma-incubator/compass/components/atom-tenants-controller
WORKDIR ${BASE_APP_DIR}

COPY go.mod go.sum ${BASE_APP_DIR}/
RUN go mod download -x

COPY . .

RUN go build -o aggregator ./cmd/main.go

RUN mkdir /app && mv ./aggregator /app/aggregator

FROM alpine:3.15.3
WORKDIR /app

RUN apk --no-cache add curl ca-certificates

COPY --from=builder /app /app

CMD /app/aggregator
