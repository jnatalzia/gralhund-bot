FROM golang:1.12

WORKDIR /go/src/app
COPY ./main.go .
COPY ./giphy ./giphy
COPY ./resizer ./resizer

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build

CMD ./app -t $AUTH_KEY