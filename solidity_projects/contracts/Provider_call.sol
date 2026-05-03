// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract ProviderCall{
    uint256 public   num  ;
    uint256 public count = 666;

    function  setNum(uint256 _num) public  {
        num = _num;
    }

    function getNum() public view returns  (uint256){
        return num;
    }

    function encode() public  pure  returns (bytes memory){
        return abi.encodeWithSignature("getNum()");
    }

    function encodeUint256() public  pure  returns (bytes memory){
        return abi.encodeWithSignature("setNum(uint256)");
    }

}