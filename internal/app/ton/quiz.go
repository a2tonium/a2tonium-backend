package ton

import (
	"fmt"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func LoadQuizFromCell(c *cell.Cell) (int, string, error) {
	cs := c.BeginParse()

	// Step 1: Read and verify 32-bit magic/type prefix
	typeID, err := cs.LoadUInt(32)
	if err != nil {
		return 0, "", fmt.Errorf("failed to load type ID: %w", err)
	}
	if typeID != 0x98756485 {
		return 0, "", fmt.Errorf("unexpected type ID: got 0x%x, want 0x98756485", typeID)
	}

	// Step 2: Read quizId (uint8)
	quizId, err := cs.LoadUInt(8)
	if err != nil {
		return 0, "", fmt.Errorf("failed to load quizId: %w", err)
	}

	// Step 3: Read answers cell reference
	answersCell, err := cs.LoadRef()
	if err != nil {
		return 0, "", fmt.Errorf("failed to load answers cell: %w", err)
	}
	answers := answersCell.MustLoadStringSnake()

	return int(quizId), answers, nil
}
