/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package dpos

import (
	"bytes"
	"testing"

	"chainmaker.org/chainmaker/utils/v2"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"

	pbdpos "chainmaker.org/chainmaker/pb-go/v2/consensus/dpos"

	commonpb "chainmaker.org/chainmaker/pb-go/v2/common"
	consensuspb "chainmaker.org/chainmaker/pb-go/v2/consensus"

	"github.com/stretchr/testify/require"
)

func TestBytesEqual(t *testing.T) {
	bz := make([]byte, 0)
	require.True(t, bytes.Equal(bz, nil))
}

func TestDPoSImpl_CreateDPoSRWSet(t *testing.T) {
	impl, fn := initTestImpl(t)
	defer fn()

	proposedBlk := &consensuspb.ProposalBlock{Block: &commonpb.Block{Header: &commonpb.BlockHeader{BlockHeight: 99}}}
	rwSet, err := impl.createDPoSRWSet(nil, proposedBlk)
	require.NoError(t, err)
	require.Nil(t, rwSet)

	proposedBlk.Block.Header.BlockHeight = 100
	rwSet, err = impl.createDPoSRWSet(nil, proposedBlk)
	require.EqualError(t, err, "not found candidates from contract")
	require.Nil(t, rwSet)
}

func TestDPoSImpl_CreateNewEpoch(t *testing.T) {
	// 1. nil candidateInfo in stateDB
	impl, fn := initTestImpl(t)
	defer fn()
	newEpoch, err := impl.createNewEpoch(9, &syscontract.Epoch{EpochId: 5}, []byte("12328423"))
	require.Error(t, err)
	require.Nil(t, newEpoch)

	// 2. have some candidate infos in stateDB
	impl, fn = initDPoSWithStore(t)
	defer fn()
	blk, blkRwSet := generateCandidateBlockAndRwSet(t, 6, 10, 1)
	require.NoError(t, impl.stateDB.PutBlock(blk, blkRwSet))

	newEpoch, err = impl.createNewEpoch(9, &syscontract.Epoch{EpochId: 5}, []byte("12328423"))
	require.NoError(t, err)
	require.EqualValues(t, 6, newEpoch.EpochId)
	require.EqualValues(t, 4, len(newEpoch.ProposerVector))
	require.EqualValues(t, 9+testEpochBlkNum, newEpoch.NextEpochCreateHeight)
}

func TestDPoSImpl_selectValidators(t *testing.T) {
	impl, fn := initTestImpl(t)
	defer fn()
	var tests = []*pbdpos.CandidateInfo{
		{PeerId: "peer0", Weight: "100"},
		{PeerId: "peer1", Weight: "50"},
		{PeerId: "peer2", Weight: "200"},
		{PeerId: "peer3", Weight: "10"},
	}
	_, err := impl.selectValidators(nil, []byte("123"))
	require.Error(t, err)

	infos, err := impl.selectValidators(tests, []byte("123"))
	require.NoError(t, err)
	for i, v := range tests {
		require.EqualValues(t, v.PeerId, infos[i].PeerId)
	}
}

func TestDPoSImpl_GetValidators(t *testing.T) {
	impl, fn := initTestImpl(t)
	defer fn()

	nodeIds, err := impl.GetValidators()
	require.NoError(t, err)
	require.EqualValues(t, 4, len(nodeIds))
}

func TestDPoSImpl_GetConsensusStateJSON(t *testing.T) {
	impl, fn := initTestImpl(t)
	defer fn()

	// 1. no validator state should fire error
	_, err := impl.GetConsensusStateJSON()
	require.Error(t, err)

	// 2. add consensus state
	impl, fn = initDPoSWithStore(t)
	defer fn()
	blk, blkRwSet := generateCandidateBlockAndRwSet(t, 6, 10, 1)
	require.NoError(t, impl.stateDB.PutBlock(blk, blkRwSet))

	_, err = impl.GetConsensusStateJSON()
	require.NoError(t, err)
}

func TestDPoSImpl_AddConsensusArgsToBlock(t *testing.T) {
	impl, fn := initTestImpl(t)
	defer fn()

	// maybe used to benchmark marshal data
	txWrite := &common.TxRWSet{}
	txWrite.TxWrites = append(txWrite.TxWrites, &common.TxWrite{
		Key:          []byte("test_key"),
		Value:        []byte("test_value"),
		ContractName: "test_contract",
	})
	blk := &common.Block{Header: &common.BlockHeader{BlockHeight: 9}}
	err := impl.addConsensusArgsToBlock(txWrite, blk)
	require.NoError(t, err)
}

func TestDPoSImpl_VerifyConsensusArgs(t *testing.T) {
	impl, fn := initDPoSWithStore(t)
	defer fn()
	blk, blkRwSet := generateCandidateBlockAndRwSet(t, 6, 10, 1)
	require.NoError(t, impl.stateDB.PutBlock(blk, blkRwSet))

	// 1. nil dpos txRWrite and not the correct height to create
	currBlk := &common.Block{Header: &common.BlockHeader{
		BlockHeight:  90,
		BlockHash:    []byte(utils.GetRandTxId()),
		PreBlockHash: []byte(utils.GetRandTxId()),
	}}
	err := impl.VerifyConsensusArgs(currBlk, nil)
	require.NoError(t, err)

	// 2. the height to create dpos rwSet, but the param is nil, should be error
	currBlk.Header.BlockHeight = 100
	err = impl.VerifyConsensusArgs(currBlk, nil)
	require.Error(t, err)

	// 3. the param and create dpos rwSet is equal.
	currBlk.Header.BlockHeight = 100
	//newEpoch := &syscontract.Epoch{
	//
	//}
	//localBz, _ := proto.Marshal(&consensuspb.BlockHeaderConsensusArgs{
	//	ConsensusType: int64(consensuspb.ConsensusType_DPOS),
	//	ConsensusData: &common.TxRWSet{
	//		TxId: moduleName,
	//		TxWrites: []*commonpb.TxWrite{
	//			{
	//				ContractName: syscontract.SystemContract_DPOS_STAKE.String(),
	//				Key:          []byte(dposmgr.KeyCurrentEpoch),
	//				Value:        bz,
	//			},
	//			{
	//				ContractName: syscontract.SystemContract_DPOS_STAKE.String(),
	//				Key:          dposmgr.ToEpochKey(fmt.Sprintf("%d", epoch.EpochId)),
	//				Value:        bz,
	//			},
	//		},
	//	},
	//})
	//currBlk.Header.ConsensusArgs = localBz
	err = impl.VerifyConsensusArgs(currBlk, nil)
	require.Error(t, err)
}
