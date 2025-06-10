package edu

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"math/big"

	"github.com/xssnick/tonutils-go/address"
)

const (
	_CertificateV1CodeHex = "b5ee9c7241020b010002fe000114ff00f4a413f4bcf2c80b01020162020902f6d001d072d721d200d200fa4021103450666f04f86102f862ed44d0d200018e27fa40810101d700d20020d70b01c30093fa40019472d7216de201d2000191d4926d01e255406c158e1afa40810101d7005902d1016d6d8147c6f8425250c705f2f47059e206925f06e004d70d1ff2e0822182102fcb26a2bae30221030401d431d33f0131f8416f2410235f037080407f543467c8552082108b7717355004cb1f12cb3f810101cf0001cf16c91034413010246d50436d03c8cf8580ca00cf8440ce01fa028069cf40025c6e016eb0935bcf819d58cf8680cf8480f400f400cf81e2f400c901fb0040340803e082105fcc3d14ba8f57313504d33f20d70b01c30093fa40019472d7216de20120d70b01c30093fa40019472d7216de201d4d2000192fa00926d01e2d2000191d4926d01e2555010235f0333812155f8425260c705f2f421c000923134e30e4430e001821098756485bae3025f06f2c08205080701a630c0018ecd347f8209312d0070fb0223206ef2d080708306708810246d50436d03c8cf8580ca00cf8440ce01fa028069cf40025c6e016eb0935bcf819d58cf8680cf8480f400f400cf81e2f400c901fb0004de06008000000000436f6e67726174756c6174696f6e732120596f752068617665207375636365737366756c6c7920636f6d706c657465642074686520636f757273652101a4d307d45932f8425260216e925b7092c705e2f2e084812e94f8416f24135f038209312d00bef2f401c8598210987564855003cb1fcb07ccc9c88258c000000000000000000000000101cb67ccc970fb004034080064c87f01ca0055405054cf1612810101cf00ca0058206e95307001cb0192cf16e2216eb3957f01ca00cc947032ca00e2c9ed5401a5a11f9fda89a1a400031c4ff481020203ae01a40041ae1603860127f4800328e5ae42dbc403a4000323a924da03c4aa80d82b1c35f481020203ae00b205a202dadb028f8df084a4a18e0be5e8e0b3c5b678d8ab0a002421206ef2d08021206ef2d0802454463028597f675a8e"
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

func NewCertificateClientFromInit(api TonApi, courseAddr *address.Address, certificateIndex int64) (*CertificateClient, error) {
	stateCell, err := tlb.ToCell(&tlb.StateInit{
		Data: getCertificateData(courseAddr, certificateIndex),
		Code: getCertificateCode(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get state: %w", err)
	}
	addr := address.NewAddress(0, byte(int8(0)), stateCell.Hash())

	return NewCertificateClient(api, addr), nil
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
