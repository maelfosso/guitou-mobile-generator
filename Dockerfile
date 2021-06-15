FROM golang:1.16-alpine3.12

ENV GOPROXY="https://proxy.golang.org"
ENV GO111MODULE="on"

ENV PORT=8000
EXPOSE 8000


WORKDIR /go/src/guitou.cm/mobile/generator

RUN apk add --no-cache git 
COPY . .

RUN go build -v -o /go/bin/generator

CMD ["/go/bin/generator"]
