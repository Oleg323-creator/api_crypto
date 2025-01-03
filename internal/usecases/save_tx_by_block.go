package usecases

import (
	"api_crypto/internal/db/repository"
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math/big"
)

func (u *Usecases) SaveTxInfoByBlock(addr string, curr string) error {
	lastBlockInDb, err := u.Repository.GetLastBlockFromDB()
	if err != nil {
		log.Fatal(err)
	}

	block, err := u.Client.BlockByNumber(context.Background(), big.NewInt(lastBlockInDb))
	if err != nil {
		log.Fatalf("Failed to fetch block: %v", err)
	}

	chainID, err := u.Client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Error getting chain ID: %v", err)
	}

	// GETTING INFO ABOUT ALL TX IN BLOCK
	for _, tx := range block.Transactions() {
		signer := types.LatestSignerForChainID(chainID)
		senderAddr, err := types.Sender(signer, tx)
		if err != nil {
			log.Printf("Error determining sender address: %v", err)
			continue
		}
		/*
			senderAddr, err := types.Sender(types.NewEIP155Signer(chainID), tx)
			if err != nil {
				log.Printf("Error with EIP155Signer: %v", err)
				senderAddr, err = types.Sender(types.HomesteadSigner{}, tx)
				if err != nil {
					log.Printf("Error with HomesteadSigner: %v", err)
					continue
				}
			}
		*/
		var toAddr string
		if tx.To() != nil {
			toAddr = tx.To().Hex()
		} else {
			toAddr = "Contract Creation"
		}

		if toAddr == addr {
			data := repository.TxData{
				Hash:     tx.Hash().Hex(),
				FromAddr: senderAddr.String(),
				ToAddr:   toAddr,
				Value:    tx.Value().String(),
				Currency: curr,
			}
			err = u.Repository.SaveTxDataToDB(data)
			if err != nil {
				log.Fatalf("Error saiving data in db: %v", err)
			}

			log.Println("Tx data saved to DB")

			//AFTER SAIVING TX DATA IN DB WE JUST MERGE ALL COINS TO ROOT ADDRESS
			err = u.MergeCoinsToRoot(data)
			if err != nil {
				return err
			}

		}
	}
	return nil
}
