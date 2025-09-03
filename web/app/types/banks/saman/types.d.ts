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
    terminalName: string,
    terminalId: string,
    website: string,
    amount: number
    expiresAt: string
}

export type SuccessErrorPair = {
    success: boolean
    error: string
}