
#build stage
#FROM golang:alpine AS builder
#WORKDIR /go/src/app
#COPY . .
#RUN apk add --no-cache git
#RUN go-wrapper download   # "go get -d -v ./..."
#RUN go-wrapper install    # "go install -v ./..."

#final stage
#FROM alpine:latest
#RUN apk --no-cache add ca-certificates
#COPY --from=builder /go/bin/app /app
#ENTRYPOINT ./app
#LABEL Name=repository Version=0.0.1
#EXPOSE 9622

FROM golang:1.9
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 9622
CMD ["app"]
