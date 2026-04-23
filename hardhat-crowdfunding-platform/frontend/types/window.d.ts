interface Window {
  ethereum?: {
    isMetaMask?: boolean;
    /** 多钱包并存时（如 MetaMask + OKX）常见 */
    providers?: unknown[];
    request: (args: { method: string; params?: any[] }) => Promise<any>;
    on: (event: string, handler: (...args: any[]) => void) => void;
    removeListener: (event: string, handler: (...args: any[]) => void) => void;
  };
}

