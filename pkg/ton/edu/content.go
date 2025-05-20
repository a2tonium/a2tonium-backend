package edu

import (
	"fmt"

	"github.com/xssnick/tonutils-go/tvm/cell"
)

type ContentOffchain struct {
	URI string
}

func ContentFromCell(c *cell.Cell) (*ContentOffchain, error) {
	return ContentFromSlice(c.BeginParse())
}

func ContentFromSlice(s *cell.Slice) (*ContentOffchain, error) {
	if s.BitsLeft() < 8 {
		if s.RefsNum() == 0 {
			return &ContentOffchain{}, nil
		}
		s = s.MustLoadRef()
	}

	typ, err := s.LoadUInt(8)
	if err != nil {
		return nil, fmt.Errorf("failed to load type: %w", err)
	}
	t := uint8(typ)

	switch t {
	case 0x01:
		str, err := s.LoadStringSnake()
		if err != nil {
			return nil, fmt.Errorf("failed to load snake offchain data: %w", err)
		}

		return &ContentOffchain{
			URI: str,
		}, nil
	default:
		// If you only expect offchain content, you can reject unknown types here
		return nil, fmt.Errorf("unknown content type: %d", t)
	}
}

func (c *ContentOffchain) ContentCell() (*cell.Cell, error) {
	return cell.BeginCell().MustStoreUInt(0x01, 8).MustStoreStringSnake(c.URI).EndCell(), nil
}
