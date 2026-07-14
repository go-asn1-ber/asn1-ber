package ber_test

import (
	"fmt"

	ber "github.com/go-asn1-ber/asn1-ber"
)

// ExampleDecodePacketErr encodes a SEQUENCE containing a single OCTET STRING,
// serializes it to BER bytes, and decodes it back into a packet tree.
func ExampleDecodePacketErr() {
	seq := ber.NewSequence("greeting")
	seq.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "Hello, world", ""))

	packet, err := ber.DecodePacketErr(seq.Bytes())
	if err != nil {
		panic(err)
	}

	fmt.Println(packet.Children[0].Value)
	// Output: Hello, world
}
