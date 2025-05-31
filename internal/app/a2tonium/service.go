package a2tonium

import (
	"github.com/a2tonium/a2tonium-backend/internal/app/ton"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
)

type YoutubeService interface {
	GiveAccessToGmail(gmails, videosId []string) error
}

type TonService interface {
	GetAllCoursesByAddress(address *address.Address) ([]*edu.Course, error)
	GetAllStudentsByCourseAddress(courseAddress *address.Address) ([]*edu.Course, error)
	GetEmitsTillTx(certificateAddress *address.Address) ([]*tlb.Transaction, lastTxLT int, error)
}

type A2Tonium struct {
	YoutubeService
	TonService
}

func (a *A2Tonium) Run() {
	// TODO: (1) Get Educator Address from Seed

	// TODO: (2) Get All Courses Which He Created
	// Course should contain the `StudentNum` and `CorrectAnswers`

	// TODO: (3) Get All `Enrolled`  Students By Course Address
	// Actually we will traverse all students

	// TODO: (4) Init state for each `Enrolled Student`
	// Create the Student struct
	// if nobody yet answerQuiz and no grade      {quizId = 0, LT = lastTx}
	// if we grade already                        {quizId = <gradedQuizId> + 1, LT = <ourGradeCommentTx>}
	// if we don't grade but the answerQuiz exist {quizId = 0, LT = <firstEmitTx>}

	// TODO: (5) Loop range `Enrolled Students`
	// Just call ProcessStudent()
	// If lastTx changed goTillLastTxn and from it parse is any emit exist and after update lastProcessedTxn

		var studentsGmail []string
	for _, c := range courses {
		for i := range c.StudentNum {
			student, isComplete, err :=	a.TonService.InitStudent()
			if !isComplete {
				studentsGmail := append(studentsGmail, student.GetGmail())
			}
			a.YoutubeService.GiveAccessToGmail(studentsGmail, c.GetVideosz)
			}
}

	for _, c := range courses {
		check is number of students in course changed
		a.TonService.UpdateCourse

		for _, s := students {
			a.TonService.ProcessStudent() {

			}
		}
	}


	// TODO: (6)
	// TODO: (7)

}
