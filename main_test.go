package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

// func TestMainFunc(t *testing.T) {
// 	go main()
// }

func TestRun(t *testing.T) {
	// net.Listenerを作成
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}
	// キャンセル可能なコンテキストを作成。
	ctx, cancel := context.WithCancel(context.Background())
	// errgroup でゴルーチンをまとめて管理
	// 別ゴルーチンで run(ctx) を呼ぶ  HTTP サーバーをバックグラウンドで起動
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx, l)
	})
	// サーバーにリクエスト
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	// どんなポート番号でリッスンしているのか確認
	t.Logf("try request to %q", url)
	rsp, err := http.Get(url)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	// 終了時に自動クローズ(リソースリーク防止)
	defer rsp.Body.Close()
	// レスポンスボディ読み込み
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	// 戻り値を検証する
	want := fmt.Sprintf("Hello, %s!", in)
	// want := fmt.Sprintf("aaaaaaaHello, %s!", in) //失敗確認用
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
	// コンテキストをキャンセル → サーバーを停止
	cancel()
	// errgroup がゴルーチンの終了を待つ
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
