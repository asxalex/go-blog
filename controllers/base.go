package controllers

import (
	"errors"
	_ "fmt"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/orm"
	_ "net/url"
	"strings"

	_ "blog/models"
)

func init() {
}

type BaseController struct {
	beego.Controller
}

//
func (c *BaseController) getFilter() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	var fields []string
	var sortby []string

	var order []string
	var query = make(map[string][]string)
	var exclude = make(map[string][]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				err := errors.New("Error: invalid query key/value pair")
				//c.Data["json"] = err
				// c.ServeJSON()
				return result, err
			}
			k, v := kv[0], kv[1]
			query[k] = strings.Split(v, " ")
		}
	}

	// exclude: k:v,k:v
	if v := c.GetString("exclude"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				err := errors.New("Error: invalid exclude key/value pair")
				//c.Data["json"] = err
				// c.ServeJSON()
				return result, err
			}
			k, v := kv[0], kv[1]
			exclude[k] = strings.Split(v, " ")
		}
	}

	result["exclude"] = exclude
	result["query"] = query
	result["fields"] = fields
	result["sortby"] = sortby
	result["order"] = order
	result["offset"] = offset
	result["limit"] = limit

	return result, nil
}

func (b *BaseController) Prepare() {
}
