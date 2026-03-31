// SPDX-License-Identifier: MIT
contract OptimizedContract {
    uint[] public numbers;

    function addNumbers(uint[] calldata _nums) external {
        // 优化1：使用calldata参数，避免拷贝
        uint len = _nums.length;  // 优化2：缓存length

        for(uint i = 0; i < len; i++) {
            numbers.push(_nums[i]);
            // 注意：push操作还是要写Storage，这个很难避免
            // 但我们已经优化了参数和length读取
        }
    }

    function getSum() public view returns (uint) {
        uint sum = 0;
        uint len = numbers.length;  // 优化3：缓存Storage的length

        for(uint i = 0; i < len; i++) {
            sum += numbers[i];
        }
        return sum;
    }
}