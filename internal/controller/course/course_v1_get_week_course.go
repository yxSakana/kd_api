package course

import (
	"context"

	"kd_api/api/course/v1"
	"kd_api/internal/logic/course"
)

func (c *ControllerV1) GetWeekCourses(
	ctx context.Context,
	req *v1.GetWeekCoursesReq,
) (res *v1.GetWeekCoursesRes, err error) {
	courses, count, err := course.GetCoursesByWeek(ctx, req.AcademicYear, req.Semester, req.Week)
	var list []v1.CourseInfo
	for _, item := range courses {
		list = append(list, v1.CourseInfo{
			AcademicYear: item.AcademicYear,
			Semester:     item.Semester,
			Name:         item.Name,
			Teacher:      item.Teacher,
			Classroom:    item.Classroom,
			Weekday:      item.Weekday,
			Weeks:        item.Weeks,
		})
	}
	return &v1.GetWeekCoursesRes{
		Items: list,
		Count: count,
	}, nil
}
