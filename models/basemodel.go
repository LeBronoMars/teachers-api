package models

import (
	"time"
)

type BaseModel struct {
	Id string `json:"id" form:"id" sql:"type:varchar(100);index;not null;unique;primary key" binding:"required"`
	CreatedAt time.Time `json:"created_at" sql:"index"`
	UpdatedAt time.Time `json:"updated_at" sql:"index"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
	IsSynced bool `json:"is_synced" form:"is_synced"`
	CreatedBy string `json:"created_by"`
}

func (b *BaseModel) BeforeUpdate() (err error) {
	b.IsSynced = true
	return
}

func (b *BaseModel) BeforeCreate() (err error) {
	b.IsSynced = true
	return
}

func (b *BaseModel) AfterUpdate() (err error) {
	b.IsSynced = true
	return
}

func (b *BaseModel) AfterCreate() (err error) {
	b.IsSynced = true
	return
}
