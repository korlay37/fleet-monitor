FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum devices.csv ./

RUN go mod download 

COPY ./cmd/ ./cmd/
COPY ./internal/ ./internal/

RUN CGO_ENABLED=0 GOOS=linux go build  -o /fleet-monitor ./cmd/api/main.go

EXPOSE 6733

CMD ["/fleet-monitor"]