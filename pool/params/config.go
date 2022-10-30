package def

import (
	"github.com/kaspanet/kaspad/util"
)

//WORK IN PROGRESS... replace "placeholders"

type config struct {

	poolConfig poolConfig
}

type poolConfig struct {
	PoolName     			string 		`json:"PoolName"`
	PoolAddress  			util.Address 	`json:"PoolAddress"`
	RegionCode  			string		`json:"RegionCode"`
	BlockCacheSize 			int		`json:"BlockCacheSize"`
	KaspadPort  			string		`json:"KaspadPort"`
	KaspadAdress 			string		`json:"KaspadAdress"`
	CommandPort             	string 		`json:"CommandPort"`
	CommandAddress          	string		`json:"CommandAddress"`
	PoolListenPort			string		`json:"PoolListenPort"`
	Nppls				int		`json:"Nppls"`
	BanDurationInSecounds   	int		`json:"BanDurationInSecounds"`
	TargetShareSubmitInSecounds  	int		`json:"TargetShareSubmitInSecounds"`
	AutoTargetShareSubmit		bool		`json:"AutoTargetShareSubmit"`
	StartDifficulty         	int		`json:"iStartDifficulty"`
	MaxConnections			int		`json:"MaxConnections"`
	MaxWorkerPerConnection  	int 		`json:"MaxWorkerPerConnection"`
	MaxPoolConnReadBytes    	int		`json:"MaxPoolConnReadBytes"`
	MaxPoolConnWriteBytes 		int		`json:"MaxPoolConnWriteBytes"`
	CleanUpIntervalSecounds		int		`json:"CleanUpIntervalSecounds"`
	ShareCutOffInSecounds   	int		`json:"ShareCutOffInSecounds"`
	ReservedNonceSpace		[2]int16	`json:"ReservedNonceSpace"`

}



func ParseJsonDefaultConfig() {

}

func OverWriteConfigWithFlags(poolConfig *poolConfig) {

}
