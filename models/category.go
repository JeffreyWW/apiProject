package models

type Category struct {
	ParentId int    `orm:"pk;auto" json:"parentId"`
	Name     string `json:"name"`
}
