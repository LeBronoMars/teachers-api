package models

type Barangay struct {
	BaseModel
	CityId string `json:"city_id" form:"city_id"`
	ZipCode string `json:"zip_code" form:"zip_code"`
	Name string `json:"name" form:"name"`
}
