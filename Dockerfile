FROM golang:1.18-alpine as builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build /app/cmd/bouncer/bouncer.go

FROM alpine:3.15

RUN apk add --no-cache ca-certificates bash
RUN addgroup -g 1000 gopher && \
    adduser -h /gopher -D -u 1000 -G gopher gopher && \
    chown gopher:gopher /gopher

WORKDIR /gopher
COPY --from=builder /app/bouncer /gopher/bouncer
USER gopher
CMD ["/gopher/bouncer", "-config-file", "/gopher/config.yaml"]
