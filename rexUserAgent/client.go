package rexUserAgent

type Device struct {
	DeviceId      string  `gorm:"index:idx_client_id;column:client_id;comment:客户端ID;type: varchar(255);" json:"client_id"` // 客户端ID
	DeviceVersion string  `gorm:"column:device_version;comment:设备版本;type: varchar(64);" json:"device_version"`             // 设备版本
	DeviceScore   float32 `gorm:"column:device_score;comment:设备评分;type: float;" json:"device_score"`                       // 设备评分
}

type Client struct {
	IP              string `gorm:"column:ip;comment:IP地址;type: text;" json:"ip"`                                                 // IP地址
	IpKeychainName  string `gorm:"column:ip_keychain_name;comment:ip的加密name;type: varchar(64);" json:"ip_keychain_name"`         // ip的加密name
	IPHash          string `gorm:"index:idx_ip_hash;column:ip_hash;comment:IP地址hash值;type: varchar(255);" json:"ip_hash"`        // IP地址hash值
	Port            string `gorm:"column:port;comment:Port;type: varchar(64);" json:"port"`                                      // Port
	UserAgent       string `gorm:"column:user_agent;comment:UserAgent;type: text;" json:"user_agent"`                            // UserAgent
	CityId          int64  `gorm:"column:city_id;comment:CityId;type: bigint;" json:"city_id"`                                   // CityId
	Country         string `gorm:"column:country;comment:Country;type: varchar(64);" json:"country"`                             // Country
	Region          string `gorm:"column:region;comment:Region;type: varchar(64);" json:"region"`                                // Region
	Province        string `gorm:"column:province;comment:Province;type: varchar(64);" json:"province"`                          // Province
	City            string `gorm:"column:city;comment:City;type: varchar(64);" json:"city"`                                      // City
	ISP             string `gorm:"column:isp;comment:ISP;type: varchar(64);" json:"isp"`                                         // ISP
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
	DeviceModel     string `gorm:"column:device_model;comment:DeviceModel;type: varchar(64);" json:"device_model"`               // DeviceModel
	Device
}
