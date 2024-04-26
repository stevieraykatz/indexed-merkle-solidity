// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.23;

import {Test, console, console2} from "forge-std/Test.sol";
import {IndexedMerkleTree} from "../src/IndexedMerkleTree.sol";

contract IndexedMerkleTreeTest is Test {

    function test_verify() public {
        (bytes32 root) = _goGetTree();
        console2.logBytes32(root);
        IndexedMerkleTree.Proof memory proof = _goGetIncProof();
        assert(IndexedMerkleTree.verify(proof));
    }

    function _goGetTree() internal returns (bytes32) {
        string[] memory inputs = new string[](2);
        inputs[0] = "test/../go/imt";
        inputs[1] = "getTreeRoot";
        return abi.decode(vm.ffi(inputs), (bytes32));
    }


    function _goGetIncProof() internal returns (IndexedMerkleTree.Proof memory) {
        string[] memory inputs = new string[](3);
        inputs[0] = "test/../go/imt";
        inputs[1] = "getInclusionProof";
        inputs[2] = "123";
        return abi.decode(vm.ffi(inputs), (IndexedMerkleTree.Proof));
    }
}
