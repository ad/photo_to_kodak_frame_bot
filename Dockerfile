FROM golang:1.22-alpine AS builder

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

ARG BUILD_VERSION

WORKDIR $GOPATH/src/app
COPY go.mod go.mod
COPY go.sum go.sum
COPY vendor vendor
COPY config config
COPY logger logger
COPY app app
COPY sender sender
COPY main.go main.go
RUN CGO_ENABLED=0 go build -mod vendor -ldflags="-w -s -X main.version=${BUILD_VERSION}" -trimpath -o /dist/app

FROM scratch
RUN mkdir /tmp && chmod 1777 /tmp
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /dist /
COPY config.json /config.json


ENTRYPOINT ["/app"]

# Build arguments
ARG BUILD_ARCH
ARG BUILD_DATE
ARG BUILD_REF
ARG BUILD_VERSION

# Labels
LABEL \
    io.hass.name="photo_to_kodak_frame_bot" \
    io.hass.description="photo_to_kodak_frame_bot" \
    io.hass.arch="${BUILD_ARCH}" \
    io.hass.version="${BUILD_VERSION}" \
    io.hass.type="addon" \
    maintainer="ad <github@apatin.ru>" \
    org.label-schema.description="photo_to_kodak_frame_bot" \
    org.label-schema.build-date=${BUILD_DATE} \
    org.label-schema.name="photo_to_kodak_frame_bot" \
    org.label-schema.schema-version="1.0" \
    org.label-schema.usage="https://gitlab.com/ad/photo_to_kodak_frame_bot/-/blob/main/README.md" \
    org.label-schema.vcs-ref=${BUILD_REF} \
    org.label-schema.vcs-url="https://github.com/ad/photo_to_kodak_frame_bot/" \
    org.label-schema.vendor="HomeAssistant add-ons by ad"
