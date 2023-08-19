FROM golang:latest

WORKDIR $GOPATH/src/MiniQ

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
ENV MINIQ-AUTH=secret
RUN go build -v -o miniq .

EXPOSE 8080

CMD [ "./miniq"]