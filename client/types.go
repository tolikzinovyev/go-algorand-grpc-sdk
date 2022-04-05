package client

import (
	"github.com/tolikzinovyev/go-algorand-grpc-sdk/types"
)

type AccountResponse struct {
	AccountData types.AccountData
	Round types.Round
	AmountWithoutPendingRewards types.MicroAlgos
	MinBalance uint64
}
