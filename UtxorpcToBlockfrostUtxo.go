package plugin

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/blinklabs-io/gouroboros/ledger/common"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/cardano"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/query"
)

func UtxorpcToBlockfrostUtxo(utxo *query.AnyUtxoData) BlockfrostUtxo {

	var address string
	if utxo.GetCardano().GetAddress() != nil {
		addr, err := common.NewAddressFromBytes(utxo.GetCardano().GetAddress())
		if err != nil {
			fmt.Println("Error converting address bytes to string:", err)
		}
		address = addr.String()
	}

	amounts := []BlockfrostAsset{
		{
			Unit:     "lovelace",
			Quantity: strconv.FormatUint(utxo.GetCardano().Coin, 10),
		},
	}
	for _, asset := range utxo.GetCardano().GetAssets() {
		amounts = append(amounts, UtxorpcMultiassetToBlockfrostAssets(asset)...)
	}

	// var block string

	var dataHash *string
	if utxo.GetCardano().GetDatum().GetHash() != nil {
		hashStr := hex.EncodeToString(utxo.GetCardano().GetDatum().GetHash())
		dataHash = &hashStr
	}

	var inlineDatum string
	if utxo.GetCardano().GetDatum().GetPayload() != nil {
		inlineDatum = hex.EncodeToString(utxo.GetCardano().GetDatum().GetOriginalCbor())
	}

	var referenceScriptHash string
	var plutusV1Script common.PlutusV1Script
	var plutusV2Script common.PlutusV2Script
	var plutusV3Script common.PlutusV3Script
	if utxo.GetCardano().GetScript().GetNative().GetNativeScript() != nil {
		referenceScriptHash = hex.EncodeToString(utxo.GetCardano().GetScript().GetNative().GetScriptPubkey())
	}
	if utxo.GetCardano().GetScript().GetPlutusV1() != nil {
		plutusV1Script = utxo.GetCardano().GetScript().GetPlutusV1()
		referenceScriptHash = plutusV1Script.Hash().String()
	}
	if utxo.GetCardano().GetScript().GetPlutusV2() != nil {
		plutusV2Script = utxo.GetCardano().GetScript().GetPlutusV2()
		referenceScriptHash = plutusV2Script.Hash().String()
	}
	if utxo.GetCardano().GetScript().GetPlutusV3() != nil {
		plutusV3Script = utxo.GetCardano().GetScript().GetPlutusV3()
		referenceScriptHash = plutusV3Script.Hash().String()
	}

	return BlockfrostUtxo{
		Address:             address,
		TxHash:              hex.EncodeToString(utxo.GetTxoRef().GetHash()),
		OutputIndex:         utxo.GetTxoRef().GetIndex(),
		Amount:              amounts,
		Block:               "",
		DataHash:            dataHash,
		InlineDatum:         &inlineDatum,
		ReferenceScriptHash: &referenceScriptHash,
	}
}

func UtxorpcMultiassetToBlockfrostAssets(asset *cardano.Multiasset) []BlockfrostAsset {
	policy := hex.EncodeToString(asset.GetPolicyId())
	var assetNames []string
	var quantities []string
	for _, a := range asset.GetAssets() {
		assetName := hex.EncodeToString(a.GetName())
		quantity := strconv.FormatUint(a.GetOutputCoin(), 10)
		assetNames = append(assetNames, assetName)
		quantities = append(quantities, quantity)
	}
	var blockfrostAssets []BlockfrostAsset
	for i := 0; i < len(assetNames); i++ {
		blockfrostAssets = append(blockfrostAssets, BlockfrostAsset{
			Unit:     policy + assetNames[i],
			Quantity: quantities[i],
		})
	}
	return blockfrostAssets
}

type BlockfrostUtxo struct {
	Address             string            `json:"address"`
	TxHash              string            `json:"tx_hash"`
	OutputIndex         uint32            `json:"output_index"`
	Amount              []BlockfrostAsset `json:"amount"`
	Block               string            `json:"block"`
	DataHash            *string           `json:"data_hash,omitempty"`
	InlineDatum         *string           `json:"inline_datum,omitempty"`
	ReferenceScriptHash *string           `json:"reference_script_hash,omitempty"`
}

type BlockfrostAsset struct {
	Unit     string `json:"unit"`
	Quantity string `json:"quantity"`
}
