// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

import {PoseidonT3} from "poseidon-solidity/PoseidonT3.sol";
import {PoseidonT4} from "poseidon-solidity/PoseidonT4.sol";

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

    function verify(Proof memory proof) internal pure returns (bool) {
        uint256 computedHash = _hashNode(proof.node);
        uint256 index = proof.index;
        for (uint256 level = proof.siblings.length; level > 0; index /= 2) {
            level--;
            uint256 sibling = proof.siblings[level];
            if (sibling != 0) {
                if (index%2 == 0) {
                    computedHash = _hashPair(sibling, computedHash);
                } else {
                    computedHash = _hashPair(computedHash, sibling);
                }
            }
        }
        computedHash = _hashPair(computedHash, proof.size);
        return bytes32(computedHash) == proof.root;
    }

    function _hashPair(uint256 a, uint256 b) private pure returns (uint256) {
        return uint256(keccak256(abi.encode(a,b)));
    }

    function _hashNode(Node memory node) private pure returns (uint256) {
        return uint256(keccak256(abi.encode(node.key, node.value, node.nextKey)));
    }

    function verifyPoseidon(Proof memory proof) internal pure returns (bool) {
        uint256 computedHash = _poseidonHashNode(proof.node);
        uint256 index = proof.index;
        for (uint256 level = proof.siblings.length; level > 0; index /= 2) {
            level--;
            uint256 sibling = proof.siblings[level];
            if (sibling != 0) {
                if (index%2 == 0) {
                    computedHash = _poseidonHashPair(sibling, computedHash);
                } else {
                    computedHash = _poseidonHashPair(computedHash, sibling);
                }
            }
        }
        computedHash = _hashPair(computedHash, proof.size);     
        return bytes32(computedHash) == proof.root;
    }

    function _poseidonHashPair(uint256 a, uint256 b) private pure returns (uint256) {
        uint256[2] memory _pair = [
            a, 
            b
        ];
        return PoseidonT3.hash(_pair);
    }

    function _poseidonHashNode(Node memory node) private pure returns (uint256) {
        uint256[3] memory _node = [
            node.key,
            node.value,
            node.nextKey
        ];
        return PoseidonT4.hash(_node);
    }
}