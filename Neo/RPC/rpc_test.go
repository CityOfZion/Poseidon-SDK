package rpc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Need to do these tests so they pass

func TestRPCSendTrans(t *testing.T) {
	node := rpc{"http://seed1.cityofzion.io:8080"}
	res := node.SendTransaction("Invalid String")
	assert.Equal(t, true, res)
}

func TestRPCGetRawTrans(t *testing.T) {
	node := rpc{"http://seed2.cityofzion.io:8080"}
	result, err := node.GetRawTransaction("c059754d44dba4d0d4cce71d4c503443ef8c2124f83c1fb760373431184823bc", 0)
	fmt.Println("RESF: ", result)
	if err != nil {

		t.Fail()
	}
	assert.Equal(t, "02000177f16f02a27d34564005427b6337ce5072732d25bd4b7c603ef60180778c23400000000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c606403000000000000860a46b5e9bc93aa726b2e838043d3ee0c053a49014140b65a715a0d4e02fa02ed87bdd03d9dd783899eb2f83a7921008ea11f846db3d9a1b8915cf48ed91c9407da2c7399bb4ead0ac7c142eecab2d28b8d86f5dbfc89232102a8d1a605a4de00c3f1d0ae6707f41d45044d9cd08077362e29ebebbb1fcfc53eac", result)
}
