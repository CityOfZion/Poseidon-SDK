package scriptBuilder

//copied from http://docs.neo.org/en-us/sc/systemfees.html

// This will be used to calculate the cost of running the script.
// Neo is going to change the gas scheme, so this is not implemented yet.
// Gas cost will by default be zero.

type GasCost int64

// In fixed 8
const (
	PUSH16Cost        GasCost = 0
	NOPCost           GasCost = 0
	APPCALLCost       GasCost = 1000000
	TAILCALLCost      GasCost = 1000000
	SHA1Cost          GasCost = 1000000
	SHA256Cost        GasCost = 1000000
	HASH160Cost       GasCost = 2000000
	HASH256Cost       GasCost = 2000000
	CHECKSIGCost      GasCost = 10000000
	CHECKMULTISIGCost GasCost = 10000000
	DEFAULTCost       GasCost = 100000
)
