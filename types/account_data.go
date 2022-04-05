package types

type Status byte

const (
	// Offline indicates that the associated account receives rewards but does not participate in the consensus.
	Offline Status = iota
	// Online indicates that the associated account participates in the consensus and receive rewards.
	Online
	// NotParticipating indicates that the associated account neither participates in the consensus, nor receives rewards.
	// Accounts that are marked as NotParticipating cannot change their status, but can receive and send Algos to other accounts.
	// Two special accounts that are defined as NotParticipating are the incentive pool (also know as rewards pool) and the fee sink.
	// These two accounts also have additional Algo transfer restrictions.
	NotParticipating
)

type AccountData struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Status     Status     `codec:"onl"`
	MicroAlgos MicroAlgos `codec:"algo"`

	// RewardsBase is used to implement rewards.
	// This is not meaningful for accounts with Status=NotParticipating.
	//
	// Every block assigns some amount of rewards (algos) to every
	// participating account.  The amount is the product of how much
	// block.RewardsLevel increased from the previous block and
	// how many whole config.Protocol.RewardUnit algos this
	// account holds.
	//
	// For performance reasons, we do not want to walk over every
	// account to apply these rewards to AccountData.MicroAlgos.  Instead,
	// we defer applying the rewards until some other transaction
	// touches that participating account, and at that point, apply all
	// of the rewards to the account's AccountData.MicroAlgos.
	//
	// For correctness, we need to be able to determine how many
	// total algos are present in the system, including deferred
	// rewards (deferred in the sense that they have not been
	// reflected in the account's AccountData.MicroAlgos, as described
	// above).  To compute this total efficiently, we avoid
	// compounding rewards (i.e., no rewards on rewards) until
	// they are applied to AccountData.MicroAlgos.
	//
	// Mechanically, RewardsBase stores the block.RewardsLevel
	// whose rewards are already reflected in AccountData.MicroAlgos.
	// If the account is Status=Offline or Status=Online, its
	// effective balance (if a transaction were to be issued
	// against this account) may be higher, as computed by
	// AccountData.Money().  That function calls
	// AccountData.WithUpdatedRewards() to apply the deferred
	// rewards to AccountData.MicroAlgos.
	RewardsBase uint64 `codec:"ebase"`

	// RewardedMicroAlgos is used to track how many algos were given
	// to this account since the account was first created.
	//
	// This field is updated along with RewardBase; note that
	// it won't answer the question "how many algos did I make in
	// the past week".
	RewardedMicroAlgos MicroAlgos `codec:"ern"`

	VoteID			 VotePK `codec:"vote"`
	SelectionID  VRFPK              `codec:"sel"`
	StateProofID MerkleVerifier        `codec:"stprf"`

	VoteFirstValid  Round  `codec:"voteFst"`
	VoteLastValid   Round  `codec:"voteLst"`
	VoteKeyDilution uint64 `codec:"voteKD"`

	// If this account created an asset, AssetParams stores
	// the parameters defining that asset.  The params are indexed
	// by the Index of the AssetID; the Creator is this account's address.
	//
	// An account with any asset in AssetParams cannot be
	// closed, until the asset is destroyed.  An asset can
	// be destroyed if this account holds AssetParams.Total units
	// of that asset (in the Assets array below).
	//
	// NOTE: do not modify this value in-place in existing AccountData
	// structs; allocate a copy and modify that instead.  AccountData
	// is expected to have copy-by-value semantics.
	AssetParams map[AssetIndex]AssetParams `codec:"apar,allocbound=encodedMaxAssetsPerAccount"`

	// Assets is the set of assets that can be held by this
	// account.  Assets (i.e., slots in this map) are explicitly
	// added and removed from an account by special transactions.
	// The map is keyed by the AssetID, which is the address of
	// the account that created the asset plus a unique counter
	// to distinguish re-created assets.
	//
	// Each asset bumps the required MinBalance in this account.
	//
	// An account that creates an asset must have its own asset
	// in the Assets map until that asset is destroyed.
	//
	// NOTE: do not modify this value in-place in existing AccountData
	// structs; allocate a copy and modify that instead.  AccountData
	// is expected to have copy-by-value semantics.
	Assets map[AssetIndex]AssetHolding `codec:"asset,allocbound=encodedMaxAssetsPerAccount"`

	// AuthAddr is the address against which signatures/multisigs/logicsigs should be checked.
	// If empty, the address of the account whose AccountData this is is used.
	// A transaction may change an account's AuthAddr to "re-key" the account.
	// This allows key rotation, changing the members in a multisig, etc.
	AuthAddr Address `codec:"spend"`

	// AppLocalStates stores the local states associated with any applications
	// that this account has opted in to.
	AppLocalStates map[AppIndex]AppLocalState `codec:"appl,allocbound=EncodedMaxAppLocalStates"`

	// AppParams stores the global parameters and state associated with any
	// applications that this account has created.
	AppParams map[AppIndex]AppParams `codec:"appp,allocbound=EncodedMaxAppParams"`

	// TotalAppSchema stores the sum of all of the LocalStateSchemas
	// and GlobalStateSchemas in this account (global for applications
	// we created local for applications we opted in to), so that we don't
	// have to iterate over all of them to compute MinBalance.
	TotalAppSchema StateSchema `codec:"tsch"`

	// TotalExtraAppPages stores the extra length in pages (MaxAppProgramLen bytes per page)
	// requested for app program by this account
	TotalExtraAppPages uint32 `codec:"teap"`
}
