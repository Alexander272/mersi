package constants

const (
	InstrumentStatusWork   string = "work"
	InstrumentStatusRepair string = "repair"
	InstrumentStatusDec    string = "decommissioning"
	InstrumentDeleted      string = "deleted"
	InstrumentDraft        string = "draft"
)

type InstrumentStatus int64

const (
	Work InstrumentStatus = iota
	Repair
	Dec
)
