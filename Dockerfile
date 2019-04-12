FROM golang:1.12

ENV IMAGEPATH=/tmp
ENV PROJECTROOT=/go/src/github.com/jnatalzia/gralhund-bot

RUN mkdir -p /tmp
RUN mkdir -p $PROJECTROOT

WORKDIR $PROJECTROOT
COPY ./main.go .
COPY ./giphy ./giphy
COPY ./resizer ./resizer
COPY ./commands ./commands
COPY ./utils ./utils
COPY ./docs ./docs

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build

CMD ./gralhund-bot -t $AUTH_KEY