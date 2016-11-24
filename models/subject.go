package models

type Subject struct {
	BaseModel
	Desciption string `json:"desciption" form:"desciption"`
	IsSynced bool `json:"is_synched" form:"is_synched"`
}