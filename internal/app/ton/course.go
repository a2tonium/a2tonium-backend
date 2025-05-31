package ton

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	certificateJsonGenerator "github.com/a2tonium/a2tonium-backend/internal/app/certificate_json_generator"
	"github.com/a2tonium/a2tonium-backend/internal/app/ipfs"
	"github.com/a2tonium/a2tonium-backend/pkg/ton/crypto"
	"github.com/a2tonium/a2tonium-backend/pkg/ton/edu"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type Course struct {
	OwnerAddress *address.Address
	Content      *edu.ContentOffchain
	courseClient *edu.CourseClient
	// StudyingStudents - Certificates where (is_initialized == false)
	StudyingStudents []*Student
	// quizCorrectAnswers - correctAnswers for each quiz (index + 1 == quizId)
	QuizCorrectAnswers []string
	// studentNum - All students number (NFT Collection next_id value)
	StudentNum *big.Int
}

func (c *Course) Process(ctx context.Context, api *ton.APIClient, w *wallet.Wallet) error {
	for _, s := range c.StudyingStudents {

		time.Sleep(5 * time.Second)
		newEmits, err := s.getNewEmits(ctx, api, c.OwnerAddress)
		if err != nil {
			return err
		}

		for i := len(newEmits) - 1; i >= 0; i-- {
			newEmit := newEmits[i]
			s.lastProcessedHash = newEmit.txHash
			s.lastProcessedLt = newEmit.txLt
			quizId, answers, err := LoadQuizFromCell(newEmit.payload)
			if err != nil {
				return err
			}

			if quizId != s.QuizId {
				continue
			}

			if len(c.QuizCorrectAnswers) == quizId {
				allGrades, err := s.GetAllGrades(ctx, api, c.OwnerAddress, quizId)
				if err != nil {
					return err
				}

				averageGrade, err := averagePercent(allGrades)
				if err != nil {
					return err
				}
				metadata, err := ipfs.FetchCourseMetadata(ctx, c.Content.URI)

				certificateJson, err := certificateJsonGenerator.GenerateCertificateJSON(certificateJsonGenerator.Certificate{
					Name:        fmt.Sprintf("%s Course Certificate", metadata.Name),
					Description: fmt.Sprintf("Certificate of completion for the %s Course. Awarded to: %s. This NFT certifies successful completion of all modules.", metadata.Name, s.IIN),
					Image:       fmt.Sprintf("ipfs://%s", metadata.CourseCompletion[0].Certificate),
					Attributes: []certificateJsonGenerator.Attribute{
						{TraitType: "	Student IIN", Value: s.IIN},
						{TraitType: "Average Grade", Value: fmt.Sprintf("%.2f%%", averageGrade)},
					},
					QuizGrades: allGrades,
				})
				if err != nil {
					return err
				}
				cid, err := ipfs.UploadJSONToPinata(certificateJson, fmt.Sprintf("certificate_%d.json", time.Now().UnixNano()))
				if err != nil {
					return err
				}

				mintData, err := c.courseClient.BuildCertificateIssuePayload(s.CertificateAddress,
					&edu.ContentOffchain{
						URI: fmt.Sprintf("ipfs://%s", cid),
					})
				if err != nil {
					panic(err)
				}

				fmt.Println("Minting NFT...")
				mint := wallet.SimpleMessage(c.courseClient.GetCourseAddress(), tlb.MustFromTON("0.02"), mintData)

				tx, block, err := w.SendWaitTransaction(context.Background(), mint)
				if err != nil {
					log.Fatalln("SendWaitTransaction err:", err.Error())
					return err
				}

				fmt.Println("tx:", tx)
				fmt.Println("block:", block)

				// TODO: make the same like in main branch
				students, err := c.GetCurrentlyStudyingStudents(ctx, api)
				if err != nil {
					return err
				}
				c.StudyingStudents = students

				continue
			}

			grade := compareStrings(answers, c.QuizCorrectAnswers[quizId])
			gradingString := fmt.Sprintf("%v | %s | %v | %s", quizId+1, grade, newEmit.txLt, base64.StdEncoding.EncodeToString(newEmit.txHash))
			fmt.Println(fmt.Sprintf("GradingString: %q", gradingString))
			transfer, err := w.BuildTransfer(s.CertificateAddress, tlb.MustFromTON("0.0001"), false, gradingString)
			if err != nil {
				log.Fatalln("Transfer err:", err.Error())
				return err
			}

			tx, block, err := w.SendWaitTransaction(ctx, transfer)
			if err != nil {
				log.Fatalln("SendWaitTransaction err:", err.Error())
				return err
			} else {
				s.QuizId += 1
			}

			fmt.Println("tx:", tx)
			fmt.Println("block:", block)
		}
	}

	return nil
}

func (c *Course) AssignQuizAnswersFromContent(ctx context.Context, recipientPrivateKey []byte) error {
	metadata, err := ipfs.FetchCourseMetadata(ctx, c.Content.URI)
	if err != nil {
		return err
	}

	log.Println("CorrectEncrypted:", metadata.QuizAnswers.EncryptedAnswers, "\nPublicKey:", metadata.QuizAnswers.SenderPublicKey)
	correctAnswersString, err := crypto.DecryptX25519AESCBCMessage(metadata.QuizAnswers.EncryptedAnswers, metadata.QuizAnswers.SenderPublicKey, recipientPrivateKey)
	if err != nil {
		return err
	}

	correctAnswers := strings.Split(correctAnswersString, " ")
	log.Println("CORRECT ANSWERS:", correctAnswers)
	c.QuizCorrectAnswers = correctAnswers

	return nil
}

func (c *Course) GetCurrentlyStudyingStudents(ctx context.Context, api edu.TonApi) ([]*Student, error) {
	var students []*Student
	for i := int64(0); ; i++ {
		certificate := edu.CertificateFromInit(ctx, api, c.courseClient.GetCourseAddress(), i)
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

		if certificateData != nil {
			if certificateData.Initialized != true {
				credentials := strings.Split(certificateData.Credentials, " | ")
				student := &Student{
					CertificateAddress: certificate.GetCertificateAddress(),
					IIN:                credentials[0],
					Gmail:              credentials[1],
				}
				students = append(students, student)
			}
		}
	}

	return students, nil
}

func GetAllCreatedCourses(ctx context.Context, api edu.TonApi, ownerAddress *address.Address) ([]*Course, error) {
	var createdCourses []*Course

	for i := int64(0); ; i++ {
		courseClient := courseFromInit(ctx, api, ownerAddress, i)
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
		if err != nil {
			panic(err)
		}

		course := &Course{
			OwnerAddress: courseData.OwnerAddress,
			courseClient: courseClient,
			StudentNum:   courseData.NextItemIndex,
			Content:      courseData.Content,
		}

		createdCourses = append(createdCourses, course)
	}

	return createdCourses, nil
}
