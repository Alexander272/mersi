package models

import "mime/multipart"

type ImportDTO struct {
	RealmId   string
	SectionId string
	BidType   string
	UserId    string
	File      *multipart.FileHeader
}

type Template struct {
	Name                      int
	DateOfReceipt             int
	Type                      int
	FactoryNumber             int
	MeasurementLimits         int
	Accuracy                  int
	StateRegister             int
	CountryOfProduce          int
	Manufacturer              int
	Responsible               int
	Inventory                 int
	YearOfIssue               int
	InterVerificationInterval int
	Notes                     int
	VerificationDate          int
	NextVerificationDate      int
	Repair                    int
	Preservation              int
	Transfer                  int
	TransferToDepartment      int
	WriteOff                  int
	Person                    int
	Place                     int
}
