// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract OrdinalsEVM is ERC20 {
    constructor() ERC20("OrdinalsEVM", "OEVM") {
        uint256 initialSupply = 1000000000 * (10 ** 18); // 10亿乘以10的18次方，以支持18位小数
        _mint(msg.sender, initialSupply);
    }
}