package ton

// Structures which will be sent to a2tonium

// (1) Certificate Issue
type CertificateIssue struct {
	AverageGrade, CompletionDate, StudentIIN, CourseName, ImageLink string
	QuizGrades                                                      []string
	CourseIndex, StudentIndex                                       int
}

type QuizGrade struct {
}

type CourseMetadata struct {
	Name             string
	QuizAnswers      QuizAnswers
	CourseCompletion []CourseCompletion
}

type QuizAnswers struct {
	EncryptedAnswers string
	SenderPublicKey  string
}

type CourseCompletion struct {
	GradeHighThan string
	Certificate   string
}
