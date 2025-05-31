package ton

import "github.com/xssnick/tonutils-go/address"

// Structures which will be sent to a2tonium

// (1) Certificate Issue
type CertificateIssue struct {
	QuizGrades      []string
	CourseAddr      *address.Address
	CertificateAddr *address.Address
}

type QuizGrade struct {
}
