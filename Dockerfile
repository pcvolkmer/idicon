FROM golang:1.22-alpine AS build-env

WORKDIR /tmp/build
ADD . /tmp/build
# -ldlflags '-s' to strip binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app -ldflags '-w -s'

###

FROM scratch

LABEL org.opencontainers.image.source="https://github.com/pcvolkmer/idicon"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.description="Simple identicon service"

COPY --from=build-env /tmp/build/app /idicon

USER 8000:8000

EXPOSE 8000

ENTRYPOINT ["/idicon"]