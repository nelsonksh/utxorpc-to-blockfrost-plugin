package main

import (
	"encoding/json"
	"fmt"

	plugin "github.com/nelsonksh/utxorpc-to-blockfrost-plugin"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/query"
	"github.com/utxorpc/go-sdk"
	utxorpc "github.com/utxorpc/go-sdk/cardano"
)

func main() {
	client := utxorpc.NewClient(sdk.WithBaseUrl("https://..."))
	utxo := readUtxo(client, "3f565260990a08d24fe56db4c3605bb142e668b64f9932679106fa270aaa6153", 3)
	// utxoJson, err := json.MarshalIndent(utxo, "", "  ")
	// if err != nil {
	// 	fmt.Println("Error marshaling UTXO to JSON:", err)
	// 	return
	// }
	// fmt.Println(string(utxoJson))
	blockfrostUtxo := plugin.UtxorpcToBlockfrostUtxo(utxo)

	blockfrostUtxoJson, err := json.MarshalIndent(blockfrostUtxo, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling Blockfrost UTXO to JSON:", err)
		return
	}
	fmt.Println(string(blockfrostUtxoJson))
}

func readUtxo(
	client *utxorpc.Client,
	txHashStr string,
	txIndex uint32,
) *query.AnyUtxoData {
	resp, err := client.GetUtxoByRef(txHashStr, txIndex)
	if err != nil {
		sdk.HandleError(err)
		return nil
	}
	return resp.Msg.Items[0]
}
