package main

import (
	"fmt"
	"github.com/logikal/chameleon/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/logikal/chameleon/Godeps/_workspace/src/github.com/kelseyhightower/envconfig"
	"github.com/logikal/chameleon/Godeps/_workspace/src/github.com/martini-contrib/render"
	"log"
	"time"
  "strconv"
)

/*
Masquerade: Service name to masquerade as
Port: The port to listen on
Healthy: Default healthy or not
MinWait: Minimum time to wait before responding
MaxWait: Maximum time to wait before responding

e.g.
export CHAMELEON_MASQUERADE=apid
export CHAMELEON_PORT=9999
export CHAMELEON_HEALTHY=false
export CHAMELEON_MINWAIT="100ms"
export CHAMELEON_MAXWAIT="1s"

*/

type Specification struct {
	Masquerade string
	Port       int
	Healthy    bool
	MinWait    time.Duration
	MaxWait    time.Duration
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

	app := martini.Classic()
	app.Use(render.Renderer())
  portformat := ":%s"

  app.RunOnAddr(fmt.Sprintf(portformat, strconv.Itoa(s.Port)))
	// run application
	app.Run()

}
