package main

import (
	"context"
	"log"
	"os"
	"time"

	coreEntities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/daoleno/uniswapv3-sdk/examples/helper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var BOBA = coreEntities.NewToken(1, common.HexToAddress("0xb0ba1b6ebadeba1a63a94445f0dfb249082b5dc1"), 9, "BOBA", "Boba")

var HARRIS = coreEntities.NewToken(1, common.HexToAddress("0x155788dd4b3ccd955a5b2d461c7d6504f83f71fa"), 9, "HARRIS", "KAMALA HARRIS")

var MOODENG = coreEntities.NewToken(1, common.HexToAddress("0x28561b8a2360f463011c16b6cc0b0cbef8dbbcad"), 9, "MOODENG", "MOODENG")

var MISHA = coreEntities.NewToken(1, common.HexToAddress("0x0ccae1bc46fb018dd396ed4c45565d4cb9d41098"), 9, "MISHA", "MISHA")

var MARS = coreEntities.NewToken(1, common.HexToAddress("0xb8d6196d71cdd7d90a053a7769a077772aaac464"), 9, "MARS", "MARS")

var HANA = coreEntities.NewToken(1, common.HexToAddress("0xb3912b20b3abc78c15e85e13ec0bf334fbb924f7"), 9, "HANA", "HANA")

var DOGEFATHER = coreEntities.NewToken(1, common.HexToAddress("0x9bcda9bd845e4bfa7e6a4c012b93f8ba821ac618"), 9, "DOGEFATHER", "DOGEFATHER")

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	ethRpc := "https://eth-beanz.rpc.blockrazor.io/"
	client, err := ethclient.Dial(ethRpc)
	if err != nil {
		log.Fatal(err)
	}
	wallet := helper.InitWallet(os.Getenv("MY_PRIVATE_KEY"))
	if wallet == nil {
		log.Fatal("init wallet failed")
	}

	token := MARS
	tx, err := helper.SwapExactETHForTokens(client, wallet, token, "0.4")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(tx.Hash())

	//token := DOGEFATHER
	//tx, err := helper.SwapExactTokensForETH(client, wallet, token, "35305313.77370976")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println(tx.Hash())

	for {
		recipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * 1)
			continue
		}
		log.Fatal(recipt.GasUsed)
	}

}
