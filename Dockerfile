# Build Container
FROM golang:latest as builder
WORKDIR /go/src/github.com/drill-shishamo-alliance/asotech_server
COPY . .
# Set Environment Variable
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
# Build
RUN go build -o app main.go

# Runtime Container
FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/drill-shishamo-alliance/asotech_server/app /app
EXPOSE 3001
ENTRYPOINT ["/app"]