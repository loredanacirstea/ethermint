package types

// The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState struct {
	PortId string
}

// // NewGenesisState creates a new genesis state.
// func NewGenesisState(params Params, pairs []TokenPair) GenesisState {
// 	return GenesisState{
// 		Params:     params,
// 		TokenPairs: pairs,
// 	}
// }

// // DefaultGenesisState sets default evm genesis state with empty accounts and
// // default params and chain config values.
// func DefaultGenesisState() *GenesisState {
// 	return &GenesisState{
// 		Params: DefaultParams(),
// 	}
// }

type Params struct {
	MessageContent []byte `protobuf:"varint,1,opt,name=message_content,json=messageContent,proto3" json:"message_content,omitempty"`
}

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PortId: PortID,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return nil
}
