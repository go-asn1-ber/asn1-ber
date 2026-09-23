package ber

import (
	"bytes"
	"testing"
)

// TestLazyDeepNestingRoundTrip builds a deeply-nested constructed packet and
// verifies it round-trips. Before lazy serialization each ancestor retained a
// full copy of its subtree (O(depth x subtree) memory); this exercises the
// depth path and confirms the encoding is unchanged.
func TestLazyDeepNestingRoundTrip(t *testing.T) {
	const depth = 500

	root := NewSequence("root")
	cur := root
	for i := 0; i < depth; i++ {
		child := NewSequence("nested")
		cur.AppendChild(child)
		cur = child
	}
	cur.AppendChild(NewString(ClassUniversal, TypePrimitive, TagOctetString, "leaf", ""))

	encoded := root.Bytes()

	decoded, err := DecodePacketErr(encoded)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if !bytes.Equal(decoded.Bytes(), encoded) {
		t.Fatalf("round-trip mismatch: re-encoded bytes differ from source")
	}
}

// TestLazyAppendChildReflectsLaterMutation verifies Bytes() reflects a child
// appended after an intermediate AppendChild. The previous eager-buffering
// implementation snapshotted the child into the parent's Data at AppendChild
// time and would have missed the later grandchild.
func TestLazyAppendChildReflectsLaterMutation(t *testing.T) {
	parent := NewSequence("parent")
	child := NewSequence("child")
	parent.AppendChild(child)

	// Append to the child after it was already appended to the parent.
	child.AppendChild(NewString(ClassUniversal, TypePrimitive, TagOctetString, "late", ""))

	want := NewSequence("parent")
	wantChild := NewSequence("child")
	wantChild.AppendChild(NewString(ClassUniversal, TypePrimitive, TagOctetString, "late", ""))
	want.AppendChild(wantChild)

	if got := parent.Bytes(); !bytes.Equal(got, want.Bytes()) {
		t.Fatalf("parent.Bytes() did not reflect a child appended after AppendChild\n got: %x\nwant: %x", got, want.Bytes())
	}
}
