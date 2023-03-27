package payment

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	IERC20 "github.com/rjman-self/Platdot/bindings/IERC20"
)

func SubscribeForTransfer(listener *IERC20.IERC20Filterer, listenChannel chan<- *IERC20.IERC20Transfer, ourAddress []common.Address) (event.Subscription, error) {
	subscription, err := listener.WatchTransfer(&bind.WatchOpts{
		Start:   nil, //last block
		Context: nil, // nil = no timeout
	}, listenChannel, nil, ourAddress,
	)
	if err != nil {
		return nil, err
	}
	return subscription, err
}

// Confirms payment. Takes in token's contract address, user address and amount of tokens
// user should send, returns true whenever tokens are received. Note that you should
// specify the recepient address and gateway websocket in .env file beforehand.
func ConfirmPayment(token common.Address, userAddress string, expectedAmount string) bool {

	ctx := context.Background()
	var channel = make(chan *IERC20.IERC20Transfer)
	ourAddress := []common.Address{common.HexToAddress(os.Getenv("OUR_ADDRESS"))}
	user := common.HexToAddress(userAddress)
	gateway := os.Getenv("GATEWAY_GOERLI_WS")

	client, err := ethclient.Dial(gateway) // load from local .env file
	if err != nil {
		log.Fatalf("could not connect to Ethereum gateway: %v\n", err)
	}
	defer client.Close()

	Event, err := IERC20.NewIERC20Filterer(token, client)
	if err != nil {
		log.Fatalf("Failed to create filterer: %v", err)
	}

	subscription, err := SubscribeForTransfer(Event, channel, ourAddress)
	fmt.Println(subscription) // this is subscription to INDEXED event. This mean we can pass what exactly value of argument we want to see
	if err != nil {
		log.Fatal(err)
	}

EventLoop:
	for {
		select {
		case <-ctx.Done():
			{
				subscription.Unsubscribe()
			}
		case eventResult := <-channel:
			{
				fmt.Println("From:", eventResult.From, "Value:", eventResult.Value)
				if eventResult.Value.String() == expectedAmount && eventResult.From == user {
					subscription.Unsubscribe()
					break EventLoop
				}
			}

		}
	}
	return true
}
