// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package course

import (
	"context"

	"kd_api/api/course/v1"
)

type ICourseV1 interface {
	GetWeekCourses(ctx context.Context, req *v1.GetWeekCoursesReq) (res *v1.GetWeekCoursesRes, err error)
	GetCurrentSemesterCourses(ctx context.Context, req *v1.GetCurrentSemesterCoursesReq) (res *v1.GetCurrentSemesterCoursesRes, err error)
}
