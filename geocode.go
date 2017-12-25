package housing

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"

	"housing/util"
)

type GeocodeNoResultsError struct {
	addr string
}

func (e *GeocodeNoResultsError) Error() string {
	return fmt.Sprintf("no geocoding results for %s", e.addr)
}

type latlngprecision struct {
	lat       float64
	lng       float64
	precision float64 // meters
}

type Geocoder struct {
	APIKey          string
	PrecisionMeters float64
	cache           map[string]latlngprecision
}

func NewGeocoder(apiKey string, precision float64) *Geocoder {
	g := Geocoder{
		APIKey:          apiKey,
		PrecisionMeters: precision,
		cache:           make(map[string]latlngprecision),
	}
	return &g
}

func (g *Geocoder) PopulateCache(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("os.Open %s", fname))
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		b := []byte(scanner.Text())
		ga := struct {
			Addr      string
			Lat       float64
			Lng       float64
			Precision float64
		}{}
		if err := json.Unmarshal(b, &ga); err != nil {
			return errors.Wrap(err, fmt.Sprintf("%s", b))
		}
		g.cache[ga.Addr] = latlngprecision{lat: ga.Lat, lng: ga.Lng, precision: ga.Precision}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (g *Geocoder) GeocodeWithRetry(addr string) (float64, float64, error) {
	var geocodeErr error
	numRetries := 5
	for i := 0; i < numRetries; i++ {
		lat, lng, err := g.Geocode(addr)
		if err == nil {
			return lat, lng, nil
		}
		if gnrErr, ok := err.(*GeocodeNoResultsError); ok {
			return -1, -1, gnrErr
		}

		geocodeErr = err

		if i < numRetries-1 {
			<-time.After(time.Duration(i) * time.Second)
		}
	}
	return -1, -1, errors.Wrap(geocodeErr, "reversegeocode")
}

func (g *Geocoder) Geocode(addr string) (float64, float64, error) {
	llp, ok := g.cache[addr]
	if ok && llp.precision < g.PrecisionMeters {
		return llp.lat, llp.lng, nil
	}

	lat, lng, err := g.geocode(addr)
	if err != nil {
		return -1, -1, err
	}

	g.cache[addr] = latlngprecision{lat: lat, lng: lng, precision: 0}
	return lat, lng, nil
}

func (g *Geocoder) geocode(addr string) (float64, float64, error) {
	v := url.Values{
		"key":     {g.APIKey},
		"address": {addr},
	}
	urlStr := "https://maps.googleapis.com/maps/api/geocode/json?" + v.Encode()
	resp := struct {
		Results []struct {
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
		Status string `json:"status"`
	}{}
	_, respBody, err := util.JSONReq3("GET", urlStr, &resp)
	if err != nil {
		return -1, -1, errors.Wrap(err, "JSONReq3")
	}
	if resp.Status != "OK" {
		if resp.Status == "ZERO_RESULTS" {
			return -1, -1, &GeocodeNoResultsError{addr: addr}
		}
		return -1, -1, fmt.Errorf("google geo code: %s", respBody)
	}

	lat := resp.Results[0].Geometry.Location.Lat
	lng := resp.Results[0].Geometry.Location.Lng
	return lat, lng, nil
}
