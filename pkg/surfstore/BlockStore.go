package surfstore

import (
	context "context"
	"fmt"
	"slices"
)

type BlockStore struct {
	BlockMap map[string]*Block
	UnimplementedBlockStoreServer
}

func (bs *BlockStore) GetBlock(ctx context.Context, blockHash *BlockHash) (*Block, error) {

	h := blockHash.GetHash()
	if h == "" {
		return &Block{}, fmt.Errorf("invalid blockHash")
	}

	b, ok := bs.BlockMap[h]
	if !ok {
		return &Block{}, fmt.Errorf("invalid accessing")
	}

	return b, nil
}

func (bs *BlockStore) PutBlock(ctx context.Context, block *Block) (*Success, error) {

	h := GetBlockHashString(block.BlockData[:block.BlockSize])
	bs.BlockMap[h] = &Block{BlockSize: block.BlockSize, BlockData: block.BlockData}

	return &Success{Flag: true}, nil
}

// Given a list of hashes “in”, returns a list containing the
// subset of in that are stored in the key-value store
// TODO optimization: some sorting and binary search to improve search efficiency
func (bs *BlockStore) HasBlocks(ctx context.Context, blockHashesIn *BlockHashes) (*BlockHashes, error) {

	keys := make([]string, 0)
	for k := range bs.BlockMap {
		keys = append(keys, k)
	}

	blockHashesOut := make([]string, 0)

	for _, blockHashIn := range blockHashesIn.Hashes {
		if slices.Contains(keys, blockHashIn) {
			blockHashesOut = append(blockHashesOut, blockHashIn)
		}
	}

	// TODO: optimizations

	return &BlockHashes{Hashes: blockHashesOut}, nil
}

// This line guarantees all method for BlockStore are implemented
var _ BlockStoreInterface = new(BlockStore)

func NewBlockStore() *BlockStore {
	return &BlockStore{
		BlockMap: map[string]*Block{},
	}
}
