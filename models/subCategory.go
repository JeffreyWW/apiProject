package models

type SubCategory struct {
	Id       int    `json:"id form:id"`
	ParentId int    `json:"parentId"`
	Name     string `json:"name"`
}
