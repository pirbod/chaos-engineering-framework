FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /platform-chaos-hub ./cmd/server

FROM alpine:3.20
RUN adduser -D -H platform
COPY --from=build /platform-chaos-hub /platform-chaos-hub
USER platform
EXPOSE 8080
ENTRYPOINT ["/platform-chaos-hub"]
