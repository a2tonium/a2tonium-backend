package ton

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/a2tonium/a2tonium-backend/pkg/logger"
	"github.com/a2tonium/a2tonium-backend/pkg/ton/crypto"
	"github.com/a2tonium/a2tonium-backend/pkg/ton/edu"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"strings"
	"time"
)

type course struct {
	OwnerAddress, CourseAddress *address.Address
	Content                     *edu.ContentOffchain
	// Students - Certificates where (is_initialized == false)
	Students []*student
	// QuizCorrectAnswers - correctAnswers for each quiz (index + 1 == quizId)
	QuizCorrectAnswers []string
	courseCompletion   []CourseCompletion
	Name               string
}

func (c *course) Process(ctx context.Context, api *ton.APIClient, w *wallet.Wallet, privateKey []byte) ([]*CertificateIssue, error) {
	certificateIssueData := make([]*CertificateIssue, 0)

	for studentIndex, s := range c.Students {
		newEmits, certificateIssueCall, err := s.getNewEmitsWithCertCallCheck(ctx, api, c.OwnerAddress)
		if err != nil {
			continue
		}

		for i := len(newEmits) - 1; i >= 0; i-- {
			newEmit := newEmits[i]

			quizId, answers, err := LoadQuizFromCell(newEmit.payload)
			if err != nil {
				s.lastProcessedHash = newEmit.txHash
				s.lastProcessedLt = newEmit.txLt
				continue
			}

			if quizId != s.QuizId {
				s.lastProcessedHash = newEmit.txHash
				s.lastProcessedLt = newEmit.txLt
				continue
			}

			encrypted := strings.Split(answers, " | ")
			answers, err = crypto.DecryptX25519AESCBCMessage(encrypted[0], encrypted[1], privateKey)
			if err != nil {
				s.lastProcessedHash = newEmit.txHash
				s.lastProcessedLt = newEmit.txLt
				continue
			}

			grade := compareStrings(answers, c.QuizCorrectAnswers[quizId])
			gradingString := fmt.Sprintf("%v | %s | %v | %s", quizId+1, grade, newEmit.txLt, base64.StdEncoding.EncodeToString(newEmit.txHash))

			transfer, err := w.BuildTransfer(s.CertificateAddress, tlb.MustFromTON("0"), false, gradingString)
			if err != nil {
				break
			}

			//tx, block, err := w.SendWaitTransaction(ctx, transfer)
			_, _, err = w.SendWaitTransaction(ctx, transfer)
			logger.InfoKV(ctx, "quiz graded", "student address", s.CertificateAddress, "quizId", quizId+1,
				"grade", grade)
			if err != nil {
				break
			}
			s.QuizId += 1
			s.lastProcessedHash, s.lastProcessedLt = newEmit.txHash, newEmit.txLt
		}

		if certificateIssueCall && s.QuizId == len(c.QuizCorrectAnswers) {
			allGrades, err := s.getAllGrades(ctx, api, c.OwnerAddress)
			if err != nil {
				continue
			}
			averageGrade, err := averagePercent(allGrades)
			if err != nil {
				continue
			}
			certificateIssueData = append(certificateIssueData, &CertificateIssue{
				CourseName:     c.Name,
				StudentIIN:     s.IIN,
				AverageGrade:   fmt.Sprintf("%.2f%%", averageGrade),
				QuizGrades:     allGrades,
				StudentIndex:   studentIndex,
				CompletionDate: time.Now().Format("2006-01-02"),
			})
			logger.InfoKV(ctx, "Certificate issue data prepared", "student address", s.CertificateAddress)
		}
	}

	return certificateIssueData, nil
}

func (c *course) setMetadata(metadata *CourseMetadata, recipientPrivateKey []byte) error {
	correctAnswersString, err := crypto.DecryptX25519AESCBCMessage(metadata.QuizAnswers.EncryptedAnswers, metadata.QuizAnswers.SenderPublicKey, recipientPrivateKey)
	if err != nil {
		return fmt.Errorf("correct answers decryp")
	}

	correctAnswers := strings.Split(correctAnswersString, " ")
	c.QuizCorrectAnswers, c.Name, c.courseCompletion = correctAnswers, metadata.Name, metadata.CourseCompletion

	return nil
}

func (c *course) getStudents(ctx context.Context, api edu.TonApi) ([]*student, error) {
	var students []*student
	for i := int64(0); ; i++ {
		certificate, err := edu.NewCertificateClientFromInit(api, c.CourseAddress, i)
		if err != nil {
			return nil, err
		}

		certificateData, err := certificate.GetCertificateData(ctx)
		if err != nil {
			// 7 is `type check error. An argument to a primitive is of incorrect value type`
			// -256 is `contract is not initialized`
			if errors.Is(err, ton.ContractExecError{7}) ||
				errors.Is(err, ton.ContractExecError{-256}) {
				break
			}

			return students, err
		}

		if certificateData != nil && !certificateData.Initialized {
			credentials := strings.Split(certificateData.Credentials, " | ")
			student := &student{
				CertificateAddress: certificate.GetCertificateAddress(),
				IIN:                credentials[0],
				Gmail:              credentials[1],
				ownerAddr:          certificateData.OwnerAddress,
			}
			students = append(students, student)
		}
	}

	return students, nil
}

func (t *TonService) getAllCreatedCourses(ctx context.Context) ([]*course, error) {
	var createdCourses []*course

	for i := int64(0); ; i++ {
		courseClient, err := edu.NewCourseClientFromInit(t.api, t.wallet.WalletAddress(), i)
		if err != nil {
			return nil, err
		}

		courseData, err := courseClient.GetCourseData(ctx)
		if err != nil {
			// 7 is `type check error. An argument to a primitive is of incorrect value type`
			// -256 is `contract is not initialized`
			if errors.Is(err, ton.ContractExecError{7}) ||
				errors.Is(err, ton.ContractExecError{-256}) {
				break
			}

			return createdCourses, err
		}

		course := &course{
			CourseAddress: courseClient.GetCourseAddress(),
			OwnerAddress:  courseData.OwnerAddress,
			Content:       courseData.Content,
		}

		createdCourses = append(createdCourses, course)
	}

	return createdCourses, nil
}
