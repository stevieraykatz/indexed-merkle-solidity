## Documentation

**NOTE: This library is a WIP and unaudited -- consume with caution.**

Inspired by this paper from Aztec: https://docs.aztec.network/learn/concepts/storage/trees/indexed_merkle_tree

Consume IMT proofs on-chain with this solidity library. Built around the specific proof structure of this go-lang lib: https://github.com/mdehoog/indexed-merkle-tree

## Usage

```solidity
import {IndexedMerkleTree} from "./IndexedMerkleTree.sol";
...
IndexedMerkleTree.Proof memory proof = _getProofSomehow(); 
bool valid = IndexedMerkleTree.verify(proof)
```

### Build

```shell
$ forge build
```

### Test

The testing infrastructure calls out to go scripts implemented in the `go/` directory. Make sure to include the `--ffi` flag when running tests.  

```shell
$ forge test --ffi --vvv
```

### Help

```shell
$ forge --help
$ anvil --help
$ cast --help
```
