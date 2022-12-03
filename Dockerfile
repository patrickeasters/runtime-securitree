FROM golang:1.19-alpine as builder

WORKDIR /app

COPY . ./
RUN go mod download

RUN go build -o securitree cmd/securitree.go

FROM alpine
COPY --from=builder /app/securitree /bin/securitree

CMD [ "/bin/securitree" ]