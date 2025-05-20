package edu

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/xssnick/tonutils-go/ton"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type TonApi interface {
	WaitForBlock(seqno uint32) ton.APIClientWrapped
	CurrentMasterchainInfo(ctx context.Context) (_ *ton.BlockIDExt, err error)
	RunGetMethod(ctx context.Context, blockInfo *ton.BlockIDExt, addr *address.Address, method string, params ...any) (*ton.ExecutionResult, error)
}

type CertificateIssuePayload struct {
	_                  tlb.Magic        `tlb:"#bc3b6585"`
	CertificateAddress *address.Address `tlb:"addr"`
	CertificateContent *cell.Cell       `tlb:"^"`
}

type CourseData struct {
	CourseIndex   *big.Int
	NextItemIndex *big.Int
	Content       *ContentOffchain
	OwnerAddress  *address.Address
	Cost          *big.Int
}

// CourseClient is a wrapper around a TON smart contract that acts as an NFT collection.
// Each NFT represents a non-transferable course certificate issued to a student.
type CourseClient struct {
	addr *address.Address
	api  TonApi
}

func NewCourseClient(api TonApi, courseAddr *address.Address) *CourseClient {
	return &CourseClient{
		addr: courseAddr,
		api:  api,
	}
}

func (c *CourseClient) GetCourseAddress() *address.Address {
	return c.addr
}

// GetCertificateAddressByIndex retrieves the address of the certificate NFT at a given index.
func (c *CourseClient) GetCertificateAddressByIndex(ctx context.Context, index *big.Int) (*address.Address, error) {
	b, err := c.api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get masterchain info: %w", err)
	}
	return c.GetCertificateAddressByIndexAtBlock(ctx, index, b)
}

func (c *CourseClient) GetCertificateAddressByIndexAtBlock(ctx context.Context, index *big.Int, b *ton.BlockIDExt) (*address.Address, error) {
	res, err := c.api.WaitForBlock(b.SeqNo).RunGetMethod(ctx, b, c.addr, methodGetNFTAddress, index)
	if err != nil {
		return nil, fmt.Errorf("failed to run %s method for index %v: %w", methodGetNFTAddress, index, err)
	}

	x, err := res.Slice(0)
	if err != nil {
		return nil, fmt.Errorf("result get err: %w", err)
	}

	addr, err := x.LoadAddr()
	if err != nil {
		return nil, fmt.Errorf("failed to load address from result slice: %w", err)
	}

	return addr, nil
}

func (c *CourseClient) GetCertificateContent(ctx context.Context, index *big.Int, individualCertificateContent *ContentOffchain) (*ContentOffchain, error) {
	b, err := c.api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get masterchain info: %w", err)
	}
	return c.GetCertificateContentAtBlock(ctx, index, individualCertificateContent, b)
}

func (c *CourseClient) GetCertificateContentAtBlock(ctx context.Context, index *big.Int, individualCertificateContent *ContentOffchain, b *ton.BlockIDExt) (*ContentOffchain, error) {
	con, err := toCertificateContent(individualCertificateContent)
	if err != nil {
		return nil, fmt.Errorf("failed to convert certificate content to cell: %w", err)
	}

	res, err := c.api.WaitForBlock(b.SeqNo).RunGetMethod(ctx, b, c.addr, methodGetNFTContent, index, con)
	if err != nil {
		return nil, fmt.Errorf("failed to run %s method for index %v: %w", methodGetNFTContent, index, err)
	}

	x, err := res.Cell(0)
	if err != nil {
		return nil, fmt.Errorf("result get err: %w", err)
	}

	cnt, err := ContentFromCell(x)
	if err != nil {
		return nil, fmt.Errorf("failed to parse content: %w", err)
	}

	return cnt, nil
}

func (c *CourseClient) GetCourseData(ctx context.Context) (*CourseData, error) {
	b, err := c.api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get masterchain info: %w", err)
	}
	return c.GetCourseDataAtBlock(ctx, b)
}

func (c *CourseClient) GetCourseDataAtBlock(ctx context.Context, b *ton.BlockIDExt) (*CourseData, error) {
	res, err := c.api.WaitForBlock(b.SeqNo).RunGetMethod(ctx, b, c.addr, methodGetCourseData)
	if err != nil {
		return nil, fmt.Errorf("failed to run %s method: %w", methodGetCourseData, err)
	}

	courseIndex, err := res.Int(0)
	if err != nil {
		return nil, fmt.Errorf("course index get err: %w", err)
	}

	nextIndex, err := res.Int(1)
	if err != nil {
		return nil, fmt.Errorf("next index get err: %w", err)
	}

	content, err := res.Cell(2)
	if err != nil {
		return nil, fmt.Errorf("content get err: %w", err)
	}

	ownerRes, err := res.Slice(3)
	if err != nil {
		return nil, fmt.Errorf("owner get err: %w", err)
	}

	cost, err := res.Int(4)
	if err != nil {
		return nil, fmt.Errorf("cost get err: %w", err)
	}

	addr, err := ownerRes.LoadAddr()
	if err != nil {
		return nil, fmt.Errorf("failed to load owner address from result slice: %w", err)
	}

	cnt, err := ContentFromCell(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse content: %w", err)
	}

	return &CourseData{
		CourseIndex:   courseIndex,
		NextItemIndex: nextIndex,
		Content:       cnt,
		OwnerAddress:  addr,
		Cost:          cost,
	}, nil
}

func (c *CourseClient) BuildCertificateIssuePayload(certificateAddress *address.Address, certificateContent *ContentOffchain) (_ *cell.Cell, err error) {
	con, err := toCertificateContent(certificateContent)
	if err != nil {
		return nil, fmt.Errorf("failed to convert certificate content to cell: %w", err)
	}

	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}

	body, err := tlb.ToCell(CertificateIssuePayload{
		CertificateAddress: certificateAddress,
		CertificateContent: con,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to convert ItemMintPayload to cell: %w", err)
	}

	return body, nil
}

func toCertificateContent(content *ContentOffchain) (*cell.Cell, error) {
	if content == nil {
		return nil, errors.New("content can't be nil") // TODO: add the errors
	}
	return cell.BeginCell().MustStoreStringSnake(content.URI).EndCell(), nil
}
