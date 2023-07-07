FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
#ENV GOPROXY https://goproxy.cn,direct

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/main ./main.go


FROM alpine


WORKDIR /app
COPY --from=builder /app/main /app/main
COPY templates /app/templates

CMD ["./main"]
