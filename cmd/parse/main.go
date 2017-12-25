package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	"housing"
)

var (
	gcpAPIKey string
	cachefile string
	dirname   string
)

func init() {
	flag.StringVar(&gcpAPIKey, "gcpAPIKey", "", "GCP API Key for Google Maps Geocoding API")
	flag.StringVar(&cachefile, "cachefile", "", "cache file for prefetched geocoding results")
	flag.StringVar(&dirname, "dirname", "", "directory containing 實價登錄 files")
}

func filterOut(fname string, rowID int, row []string) bool {
	土地區段位置或建物區門牌 := row[2]
	if 土地區段位置或建物區門牌 == "" {
		glog.Errorf("empty address %s:%d %+v", fname, rowID, row)
		return true
	}

	編號 := row[27]
	if 編號 == "RPPQMLPJNHMFFGE68CA" {
		// 2017Q3 E_lvr_land_A.CSV:499
		// Invalid 建築完成年月 1991-02-30
		// [鳳山區 房地(土地+建物) 高雄市鳳山區青年路二段181~210號 19.29 住   1060331 土地0建物1車位0 十六層 十八層 住宅大樓(11層含以上有電梯) 住家用 鋼筋混凝土造 0800230 162.05 4 2 2 有 有 4860000 29991  0.0 0  RPPQMLPJNHMFFGE68CA]
		return true
	}
	if 編號 == "RPUNMLQKOHMFFAL66CA" {
		// 2017Q3 B_lvr_land_A.CSV:5579
		// Invalid 建築完成年月 1990-02-29
		// [豐原區 房地(土地+建物) 臺中市豐原區三豐路二段241~270號 134.85 住   1060505 土地1建物1車位0 全 三層 透天厝 見其他登記事項 鋼筋混凝土造 0790229 231.38 5 2 5 有 無 13700000 59210  0.0 0  RPUNMLQKOHMFFAL66CA]
		return true
	}
	if 編號 == "RPSNMLLJPHMFFIB38CA" {
		// 2017Q3 B_lvr_land_A.CSV:6618
		// Invalid 建築完成年月 1985-02-30
		// [霧峰區 房地(土地+建物) 臺中市霧峰區文化巷1~30號 80.0 住   1060531 土地2建物1車位0 全 二層 透天厝 住商用 加強磚造 0740230 129.01 3 2 2 有 無 2000000 15503  0.0 0 親友、員工或其他特殊關係間之交易。 RPSNMLLJPHMFFIB38CA]
		return true
	}

	return false
}

func parse(fname string, rowID int, row []string, geocoder *housing.Geocoder) error {
	if filterOut(fname, rowID, row) {
		return nil
	}

	ts, err := housing.ParseRow(row, geocoder)
	if err != nil {
		if _, ok := errors.Cause(err).(*housing.GeocodeNoResultsError); ok {
			return nil
		}
		return errors.Wrap(err, fmt.Sprintf("housing.ParseRow %s %d %+v", fname, rowID, row))
	}

	b, err := json.Marshal(ts)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}
	fmt.Printf("%s\n", b)
	return nil
}

func main() {
	flag.Parse()

	// We use a large precision, since the cache already contains all attempted to geocoded all addresses.
	var precisionMeters float64 = 999999
	geocoder := housing.NewGeocoder(gcpAPIKey, precisionMeters)
	if cachefile != "" {
		geocoder.PopulateCache(cachefile)
	}

	rowFn := func(fname string, rowID int, row []string) error {
		return parse(fname, rowID, row, geocoder)
	}
	err := housing.ScanDir(dirname, rowFn)
	if err != nil {
		glog.Errorf("%+v", err)
	}
}
