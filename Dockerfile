# ============================================
# STAGE 1: DEVELOPMENT (default, dengan hot reload)
# ============================================
FROM golang:1.26-alpine AS development

RUN apk add --no-cache \
    git \
    curl \
    bash \
    make \
    gcc \
    musl-dev \
    tzdata

ENV TZ=Asia/Jakarta

# Install air (hot reload) & templ
RUN go install github.com/air-verse/air@latest && \
    go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

# Copy go.mod & go.sum (kalau ada)
COPY src/go.mod src/go.sum* ./

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off

# Download dependencies di build time di-skip
# Akan dijalankan manual di runtime setelah container up
# RUN go mod download

# Copy source
COPY src/ .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]


# ============================================
# STAGE 2: PRODUCTION (untuk build final)
# ============================================
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app
COPY src/ .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:latest AS production

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Jakarta

WORKDIR /app
COPY --from=builder /server /app/server
COPY src/public /app/public
COPY src/views /app/views

EXPOSE 8080
CMD ["/app/server"]