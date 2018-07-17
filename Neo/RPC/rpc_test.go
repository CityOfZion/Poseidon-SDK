package rpc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRPCSendTrans(t *testing.T) {

	// Invalid hex string should return false
	node := Rpc{"http://seed1.cityofzion.io:8080"}
	res := node.SendTransaction("Invalid String")
	assert.Equal(t, false, res)

}

func TestRPCGetRawTrans(t *testing.T) {
	node := Rpc{"http://seed3.aphelion-neo.com:10332"}
	result, err := node.GetRawTransaction("56d477f7cffe5f3f919be798e7c752c754faf03870d202432b8157d9da0dfc57", 0)

	if err != nil {

		t.Fail()
	}

	assert.Equal(t, "d101530880bd7e5b01000000145ff7a5ad95cf4370f84514e01b1acd8fc28f3b4f14d92f268b9bd6133c9a6dc50c04e7ddd1f0aee10b53c1087472616e736665726735f731696271626c14d6cbddaf736f50ddc8992100000000000000000220d92f268b9bd6133c9a6dc50c04e7ddd1f0aee10bf0153135333138343932353635313466613363323365610000014140251c00a67aee3dd6df96bfac34b7851306506db68d80a34ea8b3c22b96fc6f94b96a0e40947bce582be02d7508b70655d626d8ca7e75fd1dca59589914ba4ed52321032e07a10aa1f1a56900e2ec1c6255aaf53f6ddc606729b8255ccbb7b90d5724b4ac", result)
}
func TestGetBlock(t *testing.T) {
	node := Rpc{"https://seed1.neo.org:20331"}

	res, _ := node.getRawBlock(2000, 0)

	assert.NotEmpty(t, res)
}

func TestInvok(t *testing.T) {
	node := Rpc{"https://seed1.neo.org:20331"}

	res, err := node.InvokeScript("1457c4cf51f12ce6d78a585f1ea9bd1f3927c7232c51c10962616c616e63654f6667e7b132b995f43dbbddd2a3268a04a2ae081eff9a")
	if err != nil {
		fmt.Println(err, res)
	}
	assert.NotEmpty(t, res)

}
