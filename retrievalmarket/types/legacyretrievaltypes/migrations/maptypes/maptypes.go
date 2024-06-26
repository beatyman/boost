package maptypes

import (
	"github.com/filecoin-project/boost/datatransfer"
	"github.com/filecoin-project/boost/markets/piecestore"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"

	"github.com/filecoin-project/boost/retrievalmarket/types/legacyretrievaltypes"
)

//go:generate cbor-gen-for --map-encoding ClientDealState1 ProviderDealState1

// Version 1 of the ClientDealState
type ClientDealState1 struct {
	legacyretrievaltypes.DealProposal
	StoreID              *uint64
	ChannelID            datatransfer.ChannelID
	LastPaymentRequested bool
	AllBlocksReceived    bool
	TotalFunds           abi.TokenAmount
	ClientWallet         address.Address
	MinerWallet          address.Address
	PaymentInfo          *legacyretrievaltypes.PaymentInfo
	Status               legacyretrievaltypes.DealStatus
	Sender               peer.ID
	TotalReceived        uint64
	Message              string
	BytesPaidFor         uint64
	CurrentInterval      uint64
	PaymentRequested     abi.TokenAmount
	FundsSpent           abi.TokenAmount
	UnsealFundsPaid      abi.TokenAmount
	WaitMsgCID           *cid.Cid
	VoucherShortfall     abi.TokenAmount
	LegacyProtocol       bool
}

// Version 1 of the ProviderDealState
type ProviderDealState1 struct {
	legacyretrievaltypes.DealProposal
	StoreID         uint64
	ChannelID       datatransfer.ChannelID
	PieceInfo       *piecestore.PieceInfo
	Status          legacyretrievaltypes.DealStatus
	Receiver        peer.ID
	TotalSent       uint64
	FundsReceived   abi.TokenAmount
	Message         string
	CurrentInterval uint64
	LegacyProtocol  bool
}
