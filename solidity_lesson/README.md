# Solidity入门到精通系列课程

## 课程概述
本课程分为三个阶段，涵盖从Solidity基础入门到框架实战的完整学习路径。

### 第一阶段：基础入门（10门课程）
掌握Solidity核心语法和基础编程能力。包括智能合约概念、开发环境搭建、EVM存储机制、数据类型、数组、映射、结构体、函数、修饰符、控制流语句等基础知识，并完成简单代币合约项目。

### 第二阶段：进阶编程（11门课程）
深入学习合约高级特性和安全开发。包括合约继承、库合约、事件、错误处理、合约间调用、智能合约安全、Gas优化技巧、设计模式等进阶主题，并完成NFT市场实战项目。

### 第三阶段：开发框架与实战（7门课程）
掌握专业开发工具和实战项目开发。学习Hardhat和Foundry两大开发框架、单元测试、Ethers.js、MetaMask集成等开发技能，通过众筹平台、多签钱包、NFT盲盒等完整项目提升实战能力。

## 课程目录

### 第一阶段：基础入门

#### 1.1 课程导言与环境搭建
**状态**: ✅  
**学时**: 2小时  
**教学内容**:
- 智能合约基本概念
- Solidity在区块链中的作用
- Remix IDE使用
- 第一个Hello World合约

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xiaoeknow.com/p/course/video/v_6909d8d7e4b0694c5b494712)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - 合约基础与环境搭建(课件).pdf](lesson1.1/Solidity智能合约开发%20-%20合约基础与环境搭建(课件).pdf)
- 学习资料：[Solidity智能合约开发 - 合约基础与环境搭建(学习资料).md](lesson1.1/Solidity智能合约开发%20-%20合约基础与环境搭建(学习资料).md)

**代码示例**:
- [HelloWorld.sol](lesson1.1/study_code/HelloWorld.sol)

---

#### 1.2 EVM存储结构
**状态**: ✅  
**学时**: 2小时  
**教学内容**:
- EVM存储架构
- storage/memory/calldata区别
- gas消耗对比
- 存储优化基础

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xiaoeknow.com/p/course/video/v_690c4990e4b0694c5b4a9790)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - EVM存储结构（课件）.pdf](lesson1.2/Solidity智能合约开发%20-%20EVM存储结构（课件）.pdf)
- 学习资料：[Solidity智能合约开发 - EVM存储结构（学习资料）.md](lesson1.2/Solidity智能合约开发%20-%20EVM存储结构（学习资料）.md)

**代码示例**:
- [OptimizedContract.sol](lesson1.2/study_code/OptimizedContract.sol)
- [SumContract.sol](lesson1.2/study_code/SumContract.sol)
- [Task.sol](lesson1.2/study_code/Task.sol)

---

#### 2.1 数据类型基础
**状态**: ✅  
**学时**: 2.5小时  
**教学内容**:
- 值类型：bool/int/uint/address/bytes
- 引用类型：arrays/struct/mapping
- 运算符使用
- 类型转换

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xiaoeknow.com/p/course/video/v_690f308ce4b0694ca13ea640)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - 数据类型基础(课件).pdf](lesson2.1/Solidity智能合约开发-%20数据类型基础(课件).pdf)
- 学习资料：[Solidity智能合约开发 - 数据类型基础(学习资料).md](lesson2.1/Solidity智能合约开发%20-%20数据类型基础(学习资料).md)

**代码示例**:
- [SafeCounter.sol](lesson2.1/study_code/SafeCounter.sol)
- [SignatureExample.sol](lesson2.1/study_code/SignatureExample.sol)
- [VotingSystem.sol](lesson2.1/study_code/VotingSystem.sol)

---

#### 3.1 数组
**状态**: ✅  
**学时**: 2小时  
**教学内容**:
- 定长数组vs动态数组
- 数组操作：push/pop/length
- 多维数组
- 数组遍历优化

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xiaoeknow.com/p/course/video/v_69157420e4b0694ca141fe73)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - 数组(课件).pdf](lesson3.1/Solidity智能合约开发%20-%20数组(课件).pdf)
- 学习资料：[Solidity智能合约开发 - 数组(学习资料).md](lesson3.1/Solidity智能合约开发%20-%20数组(学习资料).md)

**代码示例**:
- [ArrayBasics.sol](lesson3.1/study_code/ArrayBasics.sol)
- [ArrayDeletion.sol](lesson3.1/study_code/ArrayDeletion.sol)
- [AdvancedArrays.sol](lesson3.1/study_code/AdvancedArrays.sol)
- [GasOptimization.sol](lesson3.1/study_code/GasOptimization.sol)
- [UserManager.sol](lesson3.1/study_code/UserManager.sol)

---

#### 3.2 映射和结构体
**状态**: ✅  
**学时**: 2小时  
**教学内容**:
- mapping基础和限制
- 嵌套mapping
- struct定义和使用
- struct与mapping结合

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xiaoeknow.com/p/course/video/v_69185114e4b0694c5b50a6d1)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - 映射和结构体(课件).pdf](lesson3.2/Solidity智能合约开发%20-%20映射和结构体(课件).pdf)
- 学习资料：[Solidity智能合约开发 - 映射和结构体(学习资料).md](lesson3.2/Solidity智能合约开发%20-%20映射和结构体(学习资料).md)

**代码示例**:
- [MappingBasics.sol](lesson3.2/study_code/MappingBasics.sol)
- [AdvancedMapping.sol](lesson3.2/study_code/AdvancedMapping.sol)
- [StructExample.sol](lesson3.2/study_code/StructExample.sol)
- [MappingStructCombo.sol](lesson3.2/study_code/MappingStructCombo.sol)

---

#### 3.3 函数与修饰符
**状态**: ✅  
**学时**: 2.5小时  
**教学内容**:
- 函数可见性：public/private/internal/external
- 状态修饰符：view/pure/payable
- 自定义modifier
- 函数重载

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xiaoeknow.com/p/course/video/v_691c0db6e4b0694c5b5214b9)

**文档资料**: 
- 课件：[Solidity智能合约开发 - 函数与修饰符(课件).pdf](lesson3.3/Solidity智能合约开发%20-%20函数与修饰符(课件).pdf)
- 学习资料：[Solidity智能合约开发 - 函数与修饰符(学习资料).md](lesson3.3/Solidity智能合约开发%20-%20函数与修饰符(学习资料).md)

**代码示例**:
- [VisibilityDemo.sol](lesson3.3/study_code/VisibilityDemo.sol)
- [StateModifiersDemo.sol](lesson3.3/study_code/StateModifiersDemo.sol)
- [ModifierDemo.sol](lesson3.3/study_code/ModifierDemo.sol)
- [OverloadingDemo.sol](lesson3.3/study_code/OverloadingDemo.sol)
- [RoleManagement.sol](lesson3.3/study_code/RoleManagement.sol)

---

#### 4.1 控制流语句
**状态**: ✅  
**学时**: 2小时  
**教学内容**:
- if-else条件语句
- for/while/do-while循环
- break和continue
- require/assert/revert

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xiaoeknow.com/p/course/video/v_691c3798e4b0694c5b5250ef)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - 控制流语句(课件).pdf](lesson4.1/Solidity智能合约开发%20-%20控制流语句(课件).pdf)
- 学习资料：[Solidity智能合约开发 - 控制流语句(学习资料).md](lesson4.1/Solidity智能合约开发%20-%20控制流语句(学习资料).md)

**代码示例**:
- [ConditionalDemo.sol](lesson4.1/study_code/ConditionalDemo.sol)
- [LoopDemo.sol](lesson4.1/study_code/LoopDemo.sol)
- [GasProblemDemo.sol](lesson4.1/study_code/GasProblemDemo.sol)
- [ErrorHandlingDemo.sol](lesson4.1/study_code/ErrorHandlingDemo.sol)
- [VotingSystem.sol](lesson4.1/study_code/VotingSystem.sol)
- [Crowdfunding.sol](lesson4.1/study_code/Crowdfunding.sol)

---

#### 4.2 特殊类型与全局变量
**状态**: ✅  
**学时**: 2小时  
**教学内容**:
- address和address payable
- 全局变量：msg/block/tx
- enum枚举
- constant和immutable

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xiaoeknow.com/p/course/video/v_691d7f23e4b0694ca1464b2d)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - 特殊类型与全局变量(课件).pdf](lesson4.2/Solidity智能合约开发%20-%20特殊类型与全局变量(课件).pdf)
- 学习资料：[Solidity智能合约开发 - 特殊类型与全局变量(学习资料).md](lesson4.2/Solidity智能合约开发%20-%20特殊类型与全局变量(学习资料).md)

**代码示例**:
- [AddressDemo.sol](lesson4.2/study_code/AddressDemo.sol)
- [MsgDemo.sol](lesson4.2/study_code/MsgDemo.sol)
- [BlockDemo.sol](lesson4.2/study_code/BlockDemo.sol)
- [EnumDemo.sol](lesson4.2/study_code/EnumDemo.sol)
- [ConstantDemo.sol](lesson4.2/study_code/ConstantDemo.sol)

---

#### 5.1 基础项目：简单代币合约
**状态**: ✅  
**学时**: 3小时  
**教学内容**:
- ERC20标准介绍
- 代币核心功能实现
- 铸造和销毁
- 部署和测试

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xiaoeknow.com/p/course/video/v_691ec214e4b0694ca146ebc5)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - 简单代币合约(课件).pdf](lesson5.1/Solidity智能合约开发%20-%20简单代币合约(课件).pdf)
- 学习资料：[Solidity智能合约开发 - 简单代币合约(学习资料).md](lesson5.1/Solidity智能合约开发%20-%20简单代币合约(学习资料).md)

**代码示例**:
- [MyToken.sol](lesson5.1/study_code/MyToken.sol)

---

### 第二阶段：进阶编程

#### 6.1 合约继承
**状态**: ✅  
**学时**: 2.5小时  
**教学内容**:
- 单继承和多重继承
- super关键字
- 构造函数继承
- 函数重写
- 抽象合约和接口

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xiaoeknow.com/p/course/video/v_691f003fe4b0694ca1472c60)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - 合约继承(课件).pdf](lesson6.1/Solidity智能合约开发%20-%20合约继承(课件).pdf)
- 学习资料：[Solidity智能合约开发 - 合约继承(学习资料).md](lesson6.1/Solidity智能合约开发%20-%20合约继承(学习资料).md)

**代码示例**:
- [ParentDemo.sol](lesson6.1/study_code/ParentDemo.sol)
- [Shape.sol](lesson6.1/study_code/Shape.sol)
- [Super.sol](lesson6.1/study_code/Super.sol)
- [Ownable.sol](lesson6.1/study_code/Ownable.sol)
- [ERC20.sol](lesson6.1/study_code/ERC20.sol)
- [MyCustomToken.sol](lesson6.1/study_code/MyCustomToken.sol)
- [MyConcreteToken.sol](lesson6.1/study_code/MyConcreteToken.sol)

---

#### 6.2 库合约Library
**状态**: ✅  
**学时**: 2小时  
**教学内容**:
- 库合约定义和特性
- using for语法
- 内部库和外部库
- OpenZeppelin库介绍

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xiaoeknow.com/p/course/video/v_69204b00e4b0694ca148063f)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - 库合约Library(课件).pptx](lesson6.2/Solidity智能合约开发%20-%20库合约Library(课件).pptx)
- 学习资料：[Solidity智能合约开发 - 库合约Library(学习资料).md](lesson6.2/Solidity智能合约开发%20-%20库合约Library(学习资料).md)

**代码示例**:
- [Calculator.sol](lesson6.2/study_code/Calculator.sol)
- [LibUser.sol](lesson6.2/study_code/LibUser.sol)
- [BankAccount.sol](lesson6.2/study_code/BankAccount.sol)
- [SimpleNFT.sol](lesson6.2/study_code/SimpleNFT.sol)
- [Whitelist.sol](lesson6.2/study_code/Whitelist.sol)

---

#### 7.1 事件Events
**状态**: ✅  
**学时**: 2小时  
**教学内容**:
- 事件定义和用途
- indexed参数详解
- 匿名事件
- 事件最佳实践
- 事件查询和监听
- 实际应用场景

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xiaoeknow.com/p/course/video/v_69219343e4b0694ca14889ef)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - 事件Events(课件).pdf](lesson7.1/Solidity智能合约开发%20-%20事件Events(课件).pdf)
- 学习资料：[Solidity智能合约开发 - 事件Events(学习资料).md](lesson7.1/Solidity智能合约开发%20-%20事件Events(学习资料).md)

**代码示例**:
- [EventDemo.sol](lesson7.1/study_code/EventDemo.sol)
- [AnonymousEventDemo.sol](lesson7.1/study_code/AnonymousEventDemo.sol)
- [SimpleERC20.sol](lesson7.1/study_code/SimpleERC20.sol)

---

#### 7.2 错误处理和自定义错误
**状态**: ✅  
**学时**: 2小时  
**教学内容**:
- require/assert/revert详解
- 自定义错误(0.8.4+)
- try-catch异常捕获
- 错误处理最佳实践
- Gas消耗对比分析

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xiaoeknow.com/p/course/video/v_6921a5dde4b0694ca148a830)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - 错误处理和自定义错误(课件).pdf](lesson7.2/Solidity智能合约开发%20-%20错误处理和自定义错误(课件).pdf)
- 学习资料：[Solidity智能合约开发 - 错误处理和自定义错误(学习资料).md](lesson7.2/Solidity智能合约开发%20-%20错误处理和自定义错误(学习资料).md)

**代码示例**:
- [ErrorHandlingDemo.sol](lesson7.2/study_code/ErrorHandlingDemo.sol)
- [GasComparison.sol](lesson7.2/study_code/GasComparison.sol)
- [Auction.sol](lesson7.2/study_code/Auction.sol)
- [MockToken.sol](lesson7.2/study_code/MockToken.sol)

---

#### 8.1 合约间调用
**状态**: ✅  
**学时**: 2.5小时  
**教学内容**:
- 接口调用
- call/delegatecall/staticcall
- 合约创建：new和create2
- 安全的外部调用

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xiaoeknow.com/p/course/video/v_6926a2abe4b0694ca14b5eea)  
**文档资料**: 
- 课件：[Solidity智能合约开发-合约间调用(课件).pdf](lesson8.1/Solidity智能合约开发-合约间调用(课件).pdf)
- 学习资料：[Solidity智能合约开发 - 合约间调用(学习资料).md](lesson8.1/Solidity智能合约开发%20-%20合约间调用(学习资料).md)

**代码示例**:
- [CallMethodsDemo.sol](lesson8.1/study_code/CallMethodsDemo.sol)
- [Create2Demo.sol](lesson8.1/study_code/Create2Demo.sol)
- [ProxyDemo.sol](lesson8.1/study_code/ProxyDemo.sol)
- [ReentrancyDemo.sol](lesson8.1/study_code/ReentrancyDemo.sol)
- [SimpleToken.sol](lesson8.1/study_code/SimpleToken.sol)

---

#### 8.2 智能合约安全基础
**状态**: ✅  
**学时**: 3小时  
**教学内容**:
- 重入攻击
- 整数溢出
- 权限控制漏洞
- 拒绝服务攻击
- 前端运行
- 安全检查清单

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xiaoeknow.com/p/course/video/v_6927f3b5e4b0694c5b582697)  
**文档资料**: 
- 课件：[Solidity智能合约 - 合约安全基础(课件).pdf](lesson8.2/Solidity智能合约%20-%20合约安全基础(课件).pdf)
- 学习资料：[Solidity智能合约开发 - 智能合约安全基础(学习资料).md](lesson8.2/Solidity智能合约开发%20-%20智能合约安全基础(学习资料).md)

**代码示例**:
- [CommitRevealVoting.sol](lesson8.2/study_code/CommitRevealVoting.sol)
- [OverflowProtected.sol](lesson8.2/study_code/OverflowProtected.sol)
- [SafeVaultCEI.sol](lesson8.2/study_code/SafeVaultCEI.sol)
- [SafeVaultOZ.sol](lesson8.2/study_code/SafeVaultOZ.sol)
- [SecureToken.sol](lesson8.2/study_code/SecureToken.sol)
- [UnsafeReward.sol](lesson8.2/study_code/UnsafeReward.sol)
- [VulnerableVault.sol](lesson8.2/study_code/VulnerableVault.sol)

---

#### 9.1 Gas优化技巧
**状态**: ✅  
**学时**: 2.5小时  
**教学内容**:
- Gas成本分析
- 存储优化
- 数据类型优化
- 函数优化
- 批量操作
- unchecked使用

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xiaoeknow.com/p/course/video/v_69296526e4b0694c5b591eb6)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - Gas优化技巧(课件).pdf](lesson9.1/Solidity智能合约开发%20-%20Gas优化技巧(课件).pdf)
- 学习资料：[Solidity智能合约开发 - Gas优化技巧(学习资料).md](lesson9.1/Solidity智能合约开发%20-%20Gas优化技巧(学习资料).md)

**代码示例**:
- [BatchOperations.sol](lesson9.1/study_code/BatchOperations.sol)
- [EventVsStorage.sol](lesson9.1/study_code/EventVsStorage.sol)
- [FunctionVisibility.sol](lesson9.1/study_code/FunctionVisibility.sol)
- [LocalStoragePointer.sol](lesson9.1/study_code/LocalStoragePointer.sol)
- [PackingGood.sol](lesson9.1/study_code/PackingGood.sol)
- [UncheckedDemo.sol](lesson9.1/study_code/UncheckedDemo.sol)

---

#### 9.2 智能合约设计模式
**状态**: ✅  
**学时**: 2.5小时  
**教学内容**:
- 访问控制模式
- 提现模式
- 状态机模式
- 代理模式
- 工厂模式
- 紧急停止模式

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_692e8a59e4b0694c5b5b4610)  
**文档资料**: 
- 课件：[Solidity智能合约 - 设计模式(课件资料).pdf](lesson9.2/Solidity智能合约%20-%20设计模式(课件资料).pdf)
- 学习资料：[Solidity智能合约开发 - 智能合约设计模式(学习资料).md](lesson9.2/Solidity智能合约开发%20-%20智能合约设计模式(学习资料).md)

**代码示例**:
- [AccessControlDemo.sol](lesson9.2/study_code/AccessControlDemo.sol)
- [WithdrawalPattern.sol](lesson9.2/study_code/WithdrawalPattern.sol)
- [StateMachineCrowdfunding.sol](lesson9.2/study_code/StateMachineCrowdfunding.sol)
- [ProxyPattern.sol](lesson9.2/study_code/ProxyPattern.sol)
- [FactoryPattern.sol](lesson9.2/study_code/FactoryPattern.sol)
- [PausablePattern.sol](lesson9.2/study_code/PausablePattern.sol)

---

#### 10.1 进阶项目：NFT市场
**状态**: ✅  
**学时**: 4小时  
**教学内容**:
- ERC721标准
- NFT铸造
- 市场上架/下架
- 买卖功能
- 版税系统
- 拍卖功能

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_69427527e4b0694ca15aaeeb)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - NFT市场(课件资料).pdf](lessonA10.1/Solidity智能合约开发%20-%20NFT市场(课件资料).pdf)
- 学习资料：[Solidity智能合约开发 - NFT市场(学习资料).md](lessonA10.1/Solidity智能合约开发%20-%20NFT市场(学习资料).md)

**代码示例**:
- [MyNFT.sol](lessonA10.1/study_code/MyNFT.sol)
- [NFTMarketplace.sol](lessonA10.1/study_code/NFTMarketplace.sol)
- [MyNFTWithRoyalty.sol](lessonA10.1/study_code/MyNFTWithRoyalty.sol)

---

### 第三阶段：开发框架与实战

#### 11.1 Hardhat 3 环境搭建与合约部署
**状态**: ✅  
**学时**: 4.5小时  
**教学内容**:
- Node.js和npm安装
- Hardhat 3项目初始化
- 项目结构详解
- hardhat.config.ts配置
- 网络配置
- 合约编译
- 编写部署脚本
- Ignition部署系统
- 本地网络部署
- 测试网部署
- 合约验证

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_69453d93e4b0694c5b660ef5)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - Hardhat 3 环境搭建(课件资料).pdf](lessonA11.1/Solidity智能合约开发%20-%20Hardhat%203%20环境搭建(课件资料).pdf)
- 学习资料：[Solidity智能合约开发 - Hardhat 3 环境搭建与合约部署(学习资料).md](lessonA11.1/Solidity智能合约开发%20-%20Hardhat%203%20环境搭建与合约部署(学习资料).md)

**代码示例**:
- [Counter.sol](lessonA11.1/study_code/Counter.sol)
- [deploy.ts](lessonA11.1/study_code/deploy.ts)
- [ignition/modules/Counter.ts](lessonA11.1/study_code/ignition/modules/Counter.ts)
- [hardhat.config.ts](lessonA11.1/study_code/hardhat.config.ts)

---

#### 11.2 Hardhat 3 单元测试基础与技巧
**状态**: ✅  
**学时**: 5小时  
**教学内容**:
- Mocha + Chai框架
- 测试结构
- 断言语法
- 测试合约部署和函数
- 测试事件触发
- 测试错误和回退
- 时间旅行和区块操作
- 测试覆盖率
- Gas报告
- 快照和恢复

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_694a3519e4b0694ca15e6347)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - Hardhat3单元测试基础与技巧(课件资料).pdf](lessonA11.2/Solidity智能合约开发%20-%20Hardhat3单元测试基础与技巧(课件资料).pdf)
- 学习资料：[Solidity智能合约开发 - Hardhat 3 单元测试基础与技巧(学习资料).md](lessonA11.2/Solidity智能合约开发%20-%20Hardhat%203%20单元测试基础与技巧(学习资料).md)

**代码示例**:
- [Counter.test.ts](lessonA11.2/study_code/Counter.test.ts)
- [Token.test.ts](lessonA11.2/study_code/Token.test.ts)
- [TimeLock.test.ts](lessonA11.2/study_code/TimeLock.test.ts)
- [ErrorHandling.test.ts](lessonA11.2/study_code/ErrorHandling.test.ts)

---

#### 12.1 Foundry框架学习
**状态**: ✅  
**学时**: 4.5小时  
**教学内容**:
- Foundry安装和配置
- Forge、Cast、Anvil工具介绍
- 项目初始化
- foundry.toml配置
- 网络配置
- Forge编译合约
- 编写部署脚本
- Anvil本地测试网
- 测试网部署
- 合约验证和交互
- Cast命令行工具使用

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_694d09ede4b0694c5b6944cd)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - Foundry框架学习(课件).pdf](lessonA12.1/Solidity智能合约开发%20-%20Foundry框架学习(课件).pdf)
- 学习资料：[Solidity智能合约开发 - Foundry框架学习(学习资料).md](lessonA12.1/Solidity智能合约开发%20-%20Foundry框架学习(学习资料).md)

**代码示例**:
- [foundry-demo/](lessonA12.1/study_code/foundry-demo)

---

#### 13.1 Ethers.js基础与MetaMask交互
**状态**: ✅  
**学时**: 4.5小时  
**教学内容**:
- Ethers.js vs Web3.js
- Provider连接
- 读取区块链数据
- 读取合约状态
- 监听事件
- MetaMask Provider
- 请求账户连接
- 切换网络
- 发送交易
- 签名消息

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_694d0a16e4b0694c5b6944e4)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - Ethers.js基础与MetaMask交互.pptx](lessonA13.1/Solidity智能合约开发%20-%20Ethers.js基础与MetaMask交互.pptx)
- 学习资料：[Solidity智能合约开发 - Ethers.js基础与MetaMask交互(学习资料).md](lessonA13.1/Solidity智能合约开发%20-%20Ethers.js基础与MetaMask交互(学习资料).md)

**代码示例**:
- [ethers-demo/](lessonA13.1/study_code/ethers-demo/)

---

#### 14.1 众筹平台-实战开发
**状态**: ✅  
**学时**: 6小时  
**教学内容**:
- 需求分析
- 架构设计
- 核心合约开发
- 状态机实现
- 工厂合约
- 完整测试
- 前端开发
- 部署和验证

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_6954dc17e4b0694c5b6cea8a)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - 众筹平台实战.pdf](lessonA14.1/Solidity智能合约开发%20-%20众筹平台实战.pdf)
- 学习资料：[Solidity智能合约开发 - 众筹平台实战开发(学习资料).md](lessonA14.1/Solidity智能合约开发%20-%20众筹平台实战开发(学习资料).md)

**代码示例**:
- [crowdfunding-platform/](lessonA14.1/study_code/crowdfunding-platform/)

---

#### 14.2 实战项目：多签钱包
**状态**: ✅  
**学时**: 6-8小时  
**教学内容**:
- 多重签名机制原理
- 所有者管理（添加/删除）
- 交易提案创建与管理
- 确认机制与撤销确认
- 交易执行与重入防护
- 权限控制与安全设计
- 完整测试编写
- 前端集成与事件监听

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_696df70be4b0694c5b785730)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - 多签钱包项目实战.pptx](lessonA14.2/Solidity智能合约开发%20-%20多签钱包项目实战.pptx)
- 学习资料：[Solidity智能合约开发 - 多签钱包项目实战(学习资料).md](lessonA14.2/Solidity智能合约开发%20-%20多签钱包项目实战(学习资料).md)

**代码示例**:
- [multisig-wallet/](lessonA14.2/study_code/multisig-wallet/)

---

#### 14.3 实战项目：NFT盲盒
**状态**: ✅  
**学时**: 8-10小时  
**教学内容**:
- ERC721标准实现
- Chainlink VRF可验证随机函数集成
- 随机数生成与稀有度分配算法
- 盲盒购买与延迟揭示机制
- IPFS元数据存储与管理
- 分阶段销售与白名单机制
- 模块化架构设计（SaleManager、VRFHandler）
- UUPS代理模式合约升级
- 完整测试与部署流程

**课程地址**: [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_6970a84ce4b0694c5b79db09)  
**文档资料**: 
- 课件：[Solidity智能合约开发 - NFT盲盒项目实战.pptx](lessonA14.3/Solidity智能合约开发%20-%20NFT盲盒项目实战.pptx)
- 学习资料：[Solidity智能合约开发 - NFT盲盒项目实战(学习资料).md](lessonA14.3/Solidity智能合约开发%20-%20NFT盲盒项目实战(学习资料).md)

**代码示例**:
- [nft-blindbox/](lessonA14.3/study_code/nft-blindbox/)

---

## 课程进度统计

- ✅ : 25门课程（已完成）
- 📝 : 0门课程（录制中）
- **总课程数**: 25门
- **完成进度**: 100.0%

## 状态说明

- ✅ : 课程已完成录制和剪辑，可以观看
- 📝 : 课程正在录制中

## 项目结构

```
solidity_lesson/
├── lesson1.1/          # 课程导言与环境搭建
├── lesson1.2/          # EVM存储结构
├── lesson2.1/          # 数据类型基础
├── lesson3.1/          # 数组
├── lesson3.2/          # 映射和结构体
├── lesson3.3/          # 函数与修饰符
├── lesson4.1/          # 控制流语句
├── lesson4.2/          # 特殊类型与全局变量
├── lesson5.1/          # 基础项目：简单代币合约
├── lesson6.1/          # 合约继承
├── lesson6.2/          # 库合约Library
├── lesson7.1/          # 事件Events
├── lesson7.2/          # 错误处理和自定义错误
├── lesson8.1/          # 合约间调用
├── lesson8.2/          # 智能合约安全基础
├── lesson9.1/          # Gas优化技巧
├── lesson9.2/          # 智能合约设计模式
├── lessonA10.1/        # 进阶项目：NFT市场
├── lessonA11.1/        # Hardhat 3 环境搭建与合约部署
├── lessonA11.2/        # Hardhat 3 单元测试基础与技巧
├── lessonA12.1/        # Foundry框架学习
├── lessonA13.1/        # Ethers.js基础与MetaMask交互
├── lessonA14.1/        # 众筹平台-实战开发
├── lessonA14.2/        # 多签钱包项目实战
├── lessonA14.3/        # NFT盲盒项目实战
└── README.md           # 课程说明文档
```

## 学习建议

1. **循序渐进**: 按照课程编号顺序学习，每个阶段都是下一阶段的基础
2. **动手实践**: 每节课后在Remix IDE中实践代码示例
3. **完成作业**: 认真完成课后作业，巩固所学知识
4. **项目实战**: 重点关注各阶段的实战项目，提升实际开发能力

## 联系方式

如有问题或建议，欢迎提Issue或PR。 
