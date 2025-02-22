package flow

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecordBinaryEncoding(t *testing.T) {
	// Makes sure that we read the C *packed* flow structure according
	// to the order defined in bpf/flow.h
	fr, err := ReadFrom(bytes.NewReader([]byte{
		0x01, 0x02, // u16 protocol
		0x03,                               // u16 direction
		0x04, 0x05, 0x06, 0x07, 0x08, 0x09, // data_link: u8[6] src_mac
		0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, // data_link: u8[6] dst_mac
		0x06, 0x07, 0x08, 0x09, // network: u32 src_ip
		0x0a, 0x0b, 0x0c, 0x0d, // network: u32 dst_ip
		0x0e, 0x0f, // transport: u16 src_port
		0x10, 0x11, // transport: u16 dst_port
		0x12,                                           // transport: u8protocol
		0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, // u64 bytes
	}))
	require.NoError(t, err)

	assert.Equal(t, Record{
		rawRecord: rawRecord{
			key: key{
				Protocol:  0x0201,
				Direction: 0x03,
				DataLink: DataLink{
					SrcMac: MacAddr{0x04, 0x05, 0x06, 0x07, 0x08, 0x09},
					DstMac: MacAddr{0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
				},
				Network: Network{
					SrcAddr: 0x09080706,
					DstAddr: 0x0d0c0b0a,
				},
				Transport: Transport{
					SrcPort:  0x0f0e,
					DstPort:  0x1110,
					Protocol: 0x12,
				},
			},
			Bytes: 0x1a19181716151413,
		},
	}, *fr)
}
