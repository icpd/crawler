package boot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/whoisix/subscribe2clash/app/router"
	"github.com/whoisix/subscribe2clash/library/global"
)

func initHttpServer() error {

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	router.RegisterRouter(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", global.ListenAddr, global.ListenPort),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	log.Println("server exiting")
	return nil
}
