// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

library IndexedMerkleTree {
    
    struct Node {
        uint256 key;
        uint256 value;
        uint256 nextKey;
    }

    struct Proof {
        bytes32 root;
        uint256 size;
        uint256 index;
        Node node;
        uint256[] siblings;
    }

    function verify(
        Proof memory proof,
        uint256 levels
    ) internal pure returns (bool) {
        bytes32 computedHash = _hashNode(proof.node);
        uint256 index = proof.index;
        for (uint256 level = levels; level > 0; index /= 2) {
            level--;
            uint256 sibling = proof.siblings[level];
            if (sibling != 0) {
                if (index%2 == 0) {
                    computedHash = keccak256(abi.encode(sibling, computedHash));
                } else {
                    computedHash = keccak256(abi.encode(computedHash, sibling));
                }
            }
        }
        computedHash = keccak256(abi.encode(computedHash, proof.size));
        return computedHash == proof.root;
    }

    function _hashNode(Node memory node) internal pure returns (bytes32) {
        return keccak256(abi.encode(node.key, node.value, node.nextKey));
    }
}