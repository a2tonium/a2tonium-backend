package a2tonium

import (
	"context"
	"fmt"
	"github.com/a2tonium/a2tonium-backend/internal/app/ipfs"
	jsonGenerator "github.com/a2tonium/a2tonium-backend/internal/app/json_generator"
	"github.com/a2tonium/a2tonium-backend/internal/app/ton"
	"github.com/a2tonium/a2tonium-backend/pkg/logger"
	"time"
)

type TonService interface {
	Reload(ctx context.Context) error
	GetCoursesURI() []string
	GetAllStudentsGmail() []string
	SetCoursesMetadata(ctx context.Context, coursesMetadata []*ton.CourseMetadata) error
	ProcessAllCourses(ctx context.Context) ([]*ton.CertificateIssue, error)
	CertificateIssue(ctx context.Context, courseIndex, studentIndex int, cid string) error
}

type IpfsService interface {
	UploadJSONToPinata(ctx context.Context, jsonData string) (string, error)
	FetchCourseMetadata(ctx context.Context, ipfsURI string) (*ipfs.CourseMetadata, error)
}

type JsonGeneratorService interface {
	GenerateCertificateJSON(ctx context.Context, cert *jsonGenerator.Certificate) (string, error)
}

type A2Tonium struct {
	tonService           TonService
	ipfsService          IpfsService
	jsonGeneratorService JsonGeneratorService
}

func NewA2Tonium(tonService TonService, ipfsService IpfsService, jsonGeneratorService JsonGeneratorService) *A2Tonium {
	return &A2Tonium{
		tonService,
		ipfsService,
		jsonGeneratorService,
	}
}

func (a *A2Tonium) Init(ctx context.Context) error {
	var coursesMetadata []*ton.CourseMetadata
	coursesUri := a.tonService.GetCoursesURI()
	for _, courseUri := range coursesUri {
		metadata, err := a.ipfsService.FetchCourseMetadata(ctx, courseUri)
		if err != nil {
			logger.ErrorKV(ctx, logger.Err, err)
			return fmt.Errorf("ipfsService.FetchCourseMetadta for %s course failed: %w", err)
		}

		var courseCompletion []ton.CourseCompletion
		for _, metaCourseComp := range metadata.CourseCompletion {
			courseCompletion = append(courseCompletion, ton.CourseCompletion{
				GradeHighThan: metaCourseComp.GradeHighThan,
				Certificate:   metaCourseComp.Certificate,
			})
		}

		coursesMetadata = append(coursesMetadata, &ton.CourseMetadata{
			Name: metadata.Name,
			QuizAnswers: ton.QuizAnswers{
				EncryptedAnswers: metadata.QuizAnswers.EncryptedAnswers,
				SenderPublicKey:  metadata.QuizAnswers.SenderPublicKey,
			},
			CourseCompletion: courseCompletion,
		})
	}

	if err := a.tonService.SetCoursesMetadata(ctx, coursesMetadata); err != nil {
		logger.ErrorKV(ctx, logger.Err, err)
		return fmt.Errorf("tonService.SetCoursesMetadata failed: %w", err)
	}

	return nil
}

func (a *A2Tonium) Run(ctx context.Context) error {
	for {
		time.Sleep(5 * time.Second)
		certificatesIssueData, err := a.tonService.ProcessAllCourses(ctx)
		if err != nil {
			logger.ErrorKV(ctx, logger.Err, err)
		}
		if len(certificatesIssueData) == 0 {
			continue
		}

		for _, certificateData := range certificatesIssueData {
			certificateJson, err := a.jsonGeneratorService.GenerateCertificateJSON(ctx, &jsonGenerator.Certificate{
				Name:  certificateData.CourseName,
				IIN:   certificateData.StudentIIN,
				Image: certificateData.ImageLink,
				Attributes: []jsonGenerator.Attribute{
					{TraitType: "Student IIN", Value: certificateData.StudentIIN},
					{TraitType: "Average Grade", Value: certificateData.AverageGrade},
					{TraitType: "Completion Data", Value: certificateData.CompletionDate},
				},
				QuizGrades: certificateData.QuizGrades,
			})
			if err != nil {
				continue
			}

			cid, err := a.ipfsService.UploadJSONToPinata(ctx, certificateJson)
			if err != nil {
				continue
			}

			a.tonService.CertificateIssue(ctx, certificateData.CourseIndex, certificateData.StudentIndex, cid)
		}
	}
}
