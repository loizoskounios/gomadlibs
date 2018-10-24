FROM golang:1.11.1-stretch as builder
LABEL stage=intermediate
WORKDIR /gomadlibs-src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gomadlibs .

FROM alpine:3.8
RUN apk --no-cache add ca-certificates
WORKDIR /gomadlibs
COPY --from=builder /gomadlibs-src/gomadlibs .
COPY --from=builder /gomadlibs-src/stories ./stories
COPY --from=builder /gomadlibs-src/docker-entrypoint.sh .
ENTRYPOINT ["./docker-entrypoint.sh"]
