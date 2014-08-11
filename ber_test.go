package ber

import (
	"bytes"
	"io"
	"testing"
)

func TestEncodeDecodeInterger(t *testing.T) {
	var integer uint64 = 10

	encodedInteger := EncodeInteger(integer)
	decodedInteger := DecodeInteger(encodedInteger)

	if integer != decodedInteger {
		t.Error("wrong should be equal", integer, decodedInteger)
	}

}

func TestBoolean(t *testing.T) {
	var value bool = true

	packet := NewBoolean(ClassUniversal, TypePrimitive, TagBoolean, value, "first Packet, True")

	newBoolean, ok := packet.Value.(bool)
	if !ok || newBoolean != value {
		t.Error("error during creating packet")
	}

	encodedPacket := packet.Bytes()

	newPacket := DecodePacket(encodedPacket)

	newBoolean, ok = newPacket.Value.(bool)
	if !ok || newBoolean != value {
		t.Error("error during decoding packet")
	}

}

func TestInteger(t *testing.T) {
	var value uint64 = 10

	packet := NewInteger(ClassUniversal, TypePrimitive, TagInteger, value, "Integer, 10")

	newInteger, ok := packet.Value.(uint64)
	if !ok || newInteger != value {
		t.Error("error during creating packet")
	}

	encodedPacket := packet.Bytes()

	newPacket := DecodePacket(encodedPacket)

	newInteger, ok = newPacket.Value.(uint64)
	if !ok || newInteger != value {
		t.Error("error during decoding packet")
	}

}

func TestString(t *testing.T) {
	var value string = "Hic sunt dracones"

	packet := NewString(ClassUniversal, TypePrimitive, TagOctetString, value, "String")

	newValue, ok := packet.Value.(string)
	if !ok || newValue != value {
		t.Error("error during creating packet")
	}

	encodedPacket := packet.Bytes()

	newPacket := DecodePacket(encodedPacket)

	newValue, ok = newPacket.Value.(string)
	if !ok || newValue != value {
		t.Error("error during decoding packet")
	}

}

func TestSequenceAndAppendChild(t *testing.T) {

	p1 := NewString(ClassUniversal, TypePrimitive, TagOctetString, "HIC SVNT LEONES", "String")
	p2 := NewString(ClassUniversal, TypePrimitive, TagOctetString, "HIC SVNT DRACONES", "String")
	p3 := NewString(ClassUniversal, TypePrimitive, TagOctetString, "Terra Incognita", "String")

	sequence := NewSequence("a sequence")
	sequence.AppendChild(p1)
	sequence.AppendChild(p2)
	sequence.AppendChild(p3)

	if len(sequence.Children) != 3 {
		t.Error("wrong length for children array should be three =>", len(sequence.Children))
	}

	encodedSequence := sequence.Bytes()

	decodedSequence := DecodePacket(encodedSequence)
	if len(decodedSequence.Children) != 3 {
		t.Error("wrong length for children array should be three =>", len(decodedSequence.Children))
	}

}

func TestPrint(t *testing.T) {
	p1 := NewString(ClassUniversal, TypePrimitive, TagOctetString, "Answer to the Ultimate Question of Life, the Universe, and Everything", "Question")
	p2 := NewInteger(ClassUniversal, TypePrimitive, TagInteger, 42, "Answer")
	p3 := NewBoolean(ClassUniversal, TypePrimitive, TagBoolean, true, "Validity")

	sequence := NewSequence("a sequence")
	sequence.AppendChild(p1)
	sequence.AppendChild(p2)
	sequence.AppendChild(p3)

	PrintPacket(sequence)

	encodedSequence := sequence.Bytes()
	PrintBytes(encodedSequence, "\t")
}

func TestReadPacket(t *testing.T) {
	packet := NewString(ClassUniversal, TypePrimitive, TagOctetString, "Ad impossibilia nemo tenetur", "string")
	var buffer io.ReadWriter
	buffer = new(bytes.Buffer)

	buffer.Write(packet.Bytes())

	newPacket, err := ReadPacket(buffer)
	if err != nil {
		t.Error("error during ReadPacket", err)
	}

	if !bytes.Equal(newPacket.ByteValue, packet.ByteValue) {
		t.Error("packets should be the same")
	}
}
