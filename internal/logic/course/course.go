package course

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"kd_api/internal/dao"
	"kd_api/internal/logic/user"
	"kd_api/internal/model/do"
	"kd_api/internal/model/entity"
	"kd_api/internal/spider/engine"
	"kd_api/utility"
)

func GetCurrentSemesterCourses(ctx context.Context, week uint) ([]entity.Course, uint, error) {
	now := time.Now()
	month := now.Month()
	var semester uint
	if month > 8 {
		semester = 0
	} else {
		semester = 1
	}
	academicYear := uint(now.Year()) - semester
	return GetCoursesByWeek(ctx, academicYear, semester, week)
}

func GetCoursesByWeek(
	ctx context.Context, academicYear, semester, week uint,
) ([]entity.Course, uint, error) {
	studentId, err := user.GetStudentID(ctx)
	if err != nil {
		return nil, 0, err
	}
	weekStr := strconv.Itoa(int(week))

	// TODO: From redis cache
	// From mysql db
	var courses []entity.Course
	err = dao.Course.Ctx(ctx).Where(g.Map{
		"student_id":    studentId,
		"academic_year": academicYear,
		"semester":      semester,
	}).Wheref("weeks LIKE '%%%s%%' OR weeks LIKE '%s,%%' OR weeks LIKE ',%%%s%%,' OR weeks LIKE ',%s%%'",
		weekStr, weekStr, weekStr, weekStr).Scan(&courses)
	if err != nil {
		return nil, 0, err
	}
	if len(courses) != 0 {
		return courses, uint(len(courses)), nil
	}
	// From dean office system
	//// 1) Remove old
	//_, err = dao.Course.Ctx(ctx).Where(g.Map{
	//	"student_id":    studentId,
	//	"academic_year": academicYear,
	//	"semester":      semester,
	//}).Delete()
	//if err != nil {
	//	return nil, err
	//}
	// 2) Retrieve
	items, err := FetchCoursesByCrawler(ctx, academicYear, semester)
	if err != nil {
		return nil, 0, err
	}
	// 3) Insert
	itoa := utility.IntListToStr
	doCourses := make([]do.Course, len(items))
	for i, item := range items {
		doCourses[i] = do.Course{
			StudentId:    studentId,
			AcademicYear: item.AcademicYear,
			Name:         item.Name,
			Teacher:      item.Teacher,
			Classroom:    item.Classroom,
			Weekday:      item.Weekday,
			Weeks:        strings.Join(itoa(item.Weeks), ","),
			Times:        strings.Join(itoa(item.Times), ","),
		}
	}
	_, err = dao.Course.Ctx(ctx).Data(doCourses).Insert()
	if err != nil {
		return nil, 0, err
	}

	var result []entity.Course
	for _, item := range items {
		if slices.Contains(item.Weeks, int(week)) {
			result = append(result, entity.Course{
				StudentId:    uint(studentId),
				AcademicYear: item.AcademicYear,
				Semester:     item.Semester,
				Name:         item.Name,
				Teacher:      item.Teacher,
				Classroom:    item.Classroom,
				Weekday:      item.Weekday,
				Weeks:        strings.Join(itoa(item.Weeks), ","),
				Times:        strings.Join(itoa(item.Times), ","),
			})
		}
	}
	return result, uint(len(result)), nil
}

func FetchCoursesByCrawler(ctx context.Context, academicYear, semester uint) ([]engine.ClassItem, error) {
	studentInfo, err := user.GetStudentInfo(ctx)
	if err != nil {
		return nil, err
	}
	studentId := gconv.String(studentInfo.StudentId)
	studentCode := gconv.String(studentInfo.StudentCode)

	crawler, ok := ctx.Value("spider").(*engine.KdClient)
	if !ok {
		return nil, gerror.New("Error from context to get spider")
	}

	// Query user info
	r, err := dao.User.Ctx(ctx).Where(do.User{
		StudentId: studentId,
	}).One()
	if err != nil {
		return nil, err
	}
	var userInfo entity.User
	err = r.Struct(&userInfo)
	if err != nil {
		return nil, err
	}
	pwByte, err := gbase64.DecodeString(userInfo.Password)
	if err != nil {
		return nil, err
	}
	userInfo.Password = string(pwByte)

	err = crawler.LoginSystem(studentId, userInfo.Password, 3, false)
	if err != nil {
		return nil, gerror.New("failed to login dean office")
	}

	items, err := crawler.GetClass(int(academicYear), int(semester), studentCode)
	if err != nil {
		return nil, gerror.New("failed to fetch class from dean office")
	}

	itoa := utility.IntListToStr
	for _, item := range items {
		_, err = dao.Course.Ctx(ctx).Data(do.Course{
			AcademicYear: item.AcademicYear,
			Semester:     item.Semester,
			StudentId:    studentId,
			Name:         item.Name,
			Teacher:      item.Teacher,
			Classroom:    item.Classroom,
			Weekday:      item.Weekday,
			Weeks:        strings.Join(itoa(item.Weeks), ","),
			Times:        strings.Join(itoa(item.Times), ","),
		}).Insert()
		if err != nil {
			return nil, err
		}
	}
	return items, nil
}
