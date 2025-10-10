package solana
type Block struct{ Slot uint64; BlockTime *int64; BlockHash string; ParentSlot uint64; Transactions []Tx }
type Tx struct{ Signature string }
