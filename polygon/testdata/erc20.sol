// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract BxEToken is ERC20 {
    constructor() ERC20("Test BxE Token", "BXET") {
        uint256 initialSupply = 21000000 * (10 ** 18); // 2100W Token
        _mint(msg.sender, initialSupply);
    }
}