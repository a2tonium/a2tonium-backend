package edu

import (
	"encoding/hex"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

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
