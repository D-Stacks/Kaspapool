package def

import (
	"github.com/kaspanet/kaspad/util"
	"math/big"
)

//WORK IN PROGRESS... replace "placeholders"

//Pool meta data
var Version = "0.0.1"
var Devlopers = "jwj"

//KaspadConstants
var MaxPoW, _ = new(big.Int).SetString("7EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE", 16)
var KaspaBPS = 1 
var testnetPrefix = util.Bech32PrefixKaspaTest
var mainnetPrefix = util.Bech32PrefixKaspa
var KCluster = 18
