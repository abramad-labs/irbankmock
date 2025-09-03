export const fetcher = (...args: Parameters<typeof fetch>) => fetch(...args).then(res => res.json())


export class ResponseError<T = any> extends Error {
    public response: Response
    public status: number
    public data?: T
    public constructor(message: string, response: Response) {
        super(message)
        this.response = response
        this.status = response.status
    }
}

export const fetcherWithError = (...args: Parameters<typeof fetch>) =>
    fetch(...args)
    .then(async res => {
        if(!res.ok) {
            const err = new ResponseError("error occured in fetch", res)
            err.data = await err.response.json()
            throw err
        }
        return res.json()
    })