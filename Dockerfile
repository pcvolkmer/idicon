FROM golang:1.16-alpine AS build-env

ENV USER=appuser
ENV UID=8000

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/null" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /tmp/build
ADD . /tmp/build
# -ldlflags '-s' to strip binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app -ldflags '-w -s'

###

FROM scratch

COPY --from=build-env /etc/passwd /etc/passwd
COPY --from=build-env /etc/group /etc/group
COPY --from=build-env /tmp/build/app /usr/local/bin/idicon

USER appuser:appuser

EXPOSE 8000

ENTRYPOINT ["/usr/local/bin/idicon"]