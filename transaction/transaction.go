package transaction

type Transaction struct {
	A鄉鎮市區         string  `json:"鄉鎮市區,omitempty"`
	A交易標的         string  `json:"交易標的,omitempty"`
	A土地區段位置或建物區門牌 string  `json:"土地區段位置或建物區門牌,omitempty"`
	A土地移轉總面積平方公尺  float64 `json:"土地移轉總面積平方公尺,omitempty"`
	A都市土地使用分區     string  `json:"都市土地使用分區,omitempty"`
	A非都市土地使用分區    string  `json:"非都市土地使用分區,omitempty"`
	A非都市土地使用編定    string  `json:"非都市土地使用編定,omitempty"`
	A交易年月日        int64   `json:"交易年月日,omitempty"`
	A交易筆棟數        string  `json:"交易筆棟數,omitempty"`
	A移轉層次         string  `json:"移轉層次,omitempty"`
	A總樓層數         string  `json:"總樓層數,omitempty"`
	A建物型態         string  `json:"建物型態,omitempty"`
	A主要用途         string  `json:"主要用途,omitempty"`
	A主要建材         string  `json:"主要建材,omitempty"`
	A建築完成年月       int64   `json:"建築完成年月,omitempty"`
	A建物移轉總面積平方公尺  float64 `json:"建物移轉總面積平方公尺,omitempty"`
	A建物現況格局_房     int     `json:"建物現況格局_房,omitempty"`
	A建物現況格局_廳     int     `json:"建物現況格局_廳,omitempty"`
	A建物現況格局_衛     int     `json:"建物現況格局_衛,omitempty"`
	A建物現況格局_隔間    string  `json:"建物現況格局_隔間,omitempty"`
	A有無管理組織       string  `json:"有無管理組織,omitempty"`
	A總價元          int     `json:"總價元,omitempty"`
	A單價每平方公尺      int     `json:"單價每平方公尺,omitempty"`
	A車位類別         string  `json:"車位類別,omitempty"`
	A車位移轉總面積平方公尺  float64 `json:"車位移轉總面積平方公尺,omitempty"`
	A車位總價元        int     `json:"車位總價元,omitempty"`
	A備註           string  `json:"備註,omitempty"`
	A編號           string  `json:"編號,omitempty"`

	Lat float64 `json:",omitempty"`
	Lng float64 `json:",omitempty"`
}
