FROM docker-infra-repo.adelaidegroup.fr/infra/golang:1.17-alpine as build

WORKDIR /app/pgbouncer-updater

COPY ./ ./

RUN export GOPROXY="https://repo.adelaidegroup.fr/artifactory/api/go/golang-remote" &&\ 
    go mod tidy && \
    go mod vendor

RUN cd cmd/ && \
    CGO_ENABLED=0 go build -o pgbouncer-updater

FROM docker-infra-repo.adelaidegroup.fr/infra/distroless/static:nonroot

EXPOSE 8080

COPY --from=build /app/pgbouncer-updater/cmd/pgbouncer-updater /
COPY --from=busybox:1.35.0-uclibc /bin/sh /bin/sh

USER 65532:65532

CMD [ "/bin/sh","-c","/pgbouncer-updater list --config /vault/secrets/config.yaml" ]
