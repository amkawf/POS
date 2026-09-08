# Stage 1 go binary build
FROM golang:alpine AS builder

WORKDIR /app

# copy dependency
COPY go.mod go.sum ./
RUN go mod download

# copy the rest of app
COPY . .

# compile go app to static linux binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/api-server ./cmd/api

# stage 2 miniml runtime environment
FROM alpine:3.20

WORKDIR /app

# copy completed binary from builder stage
COPY --from=builder /app/api-server .

#inform that container listen to port 8080
EXPOSE 8080

# command to execute when container start
CMD [ "./api-server" ]