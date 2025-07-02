package rexDao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type (
	Dao interface {
		GetDB() *gorm.DB
		Ping() error
		Close() error
		Create(ctx context.Context, tableName string, in interface{}) error
		Delete(ctx context.Context, tableName string, in interface{}, unscoped bool) error
		DeleteIds(ctx context.Context, tableName string, in interface{}, ids []uint, unscoped bool) error
		DeleteWhereAny(ctx context.Context, tableName string, in interface{}, unscoped bool, query interface{}, args ...interface{}) error
		First(ctx context.Context, tableName string, in interface{}, query interface{}, args ...interface{}) error
		Latest(ctx context.Context, tableName string, in interface{}, query interface{}, args ...interface{}) error
		Find(ctx context.Context, tableName string, in interface{}, query interface{}, args ...interface{}) error
		FindAndOrderByInterface(ctx context.Context, tableName string, orderBy string, in interface{}, query interface{}, args ...interface{}) error
		Update(ctx context.Context, tableName string, id uint, updates interface{}) error
		Count(ctx context.Context, tableName string, query interface{}, args ...interface{}) (int64, error)
		FindAndLimit(ctx context.Context, tableName string, limit, offset int, in interface{}, query interface{}, args ...interface{}) error
		FindAndLimitOrder(ctx context.Context, tableName, orderBy string, limit, offset int, in interface{}, query interface{}, args ...interface{}) error
		FindAndLimitAndSortInterface(ctx context.Context, tableName string, in interface{}, orderBy string, offset, limit int32, query interface{}, args ...interface{}) error
		DeleteWhereQuery(ctx context.Context, tableName string, in interface{}, unscoped bool, query interface{}, args ...interface{}) error
		UpdateWhereQuery(ctx context.Context, tableName string, updates interface{}, query interface{}, args ...interface{}) error
	}
	defaultDao struct {
		db *gorm.DB
	}
)

func NewDao(db *gorm.DB) Dao {
	return &defaultDao{
		db: db,
	}
}

func (d *defaultDao) GetDB() *gorm.DB {
	return d.db
}

func (d *defaultDao) Ping() error {
	// note: 设置5秒的超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	db, err := d.db.DB()
	if err != nil {
		return err
	}
	defer cancel()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return db.Ping()
	}
	return nil
}

func (d *defaultDao) Close() error {
	db, err := d.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (d *defaultDao) Create(ctx context.Context, tableName string, in interface{}) error {
	// todo: create one record
	if err := d.db.WithContext(ctx).Table(tableName).Create(in).Error; err != nil {
		return err
	}
	return nil
}

func (d *defaultDao) Delete(ctx context.Context, tableName string, in interface{}, unscoped bool) error {
	// todo: delete one record
	if unscoped {
		if err := d.db.WithContext(ctx).Table(tableName).Unscoped().Delete(in).Error; err != nil {
			return err
		}
	} else {
		if err := d.db.WithContext(ctx).Table(tableName).Delete(in).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *defaultDao) DeleteIds(ctx context.Context, tableName string, in interface{}, ids []uint, unscoped bool) error {
	// todo: delete one record
	if unscoped {
		if err := d.db.WithContext(ctx).Table(tableName).Where("id in ?", ids).Unscoped().Delete(in).Error; err != nil {
			return err
		}
	} else {
		if err := d.db.WithContext(ctx).Table(tableName).Where("id in ?", ids).Delete(in).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *defaultDao) DeleteWhereAny(ctx context.Context, tableName string, in interface{}, unscoped bool, query interface{}, args ...interface{}) error {
	// todo: delete one record
	if unscoped {
		if err := d.db.WithContext(ctx).Table(tableName).Where(query, args...).Unscoped().Delete(in).Error; err != nil {
			return err
		}
	} else {
		if err := d.db.WithContext(ctx).Table(tableName).Where(query, args...).Delete(in).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *defaultDao) First(ctx context.Context, tableName string, in interface{}, query interface{}, args ...interface{}) error {
	// todo: query one record
	if err := d.db.WithContext(ctx).Table(tableName).Where(query, args...).First(in).Error; err != nil {
		return err
	}
	return nil
}

func (d *defaultDao) Latest(ctx context.Context, tableName string, in interface{}, query interface{}, args ...interface{}) error {
	// todo: query one record
	if err := d.db.WithContext(ctx).Table(tableName).Order("id DESC").Where(query, args...).First(in).Error; err != nil {
		return err
	}
	return nil
}

func (d *defaultDao) Find(ctx context.Context, tableName string, in interface{}, query interface{}, args ...interface{}) error {
	// todo: query many records
	if err := d.db.Table(tableName).WithContext(ctx).Where(query, args...).Find(in).Error; err != nil {
		return err
	}
	return nil
}

func (d *defaultDao) FindAndOrderByInterface(ctx context.Context, tableName string, orderBy string, in interface{}, query interface{}, args ...interface{}) error {
	// todo: query many records
	if err := d.db.Table(tableName).WithContext(ctx).Order(orderBy).Where(query, args...).Find(in).Error; err != nil {
		return err
	}
	return nil
}

func (d *defaultDao) Update(ctx context.Context, tableName string, id uint, updates interface{}) error {
	// todo: update one record
	if err := d.db.Table(tableName).WithContext(ctx).Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func (d *defaultDao) Count(ctx context.Context, tableName string, query interface{}, args ...interface{}) (int64, error) {
	// todo: query many records
	var total int64
	if err := d.db.Table(tableName).WithContext(ctx).Where(query, args...).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (d *defaultDao) FindAndLimit(ctx context.Context, tableName string, limit, offset int, in interface{}, query interface{}, args ...interface{}) error {
	// todo: query many records
	if err := d.db.Table(tableName).WithContext(ctx).Where(query, args...).Limit(limit).Offset(offset).Find(in).Error; err != nil {
		return err
	}
	return nil
}

func (d *defaultDao) FindAndLimitOrder(ctx context.Context, tableName, orderBy string, limit, offset int, in interface{}, query interface{}, args ...interface{}) error {
	// todo: query many records
	if err := d.db.Table(tableName).WithContext(ctx).Where(query, args...).Order(orderBy).Limit(limit).Offset(offset).Find(in).Error; err != nil {
		return err
	}
	return nil
}

func (d *defaultDao) FindAndLimitAndSortInterface(ctx context.Context, tableName string, in interface{}, orderBy string, offset, limit int32, query interface{}, args ...interface{}) error {
	// todo: query many records
	if err := d.db.Table(tableName).WithContext(ctx).Offset(int(offset)).Limit(int(limit)).Order(orderBy).Where(query, args...).Find(in).Error; err != nil {
		return err
	}
	return nil
}

func (d *defaultDao) DeleteWhereQuery(ctx context.Context, tableName string, in interface{}, unscoped bool, query interface{}, args ...interface{}) error {
	// todo: delete one record
	if unscoped {
		if err := d.db.WithContext(ctx).Table(tableName).Where(query, args...).Unscoped().Delete(in).Error; err != nil {
			return err
		}
	} else {
		if err := d.db.WithContext(ctx).Table(tableName).Where(query, args...).Delete(in).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *defaultDao) UpdateWhereQuery(ctx context.Context, tableName string, updates interface{}, query interface{}, args ...interface{}) error {
	// todo: update one record
	if err := d.db.Table(tableName).WithContext(ctx).Where(query, args...).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}
