package main

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"os"

	"github.com/cockroachdb/pebble"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mdehoog/indexed-merkle-tree/db"
	"github.com/mdehoog/indexed-merkle-tree/imt"
	"github.com/mdehoog/poseidon/poseidon"
)

func main() {
	levels := uint64(16)
	temp, _ := os.MkdirTemp("", "*")
	pDb, _ := pebble.Open(temp, &pebble.Options{})
	imtDb := db.NewPebble(pDb)
	tx := imtDb.NewTransaction()
	tree := imt.NewTreeWriter(tx, levels, fr.Bytes, poseidon.Hash[*fr.Element])
	key := big.NewInt(123)
	tree.Insert(key, big.NewInt(456))

	if len(os.Args) < 2 {
		fmt.Println("Please provide a valid method name and optional key arg")
		os.Exit(1)
	}

	if os.Args[1] == "getTreeRoot" {
		getTreeRoot(tree)

	} else if os.Args[1] == "getInclusionProof" {
		key, success := new(big.Int).SetString(os.Args[2], 10)
		if !success {
			fmt.Println("Error parsing key argument:", os.Args[2])
			os.Exit(1)
		}
		getInclusionProof(key, tree)

	} else if os.Args[1] == "getExclusionProof" {
		key, success := new(big.Int).SetString(os.Args[2], 10)
		if !success {
			fmt.Println("Error parsing key argument:", os.Args[2])
			os.Exit(1)
		}
		getExcluisonProof(key, tree)
	}

}

func getTreeRoot(t *imt.TreeWriter) {
	root, _ := t.Root()
	fmt.Println(hexutil.Encode(padBytes(root.Bytes(), 32)))
}

func getInclusionProof(key *big.Int, t *imt.TreeWriter) {
	inclusionProof, _ := t.ProveInclusion(key)
	// Encode the tree-data for the proof
	// Solidity structs head memory locations
	// struct Proof {
	//     bytes32 root;			// 0x20
	//     uint256 size;			// 0x40
	//     uint256 index;			// 0x60
	//     Node node;				// 0x80-0xc0
	//     uint256[] siblings;		// 0xe0
	// }
	// struct Node {
	//     uint256 key;
	//     uint256 value;
	//     uint256 nextKey;
	// }
	// Encode the head location of the struct tuple
	tupleHead, _ := hexutil.Decode("0x20")
	head := hexutil.Encode(padBytes(tupleHead, 32))
	// Encode the data
	root := hexutil.Encode(padBytes(inclusionProof.Root.Bytes(), 32))
	size := hexutil.Encode(padBytes(Uint64ToBytes(inclusionProof.Size), 32))
	index := hexutil.Encode(padBytes(Uint64ToBytes(inclusionProof.Index), 32))
	nodeKey := hexutil.Encode(padBytes(inclusionProof.Node.Key.Bytes(), 32))
	nodeValue := hexutil.Encode(padBytes(inclusionProof.Node.Value.Bytes(), 32))
	nodeNext := hexutil.Encode(padBytes(inclusionProof.Node.NextKey.Bytes(), 32))
	result := head + root[2:] + size[2:] + index[2:] + nodeKey[2:] + nodeValue[2:] + nodeNext[2:]

	// Encode the head location of the dynamicly set data
	headLocation, _ := hexutil.Decode("0xe0")
	result += hexutil.Encode(padBytes(headLocation, 32))[2:]
	// Encode the length of the array for dyanamic var assignment w/ abi.decode()
	result += hexutil.Encode(padBytes(Uint64ToBytes(uint64(len(inclusionProof.Siblings))), 32))[2:]

	// Encode the members of the array
	for i := 0; i < len(inclusionProof.Siblings); i++ {
		result += hexutil.Encode(padBytes(inclusionProof.Siblings[i].Bytes(), 32))[2:]
	}
	// valid, _ := inclusionProof.Valid(&t.TreeReader)
	// fmt.Println(valid)
	fmt.Println(result)
}

func getExcluisonProof(key *big.Int, t *imt.TreeWriter) {
	exclusionProof, _ := t.ProveExclusion(key)
	fmt.Println(exclusionProof)
}

func padBytes(b []byte, length int) []byte {
	if len(b) >= length {
		return b
	}

	padding := make([]byte, length-len(b))
	return append(padding, b...)
}

// Uint64ToBytes converts the given uint64 value to slice of bytes.
func Uint64ToBytes(val uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, val)

	return b
}
