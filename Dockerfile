# Сборка приложения
FROM golang:1.25-alpine AS application
ARG GITHUB_REF=bundle
ADD . /bundle
WORKDIR /bundle

RUN \
    version=${GITHUB_REF} && \
    echo "Building service. Version: ${version}" && \
    go build -ldflags "-X main.build=${version}" -o /srv/app ./main.go


# Финальная сборка образа
FROM scratch

COPY --from=application /srv /srv
COPY --from=application /bundle/tmpl /srv/tmpl

ENV PORT=8080

EXPOSE 8080
WORKDIR /srv
ENTRYPOINT ["/srv/app"]
