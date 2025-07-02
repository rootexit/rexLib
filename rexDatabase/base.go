package rexDatabase

import (
	"gorm.io/gorm"
	"time"
)

// note： 兼容实现
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at;comment:创建时间;" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;comment:更新时间;" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UniqueIdAdminModel struct {
	CreatedBy string `gorm:"index:idx_created_by;column:created_by;comment:创建者;type: varchar(255)" json:"created_by"`
	UpdatedBy string `gorm:"index:idx_updated_by;column:updated_by;comment:更新者;type: varchar(255)" json:"updated_by"`
	DeletedBy string `gorm:"index:idx_deleted_by;column:deleted_by;comment:删除者;type: varchar(255)" json:"-"`
}

type BaseTenantModel struct {
	CreatedTenantBy string `gorm:"index:idx_created_tenant_by;column:created_tenant_by;comment:创建数据的租户;type: varchar(255)" json:"created_tenant_by"`
	UpdatedTenantBy string `gorm:"index:idx_updated_tenant_by;column:updated_tenant_by;comment:更新数据的租户;type: varchar(255)" json:"updated_tenant_by"`
	DeletedTenantBy string `gorm:"index:idx_deleted_tenant_by;column:deleted_tenant_by;comment:删除数据的租户;type: varchar(255)" json:"-"`
}

type BaseAdminModel struct {
	CreatedBy uint `gorm:"column:created_by;comment:创建者;type: int" json:"created_by"`
	UpdatedBy uint `gorm:"column:updated_by;comment:更新者;type: int" json:"updated_by"`
	DeletedBy uint `gorm:"column:deleted_by;comment:删除者;type: int" json:"-"`
}

// note: 审计字段新版实现
type AuditByStringModel struct {
	CreatedBy string `gorm:"index:idx_created_by;column:created_by;comment:创建者;type: varchar(255)" json:"created_by"`
	UpdatedBy string `gorm:"index:idx_updated_by;column:updated_by;comment:更新者;type: varchar(255)" json:"updated_by"`
	DeletedBy string `gorm:"index:idx_deleted_by;column:deleted_by;comment:删除者;type: varchar(255)" json:"-"`
}
type AuditByUintModel struct {
	CreatedBy uint `gorm:"index:idx_created_by;column:created_by;comment:创建者;type:int" json:"created_by"`
	UpdatedBy uint `gorm:"index:idx_updated_by;column:updated_by;comment:更新者;type:int" json:"updated_by"`
	DeletedBy uint `gorm:"index:idx_deleted_by;column:deleted_by;comment:删除者;type:int" json:"-"`
}

// note: rpc新版实现
type BaseRpcModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at;comment:创建时间;" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;comment:更新时间;" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
type AuditRpcByUintModel struct {
	CreatedBy string `gorm:"index:idx_created_by;column:created_by;comment:创建者;type:int" json:"createdBy"`
	UpdatedBy string `gorm:"index:idx_updated_by;column:updated_by;comment:更新者;type:int" json:"updatedBy"`
	DeletedBy string `gorm:"index:idx_deleted_by;column:deleted_by;comment:删除者;type:int" json:"-"`
}
type AuditRpcByStringModel struct {
	CreatedBy string `gorm:"index:idx_created_by;column:created_by;comment:创建者;type: varchar(255)" json:"createdBy"`
	UpdatedBy string `gorm:"index:idx_updated_by;column:updated_by;comment:更新者;type: varchar(255)" json:"updatedBy"`
	DeletedBy string `gorm:"index:idx_deleted_by;column:deleted_by;comment:删除者;type: varchar(255)" json:"-"`
}

// note: 审计模型租户版
type AuditTenantByStringModel struct {
	CreatedTenantBy string `gorm:"index:idx_created_tenant_by;column:created_tenant_by;comment:创建数据的租户;type: varchar(255)" json:"created_tenant_by"`
	UpdatedTenantBy string `gorm:"index:idx_updated_tenant_by;column:updated_tenant_by;comment:更新数据的租户;type: varchar(255)" json:"updated_tenant_by"`
	DeletedTenantBy string `gorm:"index:idx_deleted_tenant_by;column:deleted_tenant_by;comment:删除数据的租户;type: varchar(255)" json:"-"`
}
type AuditTenantRpcByStringModel struct {
	CreatedTenantBy string `gorm:"index:idx_created_tenant_by;column:created_tenant_by;comment:创建数据的租户;type: varchar(255)" json:"createdTenantBy"`
	UpdatedTenantBy string `gorm:"index:idx_updated_tenant_by;column:updated_tenant_by;comment:更新数据的租户;type: varchar(255)" json:"updatedTenantBy"`
	DeletedTenantBy string `gorm:"index:idx_deleted_tenant_by;column:deleted_tenant_by;comment:删除数据的租户;type: varchar(255)" json:"-"`
}
type AuditTenantByUintModel struct {
	CreatedTenantBy uint `gorm:"index:idx_created_tenant_by;column:created_tenant_by;comment:创建数据的租户;type:int" json:"created_tenant_by"`
	UpdatedTenantBy uint `gorm:"index:idx_updated_tenant_by;column:updated_tenant_by;comment:更新数据的租户;type:int" json:"updated_tenant_by"`
	DeletedTenantBy uint `gorm:"index:idx_deleted_tenant_by;column:deleted_tenant_by;comment:删除数据的租户;type:int" json:"-"`
}
type AuditTenantRpcByUintModel struct {
	CreatedTenantBy uint `gorm:"index:idx_created_tenant_by;column:created_tenant_by;comment:创建数据的租户;type:int" json:"createdTenantBy"`
	UpdatedTenantBy uint `gorm:"index:idx_updated_tenant_by;column:updated_tenant_by;comment:更新数据的租户;type:int" json:"updatedTenantBy"`
	DeletedTenantBy uint `gorm:"index:idx_deleted_tenant_by;column:deleted_tenant_by;comment:删除数据的租户;type:int" json:"-"`
}
