package v1

type CourseInfo struct {
	AcademicYear int    `json:"academicYear"`
	Semester     int    `json:"semester"`
	Name         string `json:"name"`
	Teacher      string `json:"teacher"`
	Classroom    string `json:"classroom"`
	Weekday      int    `json:"weekday"`
	Weeks        string `json:"weeks"`
	Times        string `json:"times"`
}
