package engine

import (
	"fmt"
	"os"
	"regexp"
)

type StudentInfo struct {
	SchoolCode  uint32 `json:"school_code"`
	SchoolName  string `json:"school_name"`
	StudentId   uint32 `json:"student_id"`   // 学号
	StudentCode uint64 `json:"student_code"` // 不知道含义，解密时用
	Name        string `json:"name"`
}

type ClassItem struct {
	AcademicYear int // 学年
	Semester     int // 学期 从0开始
	Name         string
	Teacher      string
	Weekday      int
	Weeks        []int
	Times        []int
	Classroom    string
}

func (ci ClassItem) String() string {
	return fmt.Sprintf("{\n"+
		"\tName: %s, Teacher: %s, Weekday: %d, Weeks: %v, Times: %v, Classroom: %s\n"+
		"}",
		ci.Name, ci.Teacher, ci.Weekday, ci.Weeks, ci.Times, ci.Classroom)
}

const (
	jsFilename = "kd_system/js/encrypt.js"

	indexUrl                      = "https://202.206.64.231"
	hbkjjwUrl                     = indexUrl + "/hbkjjw/"
	loginApi                      = "https://202.206.64.231/hbkjjw/cas/login.action"
	logonApi                      = "https://202.206.64.231/hbkjjw/cas/logon.action"
	validateCodeApi               = "https://202.206.64.231/hbkjjw/cas/genValidateCode"
	exportApi                     = "https://202.206.64.231/hbkjjw/student/wsxk.xskcb_excel10319.jsp"
	studentInfoApi                = "https://202.206.64.231/hbkjjw/custom/js/SetRootPath.jsp"
	teachingEvaluateTableApi      = "https://202.206.64.231/hbkjjw/taglib/DataTable.jsp" // ?tableId=50058&fre=1
	teachingEvaluateFormApi       = "https://202.206.64.231/hbkjjw/student/wspj_tjzbpj_wjdcb_pj.jsp"
	teachingEvaluateFormSubmitApi = "https://202.206.64.231/hbkjjw/jw/wspjZbpjWjdc/save.action"

	defaultUA = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
)

var (
	tmpDir = os.TempDir()

	desKeyReg      = regexp.MustCompile(`var _deskey = '(.*?)';`)
	studentCodeReg = regexp.MustCompile(`var G_USER_CODE = '(\d+?)';`)
	varMapReg      = regexp.MustCompile(`var(.*?)=(.*?);`)
	defaultHeaders = map[string]string{
		"User-Agent": defaultUA,
		"Referer":    "https://202.206.64.231/hbkjjw/frame/home/homepage.jsp",
	}
	cookiesPool = NewCookiesPool("localhost:6379", "", 0)
)
