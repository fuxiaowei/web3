'use client';

import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { ethers, type Eip1193Provider } from 'ethers';

/** 浏览器扩展注入的 MetaMask inpage provider（ethers 的 Eip1193Provider 不含事件 API） */
type MetaMaskInpageProvider = Eip1193Provider & {
  on?: (event: string, handler: (...args: any[]) => void) => void;
  removeListener?: (event: string, handler: (...args: any[]) => void) => void;
};

function isLikelyMetaMask(p: unknown): p is MetaMaskInpageProvider {
  if (!p || typeof p !== 'object') return false;
  const x = p as Record<string, unknown>;
  if (x.isMetaMask !== true) return false;
  if (x.isOkxWallet === true || x.isOKExWallet === true) return false;
  return typeof (p as Eip1193Provider).request === 'function';
}

/** 在 OKX 等多钱包同时注入时，避免误用顶层的 window.ethereum（常为 OKX） */
function getMetaMaskEip1193Provider(): MetaMaskInpageProvider | null {
  if (typeof window === 'undefined') return null;
  const eth = window.ethereum;
  if (!eth) return null;

  const multi = eth.providers;
  if (Array.isArray(multi) && multi.length > 0) {
    const mm = multi.find(isLikelyMetaMask);
    if (mm) return mm;
  }

  if (isLikelyMetaMask(eth)) return eth as MetaMaskInpageProvider;
  return null;
}

const CONFIG_CHAIN_ID_RAW = process.env.NEXT_PUBLIC_CHAIN_ID?.trim();
const CONFIG_RPC_URL = process.env.NEXT_PUBLIC_RPC_URL?.trim();
const CONFIG_CHAIN_NAME = process.env.NEXT_PUBLIC_CHAIN_NAME?.trim() || 'Local Network';
const CONFIG_EXPLORER_URL = process.env.NEXT_PUBLIC_BLOCK_EXPLORER_URL?.trim();

function parseExpectedChainId(): bigint | null {
  if (!CONFIG_CHAIN_ID_RAW) {
    console.warn('NEXT_PUBLIC_CHAIN_ID 未设置，将不会自动切换网络');
    return null;
  }
  try {
    return BigInt(CONFIG_CHAIN_ID_RAW);
  } catch {
    console.warn('无效的 NEXT_PUBLIC_CHAIN_ID:', CONFIG_CHAIN_ID_RAW);
    return null;
  }
}

function getExpectedChainIdHex(): string | null {
  const id = parseExpectedChainId();
  if (id === null) return null;
  return `0x${id.toString(16)}`;
}

interface Web3ContextType {
  provider: ethers.BrowserProvider | null;
  signer: ethers.JsonRpcSigner | null;
  account: string | null;
  connectWallet: () => Promise<void>;
  disconnect: () => void;
  isConnecting: boolean;
}

const Web3Context = createContext<Web3ContextType | undefined>(undefined);

export function Web3Provider({ children }: { children: ReactNode }) {
  const [provider, setProvider] = useState<ethers.BrowserProvider | null>(null);
  const [signer, setSigner] = useState<ethers.JsonRpcSigner | null>(null);
  const [account, setAccount] = useState<string | null>(null);
  const [isConnecting, setIsConnecting] = useState(false);

  const switchToConfiguredChain = async () => {
    const ethereum = getMetaMaskEip1193Provider();
    if (!ethereum) return;
    const chainIdHex = getExpectedChainIdHex();
    if (!chainIdHex || !CONFIG_RPC_URL) {
      throw new Error('请在 .env.local 中配置 NEXT_PUBLIC_CHAIN_ID 与 NEXT_PUBLIC_RPC_URL');
    }
    try {
      await ethereum.request({
        method: 'wallet_switchEthereumChain',
        params: [{ chainId: chainIdHex }],
      });
    } catch (switchError: any) {
      if (switchError.code === 4902) {
        const chainParams: Record<string, unknown> = {
          chainId: chainIdHex,
          chainName: CONFIG_CHAIN_NAME,
          nativeCurrency: {
            name: 'ETH',
            symbol: 'ETH',
            decimals: 18,
          },
          rpcUrls: [CONFIG_RPC_URL],
        };
        if (CONFIG_EXPLORER_URL) {
          chainParams.blockExplorerUrls = [CONFIG_EXPLORER_URL];
        }
        try {
          await ethereum.request({
            method: 'wallet_addEthereumChain',
            params: [chainParams],
          });
        } catch (addError) {
          console.error('添加自定义网络失败:', addError);
          throw addError;
        }
      } else {
        throw switchError;
      }
    }
  };

  useEffect(() => {
    const ethereum = getMetaMaskEip1193Provider();
    if (typeof window !== 'undefined' && ethereum) {
      const initProvider = new ethers.BrowserProvider(ethereum);
      
      // Check network and switch if needed（与 .env 中 NEXT_PUBLIC_CHAIN_ID 一致）
      initProvider.getNetwork().then(async (network) => {
        const expectedChainId = parseExpectedChainId();
        if (expectedChainId !== null && network.chainId !== expectedChainId) {
          console.log(
            `⚠️ 当前网络: ${network.chainId} (${network.name})，期望: ${expectedChainId} (${CONFIG_CHAIN_NAME})`
          );
          console.log('🔄 尝试切换到配置的网络...');
          try {
            await switchToConfiguredChain();
            // After switching, update provider without reloading
            console.log('✅ Network switched, updating provider...');
            const mm = getMetaMaskEip1193Provider();
            if (!mm) {
              console.error('❌ MetaMask provider is not available');
              return;
            }
            const newProvider = new ethers.BrowserProvider(mm);
            setProvider(newProvider);
            // Re-check accounts after network switch
            const accounts = await mm.request({ method: 'eth_accounts' });
            if (accounts.length > 0) {
              const newSigner = await newProvider.getSigner();
              setSigner(newSigner);
              setAccount(accounts[0]);
            }
            return;
          } catch (error) {
            console.error('❌ Failed to switch network:', error);
            alert(
              `请在 MetaMask 中手动切换到链 ID ${CONFIG_CHAIN_ID_RAW ?? '(未配置)'}（${CONFIG_CHAIN_NAME}），RPC: ${CONFIG_RPC_URL ?? '(未配置)'}`
            );
          }
        } else if (expectedChainId !== null) {
          console.log(`✅ 已连接到配置网络（chainId=${expectedChainId}）`);
        }
        
        setProvider(initProvider);

        // Check if already connected
        ethereum
          .request({ method: 'eth_accounts' })
          .then((accounts: string[]) => {
            if (accounts.length > 0) {
              handleAccountsChanged(accounts);
            }
          });
      });

      // Listen for account changes
      if (ethereum.on) {
        ethereum.on('accountsChanged', handleAccountsChanged);
        ethereum.on('chainChanged', async () => {
          // Update provider when chain changes, without reloading page
          console.log('🔄 Chain changed, updating provider...');
          const mm = getMetaMaskEip1193Provider();
          if (!mm) {
            console.error('❌ MetaMask provider is not available');
            return;
          }
          const newProvider = new ethers.BrowserProvider(mm);
          const network = await newProvider.getNetwork();
          console.log('✅ New network:', { chainId: Number(network.chainId), name: network.name });
          setProvider(newProvider);
          // Re-check accounts after chain change
          const accounts = await mm.request({ method: 'eth_accounts' });
          if (accounts.length > 0) {
            const newSigner = await newProvider.getSigner();
            setSigner(newSigner);
            setAccount(accounts[0]);
          } else {
            setAccount(null);
            setSigner(null);
          }
        });
      }
    }

    return () => {
      const mm = getMetaMaskEip1193Provider();
      if (mm && mm.removeListener) {
        mm.removeListener('accountsChanged', handleAccountsChanged);
      }
    };
  }, []);

  const handleAccountsChanged = async (accounts: string[]) => {
    if (accounts.length === 0) {
      setAccount(null);
      setSigner(null);
    } else {
      const newAccount = accounts[0];
      setAccount(newAccount);
      
      // 确保 provider 存在后再获取 signer
      const mm = getMetaMaskEip1193Provider();
      if (mm) {
        try {
          const currentProvider = new ethers.BrowserProvider(mm);
          const newSigner = await currentProvider.getSigner();
          setSigner(newSigner);
          // 如果 provider 状态还没有设置，也设置它
          setProvider(currentProvider);
          console.log('✅ Account changed, signer updated:', newAccount);
        } catch (error) {
          console.error('Failed to get signer after account change:', error);
        }
      }
    }
  };


  const connectWallet = async () => {
    const ethereum = getMetaMaskEip1193Provider();
    if (typeof window === 'undefined' || !ethereum) {
      alert('未检测到 MetaMask，请安装 MetaMask；若已安装，请确认未仅启用 OKX 等其它钱包为默认注入。');
      return;
    }

    try {
      setIsConnecting(true);
      const newProvider = new ethers.BrowserProvider(ethereum);
      
      const network = await newProvider.getNetwork();
      const expectedChainId = parseExpectedChainId();

      if (expectedChainId !== null && network.chainId !== expectedChainId) {
        console.log(`Current network: ${network.chainId}, expected: ${expectedChainId}`);
        await switchToConfiguredChain();
        // Update provider after network switch, without reloading
        const mmAfterSwitch = getMetaMaskEip1193Provider();
        if (!mmAfterSwitch) {
          console.error('❌ MetaMask provider is not available');
          return;
        }
        const updatedProvider = new ethers.BrowserProvider(mmAfterSwitch);
        const accounts = await updatedProvider.send('eth_requestAccounts', []);
        const newSigner = await updatedProvider.getSigner();
        setProvider(updatedProvider);
        setSigner(newSigner);
        setAccount(accounts[0]);
        return;
      }

      const accounts = await newProvider.send('eth_requestAccounts', []);
      const newSigner = await newProvider.getSigner();

      setProvider(newProvider);
      setSigner(newSigner);
      setAccount(accounts[0]);
    } catch (error: any) {
      console.error('Error connecting wallet:', error);
      alert(error.message || 'Failed to connect wallet');
    } finally {
      setIsConnecting(false);
    }
  };

  const disconnect = () => {
    setAccount(null);
    setSigner(null);
  };

  return (
    <Web3Context.Provider
      value={{
        provider,
        signer,
        account,
        connectWallet,
        disconnect,
        isConnecting,
      }}
    >
      {children}
    </Web3Context.Provider>
  );
}

export function useWeb3() {
  const context = useContext(Web3Context);
  if (context === undefined) {
    throw new Error('useWeb3 must be used within a Web3Provider');
  }
  return context;
}

