/**
 * 兼容某些 TypeScript 版本对 package exports + typesVersions 的解析差异，
 * 避免 `@nomicfoundation/hardhat-ignition/modules` 在编辑器中误报红线。
 */
declare module "@nomicfoundation/hardhat-ignition/modules" {
  export { buildModule } from "@nomicfoundation/ignition-core";
}
