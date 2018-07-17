// General components of a transaction in neo
package transaction

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	neo "multicrypt/Neo"
	"multicrypt/Utils/SliceUtils"
	uint160 "multicrypt/Utils/Uint160"
	uint256 "multicrypt/Utils/Uint256"
	writer "multicrypt/Utils/Writer"
	crypto "multicrypt/crypto"
	rfc "multicrypt/crypto/Signing"
	"multicrypt/crypto/bip32"
	"multicrypt/crypto/elliptic"
	"sort"
)

type Input struct {
	PrevHash  string
	PrevIndex uint16
	Address   string
	External  int
	Index     int // this is the derivation integer/path for this address
	// This should be saved in phone db or have a retrieve function for it
}
type Output struct {
	AssetId string
	Value   int64
	Address string
}

type Attributes struct {
	Usage uint16
	Data  []byte
}

type Witness struct {
	VerificationScript []byte
	InvocationScript   []byte
}

const (
	AttributeScript = 0x20
	AttributeRemark = 0xf0
)

type fn func(w io.Writer) error
type BasicTransaction struct {
	Type       uint8
	Version    uint8
	Attributes []Attributes
	Inputs     []Input
	Outputs    []Output
	Witnesses  []Witness
	SystemFee  float32
	NetworkFee float32
	privKeys   []bip32.Key
	F          func(w io.Writer) error
}

func (c *BasicTransaction) AddScriptAttributes() {
	// Add the empty remark so that, we do not need to send GAS when doing a transfer
	for _, input := range c.Inputs {
		scriptHash := AddressToScriptHash(input.Address)

		scriptHashAsBytes, err := hex.DecodeString(scriptHash)
		if err != nil {

			fmt.Println("Err with Input Scripthash conversion")
			return
		}
		fmt.Println(scriptHash, scriptHashAsBytes)
		c.Attributes = append(c.Attributes, Attributes{uint16(0x20), scriptHashAsBytes})
	}
}
func (c *BasicTransaction) AddAttributes(usage uint16, data []byte) {

	c.Attributes = append(c.Attributes, Attributes{usage, data})

}

func (c *BasicTransaction) AddOutput(assetId string, value int64, address string) {
	c.Outputs = append(c.Outputs, Output{assetId, value, address})
}

func (c *BasicTransaction) AddInput(prevHash string, prevIndex uint16, Address string, External, Index int) {
	c.Inputs = append(c.Inputs, Input{prevHash, prevIndex, Address, External, Index})
}

func (c *BasicTransaction) SetSystemFee(value float32) {
	c.SystemFee = 0.0
	// Contract system transaction fees are always zero
}

func (c *BasicTransaction) SetNetworkFee(value float32) {
	c.NetworkFee = value
}

type byScriptHash []Input

func (s byScriptHash) Len() int {
	return len(s)
}
func (s byScriptHash) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byScriptHash) Less(i, j int) bool {
	script1 := (AddressToScriptHash(s[i].Address))
	script2 := (AddressToScriptHash(s[j].Address))
	sh1U160, _ := uint160.Uint160DecodeString(script1)
	sh2U160, _ := uint160.Uint160DecodeString(script2)

	for i := 0; i < len(sh1U160); i++ {
		if sh1U160.BytesReverse()[i] < sh2U160.BytesReverse()[i] {
			return true
		} else if sh1U160.BytesReverse()[i] > sh2U160.BytesReverse()[i] {
			return false
		}
	}
	return true
}

func (c *BasicTransaction) OrderInputsByAddressHash() {
	sort.Slice(c.Inputs, func(i, j int) bool {
		scriptHash1 := AddressToScriptHash(c.Inputs[i].Address)
		scriptHash2 := AddressToScriptHash(c.Inputs[j].Address)

		sh1U160, _ := uint160.Uint160DecodeString(scriptHash1)
		sh2U160, _ := uint160.Uint160DecodeString(scriptHash2)

		for i := len(sh1U160.Bytes()) - 1; i >= 0; i-- {
			if sh1U160.Bytes()[i] < sh2U160.Bytes()[i] {
				return true
			} else if sh1U160.Bytes()[i] > sh2U160.Bytes()[i] {
				return false
			}
		}
		return true
	})

}

// func (c *BasicTransaction) OrderInputBySignature(){
// 	sort.Slice(c.Witnesses, func(i, j int) bool {
// 		return c.Witnesses[i]. < c.Witnesses[j]
// 	})
// }

// use the seed along with the path/index in each input to sign the transaction
func (c *BasicTransaction) SignTransaction(seed string) {

	sort.Sort(byScriptHash(c.Inputs))

	for _, input := range c.Inputs {

		c.AddPrivateKey(seed, input.External, input.Index)
	}

	lenPubKey := []byte{33}
	opCodeCheckSig := []byte{0xac}
	lenInvocScript := []byte{64}

	curve := elliptic.ChosenCurve
	curve.SetCurveSecp256r1()

	encountered := map[string]bool{}

	for _, privKey := range c.privKeys {

		//Check for duplicate. multiple inputs with same address will give duplicates
		key := hex.EncodeToString(privKey.Key)
		if encountered[key] {
			continue
		}
		encountered[key] = true

		verificationScript := append(privKey.PublicKey().Key, opCodeCheckSig...)
		verificationScript = append(lenPubKey, verificationScript...)
		buf := new(bytes.Buffer)
		err := c.SerialiseTransaction(buf, false)
		if err != nil {
			fmt.Println("Error with serialisation")
			return
		}
		signedTransaction, err := rfc.Sign(curve, buf.Bytes(), privKey)
		if err != nil {
			fmt.Println("Err with signing trans")
			return

		}

		invocationScript := append(lenInvocScript, signedTransaction...)
		c.Witnesses = append(c.Witnesses, Witness{verificationScript, invocationScript})
	}

}

func (c *BasicTransaction) AddPrivateKey(seed string, external, index int) {
	neo := neo.Coin{}
	c.privKeys = append(c.privKeys, *neo.GeneratePrivateKey(seed, external, index))
}

func (c *BasicTransaction) SerialiseTransaction(buf *bytes.Buffer, signed bool) error {

	// Order: Type, Version, special field, attributes, inputs, outputs, scripts

	if err := c.SerialiseType(buf); err != nil {
		return err
	}
	fmt.Println("Done 1")
	if err := c.SerialiseVersion(buf); err != nil {
		return err
	}
	fmt.Println("Done 2")

	// No special types to serialise. Make a serialise transaction for each type? Or we could use the type
	// and have each transaction conform to an interface which will serialise.
	if err := c.F(buf); err != nil {
		return err
	}
	fmt.Println("Done 3")
	if err := c.SerialiseAttributes(buf); err != nil {
		return err
	}
	fmt.Println("Done 4")

	if err := c.SerialiseInputs(buf); err != nil {
		return err
	}
	fmt.Println("Done 5")

	if err := c.SerialiseOutputs(buf); err != nil {
		return err
	}
	fmt.Println("Done 6")

	// if false, do not serialise scripts
	if signed && len(c.Witnesses) > 0 {
		if err := c.SerialiseScripts(buf); err != nil {
			return err
		}
		fmt.Println("Done 7")
	}
	return nil
}
func (c *BasicTransaction) SerialiseType(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, uint8(c.Type)); err != nil {
		return err
	}
	return nil
}

func (c *BasicTransaction) SerialiseVersion(w io.Writer) error {

	if err := binary.Write(w, binary.LittleEndian, c.Version); err != nil {
		return err
	}
	return nil
}

func (i *Input) Encode(w io.Writer) error {
	prevHash, err := uint256.Uint256DecodeString(i.PrevHash)
	if err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, prevHash); err != nil {
		return err
	}
	prevIndex := i.PrevIndex
	if err := binary.Write(w, binary.LittleEndian, prevIndex); err != nil {
		return err
	}
	return nil
}
func (c *BasicTransaction) SerialiseInputs(w io.Writer) error {
	if err := writer.WriteVarUint(w, uint64(len(c.Inputs))); err != nil {
		return err
	}

	for _, input := range c.Inputs {
		if err := input.Encode(w); err != nil {
			return err
		}
	}

	return nil
}

func (o *Output) Encode(w io.Writer) error {
	assetID, err := uint256.Uint256DecodeString((o.AssetId))
	scriptHash, err := uint160.Uint160DecodeString(AddressToScriptHash(o.Address))
	err = binary.Write(w, binary.LittleEndian, assetID)
	err = binary.Write(w, binary.LittleEndian, o.Value)
	err = binary.Write(w, binary.LittleEndian, scriptHash)
	if err != nil {
		return err
	}
	return nil
}
func (c *BasicTransaction) SerialiseOutputs(w io.Writer) error {

	if err := writer.WriteVarUint(w, uint64(len(c.Outputs))); err != nil {
		return err
	}

	for _, output := range c.Outputs {
		if err := output.Encode(w); err != nil {
			return err
		}
	}
	return nil

}
func (attribute *Attributes) Encode(w io.Writer) error {
	maximumSizeOfASingleAttribute := 65535 // TODO: Find out how to organise constants
	dataByte := attribute.Data
	lenData := len(dataByte)
	var err error
	if lenData > maximumSizeOfASingleAttribute {

		return errors.New("Max Data Reached")
	}

	if err = binary.Write(w, binary.LittleEndian, uint8(attribute.Usage)); err != nil { // uint8 or it will write it as 2000 for script for example
		return err
	}

	if attribute.Usage == 0x81 || (attribute.Usage >= 0xa1 && attribute.Usage <= 0xaf) {
		err = binary.Write(w, binary.LittleEndian, dataByte[:32])
	} else if attribute.Usage == 0x02 || attribute.Usage == 0x03 { //ECDH Series
		err = binary.Write(w, binary.LittleEndian, dataByte[1:33])
	} else if attribute.Usage == 0x80 || attribute.Usage == 0x81 || attribute.Usage == 0x90 || attribute.Usage >= 0xf0 {
		//TODO: Change these to proper names and Description Length is actually uint8
		err = writer.WriteVarUint(w, uint64(lenData))
		err = binary.Write(w, binary.LittleEndian, dataByte)
	} else {
		err = binary.Write(w, binary.LittleEndian, dataByte)
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *BasicTransaction) SerialiseAttributes(w io.Writer) error {

	lenAttrs := uint64(len(c.Attributes))
	if err := writer.WriteVarUint(w, lenAttrs); err != nil {
		return err
	}

	for _, attribute := range c.Attributes {
		if err := attribute.Encode(w); err != nil {
			return err
		}
	}
	return nil

}
func (witness *Witness) Encode(w io.Writer) error {

	lenInvo := uint64(len(witness.InvocationScript))
	err := writer.WriteVarUint(w, lenInvo)

	err = binary.Write(w, binary.LittleEndian, witness.InvocationScript)

	lenVerif := uint64(len(witness.VerificationScript))
	err = writer.WriteVarUint(w, lenVerif)
	err = binary.Write(w, binary.LittleEndian, witness.VerificationScript)

	if err != nil {
		return err
	}
	return nil
}

func (c *BasicTransaction) SerialiseScripts(w io.Writer) error {

	lenWitnesses := uint64(len(c.Witnesses))
	if err := writer.WriteVarUint(w, lenWitnesses); err != nil {
		return err
	}

	for _, script := range c.Witnesses {
		if err := script.Encode(w); err != nil {
			return err
		}
	}
	return nil

}

func (c *BasicTransaction) GetHash() string {
	buf := new(bytes.Buffer)
	c.SerialiseTransaction(buf, false)
	hash1, err := crypto.HashSha256(buf.Bytes())
	if err != nil {
		fmt.Println("Error with hashing of transaction")
		return ""
	}
	hash2, err := crypto.HashSha256(hash1)
	if err != nil {
		fmt.Println("Error with hashing of transaction")
		return ""
	}
	transHash := hex.EncodeToString(sliceUtils.Reverse(hash2))
	return transHash

}

// TODO: move this to a better place, like Utils/StringUtils
func ReverseString(s string) string {
	var reverse string
	for i := len(s) - 1; i >= 0; i-- {
		reverse += string(s[i])
	}
	return reverse
}

func AddressToScriptHash(address string) string {

	decodedAddressAsBytes, err := crypto.Base58Decode(address)
	if err != nil {
		return ""
	}
	decodedAddressAsHex := hex.EncodeToString(decodedAddressAsBytes)
	scriptHash := (decodedAddressAsHex[2:42])
	return scriptHash
}
