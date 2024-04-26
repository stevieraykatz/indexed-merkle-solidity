package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/cockroachdb/pebble"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/mdehoog/indexed-merkle-tree/db"
	"github.com/mdehoog/indexed-merkle-tree/imt"
	"github.com/mdehoog/poseidon/poseidon"
)

func main() {
	levels := uint64(64)

	temp, _ := os.MkdirTemp("", "*")
	pDb, _ := pebble.Open(temp, &pebble.Options{})
	imtDb := db.NewPebble(pDb)
	tx := imtDb.NewTransaction()
	tree := imt.NewTreeWriter(tx, levels, fr.Bytes, poseidon.Hash[*fr.Element])

	fmt.Println(tree.Root())

	key := big.NewInt(123)
	exclusionProof, _ := tree.ProveExclusion(key)
	fmt.Println("exclusionProof")
	fmt.Println(exclusionProof)

	success, _ := exclusionProof.Valid(&tree.TreeReader)
	fmt.Println("Success")
	fmt.Println(success)

	insertProof, _ := tree.Insert(key, big.NewInt(456))
	fmt.Println("insertProof")
	fmt.Println(insertProof)

	updateProof, _ := tree.Update(key, big.NewInt(789))
	fmt.Println("updateProof")
	fmt.Println(updateProof)

	inclusionProof, _ := tree.ProveInclusion(key)
	fmt.Println("inclusionProof")
	fmt.Println(inclusionProof)
}
