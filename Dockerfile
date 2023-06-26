FROM golang:1.19

# RUN apk update && apk upgrade && \
#     apk add --no-cache bash git openssh

WORKDIR /app

COPY crw/go.* .
RUN go mod download -x

COPY crw .

RUN go build -o app
RUN chmod +x app
CMD ["./app"]