package ton

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/a2tonium/a2tonium-backend/pkg/logger"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"strconv"
	"strings"
	"time"
)

type student struct {
	// QuizId - Last passed quiz id + 1
	QuizId                        int
	Gmail, IIN                    string
	CertificateAddress, ownerAddr *address.Address

	lastProcessedLt   uint64
	lastProcessedHash []byte
}

type emit struct {
	payload *cell.Cell
	txLt    uint64
	txHash  []byte
}

func (s *student) getAllGrades(ctx context.Context, api *ton.APIClient, courseOwner *address.Address) ([]string, error) {
	var account *tlb.Account
	lastQuizId := s.QuizId
	account, err := getAccountInstance(ctx, api, s.CertificateAddress)
	if err != nil {
		logger.ErrorKV(ctx, "student getAccountInstance() failed", logger.Err, err)
		return nil, fmt.Errorf("s.getAllGrades.getAccountInstance failed: %w", err)
	}

	grades := make([]string, lastQuizId)
	lastLt := account.LastTxLT
	lastHash := account.LastTxHash

	// Fetch transactions in batches (limit is total transactions to fetch)
	batchSize := uint32(50)
	for lastQuizId > 0 {
		var txs []*tlb.Transaction
		for i := 0; i < 3; i++ {
			txs, err = api.ListTransactions(ctx, s.CertificateAddress, batchSize, lastLt, lastHash)
			if err != nil {
				logger.WarnKV(ctx, "student api.ListTransaction failed", logger.Err, err)
				time.Sleep(1 * time.Second)
				continue
			}
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

	if lastQuizId == 0 {
		return grades, nil
	}

	return grades, errors.New("not all grades was found")
}

// Collect new emitted external-out messages (emits) for a contract address and return is the certificate issue call was from student
func (s *student) getNewEmitsWithCertCallCheck(ctx context.Context, api *ton.APIClient, courseOwner *address.Address) ([]*emit, bool, error) {
	var (
		certificateIssue bool
		err              error
	)

	account, err := getAccountInstance(ctx, api, s.CertificateAddress)
	if err != nil {
		logger.ErrorKV(ctx, "s.getAllGrades.getAccountInstance failed", logger.Err, err)
		return nil, certificateIssue, fmt.Errorf("s.getAllGrades.getAccountInstance failed: %w", err)
	}

	var emits []*emit
	lastLt, lastHash := account.LastTxLT, account.LastTxHash
	lastProcessedLt, lastProcessedHash := s.lastProcessedLt, s.lastProcessedHash

	if lastLt == lastProcessedLt && bytes.Compare(lastHash, lastProcessedHash) == 0 {
		return emits, certificateIssue, nil
	}

	// Fetch transactions in batches (limit is total transactions to fetch)
	batchSize := uint32(50)
	var stop, foundLastGrade bool
	var calls int

	for ; !stop; calls++ {
		var txs []*tlb.Transaction
		for i := 0; i < 3; i++ {
			txs, err = api.ListTransactions(ctx, s.CertificateAddress, batchSize, lastLt, lastHash)
			if err != nil {
				logger.WarnKV(ctx, "student api.ListTransaction failed", logger.Err, err)
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}

		if len(txs) == 0 {
			break
		}

		for i := len(txs) - 1; i >= 0; i-- {
			tx := txs[i]
			if tx.LT == lastProcessedLt && bytes.Compare(tx.Hash, lastProcessedHash) == 0 {
				if len(emits) == 0 {
					s.lastProcessedLt, s.lastProcessedHash = txs[len(txs)-1].LT, txs[len(txs)-1].Hash
				}
				stop = true
				break
			}

			// Traverse outgoing messages linked list
			if tx.IO.Out != nil {
				currents, err := tx.IO.Out.ToSlice()
				if err != nil {
					return nil, certificateIssue, err
				}
				for _, cur := range currents {
					if cur.Msg.DestAddr().IsAddrNone() {
						emits = append(emits, &emit{
							payload: cur.Msg.Payload(),
							txLt:    tx.LT,
							txHash:  tx.Hash,
						})
					}
				}
			}
			if !foundLastGrade && tx.IO.In != nil && tx.IO.In.Msg.SenderAddr().String() == courseOwner.String() {
				slice := tx.IO.In.Msg.Payload().BeginParse()
				q, err := slice.LoadUInt(32)
				if err != nil {
					continue
				}
				if q == 0 {
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
				}
			} else if tx.IO.In != nil && tx.IO.In.Msg.SenderAddr().String() == s.ownerAddr.String() {
				slice := tx.IO.In.Msg.Payload().BeginParse()
				q, err := slice.LoadUInt(32)
				if err != nil {
					continue
				}
				if q == 0 {
					payload := slice.MustLoadStringSnake()
					certificateIssue = strings.HasPrefix(payload, "Rating:")
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

	return emits, certificateIssue, nil
}

func getAccountInstance(ctx context.Context, api *ton.APIClient, certificateAddress *address.Address) (*tlb.Account, error) {
	var (
		account *tlb.Account
		block   *ton.BlockIDExt
		err     error
	)
	retryLimit := 3

	// Get current masterchain info for block reference
	for i := 0; i < retryLimit; i++ {
		block, err = api.CurrentMasterchainInfo(ctx)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	if err != nil {
		logger.ErrorKV(ctx, "getting current masterchain info failed", logger.Err, err)
		return nil, fmt.Errorf("getting current masterchain info failed: %w", err)
	}
	for i := 0; i < retryLimit; i++ {
		// Get account state to find last transaction LT and hash
		account, err = api.GetAccount(ctx, block, certificateAddress)
		if err == nil {
			break
		}
	}
	if err != nil {
		logger.ErrorKV(ctx, "api.GetAccount failed", logger.Err, err)
		return nil, fmt.Errorf("api.GetAccount failed: %w", err)
	}

	return account, nil
}
