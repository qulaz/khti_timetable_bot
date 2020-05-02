FROM golang:1.13.10-buster as builder

WORKDIR /build
ADD . /build/
RUN go build bot/main.go

FROM golang:1.13.10-buster as runtime
COPY --from=builder /build/main /app/
COPY --from=builder /build/bot/db/migrations /app/db/migrations/
WORKDIR /app
CMD ["./main"]