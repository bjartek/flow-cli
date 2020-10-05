package get

import (
	"fmt"
	"log"

	"github.com/onflow/flow-go-sdk"
	"github.com/psiemens/sconfig"
	"github.com/spf13/cobra"

	cli "github.com/dapperlabs/flow-cli/flow"
)

type Config struct {
	Host string `default:"127.0.0.1:3569" flag:"host" info:"Flow Observation API host address"`
	Code bool   `default:"false" flag:"code" info:"Display code deployed to the account"`
}

var conf Config

var Cmd = &cobra.Command{
	Use:   "get <address>",
	Short: "Get account info",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		acc := cli.GetAccount(conf.Host, flow.HexToAddress(args[0]))

		printAccount(acc, conf.Code)
	},
}

func init() {
	initConfig()
}

func initConfig() {
	err := sconfig.New(&conf).
		FromEnvironment(cli.EnvPrefix).
		BindFlags(Cmd.PersistentFlags()).
		Parse()
	if err != nil {
		log.Fatal(err)
	}
}

func printAccount(account *flow.Account, printCode bool) {
	fmt.Println()
	fmt.Println("Address: " + account.Address.Hex())
	fmt.Println("Balance : ", account.Balance)
	fmt.Println("Total Keys: ", len(account.Keys))
	for _, key := range account.Keys {
		fmt.Println("  ---")
		fmt.Println("  Key ID: ", key.Index)
		fmt.Printf("  PublicKey: %x\n", key.PublicKey.Encode())
		fmt.Println("  SigAlgo: ", key.SigAlgo)
		fmt.Println("  HashAlgo: ", key.HashAlgo)
		fmt.Println("  Weight: ", key.Weight)
		fmt.Println("  SequenceNumber: ", key.SequenceNumber)
	}
	fmt.Println("  ---")
	if printCode {
		fmt.Println("Code:")
		fmt.Println(string(account.Code))
	}
	fmt.Println()
}