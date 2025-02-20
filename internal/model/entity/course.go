// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Course is the golang structure for table course.
type Course struct {
	Id           int    `json:"id"           orm:"id"            description:""` //
	StudentId    uint   `json:"studentId"    orm:"student_id"    description:""` //
	AcademicYear int    `json:"academicYear" orm:"academic_year" description:""` //
	Semester     int    `json:"semester"     orm:"semester"      description:""` //
	Name         string `json:"name"         orm:"name"          description:""` //
	Teacher      string `json:"teacher"      orm:"teacher"       description:""` //
	Classroom    string `json:"classroom"    orm:"classroom"     description:""` //
	Weekday      int    `json:"weekday"      orm:"weekday"       description:""` //
	Weeks        string `json:"weeks"        orm:"weeks"         description:""` //
	Times        string `json:"times"        orm:"times"         description:""` //
}
