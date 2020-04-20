package cli

import (
	"context"
	"fmt"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/client"
	"github.com/dapperlabs/flow-go-sdk/crypto"
	"google.golang.org/grpc"
)

func SendTransaction(host string, signerAccount *Account, script []byte) {

	flowClient, err := client.New(host, grpc.WithInsecure())
	if err != nil {
		Exitf(1, "Failed to connect to host: %s", err)
	}

	signerAddress := signerAccount.Address

	fmt.Printf("Getting information for account with address 0x%s ...\n", signerAddress.Short())

	account, err := flowClient.GetAccount(context.Background(), signerAddress)
	if err != nil {
		Exitf(1, "Failed to get account with address %s: 0x%s", signerAddress.Short(), err)
	}

	signer := crypto.NewNaiveSigner(
		signerAccount.PrivateKey.PrivateKey,
		signerAccount.PrivateKey.HashAlgo,
	)

	// TODO: always use first?
	accountKey := account.Keys[0]

	tx := flow.NewTransaction().
		SetScript(script).
		SetProposalKey(signerAddress, accountKey.ID, accountKey.SequenceNumber).
		SetPayer(signerAddress).
		AddAuthorizer(signerAddress)

	err = tx.SignEnvelope(signerAddress, accountKey.ID, signer)
	if err != nil {
		Exitf(1, "Failed to sign transaction: %s", err)
	}

	fmt.Printf("Submitting transaction with ID %s ...\n", tx.ID())

	err = flowClient.SendTransaction(context.Background(), *tx)
	if err == nil {
		fmt.Printf("Successfully submitted transaction with ID %s\n", tx.ID())
	} else {
		Exitf(1, "Failed to submit transaction: %s", err)
	}
}
