import { SamanTerminal, type CreateTerminalPayload } from "@/types/banks/saman/types";
import axios from "axios";

export const createTerminal = (payload: CreateTerminalPayload) => {
    return axios.post<SamanTerminal>('/banks/saman/management/terminal', payload)
}