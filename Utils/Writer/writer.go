package util
// Taken from anthdm
import (
	"encoding/binary"
	"errors"
	"io"
)

func WriteVarUint(w io.Writer, val uint64) error {
	if val < 0 {
		return errors.New("value out of range")
	}
	if val < 0xfd {
		binary.Write(w, binary.LittleEndian, uint8(val))
		return nil
	}
	if val < 0xFFFF {
		binary.Write(w, binary.LittleEndian, byte(0xfd))
		binary.Write(w, binary.LittleEndian, uint16(val))
		return nil
	}
	if val < 0xFFFFFFFF {
		binary.Write(w, binary.LittleEndian, byte(0xfe))
		binary.Write(w, binary.LittleEndian, uint32(val))
		return nil
	}

	binary.Write(w, binary.LittleEndian, byte(0xff))
	binary.Write(w, binary.LittleEndian, val)

	return nil
}
