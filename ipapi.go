package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

// MockIPAPIResponse represents the json response of the ip-api.com API
// Documentation can be found at https://ip-api.com/docs/api:json
type MockIPAPIResponse struct {
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

// IPAPI is the http://ip-api.com/ mock handler for the /{ip} route
func IPAPI(ctx *fasthttp.RequestCtx) {
	if isFailure() {
		fmt.Fprintf(ctx, "%s", []byte(`error`))
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	cc := getRandomCountryCode()

	m := MockIPAPIResponse{
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
