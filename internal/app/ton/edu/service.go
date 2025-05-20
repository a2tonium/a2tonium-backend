package edu

import (
	"context"
	"github.com/a2tonium/a2tonium-backend/pkg/logger"
	"github.com/a2tonium/a2tonium-backend/pkg/ton/edu"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
)

// Edu - this name is like `nft`, but shows that is `edu` :)
type Edu interface {
	CourseFromInit(ctx context.Context, init CourseInit) *edu.CourseClient
	CertificateFromInit(ctx context.Context, init CertificateInit) *edu.CertificateClient
}

type EduService struct {
	api edu.TonApi
}

func NewEduService(api edu.TonApi) *EduService {
	return &EduService{api: api}
}

func (e *EduService) CourseFromInit(ctx context.Context, init CourseInit) *edu.CourseClient {
	stateCell, err := tlb.ToCell(&tlb.StateInit{
		Data: getCourseData(init.OwnerAddress, init.CourseIndex),
		Code: getCourseCode(),
	})
	if err != nil {
		logger.ErrorKV(ctx, "edu.CourseFromInit.tlb.ToCell failed:", logger.Err, err)
	}
	addr := address.NewAddress(0, byte(int8(0)), stateCell.Hash())

	return edu.NewCourseClient(e.api, addr)
}

func (e *EduService) CertificateFromInit(ctx context.Context, init CertificateInit) *edu.CertificateClient {
	stateCell, err := tlb.ToCell(&tlb.StateInit{
		Data: getCertificateData(init.CollectionAddress, init.CertificateIndex),
		Code: getCertificateCode(),
	})
	if err != nil {
		logger.ErrorKV(ctx, "edu.CertificateFromInit.tlb.ToCell failed:", logger.Err, err)
	}
	addr := address.NewAddress(0, byte(int8(0)), stateCell.Hash())

	return edu.NewCertificateClient(e.api, addr)
}
