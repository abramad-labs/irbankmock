export type SamanTerminal = {
    id: number;
    name: string;
    username: string;
    password: string;
};

export type SamanTerminalsResponse = {
    terminals: SamanTerminal[];
    endpoints: {
        paymentGateway: string;
        paymentToken: string;
        receipt: string;
        verifyTransaction: string;
        reverseTransaction: string;
    };
};

export type CreateTerminalPayload = {
    name: string;
};

export type SamanPublicTokenInfoResponse = {
    terminalName: string;
    terminalId: string;
    website: string;
    amount: number;
    expiresAt: string;
};

export type SuccessErrorPair = {
    success: boolean;
    error: string;
};

export type BankSepCancelOrFailTokenRequest = {
    token: string;
};

type BankSepSubmitTokenRequest = {
    token: string;
    cardNumber: string;
    cvv: number;
    expiryMonth: number;
    expiryYear: number;
    cardPassword: string;
    captcha: string;
};

type BankSepTokenFinalizeResponseCallbackData = {
    MID: string;
    terminalId: string;
    state: string;
    status: string;
    rrn: string;
    refNum: string;
    resNum: string;
    traceNo: string;
    amount: string;
    affectiveAmount: string;
    wage: string;
    securePan: string;
    hashedCardNumber: string;
    token: string;
};

type BankSepTokenFinalizeResponse = {
    redirectURL: string;
    callbackData: BankSepTokenFinalizeResponseCallbackData;
};
