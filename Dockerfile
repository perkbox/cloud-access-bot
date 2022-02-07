FROM golang:1.17 as gobuild

WORKDIR /
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a

FROM alpine

COPY --from=gobuild /cloud-access-bot /cloud-access-bot
RUN apk add -U --no-cache ca-certificates

CMD ["/cloud-access-bot"]

