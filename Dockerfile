ARG GO_VERSION=1.25

FROM golang:${GO_VERSION}-alpine AS app_dev
RUN go install github.com/mitranim/gow@latest
WORKDIR /srv
EXPOSE 8080
