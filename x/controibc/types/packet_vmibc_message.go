package types

import "github.com/tendermint/tendermint/crypto/tmhash"

// ValidateBasic is used for validating the packet
func (p VmibcMessagePacketData) ValidateBasic() error {

	// TODO: Validate the packet data

	return nil
}

// GetBytes is a helper for serialising
func (p VmibcMessagePacketData) GetBytes() ([]byte, error) {
	var modulePacket ControibcPacketData

	modulePacket.Packet = &ControibcPacketData_VmibcMessagePacket{&p}

	return modulePacket.Marshal()
}

// GetID returns the SHA256 hash of the ERC20 address and denomination
func (p VmibcMessagePacketData) GetID() []byte {
	id := p.Body
	return tmhash.Sum([]byte(id))
}
