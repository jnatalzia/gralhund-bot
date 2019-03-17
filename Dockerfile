FROM golang:1.12

WORKDIR /go/src/app
COPY ./main.go .
COPY ./types.go .

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build

RUN ls

CMD ./app -t $AUTH_KEY