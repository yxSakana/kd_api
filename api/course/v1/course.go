package v1

import "github.com/gogf/gf/v2/frame/g"

type GetWeekCoursesReq struct {
	g.Meta       `path:"/course/{academic_year}/{semester}/{week}" method:"get" sm:"获取指定学年的学期的周次的课程" tags:"course"`
	AcademicYear uint `json:"academic_year" v:"required" dc:"学年,现实年数-学期"`
	Semester     uint `json:"semester" v:"required" dc:"学期,从0开始"`
	Week         uint `json:"week" v:"required" dc:"周次"`
}
type GetWeekCoursesRes struct {
	Items []CourseInfo `json:"items" dc:"课程"`
	Count uint         `json:"count"`
}

type GetCurrentSemesterCoursesReq struct {
	g.Meta `path:"/course/{week}" method:"get" sm:"获取当前学年、学前的指定周次的课程" tags:"course"`
	Week   uint `json:"week" v:"required"`
}
type GetCurrentSemesterCoursesRes struct {
	Items []CourseInfo `json:"items"`
	Count uint         `json:"count"`
}
