package edu

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/a2tonium/a2tonium-backend/pkg/ton/edu"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"strconv"
	"strings"
)

type Student struct {
	// QuizId - Last passed quiz id + 1
	QuizId             int
	Gmail, IIN         string
	CertificateAddress *address.Address

	lastProcessedLt   uint64
	lastProcessedHash []byte
}

type Emit struct {
	payload *cell.Cell
	txLt    uint64
	txHash  []byte
}

func (s *Student) GetAllGrades(ctx context.Context, api *ton.APIClient, courseOwner *address.Address, lastQuizId int) ([]string, error) {
	var account *tlb.Account
	for {
		block, err := api.CurrentMasterchainInfo(ctx)
		if err != nil {
			continue
		}

		// Get account state to find last transaction LT and hash
		account, err = api.GetAccount(ctx, block, s.CertificateAddress)
		if err == nil {
			break
		}
	}

	grades := make([]string, lastQuizId)
	lastLt := account.LastTxLT
	lastHash := account.LastTxHash

	// Fetch transactions in batches (limit is total transactions to fetch)
	batchSize := uint32(50)
	for lastQuizId > 0 {
		txs, err := api.ListTransactions(ctx, s.CertificateAddress, batchSize, lastLt, lastHash)
		fmt.Println("listTransactions", txs)
		if err != nil {
			fmt.Println("CHHHEHC", err)
			break
		}

		for i := len(txs) - 1; i >= 0; i-- {
			tx := txs[i]

			if tx.IO.In != nil {
				if tx.IO.In.Msg.SenderAddr().String() == courseOwner.String() {
					slice := tx.IO.In.Msg.Payload().BeginParse()
					slice.LoadUInt(32)
					payload := slice.MustLoadStringSnake()
					payloadSlice := strings.Split(payload, " | ")
					if len(payloadSlice) != 4 {
						continue
					}

					quizId, err := strconv.Atoi(payloadSlice[0])
					if err != nil {
						continue
					}

					if quizId == lastQuizId {
						grades[quizId-1] = payloadSlice[1]
						lastQuizId--
					}
				}
			}
		}

		// Prepare for next batch: use oldest tx LT/hash from current batch
		lastTx := txs[0]
		lastLt = lastTx.PrevTxLT
		lastHash = lastTx.PrevTxHash
	}

	fmt.Println(grades, lastQuizId)
	if lastQuizId == 0 {
		return grades, nil
	}
	return grades, errors.New("not all grades was found")
}

// Collect all emitted external-out messages (emits) for a contract address
func (s *Student) getNewEmits(ctx context.Context, api *ton.APIClient, courseOwner *address.Address) ([]*Emit, error) {
	// Get current masterchain info for block reference
	var account *tlb.Account
	for {
		block, err := api.CurrentMasterchainInfo(ctx)
		if err != nil {
			continue
		}

		// Get account state to find last transaction LT and hash
		account, err = api.GetAccount(ctx, block, s.CertificateAddress)
		if err == nil {
			break
		}
	}

	var emits []*Emit
	lastLt := account.LastTxLT
	lastHash := account.LastTxHash
	fmt.Printf("Account: %s", s.CertificateAddress)

	lastProcessedHash := s.lastProcessedHash
	lastProcessedLt := s.lastProcessedLt

	if lastLt == lastProcessedLt && bytes.Compare(lastHash, lastProcessedHash) == 0 {
		return emits, nil
	}

	// Fetch transactions in batches (limit is total transactions to fetch)
	batchSize := uint32(50)
	var stop bool
	var foundLastGrade bool
	var calls int
	for ; !stop; calls++ {
		txs, err := api.ListTransactions(ctx, s.CertificateAddress, batchSize, lastLt, lastHash)
		if err != nil {
			break
		}

		for i := len(txs) - 1; i >= 0; i-- {
			tx := txs[i]
			fmt.Println("NADO VSE", fmt.Sprintf("tx.LT: %v", tx.LT), fmt.Sprintf("lastProcessedLt: %v", lastProcessedLt),
				fmt.Sprintf("tx.Hash: %x", tx.Hash), fmt.Sprintf("lastProcessedHash: %x", lastProcessedHash))
			fmt.Println(tx.LT == lastProcessedLt && bytes.Compare(tx.Hash, lastProcessedHash) == 0)
			if tx.LT == lastProcessedLt && bytes.Compare(tx.Hash, lastProcessedHash) == 0 {
				if len(emits) == 0 {
					s.lastProcessedLt = txs[len(txs)-1].LT
					s.lastProcessedHash = txs[len(txs)-1].Hash
				}
				stop = true
				break
			}
			// Traverse outgoing messages linked list
			if tx.IO.Out != nil {
				currents, err := tx.IO.Out.ToSlice()
				if err != nil {
					return nil, err
				}
				for _, cur := range currents {
					if cur.Msg.DestAddr().IsAddrNone() {
						emits = append(emits, &Emit{
							payload: cur.Msg.Payload(),
							txLt:    tx.LT,
							txHash:  tx.Hash,
						})
					}
				}
			}
			if !foundLastGrade && tx.IO.In != nil {
				if tx.IO.In.Msg.SenderAddr().String() == courseOwner.String() {
					slice := tx.IO.In.Msg.Payload().BeginParse()
					slice.LoadUInt(32)
					payload := slice.MustLoadStringSnake()
					payloadSlice := strings.Split(payload, " | ")
					if len(payloadSlice) != 4 {
						continue
					}

					quizId, err := strconv.Atoi(payloadSlice[0])
					if err != nil {
						continue
					}

					lt, err := strconv.ParseUint(payloadSlice[2], 10, 64)
					if err != nil {
						continue
					}

					hash, err := base64.StdEncoding.DecodeString(payloadSlice[3])
					if err != nil {
						continue
					}
					s.QuizId, lastProcessedLt, lastProcessedHash = quizId, lt, hash
					foundLastGrade = true
					fmt.Println("getNewEmits", quizId, lt, hash)
				}
			}
		}

		// Prepare for next batch: use oldest tx LT/hash from current batch
		lastTx := txs[0]
		lastLt = lastTx.PrevTxLT
		lastHash = lastTx.PrevTxHash
	}

	if calls == 1 && len(emits) == 0 {
		s.lastProcessedLt = account.LastTxLT
		s.lastProcessedHash = account.LastTxHash
	}

	return emits, nil
}

func certificateFromInit(ctx context.Context, api edu.TonApi, collectionAddress *address.Address, certificateIndex int64) *edu.CertificateClient {
	state := &tlb.StateInit{
		Data: getCertificateData(collectionAddress, certificateIndex),
		Code: getCertificateCode(),
	}

	stateCell, _ := tlb.ToCell(state)

	addr := address.NewAddress(0, byte(int8(0)), stateCell.Hash())
	return edu.NewCertificateClient(api, addr)
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
