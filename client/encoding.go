package client

import (
	"github.com/tolikzinovyev/go-algorand-grpc-sdk/client/internal/proto"
	"github.com/tolikzinovyev/go-algorand-grpc-sdk/types"
)

func unconvertStatus(status proto.Status) types.Status {
	switch status {
	case proto.Status_OFFLINE:
		return types.Offline
	case proto.Status_ONLINE:
		return types.Online
	case proto.Status_NOT_PARTICIPATING:
		return types.NotParticipating
	}

	return 255
}

func unconvertVotePK(votePK []byte) types.VotePK {
	var res types.VotePK
	copy(res[:], votePK)
	return res
}

func unconvertVRFPK(vrfPK []byte) types.VRFPK {
	var res types.VRFPK
	copy(res[:], vrfPK)
	return res
}

func unconvertMerkleVerifier(merkleVerifier []byte) types.MerkleVerifier {
	var res types.MerkleVerifier
	copy(res[:], merkleVerifier)
	return res
}

func unconvertAssetMetadataHash(hash []byte) [types.AssetMetadataHashLen]byte {
	var res [types.AssetMetadataHashLen]byte
	copy(res[:], hash)
	return res
}

func unconvertAddress(address []byte) types.Address {
	var res types.Address
	copy(res[:], address)
	return res
}

func unconvertAssetParams(params *proto.AssetParams) types.AssetParams {
	return types.AssetParams{
		Total: params.Total,
		Decimals: params.Decimals,
		DefaultFrozen: params.DefaultFrozen,
		UnitName: string(params.UnitName),
		AssetName: string(params.AssetName),
		URL: string(params.Url),
		MetadataHash: unconvertAssetMetadataHash(params.MetadataHash),
		Manager: unconvertAddress(params.Manager),
		Reserve: unconvertAddress(params.Reserve),
		Freeze: unconvertAddress(params.Freeze),
		Clawback: unconvertAddress(params.Clawback),
	}
}

func unconvertAssetParamsMap(m map[uint64]*proto.AssetParams) map[types.AssetIndex]types.AssetParams {
	res := make(map[types.AssetIndex]types.AssetParams, len(m))
	for k, v := range m {
		res[types.AssetIndex(k)] = unconvertAssetParams(v)
	}
	return res
}

func unconvertAssetHolding(holding *proto.AssetHolding) types.AssetHolding {
	return types.AssetHolding{
		Amount: holding.Amount,
		Frozen: holding.Frozen,
	}
}

func unconvertAssetHoldingMap(m map[uint64]*proto.AssetHolding) map[types.AssetIndex]types.AssetHolding {
	res := make(map[types.AssetIndex]types.AssetHolding, len(m))
	for k, v := range m {
		res[types.AssetIndex(k)] = unconvertAssetHolding(v)
	}
	return res
}

func unconvertStateSchema(schema *proto.StateSchema) types.StateSchema {
	return types.StateSchema{
		NumUint: schema.NumUint,
		NumByteSlice: schema.NumByteSlice,
	}
}

func unconvertTealKeyValue(tkv []*proto.TealKeyValue) types.TealKeyValue {
	res := make(types.TealKeyValue, len(tkv))

	for _, e := range tkv {
		var tv types.TealValue
		switch x := e.Value.(type) {
		case *proto.TealKeyValue_Uint:
			tv.Type = types.TealUintType
			tv.Uint = x.Uint
		case *proto.TealKeyValue_Bytes:
			tv.Type = types.TealBytesType
			tv.Bytes = string(x.Bytes)
		}
		res[string(e.Key)] = tv
	}

	return res
}

func unconvertAppLocalState(state *proto.AppLocalState) types.AppLocalState {
	return types.AppLocalState{
		Schema: unconvertStateSchema(state.Schema),
		KeyValue: unconvertTealKeyValue(state.KeyValue),
	}
}

func unconvertAppLocalStateMap(m map[uint64]*proto.AppLocalState) map[types.AppIndex]types.AppLocalState {
	res := make(map[types.AppIndex]types.AppLocalState, len(m))
	for k, v := range m {
		res[types.AppIndex(k)] = unconvertAppLocalState(v)
	}
	return res
}

func unconvertAppParams(params *proto.AppParams) types.AppParams {
	return types.AppParams{
		ApprovalProgram: params.ApprovalProgram,
		ClearStateProgram: params.ClearStateProgram,
		GlobalState: unconvertTealKeyValue(params.GlobalState),
		StateSchemas: types.StateSchemas{
			LocalStateSchema: unconvertStateSchema(params.LocalStateSchema),
			GlobalStateSchema: unconvertStateSchema(params.GlobalStateSchema),
		},
		ExtraProgramPages: params.ExtraProgramPages,
	}
}

func unconvertAppParamsMap(m map[uint64]*proto.AppParams) map[types.AppIndex]types.AppParams {
	res := make(map[types.AppIndex]types.AppParams, len(m))
	for k, v := range m {
		res[types.AppIndex(k)] = unconvertAppParams(v)
	}
	return res
}

func unconvertAccountData(accountData *proto.AccountData) types.AccountData {
	return types.AccountData{
		Status: unconvertStatus(accountData.Status),
		MicroAlgos: types.MicroAlgos(accountData.Microalgos),
		RewardsBase: accountData.RewardsBase,
		RewardedMicroAlgos: types.MicroAlgos(accountData.RewardedMicroalgos),
		VoteID: unconvertVotePK(accountData.VoteId),
		SelectionID: unconvertVRFPK(accountData.SelectionId),
		StateProofID: unconvertMerkleVerifier(accountData.StateProofId),
		VoteFirstValid: types.Round(accountData.VoteFirstValid),
		VoteLastValid: types.Round(accountData.VoteLastValid),
		VoteKeyDilution: accountData.VoteKeyDilution,
		AssetParams: unconvertAssetParamsMap(accountData.AssetParams),
		Assets: unconvertAssetHoldingMap(accountData.AssetHoldings),
		AuthAddr: unconvertAddress(accountData.AuthAddr),
		AppLocalStates: unconvertAppLocalStateMap(accountData.AppLocalStates),
		AppParams: unconvertAppParamsMap(accountData.AppParams),
		TotalAppSchema: unconvertStateSchema(accountData.TotalAppSchema),
		TotalExtraAppPages: accountData.TotalExtraAppPages,
	}
}

func unconvertAccountResponse(response *proto.AccountResponse) AccountResponse {
	return AccountResponse{
		AccountData: unconvertAccountData(response.AccountData),
		Round: types.Round(response.Round),
		AmountWithoutPendingRewards: types.MicroAlgos(response.AmountWithoutPendingRewards),
		MinBalance: response.MinBalance,
	}
}
