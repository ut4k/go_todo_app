# go 1.24板にかきなおしたやつ
# airが動かないため
# -------------------------------------------------------
# デプロイ用バイナリ作成
FROM golang:1.24-bullseye as deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app

# -------------------------------------------------------
# デプロイ用軽量コンテナ
FROM debian:bullseye-slim as deploy

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=deploy-builder /app/app .

CMD ["./app"]

# -------------------------------------------------------
# ローカル開発用ホットリロード
FROM golang:1.24 as dev

WORKDIR /app

# RUN go install github.com/cosmtrek/air/v2@latest #ふるい
RUN go install github.com/air-verse/air@latest

ENV PATH=$PATH:/go/bin

CMD ["air"]

