package services

import (
	"context"
	"crypto-wallet/blockchain"
	"crypto-wallet/config"
	"crypto-wallet/database"
	"crypto-wallet/models"
	"crypto-wallet/utils"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
)

// ===============================
// GET BLOCKCHAIN (BLOCK EXPLORER)
// ===============================
func GetBlockchainService() ([]blockchain.Block, error) {

	ctx := context.Background()

	snaps, err := config.Firestore.Collection("blocks").
		OrderBy("index", firestore.Asc).
		Documents(ctx).
		GetAll()

	if err != nil {
		return nil, err
	}

	var blocks []blockchain.Block

	for _, doc := range snaps {
		var b blockchain.Block
		doc.DataTo(&b)
		blocks = append(blocks, b)
	}

	return blocks, nil
}

// ===============================
// MINE PENDING TRANSACTIONS
// ===============================
func MinePendingTransactionsService(miner string) (*blockchain.Block, error) {

	ctx := context.Background()

	// -------------------------------------
	// 1. LOAD PENDING TXs
	// -------------------------------------
	pendCol := config.Firestore.Collection("pending_transactions")

	pendDocs, err := pendCol.Documents(ctx).GetAll()
	if err != nil {
		database.AddLog("error", "Failed to load pending transactions: "+err.Error())
		return nil, err
	}
	if len(pendDocs) == 0 {
		database.AddLog("info", "Mining attempted but no pending transactions found")
		return nil, errors.New("no pending transactions to mine")
	}

	database.AddLog("info",
		fmt.Sprintf("Mining started by %s — Pending TXs: %d", miner, len(pendDocs)),
	)

	// Convert docs → map slice
	txs := []map[string]interface{}{}
	for _, d := range pendDocs {
		txs = append(txs, d.Data())
	}

	// -------------------------------------
	// 2. LOAD LAST BLOCK
	// -------------------------------------
	blocksCol := config.Firestore.Collection("blocks")

	var index int
	prevHash := "GENESIS"

	iter := blocksCol.OrderBy("index", firestore.Desc).Limit(1).Documents(ctx)
	lastDoc, err := iter.Next()

	if err == nil {
		var lastBlock blockchain.Block
		lastDoc.DataTo(&lastBlock)
		index = lastBlock.Index + 1
		prevHash = lastBlock.Hash
	} else {
		index = 0
	}

	timestamp := utils.Now()

	// -------------------------------------
	// 3. CREATE BLOCK + ASSIGN MINER REWARD
	// -------------------------------------
	block := blockchain.NewBlock(index, timestamp, prevHash, txs)
	block.AssignMiner(miner)

	database.AddLog("info", "Starting Proof-of-Work for block "+fmt.Sprint(index))

	// PROOF OF WORK
	blockchain.MineBlock(block)

	database.AddLog("success", "Block mined! Hash: "+block.Hash)

	// -------------------------------------
	// 4. SAVE BLOCK IN DB
	// -------------------------------------
	_, _, err = blocksCol.Add(ctx, block)
	if err != nil {
		database.AddLog("error", "Failed to save block: "+err.Error())
		return nil, err
	}

	// -------------------------------------
	// 5. UPDATE UTXOs
	// -------------------------------------
	for _, tx := range txs {

		// Mark spent inputs
		if inputs, ok := tx["inputs"].([]interface{}); ok {
			for _, inp := range inputs {
				inpm := inp.(map[string]interface{})
				utxoID := inpm["id"].(string)
				_ = database.MarkUTXOSpent(utxoID)
			}
		}

		// Receiver output
		if receiver, ok := tx["receiver"].(string); ok {
			amt := tx["amount"].(float64)

			recUtxo := models.UTXO{
				WalletID: receiver,
				Amount:   amt,
				Spent:    false,
				TxID:     block.Hash,
			}
			database.AddUTXO(recUtxo)
		}

		// Change output
		if change, ok := tx["change"].(float64); ok && change > 0 {
			sender := tx["sender"].(string)

			changeUTXO := models.UTXO{
				WalletID: sender,
				Amount:   change,
				Spent:    false,
				TxID:     block.Hash,
			}
			database.AddUTXO(changeUTXO)
		}
	}

	// ----------------------------------------------
	// 6B. WRITE COMPLETED TRANSACTION HISTORY RECORDS
	// ----------------------------------------------
	for _, tx := range txs {

		sender, hasSender := tx["sender"].(string)
		receiver, hasReceiver := tx["receiver"].(string)
		if !hasSender || !hasReceiver {
			continue
		}

		amount := tx["amount"].(float64)

		note := ""
		if tx["note"] != nil {
			note = tx["note"].(string)
		}

		timestamp := int64(tx["timestamp"].(int64))

		// Special Zakat Transaction Logging
		if note == "Monthly zakat deduction" {
			database.AddLog("success",
				fmt.Sprintf("Zakat processed in block %s for wallet %s (amount %.4f)",
					block.Hash, sender, amount,
				),
			)
		}

		// SENT RECORD
		_, _, _ = config.Firestore.Collection("transactions").Add(ctx, map[string]interface{}{
			"wallet_id":    sender,
			"type":         "sent",
			"amount":       amount,
			"note":         note,
			"timestamp":    timestamp,
			"counterparty": receiver,
		})

		// RECEIVED RECORD
		_, _, _ = config.Firestore.Collection("transactions").Add(ctx, map[string]interface{}{
			"wallet_id":    receiver,
			"type":         "received",
			"amount":       amount,
			"note":         note,
			"timestamp":    timestamp,
			"counterparty": sender,
		})

		database.AddLog("success",
			fmt.Sprintf("History written: %s sent %.2f to %s", sender, amount, receiver))
	}

	// -------------------------------------
	// 7. ADD MINING REWARD
	// -------------------------------------
	rewardUTXO := models.UTXO{
		WalletID: miner,
		Amount:   50,
		Spent:    false,
		TxID:     block.Hash,
	}
	database.AddUTXO(rewardUTXO)

	database.AddLog("success",
		fmt.Sprintf("Mining reward (50 coins) sent to %s", miner),
	)

	// -------------------------------------
	// 8. DELETE PENDING TRANSACTIONS
	// -------------------------------------
	batch := config.Firestore.Batch()
	for _, d := range pendDocs {
		batch.Delete(d.Ref)
	}
	_, err = batch.Commit(ctx)

	if err != nil {
		database.AddLog("error", "Failed to clear pending tx: "+err.Error())
		return nil, err
	}

	database.AddLog("info", "Pending transactions cleared.")

	return block, nil
}
