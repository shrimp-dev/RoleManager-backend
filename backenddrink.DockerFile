FROM golang:1.17

ARG version=1.0.0-SNAPSHOT

COPY ./backenddrink /home/backenddrink

EXPOSE 3000

WORKDIR /home/backenddrink

CMD  go mod tidy & echo "start backend" & go run main.go