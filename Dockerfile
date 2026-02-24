# ---------- Stage 1: Build ----------
    FROM golang:1.22-alpine AS builder

    WORKDIR /app
    
    COPY go.mod go.sum* ./
    RUN go mod download
    
    COPY . .
    
    # билдим именно cmd
    RUN CGO_ENABLED=0 GOOS=linux go build -o hot-coffee ./cmd
    
    # ---------- Stage 2: Runtime ----------
    FROM alpine:latest
    
    RUN apk --no-cache add ca-certificates
    
    WORKDIR /root/
    
    COPY --from=builder /app/hot-coffee .
    COPY --from=builder /app/data ./data
    COPY --from=builder /app/frontend ./frontend
    
    EXPOSE 8080
    
    CMD ["./hot-coffee"]