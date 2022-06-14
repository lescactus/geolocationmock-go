package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/biter777/countries"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
)

const (
	minRandomCountryCode = 2
	maxRandomCountryCode = 997

	LatMin = -90.0
	LatMax = 90.0
	LonMin = -180.0
	LonMax = 180.0
)

var (
	provider = flag.String("provider", "ipapi", "Provier name. Available: 'ipapi', 'ipbase'")
	latency  = flag.Duration("latency", 0, "Response request latency")
	failure  = flag.Int("failure", 0, "Failure response rate")
	e2e      = flag.Bool("e2e", false, "Enable 'e2e' mode: set fixed responses instead of random")
)

// Return *failure % true
func isFailure() bool {
	if *failure == 0 {
		return false
	}
	i := rand.Intn(100)

	if i > *failure {
		return true
	} else {
		return false
	}
}

func getRandomCountryCode() countries.CountryCode {
	var c countries.CountryCode
	var r int

	for c == countries.Unknown {
		// Countries codes are somewhere between 2 and 997
		r = rand.Intn(maxRandomCountryCode-minRandomCountryCode) + minRandomCountryCode
		c = countries.ByNumeric(r)
	}
	return c
}

func getRandomCity() string {
	i := rand.Intn(len(CityList))
	return CityList[i]
}

func getRandomCoordinate(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func main() {
	flag.Parse()

	// Init seed
	rand.Seed(time.Now().UnixNano())

	// listener
	ln, err := reuseport.Listen("tcp4", "0.0.0.0:8000")
	if err != nil {
		panic(fmt.Sprintf("error in reuseport listener: %s", err))
	}

	// router and route registration
	r := router.New()
	switch *provider {
	case "ipapi":
		r.GET("/{ip}", IPAPI)
	case "ipbase":
		if *e2e {
			r.GET("/ip/1.1.1.1", IPBaseOneONeOneOne)
			r.GET("/ip/2.2.2.2", IPBaseTwoTwoTwoTwo)
			r.GET("/ip/3.3.3.3", IPBaseThreeThreeThreeThree)
		}
		r.GET("/{ip:*}", IPBase)
	default:
		r.GET("/{ip:*}", IPAPI)
	}

	// hasthttp server
	server := &fasthttp.Server{
		Concurrency: 16 * 256 * 1024,
		Handler:     r.Handler,
	}

	// listen and serve
	fmt.Printf("Starting server mocking %s provider...\n", *provider)
	if err = server.Serve(ln); err != nil {
		panic(fmt.Sprintf("error in fasthttp Server: %s", err))
	}
}
