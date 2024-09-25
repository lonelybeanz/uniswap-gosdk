package helper

import (
	"context"
	"errors"
	"fmt"
	coreEntities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/daoleno/uniswapv3-sdk/examples/contract"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"time"
)

func initSwap(client *ethclient.Client, w *Wallet) (*contract.Uniswapv2RouterV2, *bind.TransactOpts, error) {
	uniswapv2RouterV2, err := contract.NewUniswapv2RouterV2(common.HexToAddress(ContractV2SwapRouterV2), client)
	if err != nil {
		return nil, nil, err
	}
	if uniswapv2RouterV2 == nil {
		return nil, nil, errors.New("uniswapv2RouterV2 is nil")
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, err
	}
	//增加10%
	gasPrice = IntWithDecimal(gasPrice.Uint64()*15/10, 0)
	fmt.Printf("gasPrice=%d\n", gasPrice.Uint64())

	nonce, err := client.NonceAt(context.Background(), w.PublicKey, nil)
	if err != nil {
		return nil, nil, err
	}

	signer := types.LatestSignerForChainID(big.NewInt(1))
	opts := &bind.TransactOpts{
		From: w.PublicKey,
		Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != w.PublicKey {
				return nil, bind.ErrNotAuthorized
			}
			signature, err := crypto.Sign(signer.Hash(tx).Bytes(), w.PrivateKey)
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
		Context: context.Background(),
		Nonce:   big.NewInt(int64(nonce)),
	}
	opts.Value = big.NewInt(0)
	opts.GasLimit = uint64(3000000)
	opts.GasFeeCap = gasPrice //big.NewInt(18 * 1e9)
	opts.GasTipCap = big.NewInt(0.1)

	return uniswapv2RouterV2, opts, nil
}

func SwapExactTokensForETH(client *ethclient.Client, w *Wallet, token0 *coreEntities.Token, amount0 string) (*types.Transaction, error) {

	uniswapv2RouterV2, opts, err := initSwap(client, w)
	if err != nil {
		return nil, err
	}

	swapValue := FloatStringToBigInt(amount0, int(token0.Decimals()))

	//token0.Approve(client, w, uniswapv2RouterV2.Address, swapValue)

	var amountIn, amountOutMin, deadline *big.Int
	amountIn = swapValue

	deadline = big.NewInt(int64(time.Now().Add(time.Minute * 10).Unix()))

	path := make([]common.Address, 2)
	path[0] = token0.Address
	weth, err := uniswapv2RouterV2.WETH(nil)
	if err != nil {
		return nil, err
	}
	path[1] = weth

	//预估兑换额
	amountsOut, err := uniswapv2RouterV2.GetAmountsOut(nil, amountIn, path)
	if err != nil {
		return nil, err
	}
	for _, a := range amountsOut {
		fmt.Printf("amountsOut=%s\n", a)
	}
	amountOutMin = amountsOut[1]

	to := w.PublicKey

	fmt.Printf("SwapExactTokensForETH amountIn=%s amountOutMin=%s deadline=%s\n", amountIn, amountOutMin, deadline)

	tx, err := uniswapv2RouterV2.SwapExactTokensForETH(opts, amountIn, amountOutMin, path, to, deadline)
	if err != nil {
		return nil, err
	}
	fmt.Printf("SwapExactTokensForETH tx=%s\n", tx.Hash())
	return tx, nil
}

func SwapExactETHForTokens(client *ethclient.Client, w *Wallet, token0 *coreEntities.Token, amount string) (*types.Transaction, error) {

	uniswapv2RouterV2, opts, err := initSwap(client, w)
	if err != nil {
		return nil, err
	}

	swapValue := FloatStringToBigInt(amount, 18)
	opts.Value = swapValue

	var amountIn, amountOutMin, deadline *big.Int
	amountIn = swapValue
	deadline = big.NewInt(time.Now().UTC().Add(time.Minute * 10).Unix())

	path := make([]common.Address, 2)
	weth, err := uniswapv2RouterV2.WETH(nil)
	if err != nil {
		return nil, err
	}
	path[0] = weth
	path[1] = token0.Address

	//预估兑换额
	amountsOut, err := uniswapv2RouterV2.GetAmountsOut(nil, amountIn, path)
	if err != nil {
		return nil, err
	}
	for _, a := range amountsOut {
		fmt.Printf("amountsOut=%s\n", a)
	}

	amountOutMin = amountsOut[1]
	to := w.PublicKey

	fmt.Printf("SwapExactETHForTokens amountIn=%s amountOutMin=%s deadline=%s\n", amountIn, amountOutMin, deadline)

	//执行兑换
	tx, err := uniswapv2RouterV2.SwapExactETHForTokens(opts, amountOutMin, path, to, deadline)
	if err != nil {
		return nil, err
	}
	fmt.Printf("SwapExactETHForTokens tx=%s\n", tx.Hash())
	return tx, nil
	//return nil, errors.New("not support")
}
