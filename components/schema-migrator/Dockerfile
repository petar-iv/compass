FROM alpine:3.13.5

ARG MIGRATE_VER=4.5.0
ARG KUBECTL_VERSION="v1.17.3"

WORKDIR /migrate

RUN apk --no-cache add bash postgresql-client curl jq &&\
        wget https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl -O /usr/local/bin/kubectl &&\
        chmod +x /usr/local/bin/kubectl
RUN wget https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE_VER}/migrate.linux-amd64.tar.gz -O - | tar -xz
RUN mv migrate.linux-amd64 /usr/local/bin/migrate

COPY ./migrations/ ./migrations
COPY ./run.sh ./run.sh

ENTRYPOINT ["./run.sh"]
