package main

import (
	"apiProject/models"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

var key = "4740860db0d0fcd49666663ae5f97f0e"

func main() {
	dbConfig()
	requestForFood()

}
func requestForFood() {
	request := httplib.Get("http://apis.juhe.cn/cook/index")
	request.Param("key", key)
	request.Param("cid", "1")
	request.Param("pn", "0")
	request.Param("rn", "30")
	json := make(map[string]interface{})
	request.ToJSON(&json)
	println("")

}
func requestCategory() {
	request := httplib.Get("http://apis.juhe.cn/cook/category")
	request.Param("key", key)
	json := make(map[string][]map[string]interface{})
	request.ToJSON(&json)
	categories := json["result"]
	o := orm.NewOrm()
	var name interface{} = "8"
	fuck, _ := name.(int)
	println(fuck)
	for _, value := range categories {
		println(value)
		parentIdString, _ := value["parentId"].(string)
		bb := value["name"].(string)
		parentId, _ := strconv.Atoi(parentIdString)
		println(bb)
		category := models.Category{ParentId: parentId, Name: bb}
		o.Insert(&category)

		list := value["list"].([]interface{})
		for _, subCategory := range list {
			println(subCategory)
			subCategoryMap, _ := subCategory.(map[string]interface{})
			subIdString, _ := subCategoryMap["id"].(string)
			subId, _ := strconv.Atoi(subIdString)

			subPIdString, _ := subCategoryMap["parentId"].(string)
			subPId, _ := strconv.Atoi(subPIdString)
			subName := subCategoryMap["name"].(string)

			//
			subCategoryStruct := models.SubCategory{Id: subId, ParentId: subPId, Name: subName}
			_, er := o.Insert(&subCategoryStruct)
			if er != nil {
				println("1")

			}

		}

		println(list)

	}
}

func dbConfig() {
	/**默认自带,可以不写下面这句*/
	//orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:crfchina@tcp(localhost:3306)/food?charset=utf8", 30)
	orm.RegisterModel(new(models.Category))
	orm.RegisterModel(new(models.SubCategory))

	/**自动创建表,个人建议还是先建好表再处理逻辑*/
	//orm.RunSyncdb("default", false, true)
}
