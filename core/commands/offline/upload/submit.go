package upload

import (
	"fmt"

	"github.com/TRON-US/go-btfs/core/escrow"
	"github.com/TRON-US/go-btfs/core/guard"

	escrowpb "github.com/tron-us/go-btfs-common/protos/escrow"
)

func submit(rss *RenterSession, shardHashes []string, fileSize int64, offlineSigning bool) {
	rss.submit()
	bs, t, err := prepareContracts(rss, shardHashes)
	if err != nil {
		// TODO: handle error
		return
	}
	fmt.Println(len(bs), t)
	err = checkBalance(rss, offlineSigning, t)
	if err != nil {
		fmt.Println("get balance error:", err)
		// TODO: handle error
		return
	}
	req, err := NewContractRequest(rss.ctxParams.cfg, rss.ssId, bs, t)
	if err != nil {
		// TODO: handle error
		return
	}
	var amount int64 = 0
	for _, c := range req.Contract {
		amount += c.Contract.Amount
	}
	submitContractRes, err := escrow.SubmitContractToEscrow(rss.ctxParams.ctx, rss.ctxParams.cfg, req)
	if err != nil {
		//TODO: handle error
		return
	}
	fmt.Println("submitContractRes", submitContractRes)
	//TODO
	//doPay()
	return
}

func prepareContracts(rss *RenterSession, shardHashes []string) ([]*escrowpb.SignedEscrowContract, int64, error) {
	var signedContracts []*escrowpb.SignedEscrowContract
	var totalPrice int64
	for _, hash := range shardHashes {
		shard, err := GetRenterShard(rss.ctxParams, rss.ssId, hash)
		if err != nil {
			return nil, 0, err
		}
		c, err := shard.contracts()
		if err != nil {
			return nil, 0, err
		}
		sc, err := escrow.UnmarshalEscrowContract(c.SignedEscrowContract)
		if err != nil {
			return nil, 0, err
		}
		signedContracts = append(signedContracts, sc)
		guardContract, err := guard.UnmarshalGuardContract(c.SignedGuardContract)
		if err != nil {
			return nil, 0, err
		}
		totalPrice += guardContract.Amount
	}
	return signedContracts, totalPrice, nil
}
