package housing

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/traditionalchinese"

	"housing/transaction"
)

func atoiOrDashDashOne(s string) (int, error) {
	if s == "--" {
		return 1, nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return -1, errors.Wrap(err, fmt.Sprintf("strconv.Atoi %s", s))
	}
	return i, nil
}

func parseROCDate(rocDate string) (int64, error) {
	dayStr := ""
	monStr := ""
	yearStr := ""
	if len(rocDate) == 3 && rocDate[0] == '0' {
		dayStr = "--"
		monStr = "--"
		yearStr = rocDate
	} else if len(rocDate) == 5 && rocDate[0] == '0' {
		dayStr = "--"
		monStr = rocDate[len(rocDate)-2:]
		yearStr = rocDate[:len(rocDate)-2]
	} else if len(rocDate) == 6 {
		dayStr = rocDate[len(rocDate)-2:]
		monStr = rocDate[len(rocDate)-4 : len(rocDate)-2]
		yearStr = rocDate[:len(rocDate)-4]
	} else if len(rocDate) == 7 {
		dayStr = rocDate[len(rocDate)-2:]
		monStr = rocDate[len(rocDate)-4 : len(rocDate)-2]
		yearStr = rocDate[:len(rocDate)-4]
	} else {
		return -1, fmt.Errorf("invalid ROC date %s", rocDate)
	}

	dayInt, err := atoiOrDashDashOne(dayStr)
	if err != nil {
		return -1, errors.Wrap(err, "parseDay")
	}
	monInt, err := atoiOrDashDashOne(monStr)
	if err != nil {
		return -1, errors.Wrap(err, "parseMon")
	}
	yearInt, err := strconv.Atoi(yearStr)
	if err != nil {
		return -1, errors.Wrap(err, "parseYear")
	}

	shortForm := "2006-01-02"
	tryStr := fmt.Sprintf("%04d-%02d-%02d", yearInt+1911, monInt, dayInt)
	tm, err := time.Parse(shortForm, tryStr)
	if err != nil {
		return -1, errors.Wrap(err, fmt.Sprintf("time.Parse %s", rocDate))
	}

	return tm.Unix(), nil
}

type parser struct {
	err error
}

func (p *parser) Error() error {
	return p.err
}

func (p *parser) parseInt(s, desc string) int {
	if p.err != nil {
		return -1
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		p.err = errors.Wrap(err, desc)
		return -1
	}
	return i
}

func (p *parser) parseFloat(s, desc string) float64 {
	if p.err != nil {
		return -1
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		p.err = errors.Wrap(err, desc)
		return -1
	}
	return f
}

func (p *parser) parseROCDate(s, desc string) int64 {
	if p.err != nil {
		return -1
	}
	dt, err := parseROCDate(s)
	if err != nil {
		p.err = errors.Wrap(err, desc)
		return -1
	}
	return dt
}

func (p *parser) parseROCDateIfNotEmpty(s, desc string) int64 {
	if p.err != nil {
		return -1
	}
	if s == "" {
		return -1
	}
	return p.parseROCDate(s, desc)
}

func ParseRow(row []string, geocoder *Geocoder) (*transaction.Transaction, error) {
	p := &parser{}
	ts := transaction.Transaction{}
	ts.A鄉鎮市區 = row[0]
	ts.A交易標的 = row[1]
	ts.A土地區段位置或建物區門牌 = row[2]
	ts.A土地移轉總面積平方公尺 = p.parseFloat(row[3], "土地移轉總面積平方公尺")
	ts.A都市土地使用分區 = row[4]
	ts.A非都市土地使用分區 = row[5]
	ts.A非都市土地使用編定 = row[6]
	ts.A交易年月日 = p.parseROCDate(row[7], "交易年月日")
	ts.A交易筆棟數 = row[8]
	ts.A移轉層次 = row[9]
	ts.A總樓層數 = row[10]
	ts.A建物型態 = row[11]
	ts.A主要用途 = row[12]
	ts.A主要建材 = row[13]
	ts.A建築完成年月 = p.parseROCDateIfNotEmpty(row[14], "建築完成年月")
	ts.A建物移轉總面積平方公尺 = p.parseFloat(row[15], "建物移轉總面積平方公尺")
	ts.A建物現況格局_房 = p.parseInt(row[16], "建物現況格局_房")
	ts.A建物現況格局_廳 = p.parseInt(row[17], "建物現況格局_廳")
	ts.A建物現況格局_衛 = p.parseInt(row[18], "建物現況格局_衛")
	ts.A建物現況格局_隔間 = row[19]
	ts.A有無管理組織 = row[20]
	ts.A總價元 = p.parseInt(row[21], "總價元")
	ts.A單價每平方公尺 = p.parseInt(row[22], "單價每平方公尺")
	ts.A車位類別 = row[23]
	ts.A車位移轉總面積平方公尺 = p.parseFloat(row[24], "車位移轉總面積平方公尺")
	ts.A車位總價元 = p.parseInt(row[25], "車位總價元")
	ts.A備註 = row[26]
	ts.A編號 = row[27]

	lat, lng, err := geocoder.GeocodeWithRetry(row[2])
	if err != nil {
		return nil, errors.Wrap(err, "Geocode")
	}
	ts.Lat = lat
	ts.Lng = lng

	if err := p.Error(); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%+v", row))
	}
	return &ts, nil
}

func checkHeader(row []string) error {
	vals := []string{
		"鄉鎮市區",
		"交易標的",
		"土地區段位置或建物區門牌",
		"土地移轉總面積平方公尺",
		"都市土地使用分區",
		"非都市土地使用分區",
		"非都市土地使用編定",
		"交易年月日",
		"交易筆棟數",
		"移轉層次",
		"總樓層數",
		"建物型態",
		"主要用途",
		"主要建材",
		"建築完成年月",
		"建物移轉總面積平方公尺",
		"建物現況格局-房",
		"建物現況格局-廳",
		"建物現況格局-衛",
		"建物現況格局-隔間",
		"有無管理組織",
		"總價元",
		"單價每平方公尺",
		"車位類別",
		"車位移轉總面積平方公尺",
		"車位總價元",
		"備註",
		"編號",
	}
	if len(row) != len(vals) {
		return fmt.Errorf("number of columns %d not equal to %d", len(row), len(vals))
	}
	for i, col := range row {
		if col != vals[i] {
			return fmt.Errorf("column %d %s not equal to %s", i, col, vals[i])
		}
	}

	return nil
}

func ScanFile(fname string, rowFn func(fname string, rowID int, row []string) error) error {
	f, err := os.Open(fname)
	if err != nil {
		return errors.Wrap(err, "os.Open")
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.Wrap(err, "ioutil.Readall")
	}
	decoder := traditionalchinese.Big5.NewDecoder()
	big5Decoded, err := decoder.String(string(b))
	if err != nil {
		glog.Errorf("%+v", err)
		return errors.Wrap(err, "Big5.Decode")
	}
	r := csv.NewReader(bytes.NewReader([]byte(big5Decoded)))
	records, err := r.ReadAll()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("csv.Readall %s", fname))
	}

	if err := checkHeader(records[0]); err != nil {
		return errors.Wrap(err, "checkHeader")
	}
	offset := 1
	for i, rec := range records[offset:] {
		交易標的 := rec[1]
		if 交易標的 == "土地" || 交易標的 == "車位" {
			continue
		}
		單價每平方公尺 := rec[22]
		if 單價每平方公尺 == "" {
			continue
		}

		if err := rowFn(fname, i+offset, rec); err != nil {
			return err
		}
	}

	return nil
}

func ScanDir(dirname string, rowFn func(fname string, rowID int, row []string) error) error {
	counties := make(map[string]string)
	counties["C"] = "基隆市"
	counties["A"] = "臺北市"
	counties["F"] = "新北市"
	counties["H"] = "桃園縣"
	counties["O"] = "新竹市"
	counties["J"] = "新竹縣"
	counties["K"] = "苗栗縣"
	counties["B"] = "臺中市"
	counties["M"] = "南投縣"
	counties["N"] = "彰化縣"
	counties["P"] = "雲林縣"
	counties["I"] = "嘉義市"
	counties["Q"] = "嘉義縣"
	counties["D"] = "臺南市"
	counties["E"] = "高雄市"
	counties["T"] = "屏東縣"
	counties["G"] = "宜蘭縣"
	counties["U"] = "花蓮縣"
	counties["V"] = "臺東縣"
	counties["X"] = "澎湖縣"
	counties["W"] = "金門縣"
	counties["Z"] = "連江縣"

	tradeTypes := make(map[string]string)
	tradeTypes["A"] = "不動產買賣"
	// We do not handle B:預售屋買賣 since they only have 地號 instead of addresses that are geocodable.
	// We do not handle C:不動產租賃 since they are not purchases.

	for county, _ := range counties {
		for trade, _ := range tradeTypes {
			basename := fmt.Sprintf("%s_lvr_land_%s.CSV", county, trade)
			fname := filepath.Join(dirname, basename)
			if err := ScanFile(fname, rowFn); err != nil {
				return errors.Wrap(err, "ScanFile")
			}
		}
	}

	return nil
}
