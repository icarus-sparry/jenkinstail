FROM golang:stretch as builder
ENV CGO_ENABLED=0
RUN go get github.com/icarus-sparry/jenkinstail

# Small but not zero sized container for the runtime
FROM scratch
COPY --from=builder /go/bin/jenkinstail /jenkinstail
ENTRYPOINT ["/jenkinstail"]
