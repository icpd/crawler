package boot

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/whoisix/subscribe2clash/app/router"
	"github.com/whoisix/subscribe2clash/internal/global"
)

func initHttpServer() {

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	router.RegisterRouter(r)

	srv := &http.Server{
		Addr:    global.Listen,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Add(-1)
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		log.Println("shutdown server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("server shutdown: %v", err)
		}
	}()
	wg.Wait()

	log.Println("server exiting")
}
