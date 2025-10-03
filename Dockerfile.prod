FROM golang:1.24 AS build

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o svc ./cmd/server/ 

FROM scratch

WORKDIR /app
COPY --from=build /app/svc .

ENTRYPOINT ["./svc"]