// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract EnumDemo {
    enum OrderStatus {
        Pending,
        Paid,
        Shipped,
        Delivered,
        Cancelled
    }
    
    mapping(uint => OrderStatus) public orders;
    uint public orderCount;
    
    event OrderStatusChanged(uint indexed orderId, OrderStatus newStatus);
    
    // 创建订单
    function createOrder() public returns (uint) {
        uint orderId = orderCount++;
        orders[orderId] = OrderStatus.Pending;
        return orderId;
    }
    
    // 更新订单状态
    function updateStatus(uint orderId, OrderStatus newStatus) public {
        require(orders[orderId] != OrderStatus.Cancelled, "Order cancelled");
        orders[orderId] = newStatus;
        emit OrderStatusChanged(orderId, newStatus);
    }
    
    // 状态机：只能按顺序转换
    function shipOrder(uint orderId) public {
        require(orders[orderId] == OrderStatus.Paid, "Must be paid");
        orders[orderId] = OrderStatus.Shipped;
        emit OrderStatusChanged(orderId, OrderStatus.Shipped);
    }
    
    // 类型转换演示
    function getStatusValue(uint orderId) public view returns (uint) {
        return uint(orders[orderId]);
    }
}