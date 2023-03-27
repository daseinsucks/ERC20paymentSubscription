package main

import (
	"fmt"
	payment "generalPayment/lib"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
)

func main() {

	output := make(chan bool)
	_ = godotenv.Load()

	//Test token address
	USDT := common.HexToAddress("0xd1E9b088553010E4F683Bde4D28BEa4631903E34")

	//Mainnet token addresses
	/* USDC := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	USDT := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	DAI := common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F")
	WETH := common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	WBTC := common.HexToAddress("0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599")
	*/

	go payment.ConfirmPayment(USDT, "0x383A9e83E36796106EaC11E8c2Fbe8b92Ff46D3a", "2000000000000000000", output)
	result := <-output
	fmt.Println(result)

}
