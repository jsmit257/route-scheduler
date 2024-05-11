FROM golang:bookworm as build
ENV gopath=/go
COPY . /go/src/github.com/jsmit257/rs
WORKDIR /go/src/github.com/jsmit257/rs
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o /rs -a -installsuffix cgo -cover ./cmd/...

FROM python:bookworm as deploy
ENV GOCOVERDIR=/tmp
COPY ./bin /bin
COPY ./data /data
COPY --from=build /rs /rs
CMD python3 /bin/evaluateShared.py --cmd /rs --problemDir /data
