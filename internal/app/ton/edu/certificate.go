package edu

import (
	"context"
	"encoding/hex"
	"github.com/a2tonium/a2tonium-backend/pkg/ton/edu"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func getCertificateCode() *cell.Cell {
	codeCellBytes, _ := hex.DecodeString(_CertificateV1CodeHex)
	codeCell, err := cell.FromBOC(codeCellBytes)
	if err != nil {
		panic(err)
	}

	return codeCell
}

func getCertificateData(collectionAddress *address.Address, certificateIndex int64) *cell.Cell {
	return cell.BeginCell().
		MustStoreUInt(0, 1).
		MustStoreAddr(collectionAddress).
		MustStoreInt(certificateIndex, 257).
		EndCell()
}
