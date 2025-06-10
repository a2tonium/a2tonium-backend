package edu

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/xssnick/tonutils-go/ton"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

const (
	_CourseV1CodeHex = "b5ee9c724102290100097c000114ff00f4a413f4bcf2c80b01020162021103dad0eda2edfb01d072d721d200d200fa4021103450666f04f86102f862ed44d0d200019dd30fd31fd4fa40fa0055406c158eabfa40810101d7005902d10170822009184e72a000f8425240c705f2e0848200bf6023c2fff2f48810344033e206925f06e024d749c21fe30004f90127030e04d004d31f2182101c3fe32aba8fd2365b03d4fa005932f8425240c705f2e0848200bf5d22820a625a00bcf2f4811947882201f90001f900bdf2f48813154440f84201706ddb3cc87f01ca0055405045cb0f12cb1fcc01cf1601fa02c9ed54db31e021821021ce0b32ba27040507003e00000000436f757273652075706461746564207375636365737366756c6c7901f06d6d226eb3995b206ef2d0806f22019132e2f8416f24135f03f8276f1001a18209312d00b98e478209312d0070fb02102470030481008250231036552212c8cf8580ca00cf8440ce01fa028069cf40025c6e016eb0935bcf819d58cf8680cf8480f400f400cf81e2f400c901fb00e010247003048042502306006a1036552212c8cf8580ca00cf8440ce01fa028069cf40025c6e016eb0935bcf819d58cf8680cf8480f400f400cf81e2f400c901fb0003d28f5b31d40131f8416f2430328200ca2a5328bef2f42410560447135086db3c5c705920f90022f9005ad76501d76582020134c8cb17cb0fcb0fcbffcbff71f90400c87401cb0212ca07cbffc9d08209312d00707370f842410d6d016d6de0018210bc3b6585bae3020417080c02fcc8555082105fcc3d145007cb1f15cb3f5003206e95307001cb0192cf16e201206e95307001cb0192cf16e2cc216eb3977f01ca0001fa02947032ca00e2216eb3957f01ca00cc947032ca00e2c910364540103b591036453304c8cf8580ca00cf8440ce01fa028069cf40025c6e016eb0935bcf818ae2f400c901fb0002a4090a001a58cf8680cf8480f400f400cf8101bc5154a1717088103410231029036d50536d03c8cf8580ca00cf8440ce01fa028069cf40025c6e016eb0935bcf819d58cf8680cf8480f400f400cf81e2f400c901fb004404c87f01ca0055405045cb0f12cb1fcc01cf1601fa02c9ed54db310b006000000000596f7520617265207375636365737366756c6c7920656e726f6c6c656420696e2074686520636f757273652101fefa40d45932f8425260c705f2e0848209312d0070737150056d6d586d6dc8555082105fcc3d145007cb1f15cb3f5003206e95307001cb0192cf16e201206e95307001cb0192cf16e2cc216eb3977f01ca0001fa02947032ca00e2216eb3957f01ca00cc947032ca00e2c910344130146d50436d5033c8cf8580ca00cf8440ce0d008401fa028069cf40025c6e016eb0935bcf819d58cf8680cf8480f400f400cf81e2f400c901fb004034c87f01ca0055405045cb0f12cb1fcc01cf1601fa02c9ed54db31015482f0677e8cc293f605b4e5b43b49d699b421f79482f729bfebab40fe728d7a2d3d91bae3025f05f2c0820f01d4f8425240c705f2e084820a625a0070fb02708306708827553010246d50436d03c8cf8580ca00cf8440ce01fa028069cf40025c6e016eb0935bcf819d58cf8680cf8480f400f400cf81e2f400c901fb004034c87f01ca0055405045cb0f12cb1fcc01cf1601fa02c9ed5410004a000000005769746864726177616c20636f6d706c65746564207375636365737366756c6c79020120122302012013150295b8b5ded44d0d200019dd30fd31fd4fa40fa0055406c158eabfa40810101d7005902d10170822009184e72a000f8425240c705f2e0848200bf6023c2fff2f48810344033e25514db3c6c51827140002310295ba7a3ed44d0d200019dd30fd31fd4fa40fa0055406c158eabfa40810101d7005902d10170822009184e72a000f8425240c705f2e0848200bf6023c2fff2f48810344033e25504db3c6c5182716015edb3c705920f90022f9005ad76501d76582020134c8cb17cb0fcb0fcbffcbff71f90400c87401cb0212ca07cbffc9d0170126f8280188c87001ca005a59cf16810101cf00c9180114ff00f4a413f4bcf2c80b190201621a2102f6d001d072d721d200d200fa4021103450666f04f86102f862ed44d0d200018e27fa40810101d700d20020d70b01c30093fa40019472d7216de201d2000191d4926d01e255406c158e1afa40810101d7005902d1016d6d8147c6f8425250c705f2f47059e206925f06e004d70d1ff2e0822182102fcb26a2bae302211b1c01d431d33f0131f8416f2410235f037080407f543467c8552082108b7717355004cb1f12cb3f810101cf0001cf16c91034413010246d50436d03c8cf8580ca00cf8440ce01fa028069cf40025c6e016eb0935bcf819d58cf8680cf8480f400f400cf81e2f400c901fb0040342003e082105fcc3d14ba8f57313504d33f20d70b01c30093fa40019472d7216de20120d70b01c30093fa40019472d7216de201d4d2000192fa00926d01e2d2000191d4926d01e2555010235f0333812155f8425260c705f2f421c000923134e30e4430e001821098756485bae3025f06f2c0821d201f01a630c0018ecd347f8209312d0070fb0223206ef2d080708306708810246d50436d03c8cf8580ca00cf8440ce01fa028069cf40025c6e016eb0935bcf819d58cf8680cf8480f400f400cf81e2f400c901fb0004de1e008000000000436f6e67726174756c6174696f6e732120596f752068617665207375636365737366756c6c7920636f6d706c657465642074686520636f757273652101a4d307d45932f8425260216e925b7092c705e2f2e084812e94f8416f24135f038209312d00bef2f401c8598210987564855003cb1fcb07ccc9c88258c000000000000000000000000101cb67ccc970fb004034200064c87f01ca0055405054cf1612810101cf00ca0058206e95307001cb0192cf16e2216eb3957f01ca00cc947032ca00e2c9ed5401a5a11f9fda89a1a400031c4ff481020203ae01a40041ae1603860127f4800328e5ae42dbc403a4000323a924da03c4aa80d82b1c35f481020203ae00b205a202dadb028f8df084a4a18e0be5e8e0b3c5b678d8ab22002421206ef2d08021206ef2d08024544630285902014824260291b60b7da89a1a400033ba61fa63fa9f481f400aa80d82b1d57f481020203ae00b205a202e1044012309ce54001f084a4818e0be5c10904017ec04785ffe5e91020688067c5b678d8a70272500065473210291b7ae3da89a1a400033ba61fa63fa9f481f400aa80d82b1d57f481020203ae00b205a202e1044012309ce54001f084a4818e0be5c10904017ec04785ffe5e91020688067c5b678d8ab027280000000a5474325343ce408431"
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

func NewCourseClientFromInit(api TonApi, ownerAddr *address.Address, courseIndex int64) (*CourseClient, error) {
	stateCell, err := tlb.ToCell(&tlb.StateInit{
		Data: getCourseData(ownerAddr, courseIndex),
		Code: getCourseCode(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get state: %w", err)
	}
	addr := address.NewAddress(0, byte(int8(0)), stateCell.Hash())

	return NewCourseClient(api, addr), nil
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

func BuildCertificateIssuePayload(certificateAddress *address.Address, certificateContent *ContentOffchain) (_ *cell.Cell, err error) {
	content, err := certificateContent.ContentCell()
	if err != nil {
		return nil, fmt.Errorf("failed to convert certificate content to cell: %w", err)
	}

	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}

	body, err := tlb.ToCell(CertificateIssuePayload{
		CertificateAddress: certificateAddress,
		CertificateContent: content,
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
