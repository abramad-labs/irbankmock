import { BankSepCancelOrFailTokenRequest, BankSepSubmitTokenRequest, BankSepTokenFinalizeResponse, SamanTerminal, type CreateTerminalPayload } from "@/types/banks/saman/types";
import axios from "axios";

export const createTerminal = (payload: CreateTerminalPayload) => {
    return axios.post<SamanTerminal>('/banks/saman/management/terminal', payload)
}

export const submitToken = (payload: BankSepSubmitTokenRequest) => {
    return axios.post<BankSepTokenFinalizeResponse>('/banks/saman/management/token/submit', payload)
}

export const failToken = (payload: BankSepCancelOrFailTokenRequest) => {
    return axios.post<BankSepTokenFinalizeResponse>('/banks/saman/management/token/fail', payload)
}

export const cancelToken = (payload: BankSepCancelOrFailTokenRequest) => {
    return axios.post<BankSepTokenFinalizeResponse>('/banks/saman/management/token/cancel', payload)
}