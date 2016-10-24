package models

import (
	"errors"
	"fmt"
	_ "reflect"
	"strings"

	"bufio"
	"crypto/md5"
	"github.com/astaxie/beego/orm"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Post struct {
	Id       int64
	Title    string
	Author   string
	Date     time.Time `orm:"type(datetime)"`
	Summary  string
	Viewed   int
	Tags     []*Tag `orm:"rel(m2m)"`
	Body     string `orm:"type(text)"`
	Md5Value string `orm:"size(50)"`
}

type Tag struct {
	Id    int64
	Name  string
	Posts []*Post `orm:"reverse(many)"`
}

func getFilelist(path string) []string {
	var paths []string
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		println(path)
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}

	return paths
}

func ParseFile(filename string) *Post {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer f.Close()

	md5h := md5.New()
	io.Copy(md5h, f)
	f.Seek(0, 0)

	buf := bufio.NewReader(f)
	post := new(Post)
	checked_md5 := fmt.Sprintf("%x", string(md5h.Sum([]byte(""))))
	post.Md5Value = checked_md5
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err.Error())
			break
		}
		if strings.HasPrefix(line, "x") {
			var content = make([]byte, 1024)
			var body string
			for {
				n, err := buf.Read(content)
				if err != nil {
					if err == io.EOF {
						post.Body = body
						break
					}
					fmt.Println(err)
					break
				}
				body = body + string(content[:n])
			}
		}
		switch {
		case strings.HasPrefix(line, "@title"):
			v := strings.Split(line, ":")
			title := strings.TrimSpace(v[1])
			post.Title = title
		case strings.HasPrefix(line, "@tags"):
			v := strings.Split(line, ":")
			tags := strings.Split(strings.TrimSpace(v[1]), ",")
			var post_tags []*Tag
			for _, t := range tags {
				single_tag := &Tag{Name: strings.TrimSpace(t)}
				post_tags = append(post_tags, single_tag)
			}
			post.Tags = post_tags
		case strings.HasPrefix(line, "@summary"):
			v := strings.Split(line, ":")
			summary := strings.TrimSpace(v[1])
			post.Summary = summary
		case strings.HasPrefix(line, "@author"):
			v := strings.Split(line, ":")
			author := strings.TrimSpace(v[1])
			post.Author = author
		case strings.HasPrefix(line, "@date"):
			v := strings.Split(line, ":")
			date := strings.TrimSpace(v[1])
			//fmt.Println(date)
			tm2, err := time.Parse("2006-01-02 15.04.05", date)
			if err != nil {
				fmt.Println(err)
			}
			post.Date = tm2
		}
	}
	return post
}

func File2DB() {
	//filename := "posts/test.md"
	filenames := getFilelist("posts/")
	for _, filename := range filenames {
		if !strings.HasSuffix(filename, ".md") {
			continue
		}
		post := ParseFile(filename)
		if post == nil {
			fmt.Println("parse file " + filename + " error")
			return
		}

		var posts []*Post
		o := orm.NewOrm()
		qs := o.QueryTable("post")
		qs.Filter("title", post.Title).All(&posts)
		if len(posts) > 1 {
			fmt.Println("len(post) error")
		} else if len(posts) == 1 {
			post.Id = posts[0].Id
			if posts[0].Md5Value != post.Md5Value {
				post.Viewed = posts[0].Viewed
				o.Update(post)
			}
		} else {
			id, _ := o.Insert(post)
			post.Id = id
		}

		//fmt.Println(post)

		tags := post.Tags
		tagmap := make(map[string]*Tag)
		for _, t := range tags {
			tagmap[t.Name] = t
		}
		var tagname []interface{}
		for k, _ := range tagmap {
			tagname = append(tagname, k)
		}

		//fmt.Println(tagname)

		qs = o.QueryTable("tag")
		var query_tags []*Tag
		qs.Filter("name__in", tagname...).All(&query_tags)
		//fmt.Println(len(query_tags))

		var in_flag bool
		var insert_tags []*Tag
		for k, v := range tagmap {
			in_flag = false
			for _, v1 := range query_tags {
				if k == v1.Name {
					in_flag = true
					break
				}
			}
			if !in_flag {
				insert_tags = append(insert_tags, v)
			}
		}

		//fmt.Println(len(insert_tags))
		if len(insert_tags) != 0 {
			o.InsertMulti(len(insert_tags), insert_tags)

			qs = o.QueryTable("tag")
			qs.Filter("name__in", tagname...).All(&query_tags)
		}

		// fmt.Println(query_tags)

		new_post := &Post{Id: post.Id}
		m2m := o.QueryM2M(new_post, "Tags")

		// clear all and re-insert
		m2m.Clear()
		m2m.Add(query_tags)
	}
}

func init() {
}

// GetPostById retrieves Post by Id. Returns error if
// Id doesn't exist
func GetPostById(id string) (v *Post, err error) {
	var posts []*Post
	o := orm.NewOrm()
	qs := o.QueryTable(new(Post))
	qs.Filter("title", id).All(&posts)
	if len(posts) != 1 {
		fmt.Println(err)
		return nil, errors.New("too many posts")
	}

	new_post := posts[0]
	new_post.Viewed = new_post.Viewed + 1
	o.Update(new_post)

	post := *posts[0]
	o.LoadRelated(&post, "Tags")
	return &post, nil
}

// GetAllTag retrieves all tags
func GetAllTag() (tags []*Tag, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Tag))
	_, err = qs.All(&tags)
	if err != nil {
		return nil, err
	}
	return
}

// GetAllPost retrieves all Post matches certain condition. Returns empty list if
// no records exist
func GetAllPost(exclude map[string][]string, query map[string][]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Post, count int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Post))
	count, _ = qs.Count()

	// exclude k=[v1, v2]
	for k, v := range exclude {
		k = strings.Replace(k, ".", "__", -1)
		for _, tv := range v {
			qs = qs.Exclude(k, tv)
		}
	}

	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		k = k + "__in"
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, 0, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, 0, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, 0, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, 0, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Post
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		return l, count, nil
	}
	return nil, 0, err
}

// UpdatePost updates Post by Id and returns error if
// the record to be updated doesn't exist
func UpdatePostById(m *Post) (err error) {
	o := orm.NewOrm()
	v := Post{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func GetTagsByName(tagname string) (*Tag, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Tag))
	var tags []*Tag
	qs.Filter("Name", tagname).All(&tags)
	if len(tags) != 1 {
		return nil, errors.New("invalid number of tag")
	}

	tag := tags[0]

	o.LoadRelated(tag, "Posts")
	return tag, nil
}

// DeletePost deletes Post by Id and returns error if
// the record to be deleted doesn't exist
func DeletePost(id int64) (err error) {
	o := orm.NewOrm()
	v := Post{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Post{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
