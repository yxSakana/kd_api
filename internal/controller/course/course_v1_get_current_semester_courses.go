package course

import (
	"context"

	"kd_api/api/course/v1"
	"kd_api/internal/logic/course"
)

func (c *ControllerV1) GetCurrentSemesterCourses(
	ctx context.Context, req *v1.GetCurrentSemesterCoursesReq,
) (res *v1.GetCurrentSemesterCoursesRes, err error) {
	courses, count, err := course.GetCurrentSemesterCourses(ctx, req.Week)
	if err != nil {
		return nil, err
	}

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
			Times:        item.Times,
		})
	}
	return &v1.GetCurrentSemesterCoursesRes{
		Items: list,
		Count: count,
	}, nil
}
