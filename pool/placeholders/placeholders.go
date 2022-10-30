package placeholders

import (
	"math/big"
	"time"

	"github.com/kaspanet/kaspad/util"
)

const PoolAddress = "kaspa:wemewemowemweomfweofmio"
const TargetShareSubmit = time.Second * 5
//const CalibrateInterval = time.Second * 10 //recalculate vardiff and Re 
var StartDiff = big.NewRat(20, 1)

const MaxConnectionReadSizeInBytes = 256 //tcp sniffing shows 155 bytes to be the max for mining submit
const MaxConnectionWriteSizeInByes = 256 //tcp sniffing shows 143 bytes to be the max for mining submit

const StrikeDuration = time.Second * 30

const NPplns = 100

const NetworkPrefix = util.Bech32PrefixKaspaTest

const ShareDurationCutOff = TargetShareSubmit * NPplns