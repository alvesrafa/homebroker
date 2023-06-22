FROM golang:1.19-alpine as bigBuild

WORKDIR /app

RUN go mod init challenge01

COPY *.go ./

RUN go build main.go

RUN rm -rf go.mod main.go

FROM scratch

WORKDIR /app

COPY --from=bigBuild /app/main .

CMD ["./main"]