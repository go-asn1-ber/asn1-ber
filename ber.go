package ber

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
)

type Packet struct {
	ClassType   uint8
	TagType     uint8
	Tag         uint8
	Value       interface{}
	ByteValue   []byte
	Data        *bytes.Buffer
	Children    []*Packet
	Description string
}

const (
	TagEOC              = 0x00
	TagBoolean          = 0x01
	TagInteger          = 0x02
	TagBitString        = 0x03
	TagOctetString      = 0x04
	TagNULL             = 0x05
	TagObjectIdentifier = 0x06
	TagObjectDescriptor = 0x07
	TagExternal         = 0x08
	TagRealFloat        = 0x09
	TagEnumerated       = 0x0a
	TagEmbeddedPDV      = 0x0b
	TagUTF8String       = 0x0c
	TagRelativeOID      = 0x0d
	TagSequence         = 0x10
	TagSet              = 0x11
	TagNumericString    = 0x12
	TagPrintableString  = 0x13
	TagT61String        = 0x14
	TagVideotexString   = 0x15
	TagIA5String        = 0x16
	TagUTCTime          = 0x17
	TagGeneralizedTime  = 0x18
	TagGraphicString    = 0x19
	TagVisibleString    = 0x1a
	TagGeneralString    = 0x1b
	TagUniversalString  = 0x1c
	TagCharacterString  = 0x1d
	TagBMPString        = 0x1e
	TagBitmask          = 0x1f // xxx11111b
)

var tagMap = map[uint8]string{
	TagEOC:              "EOC (End-of-Content)",
	TagBoolean:          "Boolean",
	TagInteger:          "Integer",
	TagBitString:        "Bit String",
	TagOctetString:      "Octet String",
	TagNULL:             "NULL",
	TagObjectIdentifier: "Object Identifier",
	TagObjectDescriptor: "Object Descriptor",
	TagExternal:         "External",
	TagRealFloat:        "Real (float)",
	TagEnumerated:       "Enumerated",
	TagEmbeddedPDV:      "Embedded PDV",
	TagUTF8String:       "UTF8 String",
	TagRelativeOID:      "Relative-OID",
	TagSequence:         "Sequence and Sequence of",
	TagSet:              "Set and Set OF",
	TagNumericString:    "Numeric String",
	TagPrintableString:  "Printable String",
	TagT61String:        "T61 String",
	TagVideotexString:   "Videotex String",
	TagIA5String:        "IA5 String",
	TagUTCTime:          "UTC Time",
	TagGeneralizedTime:  "Generalized Time",
	TagGraphicString:    "Graphic String",
	TagVisibleString:    "Visible String",
	TagGeneralString:    "General String",
	TagUniversalString:  "Universal String",
	TagCharacterString:  "Character String",
	TagBMPString:        "BMP String",
}

const (
	ClassUniversal   = 0   // 00xxxxxxb
	ClassApplication = 64  // 01xxxxxxb
	ClassContext     = 128 // 10xxxxxxb
	ClassPrivate     = 192 // 11xxxxxxb
	ClassBitmask     = 192 // 11xxxxxxb
)

var ClassMap = map[uint8]string{
	ClassUniversal:   "Universal",
	ClassApplication: "Application",
	ClassContext:     "Context",
	ClassPrivate:     "Private",
}

const (
	TypePrimitive   = 0  // xx0xxxxxb
	TypeConstructed = 32 // xx1xxxxxb
	TypeBitmask     = 32 // xx1xxxxxb
)

var TypeMap = map[uint8]string{
	TypePrimitive:   "Primitive",
	TypeConstructed: "Constructed",
}

var Debug bool = false

func PrintBytes(out io.Writer, buf []byte, indent string) {
	data_lines := make([]string, (len(buf)/30)+1)
	num_lines := make([]string, (len(buf)/30)+1)

	for i, b := range buf {
		data_lines[i/30] += fmt.Sprintf("%02x ", b)
		num_lines[i/30] += fmt.Sprintf("%02d ", (i+1)%100)
	}

	for i := 0; i < len(data_lines); i++ {
		out.Write([]byte(indent + data_lines[i] + "\n"))
		out.Write([]byte(indent + num_lines[i] + "\n\n"))
	}
}

func PrintPacket(p *Packet) {
	printPacket(os.Stdout, p, 0, false)
}

func printPacket(out io.Writer, p *Packet, indent int, printBytes bool) {
	indent_str := ""

	for len(indent_str) != indent {
		indent_str += " "
	}

	class_str := ClassMap[p.ClassType]

	tagtype_str := TypeMap[p.TagType]

	tag_str := fmt.Sprintf("0x%02X", p.Tag)

	if p.ClassType == ClassUniversal {
		tag_str = tagMap[p.Tag]
	}

	value := fmt.Sprint(p.Value)
	description := ""

	if p.Description != "" {
		description = p.Description + ": "
	}

	fmt.Fprintf(out, "%s%s(%s, %s, %s) Len=%d %q\n", indent_str, description, class_str, tagtype_str, tag_str, p.Data.Len(), value)

	if printBytes {
		PrintBytes(out, p.Bytes(), indent_str)
	}

	for _, child := range p.Children {
		printPacket(out, child, indent+1, printBytes)
	}
}

func resizeBuffer(in []byte, new_size int) (out []byte) {
	out = make([]byte, new_size)

	copy(out, in)

	return
}

func ReadPacket(reader io.Reader) (*Packet, error) {
	var header [2]byte
	buf := header[:]
	_, err := io.ReadFull(reader, buf)

	if err != nil {
		return nil, err
	}

	idx := 2
	var datalen int
	l := buf[1]

	if l&0x80 == 0 {
		// The length is encoded in the bottom 7 bits.
		datalen = int(l & 0x7f)
		if Debug {
			fmt.Printf("Read: datalen = %d len(buf) = %d\n  ", l, len(buf))

			for _, b := range buf {
				fmt.Printf("%02X ", b)
			}

			fmt.Printf("\n")
		}
	} else {
		// Bottom 7 bits give the number of length bytes to follow.
		numBytes := int(l & 0x7f)
		if numBytes == 0 {
			return nil, fmt.Errorf("invalid length found")
		}
		idx += numBytes
		buf = resizeBuffer(buf, 2+numBytes)
		_, err := io.ReadFull(reader, buf[2:])

		if err != nil {
			return nil, err
		}
		datalen = 0
		for i := 0; i < numBytes; i++ {
			b := buf[2+i]
			datalen <<= 8
			datalen |= int(b)
		}

		if Debug {
			fmt.Printf("Read: datalen = %d numbytes=%d len(buf) = %d\n  ", datalen, numBytes, len(buf))

			for _, b := range buf {
				fmt.Printf("%02X ", b)
			}

			fmt.Printf("\n")
		}
	}

	buf = resizeBuffer(buf, idx+datalen)
	_, err = io.ReadFull(reader, buf[idx:])

	if err != nil {
		return nil, err
	}

	if Debug {
		fmt.Printf("Read: len( buf ) = %d  idx=%d datalen=%d idx+datalen=%d\n  ", len(buf), idx, datalen, idx+datalen)

		for _, b := range buf {
			fmt.Printf("%02X ", b)
		}
	}

	p, _ := decodePacket(buf)

	return p, nil
}

func DecodeString(data []byte) string {
	return string(data)
}

func parseInt64(bytes []byte) (ret int64, err error) {
	if len(bytes) > 8 {
		// We'll overflow an int64 in this case.
		err = fmt.Errorf("integer too large")
		return
	}
	for bytesRead := 0; bytesRead < len(bytes); bytesRead++ {
		ret <<= 8
		ret |= int64(bytes[bytesRead])
	}

	// Shift up and down in order to sign extend the result.
	ret <<= 64 - uint8(len(bytes))*8
	ret >>= 64 - uint8(len(bytes))*8
	return
}

func encodeInteger(i int64) []byte {
	n := int64Length(i)
	out := make([]byte, n)

	var j int
	for ; n > 0; n-- {
		out[j] = (byte(i >> uint((n-1)*8)))
		j++
	}

	return out
}

func int64Length(i int64) (numBytes int) {
	numBytes = 1

	for i > 127 {
		numBytes++
		i >>= 8
	}

	for i < -128 {
		numBytes++
		i >>= 8
	}

	return
}

func DecodePacket(data []byte) *Packet {
	p, _ := decodePacket(data)

	return p
}

func decodePacket(data []byte) (*Packet, []byte) {
	if Debug {
		fmt.Printf("decodePacket: enter %d\n", len(data))
	}

	p := new(Packet)

	p.ClassType = data[0] & ClassBitmask
	p.TagType = data[0] & TypeBitmask
	p.Tag = data[0] & TagBitmask

	var datalen int
	l := data[1]
	datapos := 2
	if l&0x80 == 0 {
		// The length is encoded in the bottom 7 bits.
		datalen = int(l & 0x7f)
	} else {
		// Bottom 7 bits give the number of length bytes to follow.
		numBytes := int(l & 0x7f)
		if numBytes == 0 {
			return nil, nil
		}
		datapos += numBytes
		datalen = 0
		for i := 0; i < numBytes; i++ {
			b := data[2+i]
			datalen <<= 8
			datalen |= int(b)
		}
	}

	p.Data = new(bytes.Buffer)

	p.Children = make([]*Packet, 0, 2)

	p.Value = nil

	value_data := data[datapos : datapos+datalen]

	if p.TagType == TypeConstructed {
		for len(value_data) != 0 {
			var child *Packet

			child, value_data = decodePacket(value_data)
			p.AppendChild(child)
		}
	} else if p.ClassType == ClassUniversal {
		p.Data.Write(data[datapos : datapos+datalen])
		p.ByteValue = value_data

		switch p.Tag {
		case TagEOC:
		case TagBoolean:
			val, _ := parseInt64(value_data)

			p.Value = val != 0
		case TagInteger:
			p.Value, _ = parseInt64(value_data)
		case TagBitString:
		case TagOctetString:
			// the actual string encoding is not known here
			// (e.g. for LDAP value_data is already an UTF8-encoded
			// string). Return the data without further processing
			p.Value = string(value_data)
		case TagNULL:
		case TagObjectIdentifier:
		case TagObjectDescriptor:
		case TagExternal:
		case TagRealFloat:
		case TagEnumerated:
			p.Value, _ = parseInt64(value_data)
		case TagEmbeddedPDV:
		case TagUTF8String:
		case TagRelativeOID:
		case TagSequence:
		case TagSet:
		case TagNumericString:
		case TagPrintableString:
			p.Value, _ = parseInt64(value_data)
		case TagT61String:
		case TagVideotexString:
		case TagIA5String:
		case TagUTCTime:
		case TagGeneralizedTime:
		case TagGraphicString:
		case TagVisibleString:
		case TagGeneralString:
		case TagUniversalString:
		case TagCharacterString:
		case TagBMPString:
		}
	} else {
		p.Data.Write(data[datapos : datapos+datalen])
	}

	return p, data[datapos+datalen:]
}

func (p *Packet) Bytes() []byte {
	var out bytes.Buffer

	out.Write([]byte{p.ClassType | p.TagType | p.Tag})
	packet_length := encodeInteger(int64(p.Data.Len()))

	if p.Data.Len() > 127 || len(packet_length) > 1 {
		out.Write([]byte{byte(len(packet_length) | 128)})
		out.Write(packet_length)
	} else {
		out.Write(packet_length)
	}

	out.Write(p.Data.Bytes())

	return out.Bytes()
}

func (p *Packet) AppendChild(child *Packet) {
	p.Data.Write(child.Bytes())
	p.Children = append(p.Children, child)
}

func Encode(ClassType, TagType, Tag uint8, Value interface{}, Description string) *Packet {
	p := new(Packet)

	p.ClassType = ClassType
	p.TagType = TagType
	p.Tag = Tag
	p.Data = new(bytes.Buffer)

	p.Children = make([]*Packet, 0, 2)

	p.Value = Value
	p.Description = Description

	if Value != nil {
		v := reflect.ValueOf(Value)

		if ClassType == ClassUniversal {
			switch Tag {
			case TagOctetString:
				sv, ok := v.Interface().(string)

				if ok {
					p.Data.Write([]byte(sv))
				}
			}
		}
	}

	return p
}

func NewSequence(Description string) *Packet {
	return Encode(ClassUniversal, TypeConstructed, TagSequence, nil, Description)
}

func NewBoolean(ClassType, TagType, Tag uint8, Value bool, Description string) *Packet {
	intValue := int64(0)

	if Value {
		intValue = 1
	}

	p := Encode(ClassType, TagType, Tag, nil, Description)

	p.Value = Value
	p.Data.Write(encodeInteger(intValue))

	return p
}

func NewInteger(ClassType, TagType, Tag uint8, Value int64, Description string) *Packet {
	p := Encode(ClassType, TagType, Tag, nil, Description)

	p.Value = Value
	p.Data.Write(encodeInteger(Value))

	return p
}

func NewString(ClassType, TagType, Tag uint8, Value, Description string) *Packet {
	p := Encode(ClassType, TagType, Tag, nil, Description)

	p.Value = Value
	p.Data.Write([]byte(Value))

	return p
}
