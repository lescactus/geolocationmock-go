package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
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
	latency = flag.Duration("latency", 0, "Response request latency")
	failure = flag.Int("failure", 0, "Failure response rate")
)

// MockResponse represents the json response of the ip-api.com API
// Documentation can be found at https://ip-api.com/docs/api:json
type MockResponse struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

// IP is the handler for the /{ip} route
func IP(ctx *fasthttp.RequestCtx) {
	if isFailure() {
		fmt.Fprintf(ctx, "%s", []byte(`error`))
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	cc := getRandomCountryCode()

	m := MockResponse{
		Query:       fmt.Sprintf("%s", ctx.UserValue("ip")),
		CountryCode: cc.Alpha2(),
		Country:     cc.String(),
		City:        getRandomCity(),
		Lat:         getRandomCoordinate(LatMin, LatMax),
		Lon:         getRandomCoordinate(LonMin, LonMax),
	}

	resp, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	time.Sleep(*latency)
	fmt.Println(string(resp)) // just for a basic logging
	fmt.Fprintf(ctx, "%s", []byte(resp))
}

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
	r.GET("/{ip}", IP)

	// hasthttp server
	server := &fasthttp.Server{
		Concurrency: 16 * 256 * 1024,
		Handler:     r.Handler,
	}

	// listen and serve
	if err = server.Serve(ln); err != nil {
		panic(fmt.Sprintf("error in fasthttp Server: %s", err))
	}
}
