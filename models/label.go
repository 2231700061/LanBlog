package models

import (
	"github.com/astaxie/beego/orm"
	"strings"
	"reflect"
	"fmt"
	"errors"
)

/**
type Category struct {
	Id     int
	Name   string	`json:"cate"`
	Topics []*Topic `json:"-" orm:"reverse(many)"`
}

type Topic struct {
	Id       int
	Title    string
	Category *Category `orm:"rel(fk)"`
	Labels   []*Label  `orm:"rel(m2m)"`
}

type Label struct {
	Id     int
	Name   string	`json:label`
	Topics []*Topic `json:"-" orm:"reverse(many)"`
}

 */

type Label struct {
	Id          int    `json:"id"  orm:"column(id);auto"`
	Name        string `json:"name"  orm:"column(name);size(100)"`
	Description string `json:"desc"  orm:"column(description);size(500)"`
	Topics []*Topic		`json:"-" orm:"reverse(many)"`
}

func (t *Label) TableName() string {
	return "label"
}

//func init() {
//	orm.RegisterModel(new(Label))
//}

// AddLabel insert a new Label into database and returns
// last inserted Id on success.
func AddLabel(m *Label) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetLabelById retrieves Label by Id. Returns error if
// Id doesn't exist
func GetLabelById(id int) (v *Label, err error) {
	o := orm.NewOrm()
	v = &Label{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllLabel retrieves all Label matches certain condition. Returns empty list if
// no records exist
func GetAllLabel(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Label))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
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
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
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
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Label
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateLabel updates Label by Id and returns error if
// the record to be updated doesn't exist
func UpdateLabelById(m *Label) (err error) {
	o := orm.NewOrm()
	v := Label{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteLabel deletes Label by Id and returns error if
// the record to be deleted doesn't exist
func DeleteLabel(id int) (err error) {
	o := orm.NewOrm()
	v := Label{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Label{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
