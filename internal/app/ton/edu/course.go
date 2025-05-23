package edu

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	certificateJsonGenerator "github.com/a2tonium/a2tonium-backend/internal/app/certificate_json_generator"
	"github.com/a2tonium/a2tonium-backend/internal/app/ipfs"
	"github.com/a2tonium/a2tonium-backend/pkg/logger"
	"github.com/a2tonium/a2tonium-backend/pkg/ton/crypto"
	"github.com/a2tonium/a2tonium-backend/pkg/ton/edu"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
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

func (c *Course) Process(ctx context.Context, api *ton.APIClient, w *wallet.Wallet, privateKey []byte) error {
	for iteration := 0; iteration < 20; iteration++ {
		time.Sleep(5 * time.Second)
		for studentIndex := len(c.StudyingStudents) - 1; studentIndex >= 0; studentIndex-- {
			s := c.StudyingStudents[studentIndex]
			newEmits, err := s.getNewEmits(ctx, api, c.OwnerAddress)
			if err != nil {
				return err
			}
			fmt.Println("len(newEmits)", len(newEmits))
			for i := len(newEmits) - 1; i >= 0; i-- {
				newEmit := newEmits[i]
				s.lastProcessedHash = newEmit.txHash
				s.lastProcessedLt = newEmit.txLt
				//newEmit.payload.BeginParse().LoadUInt(32) ==
				quizId, answers, err := LoadQuizFromCell(newEmit.payload)
				if err != nil {
					return err
				}
				fmt.Println(quizId, s.QuizId, quizId != s.QuizId)
				if quizId != s.QuizId {
					continue
				}

				if len(c.QuizCorrectAnswers) == s.QuizId {
					c.StudyingStudents = append(c.StudyingStudents[:studentIndex], c.StudyingStudents[studentIndex+1:]...)
					continue
				}

				if len(c.QuizCorrectAnswers) == quizId {
					allGrades, err := s.GetAllGrades(ctx, api, c.OwnerAddress, quizId)
					if err != nil {
						continue
					}

					averageGrade, err := averagePercent(allGrades)
					if err != nil {
						return err
					}
					metadata, err := ipfs.FetchQuizAndCompletion(c.Content.URI)
					formattedDate := time.Now().Format("2006-01-02")

					certificateJson, err := certificateJsonGenerator.GenerateCertificateJSON(certificateJsonGenerator.Certificate{
						Name:        fmt.Sprintf("`%s` Course Certificate", metadata.Name),
						Description: fmt.Sprintf("Certificate of completion for the `%s` Course. Awarded to: `%s`. This NFT certifies successful completion of all modules.", metadata.Name, s.IIN),
						Image:       fmt.Sprintf("%s", metadata.CourseCompletion[0].Certificate),
						Attributes: []certificateJsonGenerator.Attribute{
							{TraitType: "Student IIN", Value: s.IIN},
							{TraitType: "Average Grade", Value: fmt.Sprintf("%.2f%%", averageGrade)},
							{TraitType: "Completion Date", Value: formattedDate},
						},
						QuizGrades: allGrades,
					})
					if err != nil {
						return err
					}
					cid, err := ipfs.UploadJSONToPinata(certificateJson, fmt.Sprintf("certificate_%d.json", time.Now().UnixNano()))
					if err != nil {
						fmt.Println(err)
						return err
					}

					mintData, err := c.courseClient.BuildCertificateIssuePayload(s.CertificateAddress,
						&edu.ContentOffchain{
							URI: fmt.Sprintf("ipfs://%s", cid),
						})
					if err != nil {
						panic(err)
					}

					mint := wallet.SimpleMessage(c.courseClient.GetCourseAddress(), tlb.MustFromTON("0.02"), mintData)

					rateReview := strings.Split(answers, " | ")
					review := fmt.Sprintf("Rating: %s Review: %s", rateReview[0], rateReview[1])
					transfer, err := w.BuildTransfer(s.CertificateAddress, tlb.MustFromTON("0"), false, review)
					if err != nil {
						log.Fatalln("Transfer err:", err.Error())
						return err
					}

					fmt.Println("Minting NFT && Review...")
					for {
						err = w.SendMany(ctx, []*wallet.Message{mint, transfer}, false)
						if err == nil {
							break
						}
						fmt.Println(err)
					}

					c.StudyingStudents = append(c.StudyingStudents[:studentIndex], c.StudyingStudents[studentIndex+1:]...)
					continue
				}

				encrypted := strings.Split(answers, " | ")
				answers, err = crypto.DecryptX25519AESCBCMessage(encrypted[0], encrypted[1], privateKey)
				if err != nil {
					fmt.Printf("%q\n", encrypted)
					fmt.Println(answers, err)
					continue
				}

				grade := compareStrings(answers, c.QuizCorrectAnswers[quizId])
				gradingString := fmt.Sprintf("%v | %s | %v | %s", quizId+1, grade, newEmit.txLt, base64.StdEncoding.EncodeToString(newEmit.txHash))
				fmt.Println(fmt.Sprintf("GradingString: %q", gradingString))
				transfer, err := w.BuildTransfer(s.CertificateAddress, tlb.MustFromTON("0"), false, gradingString)
				if err != nil {
					log.Fatalln("Transfer err:", err.Error())
					return err
				}

				ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Minute*2)
				tx, block, err := w.SendWaitTransaction(ctxWithTimeout, transfer)
				cancel()
				if err != nil {
					logger.ErrorKV(ctx, "SendWaitTransaction err:", logger.Err, err.Error())
					select {
					case <-ctxWithTimeout.Done():
						return fmt.Errorf("transaction timed out after %s: %w", "2min", err)
					default:
						return fmt.Errorf("transaction failed: %w", err)
					}
				} else {
					s.QuizId += 1
				}

				fmt.Println("tx:", tx)
				fmt.Println("block:", block)
			}
		}
	}
	return nil
}

func averagePercent(percentStrings []string) (float64, error) {
	var sum float64
	for _, p := range percentStrings {
		// Remove trailing '%' sign
		trimmed := strings.TrimSuffix(p, "%")
		// Parse to float64
		val, err := strconv.ParseFloat(trimmed, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse %q: %w", p, err)
		}
		sum += val
	}
	if len(percentStrings) == 0 {
		return 0, nil // avoid division by zero
	}
	return sum / float64(len(percentStrings)), nil
}
func compareStrings(s1, s2 string) string {
	// Determine the minimum length to avoid index out of range
	minLen := len(s1)
	if len(s2) < minLen {
		minLen = len(s2)
	}

	// Count matching characters at the same positions
	matchCount := 0
	for i := 0; i < minLen; i++ {
		if s1[i] == s2[i] {
			matchCount++
		}
	}

	// Use the length of the longer string as the denominator
	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
	}

	// Calculate the percentage
	percent := (float64(matchCount) / float64(maxLen)) * 100

	// Format with two decimal places
	return fmt.Sprintf("%.2f%%", percent)
}

func (c *Course) AssignQuizAnswersFromContent(ctx context.Context, recipientPrivateKey []byte) error {
	metadata, err := ipfs.FetchQuizAndCompletion(c.Content.URI)
	if err != nil {
		return err
	}

	log.Println("CorrectEncrypted:", metadata.QuizAnswers.EncryptedAnswers, "\nPublicKey:", metadata.QuizAnswers.SenderPublicKey)
	correctAnswersString, err := crypto.DecryptX25519AESCBCMessage(metadata.QuizAnswers.EncryptedAnswers, metadata.QuizAnswers.SenderPublicKey, recipientPrivateKey)
	if err != nil {
		return err
	}
	fmt.Println("CORRECTANSWERSTRING", correctAnswersString)

	correctAnswers := strings.Split(correctAnswersString, " ")
	log.Println("CORRECT ANSWERS:", correctAnswers)
	c.QuizCorrectAnswers = correctAnswers

	return nil
}

func (c *Course) GetCurrentlyStudyingStudents(ctx context.Context, api edu.TonApi) ([]*Student, error) {
	var students []*Student
	for i := int64(0); ; i++ {
		certificate := certificateFromInit(ctx, api, c.courseClient.GetCourseAddress(), i)
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

func courseFromInit(ctx context.Context, api edu.TonApi, ownerAddress *address.Address, courseId int64) *edu.CourseClient {
	code := getCourseCode()
	data := getCourseData(ownerAddress, courseId)
	state := &tlb.StateInit{
		Data: data,
		Code: code,
	}

	stateCell, err := tlb.ToCell(state)
	if err != nil {
		log.Fatal(err)
	}
	addr := address.NewAddress(0, byte(int8(0)), stateCell.Hash())

	return edu.NewCourseClient(api, addr)
}

func getCourseCode() *cell.Cell {
	codeCellBytes, _ := hex.DecodeString(_CourseV1CodeHex)

	codeCell, err := cell.FromBOC(codeCellBytes)
	if err != nil {
		panic(err)
	}

	return codeCell
}

func getCourseData(ownerAddress *address.Address, courseId int64) *cell.Cell {
	return cell.BeginCell().
		MustStoreUInt(0, 1).
		MustStoreAddr(ownerAddress).
		MustStoreInt(courseId, 257).
		EndCell()
}
