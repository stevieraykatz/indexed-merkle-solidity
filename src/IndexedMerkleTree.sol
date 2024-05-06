// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

import {PoseidonT3} from "poseidon-solidity/PoseidonT3.sol";
import {PoseidonT4} from "poseidon-solidity/PoseidonT4.sol";

/// @title Indexed Merkle Tree Verifier
///
/// @notice This library is strictly tied to the IMT implementation found in this
/// library: https://github.com/mdehoog/indexed-merkle-tree. All proofs must be formated
/// according to the structures provided.
/// Be sure to strictly pair trees and proofs with a specific hashing algorithm. This library implements
/// a mechanism for validating trees hashed using keccak256 or Poseidon as the hashing function.
library IndexedMerkleTree {
    /// @notice The supported hashing types.
    enum HashType {
        Keccak256,
        Poseidon
    }

    /// @notice A Node represents a single leaf in an IMT
    struct Node {
        /// @dev The unique identifier in the tree upon which the tree is sorted
        uint256 key;
        /// @dev The node index.
        uint256 index;
        /// @dev The value associated with this key
        uint256 value;
        /// @dev The linked-list property of an IMT depends on each node knowing the next
        /// sorted value.
        uint256 nextKey;
    }

    /// @notice The data struct containing all associated fields necessary for performing
    /// inclusion/exclusion proofs
    struct Proof {
        /// @dev The tree's root hash
        bytes32 root;
        /// @dev The tree's size
        uint256 size;
        /// @dev The Node struct containing the details for the leaf node
        Node node;
        ///@dev The sibling proof list, similar to other Merkle Tree proofs
        uint256[] siblings;
    }

    /// @notice Verifies the given inclusion proof using the provided hashing type.
    ///
    /// @param proof The proof to verify inclusion for.
    /// @param ht The hash type to use.
    function verifyInclusionProof(Proof memory proof, HashType ht) internal pure returns (bool) {
        // Ensure the proof's node is included in the tree.
        return _verify(proof, ht);
    }

    /// @notice Verifies the given exclusion proof for the provided `key` using the provided hashing type.
    ///
    /// @param key The key to verify exclusion for.
    /// @param proof The proof of the low nullifier to verify.
    /// @param ht The hash type to use.
    function verifyExclusionProof(uint256 key, Proof memory proof, HashType ht) internal pure returns (bool) {
        // Ensure that the proof's node is the key's low nullifier and check that this low nullifier is included in the tree.
        return proof.node.key < key && (proof.node.nextKey > key || proof.node.nextKey == 0) && _verify(proof, ht);
    }

    /// @notice Verifies that a proof is valid.
    ///
    /// @dev Ensure proofs are generated using the golang lib linked above
    /// @dev This method is generic over the hashing alr
    ///
    /// @param proof A complete proof struct which describes the tree and the associated
    /// data necessary to validate a provided node.
    /// @param ht The hash type to use.
    function _verify(Proof memory proof, HashType ht) private pure returns (bool) {
        (
            function (Node memory) pure returns (uint256) hashNode,
            function (uint256, uint256) pure returns (uint256) hashPair
        ) = ht == HashType.Keccak256 ? (_hashNode, _hashPair) : (_poseidonHashNode, _poseidonHashPair);

        uint256 computedHash = hashNode(proof.node);

        uint256 index = proof.node.index;
        for (uint256 level = proof.siblings.length; level > 0; index /= 2) {
            level--;

            uint256 sibling = proof.siblings[level];
            if (sibling != 0) {
                (uint256 l, uint256 r) = index % 2 == 0 ? (sibling, computedHash) : (computedHash, sibling);
                computedHash = hashPair(l, r);
            }
        }

        computedHash = _hashPair(computedHash, proof.size);
        return bytes32(computedHash) == proof.root;
    }

    function _hashNode(Node memory node) private pure returns (uint256) {
        return uint256(keccak256(abi.encode(node.key, node.value, node.nextKey)));
    }

    function _hashPair(uint256 a, uint256 b) private pure returns (uint256) {
        return uint256(keccak256(abi.encode(a, b)));
    }

    function _poseidonHashNode(Node memory node) private pure returns (uint256) {
        uint256[3] memory _node = [node.key, node.value, node.nextKey];
        return PoseidonT4.hash(_node);
    }

    function _poseidonHashPair(uint256 a, uint256 b) private pure returns (uint256) {
        uint256[2] memory _pair = [a, b];
        return PoseidonT3.hash(_pair);
    }
}
