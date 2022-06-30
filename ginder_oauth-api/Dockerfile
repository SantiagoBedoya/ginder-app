FROM golang:1.17.7 as BUILDER

ENV GOOS linux
ENV CGO_ENABLED 0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app

FROM alpine:3.14 as PRODUCTION

RUN apk add --no-cache ca-certificates

COPY --from=BUILDER app .

EXPOSE 3000

CMD ./app