FROM golang:1.16-buster as builder
WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download
COPY ./ ./
RUN go build  /app/cmd/server


FROM gcr.io/distroless/base-debian10

COPY --from=builder /app/server /app/server

WORKDIR /app
CMD ["/app/server"]
