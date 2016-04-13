package main

import (
	"fmt"
	"github.com/logikal/chameleon/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/logikal/chameleon/Godeps/_workspace/src/github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

/*
export CHAMELEON_MASQUERADE=apid
export CHAMELEON_PORT=9999
export CHAMELEON_HEALTHY=false
export CHAMELEON_MINWAIT="100ms"
export CHAMELEON_MAXWAIT="1s"
*/

type Specification struct {
	// Masquerade: Service name to masquerade as
	Masquerade string

	// Port: The port to listen on
	Port int

	// Healthy: Default healthy or not
	Healthy bool

	// MinWait: Minimum time to wait before responding
	MinWait time.Duration

	// MaxWait: Maximum time to wait before responding
	MaxWait time.Duration
}

type HealthCheck struct {
	Name    string "json:s.Masquerade"
	Healthy string "json:s.Healthy"
}

func main() {
	var s Specification
	err := envconfig.Process("chameleon", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	format := "Masquerade: %s\nPort: %d\nHealthy: %v\nMinWait: %s\nMaxWait: %s\n"
	_, err = fmt.Printf(format, s.Masquerade, s.Port, s.Healthy, s.MinWait, s.MaxWait)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Also grab the hostname for later, since it's required in our healthcheck.
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println(hostname)
	// set up a channel to receive signals on.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGUSR1)

	// If we get SIGUSR1, toggle the Healthy state.
	go func() {
		for {
			select {
			case <-sig:
				log.Println("Received USR1 signal")
				s.Healthy = !s.Healthy
				_, err = fmt.Printf(format, s.Masquerade, s.Port, s.Healthy, s.MinWait, s.MaxWait)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		}
	}()

	// set up the webserver
	r := gin.Default()

	// Routes for everyone!

	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"host":     hostname,
			"app_name": s.Masquerade,
			"message":  "pong",
			"results": gin.H{
				s.Masquerade: gin.H{
					"ok": s.Healthy,
					"message": "null",
				},
			},
		"healthcheck_version": 1,
		})
	})

	// start it up!
	portformat := ":%s"
	r.Run(fmt.Sprintf(portformat, strconv.Itoa(s.Port)))

}
