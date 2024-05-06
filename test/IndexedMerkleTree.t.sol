// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.23;

import {Test, console, console2} from "forge-std/Test.sol";
import {IndexedMerkleTree} from "../src/IndexedMerkleTree.sol";
import {Strings} from "openzeppelin-contracts/contracts/utils/Strings.sol";

contract IndexedMerkleTreeTest is Test {
    using Strings for uint256;

    function test_verifyInclusion() public {
        (bytes32 root) = _goGetTree();
        console2.logBytes32(root);
        IndexedMerkleTree.Proof memory proof = _goGetIncProof(123);
        assertTrue(IndexedMerkleTree.verifyInclusionProof(proof, IndexedMerkleTree.HashType.Keccak256));
    }

    function test_verifyExclusion() public {
        (bytes32 root) = _goGetTree();
        console2.logBytes32(root);
        IndexedMerkleTree.Proof memory proof = _goGetExcProof(42);
        assertTrue(IndexedMerkleTree.verifyExclusionProof(42, proof, IndexedMerkleTree.HashType.Keccak256));
        assertFalse(IndexedMerkleTree.verifyExclusionProof(123, proof, IndexedMerkleTree.HashType.Keccak256));
        assertFalse(IndexedMerkleTree.verifyExclusionProof(1337, proof, IndexedMerkleTree.HashType.Keccak256));
    }

    function _goGetTree() internal returns (bytes32) {
        string[] memory inputs = new string[](2);
        string memory rootdir = vm.projectRoot();
        inputs[0] = string.concat(rootdir, "/go/imt");
        inputs[1] = "getTreeRoot";
        return abi.decode(vm.ffi(inputs), (bytes32));
    }

    function _goGetIncProof(uint256 key) internal returns (IndexedMerkleTree.Proof memory) {
        string[] memory inputs = new string[](3);
        string memory rootdir = vm.projectRoot();
        inputs[0] = string.concat(rootdir, "/go/imt");
        inputs[1] = "getInclusionProof";
        inputs[2] = key.toString();
        return abi.decode(vm.ffi(inputs), (IndexedMerkleTree.Proof));
    }

    function _goGetExcProof(uint256 key) internal returns (IndexedMerkleTree.Proof memory) {
        string[] memory inputs = new string[](3);
        string memory rootdir = vm.projectRoot();
        inputs[0] = string.concat(rootdir, "/go/imt");
        inputs[1] = "getExclusionProof";
        inputs[2] = key.toString();
        return abi.decode(vm.ffi(inputs), (IndexedMerkleTree.Proof));
    }
}
