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
	//requestForFood()
	requestCategory()

}
func requestForFood(subId int, pn int) {
	println("抓取subId：" + strconv.Itoa(subId) + ",  pn:" + strconv.Itoa(pn))
	request := httplib.Get("http://apis.juhe.cn/cook/index")
	request.Param("key", key)
	subIdString := strconv.Itoa(subId)
	request.Param("cid", subIdString)
	pnString := strconv.Itoa(pn)
	request.Param("pn", pnString)
	request.Param("rn", "30")
	json := make(map[string]interface{})
	request.ToJSON(&json)
	result, _ := json["result"].(map[string]interface{})
	data, _ := result["data"].([]interface{})
	o := orm.NewOrm()
	for _, menu := range data {
		menuMap, _ := menu.(map[string]interface{})
		id, _ := strconv.Atoi(menuMap["id"].(string))
		title, _ := menuMap["title"].(string)
		tags, _ := menuMap["tags"].(string)
		imtro, _ := menuMap["imtro"].(string)
		ingredients, _ := menuMap["ingredients"].(string)
		burden, _ := menuMap["burden"].(string)
		albums, _ := menuMap["albums"].([]interface{})
		album, _ := albums[0].(string)
		steps, _ := menuMap["steps"].([]interface{})
		menuStruct := models.Menu{Id: id, Title: title, Tags: tags, Imtro: imtro, Ingredients: ingredients, Burden: burden, Album: album}
		_, erM := o.Insert(&menuStruct)
		if erM != nil {
			println(erM.Error())
		}

		for _, step := range steps {
			stepMap, _ := step.(map[string]interface{})
			img, _ := stepMap["img"].(string)
			step, _ := stepMap["step"].(string)
			stepStruct := models.Step{MenuId: id, Img: img, Step: step}
			_, er := o.Insert(&stepStruct)
			if er != nil {
				println(er.Error())
			}

		}
	}
	if len(data) == 30 {
		requestForFood(subId, pn+30)
	}
}
func requestCategory() {
	request := httplib.Get("http://apis.juhe.cn/cook/category")
	request.Param("key", key)
	json := make(map[string][]map[string]interface{})
	request.ToJSON(&json)
	categories := json["result"]
	o := orm.NewOrm()
	for _, value := range categories {
		parentIdString, _ := value["parentId"].(string)
		bb := value["name"].(string)
		parentId, _ := strconv.Atoi(parentIdString)
		category := models.Category{ParentId: parentId, Name: bb}
		o.Insert(&category)
		list := value["list"].([]interface{})
		for _, subCategory := range list {
			subCategoryMap, _ := subCategory.(map[string]interface{})
			subIdString, _ := subCategoryMap["id"].(string)
			subId, _ := strconv.Atoi(subIdString)
			subPIdString, _ := subCategoryMap["parentId"].(string)
			subPId, _ := strconv.Atoi(subPIdString)
			subName := subCategoryMap["name"].(string)
			subCategoryStruct := models.SubCategory{Id: subId, ParentId: subPId, Name: subName}
			o.Insert(&subCategoryStruct)
			requestForFood(subId, 0)
		}
	}
}

func dbConfig() {
	/**默认自带,可以不写下面这句*/
	//orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(localhost:3306)/food?charset=utf8", 30)
	orm.RegisterModel(new(models.Category))
	orm.RegisterModel(new(models.SubCategory))
	orm.RegisterModel(new(models.Menu))
	orm.RegisterModel(new(models.Step))

	/**自动创建表,个人建议还是先建好表再处理逻辑*/
	//orm.RunSyncdb("default", false, true)
}
