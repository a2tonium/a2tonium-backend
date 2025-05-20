package edu

import (
	"context"
	"fmt"
	"github.com/xssnick/tonutils-go/ton"
	"math/big"

	"github.com/xssnick/tonutils-go/address"
)

type CertificateData struct {
	Initialized   bool
	Index         *big.Int
	CourseAddress *address.Address
	OwnerAddress  *address.Address
	Content       *ContentOffchain // after certificateIssue
	Credentials   string           // from enroll till certificateIssue
}

// CertificateClient is a wrapper around a TON smart contract that acts as an NFT Item.
// Each NFT represents a non-transferable course certificate issued to a student.
type CertificateClient struct {
	addr *address.Address
	api  TonApi
}

func NewCertificateClient(api TonApi, certificateAddr *address.Address) *CertificateClient {
	return &CertificateClient{
		addr: certificateAddr,
		api:  api,
	}
}

func (c *CertificateClient) GetCertificateAddress() *address.Address {
	return c.addr
}

func (c *CertificateClient) GetCertificateData(ctx context.Context) (*CertificateData, error) {
	b, err := c.api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get masterchain info: %w", err)
	}
	return c.GetCertificateDataAtBlock(ctx, b)
}

// GetCertificateDataAtBlock retrieves the data of the certificate.
func (c *CertificateClient) GetCertificateDataAtBlock(ctx context.Context, b *ton.BlockIDExt) (*CertificateData, error) {
	res, err := c.api.WaitForBlock(b.SeqNo).RunGetMethod(ctx, b, c.addr, methodGetNFTData)
	if err != nil {
		return nil, fmt.Errorf("failed to run %s method: %w", methodGetNFTData, err)
	}

	init, err := res.Int(0)
	if err != nil {
		return nil, fmt.Errorf("err get init value: %w", err)
	}

	index, err := res.Int(1)
	if err != nil {
		return nil, fmt.Errorf("err get index value: %w", err)
	}

	courseRes, err := res.Slice(2)
	if err != nil {
		return nil, fmt.Errorf("err get course slice value: %w", err)
	}

	courseAddr, err := courseRes.LoadAddr()
	if err != nil {
		return nil, fmt.Errorf("failed to load course address from result slice: %w", err)
	}

	var ownerAddr *address.Address

	nilOwner, err := res.IsNil(3)
	if err != nil {
		return nil, fmt.Errorf("err check for nil owner slice value: %w", err)
	}

	if !nilOwner {
		ownerRes, err := res.Slice(3)
		if err != nil {
			return nil, fmt.Errorf("err get owner slice value: %w", err)
		}

		ownerAddr, err = ownerRes.LoadAddr()
		if err != nil {
			return nil, fmt.Errorf("failed to load owner address from result slice: %w", err)
		}
	} else {
		ownerAddr = address.NewAddressNone()
	}

	var (
		cnt         *ContentOffchain
		credentials string
	)

	nilContent, err := res.IsNil(4)
	if err != nil {
		return nil, fmt.Errorf("err check for nil content cell value: %w", err)
	}

	if !nilContent {
		content, err := res.Cell(4)
		if err != nil {
			return nil, fmt.Errorf("err get content cell value: %w", err)
		}

		cnt, err = ContentFromCell(content)
		if err != nil {
			credentials, err = content.BeginParse().LoadStringSnake()
			if err != nil {
				return nil, fmt.Errorf("failed to parse content: %w", err)
			}
		}
	}

	certificateData := &CertificateData{
		Initialized:   init.Cmp(big.NewInt(0)) != 0,
		Index:         index,
		CourseAddress: courseAddr,
		OwnerAddress:  ownerAddr,
	}
	if credentials != "" {
		certificateData.Credentials = credentials
	} else {
		certificateData.Content = cnt
	}
	
	return certificateData, nil
}
