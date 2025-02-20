// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// User is the golang structure for table user.
type User struct {
	StudentId   uint   `json:"studentId"   orm:"student_id"   description:""` //
	Password    string `json:"password"    orm:"password"     description:""` //
	StudentCode uint64 `json:"studentCode" orm:"student_code" description:""` //
	Name        string `json:"name"        orm:"name"         description:""` //
	SchoolCode  uint   `json:"schoolCode"  orm:"school_code"  description:""` //
	SchoolName  string `json:"schoolName"  orm:"school_name"  description:""` //
}
