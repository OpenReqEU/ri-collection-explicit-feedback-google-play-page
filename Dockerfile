FROM golang:1.9

COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 9622
CMD ["app"]
