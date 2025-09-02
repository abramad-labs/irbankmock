export type SamanTerminal = {
    id: number,
    name: string,
    username: string,
    password: string
}

export type SamanTerminalsResponse = {
    terminals: SamanTerminal[]
}

export type CreateTerminalPayload = {
    name: string
}