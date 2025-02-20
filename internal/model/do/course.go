// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Course is the golang structure of table course for DAO operations like Where/Data.
type Course struct {
	g.Meta       `orm:"table:course, do:true"`
	Id           interface{} //
	StudentId    interface{} //
	AcademicYear interface{} //
	Semester     interface{} //
	Name         interface{} //
	Teacher      interface{} //
	Classroom    interface{} //
	Weekday      interface{} //
	Weeks        interface{} //
	Times        interface{} //
}
