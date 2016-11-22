package models

import (
	"time"
)

type BaseModel struct {
	Id string `json:"id" form:"id" sql:"type:varchar(100);index;not null;unique;primary key"`
	CreatedAt time.Time `json:"created_at,omitempty" sql:"index"`
	UpdatedAt time.Time `json:"created_at,omitempty" sql:"index"`
	DeletedAt *time.Time `json:"created_at,omitempty" sql:"index"`
}