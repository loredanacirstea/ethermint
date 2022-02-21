package types

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
