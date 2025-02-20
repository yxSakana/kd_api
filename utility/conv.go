package utility

import (
	"kd_api/internal/model/entity"
	"kd_api/internal/spider/engine"
	"strconv"
	"strings"
)

func IntListToStr(in []int) []string {
	out := make([]string, len(in))
	for i, item := range in {
		out[i] = strconv.Itoa(item)
	}
	return out
}

func CrawlerClassInfoToDb(item engine.ClassItem) entity.Course {
	return entity.Course{
		AcademicYear: item.AcademicYear,
		Semester:     item.Semester,
		Name:         item.Name,
		Teacher:      item.Teacher,
		Classroom:    item.Classroom,
		Weekday:      item.Weekday,
		Weeks:        strings.Join(IntListToStr(item.Weeks), ","),
		Times:        strings.Join(IntListToStr(item.Times), ","),
	}
}

func CrawlerClassInfoListToDb(items []engine.ClassItem) []entity.Course {
	courses := make([]entity.Course, len(items))
	for i, item := range items {
		courses[i] = CrawlerClassInfoToDb(item)
	}
	return courses
}
