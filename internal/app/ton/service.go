package ton

import (
	"context"
	"fmt"
	"github.com/a2tonium/a2tonium-backend/pkg/logger"
	"github.com/a2tonium/a2tonium-backend/pkg/ton/crypto"
	"github.com/a2tonium/a2tonium-backend/pkg/ton/edu"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"strings"
)

type TonService struct {
	api     *ton.APIClient
	keyPair crypto.KeyPair
	courses []*course
	wallet  *wallet.Wallet
}

func NewTonService() *TonService {
	return &TonService{}
}

// Init - will retrieve all created courses along with their enrolled students and assign them accordingly
func (t *TonService) Init(ctx context.Context, mnemonic string) error {
	client := liteclient.NewConnectionPool()

	configUrl := "https://ton-blockchain.github.io/testnet-global.config.json"
	err := client.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	if err != nil {
		panic(err)
	}
	t.api = ton.NewAPIClient(client)

	seed := strings.Split(mnemonic, " ")
	w, err := wallet.FromSeed(t.api, seed, wallet.ConfigV5R1Final{
		NetworkGlobalID: wallet.TestnetGlobalID,
	})
	if err != nil {
		return fmt.Errorf("error creating wallet: %v", err)
	}

	keyPair, err := crypto.MnemonicToX25519KeyPair(mnemonic)
	if err != nil {
		return fmt.Errorf("error generating X25519 key pair: %v", err)
	}

	t.api, t.keyPair, t.wallet = t.api, keyPair, w

	courses, err := t.getAllCreatedCourses(ctx)
	if err != nil {
		return err
	}
	t.courses = courses

	for _, course := range courses {
		students, err := course.getStudents(ctx, t.api)
		if err != nil {
			return err
		}

		course.Students = students
	}

	return nil
}

func (t *TonService) Reload(ctx context.Context) error {
	courses, err := t.getAllCreatedCourses(ctx)
	if err != nil {
		return err
	}
	t.courses = courses

	for _, course := range courses {
		students, err := course.getStudents(ctx, t.api)
		if err != nil {
			return err
		}

		course.Students = students
	}

	return nil
}

func (t *TonService) UpdateStudents(ctx context.Context) error {
	for _, course := range t.courses {
		students, err := course.getStudents(ctx, t.api)
		if err != nil {
			return err
		}

		course.Students = students
	}

	return nil
}

// GetCoursesURIs - will return the content URIs for all courses
func (t *TonService) GetCoursesURI() []string {
	uris := make([]string, 0, len(t.courses))
	for _, c := range t.courses {
		uris = append(uris, c.Content.URI)
	}

	return uris
}

func (t *TonService) GetAllStudentsGmail() []string {
	uris := make([]string, 0, len(t.courses))
	for _, c := range t.courses {
		for _, s := range c.Students {
			uris = append(uris, s.Gmail)
		}
	}

	return uris
}

// SetCoursesMetadata updates the metadata for all courses managed by the TonService.
// It expects a slice of CourseMetadata with the same length as the number of courses.
// Returns an error if the lengths don't match or if updating metadata for any course fails.
func (t *TonService) SetCoursesMetadata(ctx context.Context, coursesMetadata []*CourseMetadata) error {
	if len(coursesMetadata) != len(t.courses) {
		return fmt.Errorf("number of courses does not match the number of coursesMetadata")
	}

	for i, course := range t.courses {
		err := course.setMetadata(coursesMetadata[i], t.keyPair.PrivateKey)
		if err != nil {
			logger.ErrorKV(ctx, "courses metadata error", "courseAddress", course.CourseAddress,
				logger.Err, err)
			return fmt.Errorf("set metadata failed: %v", err)
		}
	}

	return nil
}

// TODO: function which will operate all courses and will return the certificateIssue
func (t *TonService) ProcessAllCourses(ctx context.Context) ([]*CertificateIssue, error) {
	certificateIssueData := make([]*CertificateIssue, 0)
	for i, course := range t.courses {
		certificateIssueD, err := course.Process(ctx, t.api, t.wallet, t.keyPair.PrivateKey)
		if err != nil {
			logger.ErrorKV(ctx, "tonService.ProcessAllCourses course.Process failed", logger.Err, err)
			return certificateIssueData, err
		}

		for _, certificateIssue := range certificateIssueD {
			certificateIssue.CourseIndex = i
			certificateIssue.ImageLink = course.courseCompletion[0].Certificate
			certificateIssueData = append(certificateIssueData, certificateIssue)
		}
	}

	return certificateIssueData, nil
}

func (t *TonService) CertificateIssue(ctx context.Context, courseIndex, studentIndex int, cid string) error {
	var (
		course  = t.courses[courseIndex]
		student = course.Students[studentIndex]
	)

	mintData, err := edu.BuildCertificateIssuePayload(student.CertificateAddress,
		&edu.ContentOffchain{
			URI: fmt.Sprintf("ipfs://%s", cid),
		})
	if err != nil {
		logger.ErrorKV(ctx, "tonService.CertificateIssue edu.BuildCertificateIssuePayload failed", logger.Err, err)
		return fmt.Errorf("build certifcate issue payload failed: %w", err)
	}

	mint := wallet.SimpleMessage(course.CourseAddress, tlb.MustFromTON("0.02"), mintData)

	for i := 0; i < 3; i++ {
		_, _, err = t.wallet.SendWaitTransaction(ctx, mint)
		if err == nil {
			break
		}
	}
	if err != nil {
		logger.ErrorKV(ctx, "tonService.CertificateIssue t.wallet.SendWaitTransaction failed", logger.Err, err)
		return fmt.Errorf("sendWaitTransaction failed: %w", err)
	}

	course.Students = append(course.Students[:studentIndex], course.Students[studentIndex+1:]...)
	logger.Info(ctx, "Certificate successfully issue")
	return nil
}
