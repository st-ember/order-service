package enum

type PurchaseStatus int

const (
	Pending PurchaseStatus = iota
	Processing
	Completed
	Canceled
)

func (ps PurchaseStatus) String() string {
	switch ps {
	case Pending:
		return "pending"
	case Processing:
		return "processing"
	case Completed:
		return "completed"
	case Canceled:
		return "canceled"
	default:
		return "unknown"
	}
}
