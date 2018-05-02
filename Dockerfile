FROM golang:1.9 as builder

WORKDIR /go/src/github.com/uphy/drone-image-copy-plugin
COPY . .
RUN CGO_ENABLED=0 go build -o /drone-image-copy-plugin .

FROM docker:18.04.0-ce-dind

COPY --from=builder /drone-image-copy-plugin /bin/drone-image-copy-plugin
COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /bin/drone-image-copy-plugin /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
