package basemodule

type AddressDecoder struct {
	UnknownSlot       uint8
	FrontTopCenter    uint8
	FrontTopRight     uint8
	FrontBottomLeft   uint8
	FrontBottomCenter uint8
	FrontBottomRight  uint8
	BackTopCenter     uint8
	BackTopRight      uint8
	BackBottomLeft    uint8
	BackBottomCenter  uint8
	BackBottomRight   uint8
}

func newAddressDecoder() *AddressDecoder {
	return &AddressDecoder{
		UnknownSlot:       0x29,
		FrontTopCenter:    0x30,
		FrontTopRight:     0x31,
		FrontBottomLeft:   0x32,
		FrontBottomCenter: 0x33,
		FrontBottomRight:  0x34,
		BackTopCenter:     0x41,
		BackTopRight:      0x42,
		BackBottomLeft:    0x43,
		BackBottomCenter:  0x44,
		BackBottomRight:   0x45,
	}
}

func (addr *AddressDecoder) convertToAddress(addressVoltage uint16) uint8 {
	const (
		VoltageMax = 1024
	)
	var (
		// AddressList maps the address voltage ranges to the corresponding slot with
		// the minimum and maximum values from the ADC
		AddressList = map[uint8][]uint16{
			addr.FrontTopCenter:    {93, 186},
			addr.FrontTopRight:     {186, 279},
			addr.FrontBottomLeft:   {279, 372},
			addr.FrontBottomCenter: {372, 465},
			addr.FrontBottomRight:  {465, 558},
			addr.BackTopCenter:     {558, 651},
			addr.BackTopRight:      {651, 744},
			addr.BackBottomLeft:    {744, 837},
			addr.BackBottomCenter:  {837, 930},
			addr.BackBottomRight:   {930, 1024},
		}
	)

	if addressVoltage > 1024 {
		return addr.UnknownSlot
	}

	for key, value := range AddressList {
		if addressVoltage >= value[0] && addressVoltage < value[1] {
			return key
		}
	}

	return addr.UnknownSlot
}
