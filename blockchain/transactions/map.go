package transactions

type Map struct {
	txs map[string]*Transaction
}

func NewMap() *Map {
	return &Map{txs: make(map[string]*Transaction)}
}

func (m *Map) Set(txID string, tx *Transaction) {
	m.txs[txID] = tx
}

func (m *Map) Get(txID string) *Transaction {
	return m.txs[txID]
}

func (m *Map) GetAll() map[string]*Transaction {
	return m.txs
}
