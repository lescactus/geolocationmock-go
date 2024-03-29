package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

// MockIPBaseResponse represents the json response of the ipbase.com API
// Documentation can be found at https://ipbase.com/docs/info
type MockIPBaseResponse struct {
	Data MockIPBaseResponseData `json:"data"`
}

type MockIPBaseResponseData struct {
	Timezone   MockIPBaseResponseDataTimeZone   `json:"timezone"`
	IP         string                           `json:"ip"`
	Type       string                           `json:"type"`
	Connection MockIPBaseResponseDataConnection `json:"connection"`
	Location   MockIPBaseResponseDataLocation   `json:"location"`
}

type MockIPBaseResponseDataTimeZone struct {
	ID               string    `json:"id"`
	CurrentTime      time.Time `json:"current_time"`
	Code             string    `json:"code"`
	IsDaylightSaving bool      `json:"is_daylight_saving"`
	GmtOffset        int       `json:"gmt_offset"`
}

type MockIPBaseResponseDataConnection struct {
	Asn          int    `json:"asn"`
	Organization string `json:"organization"`
	Isp          string `json:"isp"`
}

type MockIPBaseResponseDataLocation struct {
	GeonamesID int                                     `json:"geonames_id"`
	Latitude   float64                                 `json:"latitude"`
	Longitude  float64                                 `json:"longitude"`
	Zip        string                                  `json:"zip"`
	Continent  MockIPBaseResponseDataLocationContinent `json:"continent"`
	Country    MockIPBaseResponseDataLocationCountry   `json:"country"`
	City       MockIPBaseResponseDataLocationCity      `json:"city"`
	Region     MockIPBaseResponseDataLocationRegion    `json:"region"`
}

type MockIPBaseResponseDataLocationContinent struct {
	Code           string `json:"code"`
	Name           string `json:"name"`
	NameTranslated string `json:"name_translated"`
}

type MockIPBaseResponseDataLocationCity struct {
	Name           string `json:"name"`
	NameTranslated string `json:"name_translated"`
}

type MockIPBaseResponseDataLocationRegion struct {
	Fips           string `json:"fips"`
	Alpha2         string `json:"alpha2"`
	Name           string `json:"name"`
	NameTranslated string `json:"name_translated"`
}

type MockIPBaseResponseDataLocationCountry struct {
	Alpha2       string   `json:"alpha2"`
	Alpha3       string   `json:"alpha3"`
	CallingCodes []string `json:"calling_codes"`
	Currencies   []struct {
		Symbol        string `json:"symbol"`
		Name          string `json:"name"`
		SymbolNative  string `json:"symbol_native"`
		DecimalDigits int    `json:"decimal_digits"`
		Rounding      int    `json:"rounding"`
		Code          string `json:"code"`
		NamePlural    string `json:"name_plural"`
	} `json:"currencies"`
	Emoji     string `json:"emoji"`
	Ioc       string `json:"ioc"`
	Languages []struct {
		Name       string `json:"name"`
		NameNative string `json:"name_native"`
	} `json:"languages"`
	Name              string   `json:"name"`
	NameTranslated    string   `json:"name_translated"`
	Timezones         []string `json:"timezones"`
	IsInEuropeanUnion bool     `json:"is_in_european_union"`
}

var (
	OneOneOneOne = MockIPBaseResponse{
		Data: MockIPBaseResponseData{
			IP: "1.1.1.1",
			Location: MockIPBaseResponseDataLocation{
				Latitude:  -27.4766,
				Longitude: 153.0166,
				City: MockIPBaseResponseDataLocationCity{
					Name: "South Brisbane",
				},
				Country: MockIPBaseResponseDataLocationCountry{
					Name:   "Australia",
					Alpha2: "AU",
				},
			},
		},
	}

	TwoTwoTwoTwo = MockIPBaseResponse{
		Data: MockIPBaseResponseData{
			IP: "2.2.2.2",
			Location: MockIPBaseResponseDataLocation{
				Latitude:  48.8566,
				Longitude: 2.35222,
				City: MockIPBaseResponseDataLocationCity{
					Name: "Paris",
				},
				Country: MockIPBaseResponseDataLocationCountry{
					Name:   "France",
					Alpha2: "FR",
				},
			},
		},
	}

	ThreeThreeThreeThree = MockIPBaseResponse{
		Data: MockIPBaseResponseData{
			IP: "3.3.3.3",
			Location: MockIPBaseResponseDataLocation{
				Latitude:  39.0469,
				Longitude: -77.4903,
				City: MockIPBaseResponseDataLocationCity{
					Name: "Ashburn",
				},
				Country: MockIPBaseResponseDataLocationCountry{
					Name:   "United States",
					Alpha2: "US",
				},
			},
		},
	}
)

func IPBase(ctx *fasthttp.RequestCtx) {
	if isFailure() {
		fmt.Fprintf(ctx, "%s", []byte(`error`))
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	cc := getRandomCountryCode()

	m := MockIPBaseResponse{
		Data: MockIPBaseResponseData{
			IP: fmt.Sprintf("%s", ctx.UserValue("ip")),
			Location: MockIPBaseResponseDataLocation{
				Latitude:  getRandomCoordinate(LatMin, LatMax),
				Longitude: getRandomCoordinate(LonMin, LonMax),
				City: MockIPBaseResponseDataLocationCity{
					Name: getRandomCity(),
				},
				Country: MockIPBaseResponseDataLocationCountry{
					Name:   cc.String(),
					Alpha2: cc.Alpha2(),
				},
			},
		},
	}

	resp, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	time.Sleep(*latency)
	fmt.Println(string(resp)) // just for a basic logging
	fmt.Fprintf(ctx, "%s", []byte(resp))
}

func IPBaseOneONeOneOne(ctx *fasthttp.RequestCtx) {
	resp, err := json.Marshal(OneOneOneOne)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	fmt.Println(string(resp)) // just for a basic logging
	fmt.Fprintf(ctx, "%s", []byte(resp))
}

func IPBaseTwoTwoTwoTwo(ctx *fasthttp.RequestCtx) {
	resp, err := json.Marshal(TwoTwoTwoTwo)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	fmt.Println(string(resp)) // just for a basic logging
	fmt.Fprintf(ctx, "%s", []byte(resp))
}

func IPBaseThreeThreeThreeThree(ctx *fasthttp.RequestCtx) {
	resp, err := json.Marshal(ThreeThreeThreeThree)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	fmt.Println(string(resp)) // just for a basic logging
	fmt.Fprintf(ctx, "%s", []byte(resp))
}
