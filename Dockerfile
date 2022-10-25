FROM golang:1.18-alpine3.16 as builder
COPY . /src/
WORKDIR /src/
RUN go mod download
RUN GOOS=linux go build -ldflags "-X 'main.version=0.0.1'" -o ./.bin/app ./cmd/api/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /src/.bin/app .

CMD [ "./app" ]