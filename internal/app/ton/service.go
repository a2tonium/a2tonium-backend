package ton

import (
	"context"
	"fmt"
	"github.com/a2tonium/a2tonium-backend/pkg/ton/crypto"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"log"
	"strings"
)

type TonServiceInterface interface {
	GetAllCourses(ownerAddr *address.Address) []*Course
}

type TonService struct {
	api     *ton.APIClient
	keyPair crypto.KeyPair
	Courses []*Course
	wallet  *wallet.Wallet
}

func NewTonService(api *ton.APIClient, mnemonic string) *TonService {
	log.Println("Creating Wallet ...")
	seed := strings.Split(mnemonic, " ")
	w, err := wallet.FromSeed(api, seed, wallet.ConfigV5R1Final{
		NetworkGlobalID: wallet.TestnetGlobalID,
	})
	log.Println(w.WalletAddress())
	if err != nil {
		panic(err)
	}
	log.Println("X25519 Key Pair ...")
	keyPair, err := crypto.MnemonicToX25519KeyPair(mnemonic)
	if err != nil {
		panic(err)
	}

	return &TonService{
		api:     api,
		keyPair: keyPair,
		wallet:  w,
	}
}

func (t *TonService) Run(ctx context.Context) error {
	for i, course := range t.Courses {
		fmt.Println("Processing Course", i)
		err := course.Process(ctx, t.api, t.wallet, t.keyPair.PrivateKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TonService) Init(ctx context.Context) error {
	var courses []*Course
	courses, err := GetAllCreatedCourses(ctx, t.api, t.wallet.WalletAddress())
	if err != nil {
		return err
	}
	t.Courses = courses
	log.Println("Initializing TonService ...")
	for _, course := range courses {

		err := course.AssignQuizAnswersFromContent(ctx, t.keyPair.PrivateKey)
		if err != nil {
			return err
		}

		students, err := course.GetCurrentlyStudyingStudents(ctx, t.api)
		if err != nil {
			return err
		}
		course.StudyingStudents = students
		log.Println(course.StudentNum, course.Content, course)
	}

	return nil
}

func (t *TonService) Show() {
	fmt.Println("API:", t.api)
	fmt.Println("X255219 Private Key:", string(t.keyPair.PrivateKey))
	fmt.Println("X255219 Public Key:", string(t.keyPair.PublicKey))
	fmt.Println("Wallet Address:", t.wallet.WalletAddress())
	for _, course := range t.Courses {
		fmt.Println()
		fmt.Println("Courses:")
		fmt.Println("\tStudentNum:", course.StudentNum)
		fmt.Println("\tQuiz Correct Answers:", course.QuizCorrectAnswers)
		fmt.Println("\tContent:", course.Content)
		fmt.Println()
		fmt.Println("\tStudents:")
		for _, student := range course.StudyingStudents {
			fmt.Println("\t\tStudent Wallet Address:", student.CertificateAddress)
			fmt.Println("\t\tQuiz Id:", student.QuizId)
			fmt.Println("\t\tGmail:", student.Gmail)
			fmt.Println("\t\tIIN:", student.IIN)
		}
	}
}

//
//func NewTonService(keyPair crypto.KeyPair) *TonService {
//	return &TonService{keyPair: keyPair}
//}
//
//func (t *TonService) GetCourse(ctx context.Context, courseId uint16) (*Course, error)   {}
//func (t *TonService) GetTxnTillLT(ctx context.Context, address address.Address, lt int) {}
//func (t *TonService) GetGradeOrAllEmits(ctx context.Context, address address.Address)   {}
//func (t *TonService) GetAllGrades(ctx context.Context, address address.Address)         {}
//func (t *TonService) GetAllGradesByCourse(ctx context.Context, address address.Address) {}
//func (t *TonService) GetCertificateByCourse(ctx context.Context, address address.Address, course *Course) {
//}
//
//type TonServiceI interface {
//	// GetCourse
//	GetCourse(ctx context.Context, courseId uint16) (*Course, error)
//}
//
//func getWallet(api ton.APIClientWrapped) *wallet.Wallet {
//	words := strings.Split("birth pattern then forest walnut then phrase walnut fan pumpkin pattern then cluster blossom verify then forest velvet pond fiction pattern collect then then", " ")
//	w, err := wallet.FromSeed(api, words, wallet.V4R2)
//	if err != nil {
//		panic(err)
//	}
//	return w
//}
