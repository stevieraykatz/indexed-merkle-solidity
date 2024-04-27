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
	"github.com/wealdtech/go-merkletree/v2/keccak256"
)

func main() {
	levels := uint64(16)
	temp, _ := os.MkdirTemp("", "*")
	pDb, _ := pebble.Open(temp, &pebble.Options{})
	imtDb := db.NewPebble(pDb)
	tx := imtDb.NewTransaction()
	tree := imt.NewTreeWriter(tx, levels, fr.Bytes, hashFn)
	// tree := imt.NewTreeWriter(tx, levels, fr.Bytes, poseidon.Hash[*fr.Element])
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
	} else if os.Args[1] == "testHash" {
		a, successA := new(big.Int).SetString(os.Args[2], 10)
		if !successA {
			fmt.Println("Error parsing key argument:", os.Args[2])
			os.Exit(1)
		}
		b, successB := new(big.Int).SetString(os.Args[3], 10)
		if !successB {
			fmt.Println("Error parsing key argument:", os.Args[2])
			os.Exit(1)
		}
		testHash(a, b)
	}
}

func getTreeRoot(t *imt.TreeWriter) {
	root, _ := t.Root()
	fmt.Println(hexutil.Encode(padBytes(root.Bytes(), 32)))
}

func getInclusionProof(key *big.Int, t *imt.TreeWriter) {
	inclusionProof, _ := t.ProveInclusion(key)
	_encodeForFoundry(inclusionProof)
}

func getExcluisonProof(key *big.Int, t *imt.TreeWriter) {
	exclusionProof, _ := t.ProveExclusion(key)
	_encodeForFoundry(exclusionProof)
}

func testHash(a *big.Int, b *big.Int) {
	I := make([]*big.Int, 2)
	I[0] = a
	I[1] = b
	h, _ := hashFn(I)
	fmt.Println(h)
}

func _encodeForFoundry(p *imt.Proof) {
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
	root := hexutil.Encode(padBytes(p.Root.Bytes(), 32))
	size := hexutil.Encode(padBytes(Uint64ToBytes(p.Size), 32))
	index := hexutil.Encode(padBytes(Uint64ToBytes(p.Index), 32))
	nodeKey := hexutil.Encode(padBytes(p.Node.Key.Bytes(), 32))
	nodeValue := hexutil.Encode(padBytes(p.Node.Value.Bytes(), 32))
	nodeNext := hexutil.Encode(padBytes(p.Node.NextKey.Bytes(), 32))
	result := head + root[2:] + size[2:] + index[2:] + nodeKey[2:] + nodeValue[2:] + nodeNext[2:]

	// Encode the head location of the dynamicly set data
	headLocation, _ := hexutil.Decode("0xe0")
	result += hexutil.Encode(padBytes(headLocation, 32))[2:]
	// Encode the length of the array for dyanamic var assignment w/ abi.decode()
	result += hexutil.Encode(padBytes(Uint64ToBytes(uint64(len(p.Siblings))), 32))[2:]

	// Encode the members of the array
	for i := 0; i < len(p.Siblings); i++ {
		result += hexutil.Encode(padBytes(p.Siblings[i].Bytes(), 32))[2:]
	}
	// valid, _ := inclusionProof.Valid(&t.TreeReader)
	// fmt.Println(valid)
	fmt.Println(result)
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

// Wrapper for passing keccak256 into Tree instantiation
func hashFn(E []*big.Int) (*big.Int, error) {
	k := keccak256.New()
	b := make([][]byte, len(E))
	for i := 0; i < len(E); i++ {
		b[i] = padBytes(E[i].Bytes(), 32)
	}
	h := k.Hash(b...)
	I := big.Int{}
	I.SetBytes(h)
	return &I, nil
}
