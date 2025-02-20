package engine

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"kd_api/internal/spider/goreq"
)

type KdClient struct {
	goreq.Session
}

func NewKdClient() *KdClient {
	return &KdClient{Session: *goreq.NewSession(true)}
}

func (c *KdClient) LoginSystem(username, password string, retryCount int, forceLogin bool) error {
	if retryCount > 3 {
		retryCount = 3
	}
	var err error
	for i := 0; i < retryCount; i++ {
		_, err = c.loginSystem(username, password, forceLogin)
		if err == nil {
			return nil
		}
	}
	return err
}

func (c *KdClient) GetStudentInfo() (*StudentInfo, error) {
	return c.studentInfoReq()
}

// GetClass the semester param is start from 0, studentCode is StudentInfo.StudentCode
// 如果是上半学期(2024.9): 2024,0; 如果是下半学期(2025.2): 2024,1
func (c *KdClient) GetClass(academicYear, semester int, studentCode string) ([]ClassItem, error) {
	a := strconv.Itoa(academicYear)
	s := strconv.Itoa(semester)
	params := encryptClassSchedule(a, s, studentCode)
	api := exportApi + "?params=" + params

	res, err := c.Get(api, defaultHeaders, nil)
	if err != nil {
		return nil, err
	}

	doc, err := GetDocFromRes(res)
	tables, err := ExpandTableFromHtml(doc)
	if err != nil {
		return nil, err
	}
	table := tables[1]
	classSchedule := ExpandClassScheduleInfo(table)

	for i := range classSchedule {
		classSchedule[i].AcademicYear = academicYear
		classSchedule[i].Semester = semester
	}

	return classSchedule, nil
}

func (c *KdClient) TeachingEvaluate() {
	err := c.fetchTeachingEvaluateTablesReq(50058, 1)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *KdClient) loginSystem(username, password string, forceLogin bool) (string, error) {
	// Get cookies from db
	if !forceLogin {
		err := cookiesPool.SyncToCookieJar(username, c.Jar, hbkjjwUrl)
		if err == nil {
			fmt.Println("load cookies from db")
			return "", nil
		} else {
			fmt.Println(err.Error())
		}
	}
	fmt.Println("cookies is not exits, login...")

	// login
	deskeyUrl, err := c.loginReq()
	if err != nil {
		return "", err
	}

	// fetchKeyReq
	deskeyUrl = indexUrl + deskeyUrl
	desKey, sessionId, err := c.fetchKeyReq(deskeyUrl)
	if err != nil {
		return "", err
	}

	// validate code
	code, err := c.getValidateCode()
	if err != nil {
		return "", err
	}
	//time.Sleep(1 * time.Second)

	// encryptLogonParams
	params, err := encryptLogonParams(username, password, sessionId, desKey, code)
	if err != nil {
		return "", err
	}

	// logon
	homeUrl, err := c.logonReq(params)
	if err != nil {
		return "", err
	}

	// home page
	err = c.entryIndex(homeUrl)
	if err != nil {
		return "", err
	}

	// Save cookies to redis
	u, err := url.Parse(hbkjjwUrl)
	err = cookiesPool.AddCookies(username, c.Jar.Cookies(u), 10*time.Minute)
	if err != nil {
		return "", err
	}
	return homeUrl, nil
}
