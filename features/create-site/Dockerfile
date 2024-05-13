# 使用帶有 Go 環境的官方基礎映像
FROM golang:1.18 as builder

# 安裝 curl 和其他必要工具
RUN apt-get update && apt-get install -y curl

# 下載並安裝 kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
    && install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# 下載並安裝 Helm
RUN curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# 設定工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum 文件，然後下載依賴
COPY go.mod go.sum ./
RUN go mod download

# 複製 Go 程式碼
COPY . .

# 編譯 Go 程式
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp .

# 使用較小的基礎映像來部署
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 從建構階段複製編譯好的應用和 Helm/kubectl
COPY --from=builder /app/myapp .
COPY --from=builder /usr/local/bin/helm /usr/local/bin/helm
COPY --from=builder /usr/local/bin/kubectl /usr/local/bin/kubectl

# 當容器啟動時運行 Go 程式
CMD ["./myapp"]
