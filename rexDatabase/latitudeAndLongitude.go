package rexDatabase

type CommonLocationModel struct {
	LocationName string `gorm:"column:location_name;comment:位置名称;type:varchar(255)" json:"location_name"`
	LonLat
}

// lon,lat（经度,纬度
type LonLat struct {
	Longitude float64 `gorm:"column:longitude;comment:经度;type:decimal(10,6)" json:"longitude"`
	Latitude  float64 `gorm:"column:latitude;comment:纬度;type:decimal(10,6)" json:"latitude"`
}
