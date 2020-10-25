
export default interface Receipt {
    id: string
    vendor_name: string
    receipt_number: string
    tax_id: string
    date: string
    total: number
    vat: number
    verified?: Date
    added: Date
}

// checks if any data is missing and renders exclamation point
export function isWarning(r: Receipt): boolean {
    return r.total === 0 
        || r.vat === 0
        || !r.vendor_name
        || !r.tax_id
        || !r.receipt_number
        || !r.date
        || !r.total
        || !r.vat ?
        true
        : false
}

// checks if the receipt has been verified
export function isVerified(r: Receipt): boolean {
    return !!r.verified
}