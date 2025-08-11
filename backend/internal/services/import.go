package services

type ImportService struct {
	instrument   Instrument
	verification Verification
}

type ImportDeps struct {
	Instrument   Instrument
	Verification Verification
}

func NewImportService(deps *ImportDeps) *ImportService {
	return &ImportService{
		instrument:   deps.Instrument,
		verification: deps.Verification,
	}
}

type Import interface{}
