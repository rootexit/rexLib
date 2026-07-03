package rexDatabase

type Device struct {
	DeviceId      string `gorm:"column:device_id;comment:客户端ID;type:varchar(255);" json:"device_id"`         // 客户端ID
	DeviceVersion string `gorm:"column:device_version;comment:设备版本;type:varchar(64);" json:"device_version"` // 设备版本
	DeviceScore   string `gorm:"column:device_score;comment:设备评分;type:varchar(8);" json:"device_score"`      // 设备评分
}

type ClientNetwork struct {
	IpAddress string `gorm:"column:ip_address;comment:ip地址;type: varchar(255);" json:"ip_address"` // IpAddress
	Port      string `gorm:"column:port;comment:端口;type: varchar(64);" json:"port"`                // Port
	Network   string `gorm:"column:network;comment:网关数据;type: varchar(255);" json:"network"`       // Network
	Isp       string `gorm:"column:isp;comment:ISP;type: varchar(64);" json:"isp"`                 // ISP
	// AutonomousSystemOrganization for the registered autonomous system number.
	AutonomousSystemOrganization string `gorm:"column:autonomous_system_organization;comment:针对已注册自治系统编号的自治系统组织。;type: varchar(255);" json:"autonomous_system_organization"` //nolint:lll
	// AutonomousSystemNumber for the IP address.
	AutonomousSystemNumber uint `gorm:"column:autonomous_system_number;comment:该 IP 地址的自治系统号;type: varchar(64);"  json:"autonomous_system_number"` //nolint:lll
}

type ClientLocation struct {
	Continent     string `gorm:"column:continent;comment:洲;type: varchar(255);" json:"continent"`             // Continent
	ContinentCode string `gorm:"column:continent_code;comment:洲代码;type: varchar(255);" json:"continent_code"` // ContinentCode
	Country       string `gorm:"column:country;comment:国家;type: varchar(255);" json:"country"`                // CountryCode
	CountryCode   string `gorm:"column:country_code;comment:国家代码;type: varchar(255);" json:"country_code"`    // CountryCode
	Province      string `gorm:"column:province;comment:行政区划;type: varchar(255);" json:"province"`            // Province
	ProvinceCode  string `gorm:"column:province_code;comment:行政区划;type: varchar(255);" json:"province_code"`  // Province
	City          string `gorm:"column:city;comment:城市;type: varchar(255);" json:"city"`                      // City
	CityNameID    uint   `gorm:"column:city_name_id;comment:城市代码;type: varchar(255);" json:"city_name_id"`    // CityNameId

	Longitude float64 `gorm:"column:longitude;comment:经度;type:double precision;" json:"longitude"` // Longitude
	Latitude  float64 `gorm:"column:latitude;comment:纬度;type:double precision;" json:"latitude"`   // latitude
	TimeZone  string  `gorm:"column:time_zone;comment:时区;type:varchar(255);" json:"time_zone"`     // TimeZone

	AccuracyRadius uint16 `gorm:"column:accuracy_radius;comment:精度半径;type:int;" json:"accuracy_radius"` // AccuracyRadius
}

type ClientUa struct {
	UserAgent       string `gorm:"column:user_agent;comment:用户代理;type: varchar(255);" json:"user_agent"`                         // UserAgent
	UserAgentFamily string `gorm:"column:user_agent_family;comment:UserAgentFamily;type: varchar(64);" json:"user_agent_family"` // UserAgentFamily
	UserAgentMajor  string `gorm:"column:user_agent_major;comment:UserAgentMajor;type: varchar(64);" json:"user_agent_major"`    // UserAgentMajor
	UserAgentMinor  string `gorm:"column:user_agent_minor;comment:UserAgentMinor;type: varchar(64);" json:"user_agent_minor"`    // UserAgentMinor
	UserAgentPatch  string `gorm:"column:user_agent_patch;comment:UserAgentPatch;type: varchar(64);" json:"user_agent_patch"`    // UserAgentPatch
	OsFamily        string `gorm:"column:os_family;comment:OsFamily;type: varchar(64);" json:"os_family"`                        // OsFamily
	OsMajor         string `gorm:"column:os_major;comment:OsMajor;type: varchar(64);" json:"os_major"`                           // OsMajor
	OsMinor         string `gorm:"column:os_minor;comment:OsMinor;type: varchar(64);" json:"os_minor"`                           // OsMinor
	OsPatch         string `gorm:"column:os_patch;comment:OsPatch;type: varchar(64);" json:"os_patch"`                           // OsPatch
	OsPatchMinor    string `gorm:"column:os_patch_minor;comment:OsPatchMinor;type: varchar(64);" json:"os_patch_minor"`          // OsPatchMinor
	DeviceFamily    string `gorm:"column:device_family;comment:DeviceFamily;type: varchar(64);" json:"device_family"`            // DeviceFamily
	DeviceBrand     string `gorm:"column:device_brand;comment:DeviceBrand;type: varchar(64);" json:"device_brand"`               // DeviceBrand
	DeviceModel     string `gorm:"column:device_model;comment:DeviceModel;type: varchar(64);" json:"device_model"`
}

type ClientBot struct {
	IsBot       bool   `gorm:"column:is_bot;comment:是否机器人;type: boolean;" json:"is_bot"`
	BotName     string `gorm:"column:bot_name;comment:机器人名称;type: varchar(255);" json:"bot_name"`
	BotCategory string `gorm:"column:bot_category;comment:机器人类别;type: varchar(255);" json:"bot_category"`
	BotReason   string `gorm:"column:bot_reason;comment:机器人原因;type: varchar(255);" json:"bot_reason"`
}

type Client struct {
	ClientNetwork
	ClientLocation
	ClientUa
	ClientBot
}
