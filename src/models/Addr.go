package models

import "github.com/jinzhu/gorm"

type Addr struct {
	gorm.Model
	LongAddr string`gorm:"type:varchar(256);not null; index: la_idx"`
	ShortAddr string`gorm:"type:varchar(256);not null; index: sa_idx"`
}
