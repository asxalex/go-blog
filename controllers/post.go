package controllers

import (
	"github.com/astaxie/beego"

	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"myblog/models"
	_ "os"
)

type PostController struct {
	BaseController
}

// URLMapping ...
func (c *PostController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetIndex", c.GetIndex)
}

// Get ...
// @Title Create
// @Description create Posts
// @Param	body		body 	models.Posts	true		"body for Posts content"
// @Success 201 {object} models.Posts
// @Failure 403 body is empty
// @router / [get]
func (c *PostController) GetIndex() {
	c.Ctx.Redirect(302, "/posts/")
}

// Get ...
// @Title Create
// @Description create Posts
// @Param	body		body 	models.Posts	true		"body for Posts content"
// @Success 201 {object} models.Posts
// @Failure 403 body is empty
// @router /tags/:tagname [get]
func (c *PostController) GetSpecifiedTags() {
	tagname := c.Ctx.Input.Param(":tagname")
	tag, err := models.GetTagsByName(tagname)
	fmt.Println("tag ============", tag)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	c.Data["Tag"] = tag

	c.TplName = "categorys/show_single.html"
	c.Layout = "layouts/application.html"
}

// Get ...
// @Title Create
// @Description create Posts
// @Param	body		body 	models.Posts	true		"body for Posts content"
// @Success 201 {object} models.Posts
// @Failure 403 body is empty
// @router /posts/ [get]
func (c *PostController) GetAll() {

	filter, err := c.getFilter()
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	var fields []string = filter["fields"].([]string)
	var sortby []string = filter["sortby"].([]string)

	var order []string = filter["order"].([]string)
	var query = filter["query"].(map[string][]string)
	var exclude = filter["exclude"].(map[string][]string)
	var limit int64 = filter["limit"].(int64)
	var offset int64 = filter["offset"].(int64)
	ml, count, _ := models.GetAllPost(exclude, query, fields, sortby, order, offset, limit)
	tags, err := models.GetAllTag()

	fmt.Println(ml)
	c.Data["Posts"] = ml
	c.Data["Tags"] = tags

	var tabs []int64
	tabs = append(tabs, 0)
	var i int64
	for i = 0; i*10 < count; i++ {
		tabs = append(tabs, (i)*10)
	}
	c.Data["Tabs"] = tabs
	c.Data["Current"] = offset/10 + 1

	c.TplName = "categorys/show.html"
	c.Layout = "layouts/application.html"
}

func markDowner(args ...interface{}) template.HTML {
	s := blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}

// length is the page number
// current is the current page number
// to be finished
func paginate(length, current int64) template.HTML {
	h := `<div id="pagenavi" class="noselect">`
	h = h + `<ul class="pagination pagination">`
	if current == 1 {

	}
	return template.HTML(h)
}

// GetOne ...
// @Title Get
// @Description create Posts
// @Param	body		body 	models.Posts	true		"body for Posts content"
// @Success 201 {object} models.Posts
// @Failure 403 body is empty
// @router /posts/:id [get]
func (c *PostController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")

	fmt.Println("id str =", idStr)
	post, err := models.GetPostById(idStr)
	if err != nil {
		c.Data["json"] = err.Error()
		fmt.Println("err 1")
		c.ServeJSON()
		return
	}

	//tags, err = models.GetAllTag()
	if err != nil {
		c.Data["json"] = err.Error()
		fmt.Println("err 2")
		c.ServeJSON()
		return
	}

	beego.AddFuncMap("markDown", markDowner)
	c.Data["post"] = post
	c.TplName = "articles/show.html"
	c.Layout = "layouts/application.html"
}
