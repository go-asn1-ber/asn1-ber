package ber

import "testing"

func TestNewIntegerErr(t *testing.T) {
	if _, err := NewIntegerErr(ClassUniversal, TypePrimitive, TagInteger, 42, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := NewIntegerErr(ClassUniversal, TypePrimitive, TagInteger, "nope", ""); err == nil {
		t.Fatal("expected error for non-integer value")
	}
}

func TestNewRealErr(t *testing.T) {
	if _, err := NewRealErr(ClassUniversal, TypePrimitive, TagRealFloat, 1.5, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := NewRealErr(ClassUniversal, TypePrimitive, TagRealFloat, 3, ""); err == nil {
		t.Fatal("expected error for non-float value")
	}
}

func TestNewOIDErr(t *testing.T) {
	p, err := NewOIDErr(ClassUniversal, TypePrimitive, TagObjectIdentifier, "1.2.840.113549", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	decoded, err := DecodePacketErr(p.Bytes())
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if decoded.Value != "1.2.840.113549" {
		t.Fatalf("round-trip mismatch: got %v", decoded.Value)
	}

	// Structurally invalid OID string: an error, not a panic.
	if _, err := NewOIDErr(ClassUniversal, TypePrimitive, TagObjectIdentifier, "1", ""); err == nil {
		t.Fatal("expected error for structurally invalid OID")
	}
	// Non-string value: an error, not a panic.
	if _, err := NewOIDErr(ClassUniversal, TypePrimitive, TagObjectIdentifier, 123, ""); err == nil {
		t.Fatal("expected error for non-string value")
	}
}

// TestNewOIDNoPanicOnInvalidString confirms the deprecated NewOID returns nil
// (rather than panicking) on a structurally invalid OID string, now that
// encodeOID returns an error instead of panicking.
func TestNewOIDNoPanicOnInvalidString(t *testing.T) {
	if p := NewOID(ClassUniversal, TypePrimitive, TagObjectIdentifier, "1", ""); p != nil {
		t.Fatalf("expected nil for invalid OID string, got %v", p)
	}
}

func TestNewRelativeOIDErr(t *testing.T) {
	if _, err := NewRelativeOIDErr(ClassUniversal, TypePrimitive, TagRelativeOID, "8571.3.2", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := NewRelativeOIDErr(ClassUniversal, TypePrimitive, TagRelativeOID, "not.a.number", ""); err == nil {
		t.Fatal("expected error for invalid relative OID")
	}
}

// TestDecodePacketErrNilOnContentError confirms a content-parse failure yields
// a nil packet alongside the error, never a partially-populated one. The bytes
// are a UTF8String (0x0c) of length 1 containing invalid UTF-8 (0xff).
func TestDecodePacketErrNilOnContentError(t *testing.T) {
	p, err := DecodePacketErr([]byte{0x0c, 0x01, 0xff})
	if err == nil {
		t.Fatal("expected error for invalid UTF-8")
	}
	if p != nil {
		t.Fatalf("expected nil packet on error, got %v", p)
	}
}
