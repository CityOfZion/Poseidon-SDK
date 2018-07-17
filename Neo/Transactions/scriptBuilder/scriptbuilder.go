// Copied and modified from https://github.com/CityOfZion/neo-go/blob/master/pkg/vm/emit.go

package scriptBuilder

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"math/big"
	sliceUtils "multicrypt/Utils/SliceUtils"
	writer "multicrypt/Utils/Writer"
)

type Script struct {
	Operation  string
	Args       []interface{}
	ScriptHash string // This is the scriptHash of contract we want to invoke
	gasCost    GasCost
}

func (s *Script) Bytes(w io.Writer) { //This will serialise the script, alone and not the Gas Cost

	s.EmitArr(w, s.Args)
	s.EmitString(w, s.Operation)
	s.EmitAppCall(w)
	s.EmitScriptHash(w, s.ScriptHash)

}

// TODO: use the gas cost file, to add up the gas cost whenever an operation is done
func (s *Script) GasCost() int64 {
	return int64(0)
}
func (s *Script) SerialiseScript(w io.Writer) {
	// TODO: optimise this, we do not need to do s.Bytes both times.
	b := new(bytes.Buffer)
	s.Bytes(b)

	writer.WriteVarUint(w, uint64(len(b.Bytes())))
	s.Bytes(w)
	binary.Write(w, binary.LittleEndian, s.GasCost()) // gas cost should already be in fixed8 format
}

func (s *Script) GetLength() uint64 {
	b := new(bytes.Buffer)
	s.Bytes(b)

	return uint64(len(b.Bytes()))
}

func (s *Script) Hex() string {
	b := new(bytes.Buffer)
	s.Bytes(b)
	return hex.EncodeToString(b.Bytes())
}

func (s *Script) SetOperation(operation string) {
	s.Operation = operation
}
func (s *Script) SetScriptHash(scriptHash string) {
	s.ScriptHash = scriptHash
}
func (s *Script) AddArgument(arg interface{}) {
	s.Args = append(s.Args, arg)
}

func (s *Script) EmitArr(w io.Writer, arr []interface{}) {
	if len(arr) == 0 || arr == nil {
		binary.Write(w, binary.LittleEndian, byte(PUSHF))
		binary.Write(w, binary.LittleEndian, byte(PACK))
		return
	}
	for _, element := range arr {
		switch element.(type) {
		case int:
			if value, isInt := element.(int); isInt {
				s.EmitInt(w, int64(value))
			}
		case int8:
			if value, isInt8 := element.(int8); isInt8 {
				s.EmitInt(w, int64(value))
			}
		case int16:
			if value, isInt16 := element.(int16); isInt16 {
				s.EmitInt(w, int64(value))
			}
		case int32:
			if value, isInt32 := element.(int32); isInt32 {
				s.EmitInt(w, int64(value))
			}
		case int64:
			if value, isInt64 := element.(int64); isInt64 {

				s.EmitInt(w, value)
			}
		case float64:
			if value, isFloat64 := element.(float64); isFloat64 {
				s.EmitFloat(w, value)
			}
		case string:
			if value, isString := element.(string); isString {
				s.EmitString(w, value)
			}
		case bool:
			if value, isBool := element.(bool); isBool {
				s.EmitBool(w, value)
			}
		case []byte:
			if value, isByteArr := element.([]byte); isByteArr {
				s.EmitBytes(w, value)
			}
		case []interface{}: // an array
			if value, isArr := element.([]interface{}); isArr {
				s.EmitArr(w, value)
			}
		default:
			// throw error
		}
	}
	if len(arr) <= 16 {
		i := len(arr)
		s.EmitIntLiteral(w, int64(i))
		// binary.Write(w, binary.LittleEndian, byte(len(arr))+byte(0x50))
	} else {
		s.EmitInt(w, int64(len(arr)))
		//TODO: Add Two numbers together using OP_ADD or OP_MULTIPLY
		// use mod16 to find the best possible combination
	}

	binary.Write(w, binary.LittleEndian, byte(PACK))

}

func (s *Script) EmitFloat(w io.Writer, f float64) {
	i := int64(f * 100000000)
	s.EmitInt(w, i)
}

func (s *Script) EmitBool(w io.Writer, isTrue bool) {
	if isTrue {
		s.EmitOpCode(w, PUSHT)
		return
	}
	s.EmitOpCode(w, PUSHF)
}

func (s *Script) EmitString(w io.Writer, str string) {

	s.EmitBytes(w, []byte(str)) // done

}

func (s *Script) EmitOpCode(w io.Writer, op OpCode) {
	binary.Write(w, binary.LittleEndian, byte(op))

}

func (s *Script) EmitInt(w io.Writer, i int64) { // Should aready be in fixed 8 here

	bInt := big.NewInt(i)

	val := sliceUtils.Reverse(bInt.Bytes())
	s.EmitBytes(w, val)

}

// This is to emit the literal values of numbers. E.g. input = 10, output =PUSH10 upto PUSH16
func (s *Script) EmitIntLiteral(w io.Writer, i int64) {
	if i == -1 {

		s.EmitOpCode(w, PUSHM1)

	} else if i == 0 {

		s.EmitOpCode(w, PUSHF)

	} else if i > 0 && i <= 16 {
		val := OpCode((int(PUSH1) - 1 + int(i)))
		s.EmitOpCode(w, val)
	} else {
		s.EmitInt(w, i)
	}
}

func (s *Script) EmitBytes(w io.Writer, bAppend []byte) {

	var n = len(bAppend)

	if n <= int(PUSHBYTES75) {
		binary.Write(w, binary.LittleEndian, byte(OpCode(n)))
		binary.Write(w, binary.LittleEndian, bAppend)
	} else if n < 0x100 {

		binary.Write(w, binary.LittleEndian, byte(PUSHDATA1))
		binary.Write(w, binary.LittleEndian, byte(n))
		binary.Write(w, binary.LittleEndian, bAppend)

	} else if n < 0x10000 {

		buf := make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, uint16(n))
		binary.Write(w, binary.LittleEndian, byte(PUSHDATA2))
		binary.Write(w, binary.LittleEndian, buf)
		binary.Write(w, binary.LittleEndian, bAppend)

	} else {

		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(n))

		binary.Write(w, binary.LittleEndian, byte(PUSHDATA4))
		binary.Write(w, binary.LittleEndian, buf)
		binary.Write(w, binary.LittleEndian, bAppend)
	}

}

func (s *Script) EmitScriptHash(w io.Writer, script string) {

	scriptBytes, _ := hex.DecodeString(script)
	val := sliceUtils.Reverse(scriptBytes)
	binary.Write(w, binary.LittleEndian, val)

}

func (s *Script) EmitAppCall(w io.Writer) {

	s.EmitOpCode(w, APPCALL)

}
