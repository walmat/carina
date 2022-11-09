interface Window {
  RPCAction: (action: string, input?: any[]) => Promise<any>;
}

interface RcpResponse {
  success: boolean;
  retData: any;
}

interface RcpRequest {
  type: string;
  data?: string;
}

type Callback = (...args: any[]) => void

declare var astilectron: {
  sendMessage: ({ type, data }: RcpRequest, callback?: Callback) => null | RcpResponse;
  onMessage: (callback: Callback) => any;
};