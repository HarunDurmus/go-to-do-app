FROM golang:1.20

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .
COPY ./.config ./.config
RUN go build -o go-to-do-app

CMD ["./go-to-do-app"]
