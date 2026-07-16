package ber

import (
	"bytes"
	"io"
	"math"
	"strings"
	"testing"
)

func TestEncodeDecodeInteger(t *testing.T) {
	for _, v := range []int64{0, 10, 128, 1024, math.MaxInt64, -1, -100, -128, -1024, math.MinInt64} {
		enc := encodeInteger(v)
		dec, err := ParseInt64(enc)
		if err != nil {
			t.Fatalf("Error decoding %d : %s", v, err)
		}
		if v != dec {
			t.Errorf("TestEncodeDecodeInteger failed for %d (got %d)", v, dec)
		}
	}
}

func TestBoolean(t *testing.T) {
	packet := NewBoolean(ClassUniversal, TypePrimitive, TagBoolean, true, "first Packet, True")

	newBoolean, ok := packet.Value.(bool)
	if !ok || newBoolean != true {
		t.Error("error during creating packet")
	}

	encodedPacket := packet.Bytes()

	newPacket := DecodePacket(encodedPacket)

	newBoolean, ok = newPacket.Value.(bool)
	if !ok || newBoolean != true {
		t.Error("error during decoding packet")
	}
}

func TestLDAPBoolean(t *testing.T) {
	packet := NewLDAPBoolean(ClassUniversal, TypePrimitive, TagBoolean, true, "first Packet, True")

	newBoolean, ok := packet.Value.(bool)
	if !ok || newBoolean != true {
		t.Error("error during creating packet")
	}

	encodedPacket := packet.Bytes()

	newPacket := DecodePacket(encodedPacket)

	newBoolean, ok = newPacket.Value.(bool)
	if !ok || newBoolean != true {
		t.Error("error during decoding packet")
	}
}

func TestInteger(t *testing.T) {
	var value int64 = 10

	packet := NewInteger(ClassUniversal, TypePrimitive, TagInteger, value, "Integer, 10")

	{
		newInteger, ok := packet.Value.(int64)
		if !ok || newInteger != value {
			t.Error("error creating packet")
		}
	}

	encodedPacket := packet.Bytes()

	newPacket := DecodePacket(encodedPacket)

	{
		newInteger, ok := newPacket.Value.(int64)
		if !ok || newInteger != value {
			t.Error("error decoding packet")
		}
	}
}

func TestString(t *testing.T) {
	value := "Hic sunt dracones"

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

func TestEncodeDecodeOID(t *testing.T) {
	for _, v := range []string{"0.1", "2.981", "2.3", "0.4", "0.4.5.1888", "0.10.5.1888.234.324234"} {
		enc, err := encodeOID(v)
		if err != nil {
			t.Errorf("error on encoding object identifier when encoding %s: %v", v, err)
		}
		parsed, err := parseObjectIdentifier(enc)
		if err != nil {
			t.Errorf("error on parsing object identifier when parsing %s: %v", v, err)
		}
		t.Log(enc)
		t.Log(OIDToString(parsed))
		if v != OIDToString(parsed) {
			t.Error("encoded object identifier did not match parsed")
		}
	}
}

func TestSequenceAndAppendChild(t *testing.T) {
	values := []string{
		"HIC SVNT LEONES",
		"Iñtërnâtiônàlizætiøn",
		"Terra Incognita",
	}

	sequence := NewSequence("a sequence")
	for _, s := range values {
		sequence.AppendChild(NewString(ClassUniversal, TypePrimitive, TagOctetString, s, "String"))
	}

	if len(sequence.Children) != len(values) {
		t.Errorf("wrong length for children array should be %d, got %d", len(values), len(sequence.Children))
	}

	encodedSequence := sequence.Bytes()

	decodedSequence := DecodePacket(encodedSequence)
	if len(decodedSequence.Children) != len(values) {
		t.Errorf("wrong length for children array should be %d => %d", len(values), len(decodedSequence.Children))
	}

	for i, s := range values {
		if decodedSequence.Children[i].Value.(string) != s {
			t.Errorf("expected %d to be %q, got %q", i, s, decodedSequence.Children[i].Value.(string))
		}
	}
}

func TestReadPacket(t *testing.T) {
	packet := NewString(ClassUniversal, TypePrimitive, TagOctetString, "Ad impossibilia nemo tenetur", "string")
	var buffer io.ReadWriter = new(bytes.Buffer)

	if _, err := buffer.Write(packet.Bytes()); err != nil {
		t.Error("error writing packet", err)
	}

	newPacket, err := ReadPacket(buffer)
	if err != nil {
		t.Error("error during ReadPacket", err)
	}
	newPacket.ByteValue = nil
	if !bytes.Equal(newPacket.ByteValue, packet.ByteValue) {
		t.Error("packets should be the same")
	}
}

func TestBinaryInteger(t *testing.T) {
	// data src : http://luca.ntop.org/Teaching/Appunti/asn1.html 5.7
	data := []struct {
		v int64
		e []byte
	}{
		{v: 0, e: []byte{0x02, 0x01, 0x00}},
		{v: 127, e: []byte{0x02, 0x01, 0x7F}},
		{v: 128, e: []byte{0x02, 0x02, 0x00, 0x80}},
		{v: 256, e: []byte{0x02, 0x02, 0x01, 0x00}},
		{v: -128, e: []byte{0x02, 0x01, 0x80}},
		{v: -129, e: []byte{0x02, 0x02, 0xFF, 0x7F}},
		{v: math.MaxInt64, e: []byte{0x02, 0x08, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
		{v: math.MinInt64, e: []byte{0x02, 0x08, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}

	for _, d := range data {
		if b := NewInteger(ClassUniversal, TypePrimitive, TagInteger, d.v, "").Bytes(); !bytes.Equal(d.e, b) {
			t.Errorf("Wrong binary generated for %d : got % X, expected % X", d.v, b, d.e)
		}
	}
}

func TestBinaryOctetString(t *testing.T) {
	// data src : http://luca.ntop.org/Teaching/Appunti/asn1.html 5.10

	if !bytes.Equal([]byte{0x04, 0x08, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}, NewString(ClassUniversal, TypePrimitive, TagOctetString, "\x01\x23\x45\x67\x89\xab\xcd\xef", "").Bytes()) {
		t.Error("wrong binary generated")
	}
}

// buff is an alias to build a bytes.Reader from an explicit sequence of bytes
func buff(bs ...byte) *bytes.Reader {
	return bytes.NewReader(bs)
}

func TestEOF(t *testing.T) {
	_, err := ReadPacket(buff())
	if err != io.EOF {
		t.Errorf("empty buffer: expected EOF, got %s", err)
	}

	// testCases for EOF
	testCases := []struct {
		name string
		buf  *bytes.Reader
	}{
		{"primitive", buff(0x04, 0x0a, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9)},
		{"constructed", buff(0x30, 0x06, 0x02, 0x01, 0x01, 0x02, 0x01, 0x02)},
		{"constructed indefinite length", buff(0x30, 0x80, 0x02, 0x01, 0x01, 0x02, 0x01, 0x02, 0x00, 0x00)},
	}
	for _, tc := range testCases {
		_, err := ReadPacket(tc.buf)
		if err != nil {
			t.Errorf("%s: expected no error, got %s", tc.name, err)
		}

		_, err = ReadPacket(tc.buf)
		if err != io.EOF {
			t.Errorf("%s: expected EOF, got %s", tc.name, err)
		}
	}

	// testCases for UnexpectedEOF :
	testCases = []struct {
		name string
		buf  *bytes.Reader
	}{
		{"truncated tag", buff(0x1f, 0xff)},
		{"tag and no length", buff(0x04)},
		{"truncated length", buff(0x04, 0x82, 0x02)},
		{"header with no content", buff(0x04, 0x0a)},
		{"header with truncated content", buff(0x04, 0x0a, 0, 1, 2)},

		{"constructed missing content", buff(0x30, 0x06)},
		{"constructed only first child", buff(0x30, 0x06, 0x02, 0x01, 0x01)},
		{"constructed truncated", buff(0x30, 0x06, 0x02, 0x01, 0x01, 0x02, 0x01)},

		{"indefinite missing eoc", buff(0x30, 0x80, 0x02, 0x01, 0x01, 0x02, 0x01, 0x02)},
		{"indefinite truncated eoc", buff(0x30, 0x80, 0x02, 0x01, 0x01, 0x02, 0x01, 0x02, 0x00)},
	}
	for _, tc := range testCases {
		_, err := ReadPacket(tc.buf)
		if err != io.ErrUnexpectedEOF {
			t.Errorf("%s: expected UnexpectedEOF, got %s", tc.name, err)
		}
	}
}

// buildNestedSequence builds a BER-encoded SEQUENCE nested depth levels deep.
// Each level wraps the previous; the innermost is an empty SEQUENCE.
// Only valid for small depth where each length fits in one byte (<= 127 content bytes).
func buildNestedSequence(depth int) []byte {
	inner := []byte{0x30, 0x00} // empty SEQUENCE
	for i := 1; i < depth; i++ {
		wrapped := make([]byte, 0, 2+len(inner))
		wrapped = append(wrapped, 0x30, byte(len(inner)))
		wrapped = append(wrapped, inner...)
		inner = wrapped
	}
	return inner
}

func TestLongFormLengthSentinelCollision(t *testing.T) {
	// A definite-form length of 2^64-1 must not be reinterpreted as indefinite form.
	data := []byte{
		0x30,                                                 // SEQUENCE, constructed
		0x88, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // "length" = 2^64-1
		0x05, 0x00, // NULL child
		0x00, 0x00, // EOC
	}
	if _, err := DecodePacketErr(data); err == nil {
		t.Error("expected error for definite-form length 2^64-1, got nil")
	}
}

func TestMaxNestingDepth(t *testing.T) {
	old := MaxNestingDepth
	defer func() { MaxNestingDepth = old }()

	MaxNestingDepth = 5

	// depth=5: outermost at depth 0, innermost at depth 4 — should succeed
	data5 := buildNestedSequence(5)
	_, err := DecodePacketErr(data5)
	if err != nil {
		t.Errorf("5 levels with MaxNestingDepth=5: expected success, got %v", err)
	}

	// depth=6: requires depth 5 which equals MaxNestingDepth — should fail
	data6 := buildNestedSequence(6)
	_, err = DecodePacketErr(data6)
	if err == nil {
		t.Error("6 levels with MaxNestingDepth=5: expected error, got nil")
	}
}

func TestMaxNestingDepthUnlimited(t *testing.T) {
	old := MaxNestingDepth
	defer func() { MaxNestingDepth = old }()

	MaxNestingDepth = 0

	// 50-deep nesting with single-byte lengths — should succeed with no limit
	data := buildNestedSequence(50)
	_, err := DecodePacketErr(data)
	if err != nil {
		t.Errorf("50 levels with MaxNestingDepth=0: unexpected error %v", err)
	}
}

func TestMaxNestingDepthReadPacket(t *testing.T) {
	old := MaxNestingDepth
	defer func() { MaxNestingDepth = old }()

	MaxNestingDepth = 5

	// depth=5: should succeed
	data5 := buildNestedSequence(5)
	_, err := ReadPacket(bytes.NewReader(data5))
	if err != nil {
		t.Errorf("5 levels with MaxNestingDepth=5: expected success, got %v", err)
	}

	// depth=6: should fail
	data6 := buildNestedSequence(6)
	_, err = ReadPacket(bytes.NewReader(data6))
	if err == nil {
		t.Error("6 levels with MaxNestingDepth=5: expected error, got nil")
	}
}

func TestMaxNestingDepthReadPacketUnlimited(t *testing.T) {
	old := MaxNestingDepth
	defer func() { MaxNestingDepth = old }()

	MaxNestingDepth = 0

	data := buildNestedSequence(50)
	_, err := ReadPacket(bytes.NewReader(data))
	if err != nil {
		t.Errorf("50 levels with MaxNestingDepth=0: unexpected error %v", err)
	}
}

// nulls returns n bytes worth of NULL (0x05 0x00) primitives. n must be even.
func nulls(n int) []byte {
	if n%2 != 0 {
		panic("nulls: n must be even")
	}
	buf := make([]byte, 0, n)
	for i := 0; i < n/2; i++ {
		buf = append(buf, 0x05, 0x00)
	}
	return buf
}

// buildDefinite wraps content in a definite-length packet with the given identifier byte.
func buildDefinite(tag byte, content []byte) []byte {
	out := append([]byte{tag}, encodeLength(len(content))...)
	return append(out, content...)
}

// buildIndefinite wraps content in an indefinite-length packet terminated by an EOC marker.
// The content size seen by the parser is len(content) + 2, counting the EOC.
func buildIndefinite(tag byte, content []byte) []byte {
	out := append([]byte{tag, 0x80}, content...)
	return append(out, 0x00, 0x00) // EOC
}

// buildConstructedIndefinite builds an indefinite-length SEQUENCE containing children
// NULL (0x05 0x00) primitives, terminated by an EOC marker.
func buildConstructedIndefinite(children int) []byte {
	return buildIndefinite(0x30, nulls(children*2))
}

// buildNestedIndefinite builds depth nested indefinite-length SEQUENCEs, each level holding
// perLevel NULL primitives alongside the next level down. Spreading the content across every
// level means no single level exceeds a limit that the tree as a whole does.
func buildNestedIndefinite(depth, perLevel int) []byte {
	inner := buildIndefinite(0x30, nulls(perLevel*2))
	for i := 1; i < depth; i++ {
		inner = buildIndefinite(0x30, append(nulls(perLevel*2), inner...))
	}
	return inner
}

func TestMaxPacketLengthBytes(t *testing.T) {
	old := MaxPacketLengthBytes
	defer func() { MaxPacketLengthBytes = old }()

	decoders := map[string]func([]byte) (*Packet, error){
		"DecodePacketErr": DecodePacketErr,
		"ReadPacket":      func(b []byte) (*Packet, error) { return ReadPacket(bytes.NewReader(b)) },
	}

	tests := []struct {
		name         string
		limit        int64
		data         []byte
		wantChildren int // ignored when wantErr
		wantErr      bool
	}{
		// Indefinite length carries no declared bound; only the aggregate byte count catches it.
		{"indefinite over limit", 1024, buildConstructedIndefinite(2000), 0, true},
		// SEQUENCE declaring definite long-form length 0x2000 (8192) — rejected up front.
		{"definite over limit", 1024, []byte{0x30, 0x82, 0x20, 0x00}, 0, true},
		{"within limit", 1024, buildConstructedIndefinite(10), 10, false},
		{"unlimited", 0, buildConstructedIndefinite(5000), 5000, false},

		// The limit counts content bytes, so a packet declaring exactly the limit is accepted
		// and one declaring a single byte more is not — identically for both length forms.
		{"definite primitive at limit", 1024, buildDefinite(0x04, make([]byte, 1024)), 0, false},
		{"definite primitive over limit", 1024, buildDefinite(0x04, make([]byte, 1025)), 0, true},
		{"definite constructed at limit", 1024, buildDefinite(0x30, nulls(1024)), 512, false},
		{"definite constructed over limit", 1024, buildDefinite(0x30, nulls(1026)), 0, true},
		// 1022 content bytes plus the 2-byte EOC lands exactly on the limit.
		{"indefinite at limit", 1024, buildIndefinite(0x30, nulls(1022)), 511, false},
		{"indefinite just over limit", 1024, buildIndefinite(0x30, nulls(1024)), 0, true},
	}

	for _, tt := range tests {
		for decName, decode := range decoders {
			t.Run(tt.name+"/"+decName, func(t *testing.T) {
				MaxPacketLengthBytes = tt.limit
				p, err := decode(tt.data)
				switch {
				case tt.wantErr && err == nil:
					t.Fatal("expected error, got nil")
				case !tt.wantErr && err != nil:
					t.Fatalf("unexpected error: %v", err)
				case !tt.wantErr && len(p.Children) != tt.wantChildren:
					t.Errorf("expected %d children, got %d", tt.wantChildren, len(p.Children))
				}
			})
		}
	}
}

func TestMaxPacketLengthBytesNested(t *testing.T) {
	old := MaxPacketLengthBytes
	defer func() { MaxPacketLengthBytes = old }()

	decoders := map[string]func([]byte) (*Packet, error){
		"DecodePacketErr": DecodePacketErr,
		"ReadPacket":      func(b []byte) (*Packet, error) { return ReadPacket(bytes.NewReader(b)) },
	}

	// A definite SEQUENCE understating its length, followed by an oversized indefinite child.
	// The outer header passes the limit check, so only the child's own accounting can stop it.
	understatedOuter := append([]byte{0x30, 0x64}, buildIndefinite(0x30, nulls(1100))...)

	// Definite children each well under the limit, but summing past it inside an indefinite parent.
	var definiteChildren []byte
	for i := 0; i < 4; i++ {
		definiteChildren = append(definiteChildren, buildDefinite(0x04, make([]byte, 400))...)
	}

	tests := []struct {
		name    string
		limit   int64
		data    []byte
		wantErr bool
	}{
		// No level holds more than ~200 bytes of its own, so this only trips if a child's size
		// propagates up into its ancestors' accounting.
		{"aggregate across levels", 1024, buildNestedIndefinite(8, 100), true},
		{"oversized indefinite inside definite", 1024, understatedOuter, true},
		{"definite children summing over limit", 1024, buildIndefinite(0x30, definiteChildren), true},
		{"unlimited", 0, buildNestedIndefinite(8, 100), false},
	}

	for _, tt := range tests {
		for decName, decode := range decoders {
			t.Run(tt.name+"/"+decName, func(t *testing.T) {
				MaxPacketLengthBytes = tt.limit
				_, err := decode(tt.data)
				switch {
				case tt.wantErr && err == nil:
					t.Fatal("expected error, got nil")
				case tt.wantErr && !strings.Contains(err.Error(), "greater than maximum"):
					// Guard against passing for the wrong reason, e.g. a truncation error
					t.Fatalf("expected a length limit error, got %v", err)
				case !tt.wantErr && err != nil:
					t.Fatalf("unexpected error: %v", err)
				}
			})
		}
	}
}

func TestNestedIndefiniteWithinLimit(t *testing.T) {
	old := MaxPacketLengthBytes
	defer func() { MaxPacketLengthBytes = old }()

	MaxPacketLengthBytes = 1024

	const depth, perLevel = 4, 5

	p, err := DecodePacketErr(buildNestedIndefinite(depth, perLevel))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Every level but the innermost holds perLevel NULLs plus the nested SEQUENCE
	for level := 1; level <= depth; level++ {
		want := perLevel + 1
		if level == depth {
			want = perLevel
		}
		if len(p.Children) != want {
			t.Fatalf("level %d: expected %d children, got %d", level, want, len(p.Children))
		}
		if level < depth {
			p = p.Children[perLevel]
		}
	}
}
