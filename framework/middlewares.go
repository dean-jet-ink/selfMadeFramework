package framework

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func TimeoutMiddleware(ctx *MyContext) {
	ch := make(chan struct{})
	panicCh := make(chan struct{})

	durationContext, cancel := context.WithTimeout(ctx.r.Context(), time.Second*10)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				panicCh <- struct{}{}
			}
		}()
		ctx.Next()
		ch <- struct{}{}
	}()

	defer cancel()

	select {
	case <-durationContext.Done():
		ctx.SetTimeout(true)
		ctx.WriteString("timeout")
	case <-ch:
		fmt.Println("success")
	case <-panicCh:
		ctx.WriteString("panic")
	}
}

func AuthUserMiddleware(ctx *MyContext) {
	ctx.SetValue("AuthUser", "test")
}

func TimeCostMiddleware(ctx *MyContext) {
	now := time.Now()
	ctx.Next()
	fmt.Println("time cost:", time.Since(now).Milliseconds())
}

func StaticFileMiddleware(ctx *MyContext) {
	pathName := strings.TrimSuffix(ctx.r.URL.Path, "/")
	pathName = path.Join("static", pathName)
	fileInfo, err := os.Stat(pathName)

	if err == nil && !fileInfo.IsDir() {
		fileServer := http.FileServer(http.Dir("static"))
		fileServer.ServeHTTP(ctx.w, ctx.r)
		ctx.Abort()
	}
}
