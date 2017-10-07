package ber

import (
	"errors"
	"io"
)

func readHeader(reader io.Reader) (identifier Identifier, length int64, read int64, err error) {
	if i, c, err := readIdentifier(reader); err != nil {
		return Identifier{}, 0, read, err
	} else {
		identifier = i
		read += int64(c)
	}

	if l, c, err := readLength(reader); err != nil {
		return Identifier{}, 0, read, err
	} else {
		length = l
		read += int64(c)
	}

	// Validate length type with identifier (x.600, 8.1.3.2.a)
	if length == LengthIndefinite && identifier.TagType == TypePrimitive {
		return Identifier{}, 0, read, errors.New("indefinite length used with primitive type")
	}

	return identifier, length, read, nil
}
