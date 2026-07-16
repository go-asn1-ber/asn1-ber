package ber

import "io"

func readByte(reader io.Reader) (byte, error) {
	if br, ok := reader.(io.ByteReader); ok {
		return br.ReadByte()
	}
	var b [1]byte
	if _, err := io.ReadFull(reader, b[:]); err != nil {
		return 0, err
	}
	return b[0], nil
}

func unexpectedEOF(err error) error {
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	return err
}

func isEOCPacket(p *Packet) bool {
	return p != nil &&
		p.Tag == TagEOC &&
		p.ClassType == ClassUniversal &&
		p.TagType == TypePrimitive &&
		len(p.ByteValue) == 0 &&
		len(p.Children) == 0
}
