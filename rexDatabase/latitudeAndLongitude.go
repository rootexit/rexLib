package rexDatabase

type CommonLocationModel struct {
	FormattedAddress string `gorm:"column:formatted_address;comment:格式化后的完整地址;type:varchar(255)" json:"formatted_address"`
	Country          string `gorm:"column:country;comment:国家;type: varchar(64);" json:"country"`           // Country
	Province         string `gorm:"column:province;comment:省份;type: varchar(64);" json:"province"`         // Province
	City             string `gorm:"column:city;comment:城市;type: varchar(64);" json:"city"`                 // City
	CityCode         string `gorm:"column:city_code;comment:城市行政区划代码;type: varchar(64);" json:"city_code"` // CityCode
	District         string `gorm:"column:district;comment:区域;type: varchar(64);" json:"district"`         // District
	Street           string `gorm:"column:street;comment:街道;type: varchar(255);" json:"street"`            // Street
	Number           string `gorm:"column:number;comment:门牌号;type: varchar(64);" json:"number"`            // Number
	LonLat
}

// lon,lat（经度,纬度
type LonLat struct {
	Longitude float64 `gorm:"column:longitude;comment:经度;type:decimal(10,6)" json:"longitude"`
	Latitude  float64 `gorm:"column:latitude;comment:纬度;type:decimal(10,6)" json:"latitude"`
}
