package services

import (
	"context"
	"time"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
	"github.com/Alexander272/mersi/backend/internal/services/most"
)

type fakeTx struct {
	postgres.Tx
	committed  bool
	rolledBack bool
}

func (f *fakeTx) Commit(context.Context) error   { f.committed = true; return nil }
func (f *fakeTx) Rollback(context.Context) error { f.rolledBack = true; return nil }

type fakeTxManager struct {
	TransactionManager
	execute func(ctx context.Context, fn func(tx postgres.Tx) error) error
}

func (f *fakeTxManager) ExecuteInTx(ctx context.Context, fn func(tx postgres.Tx) error) error {
	if f.execute != nil {
		return f.execute(ctx, fn)
	}
	return fn(&fakeTx{})
}

type fakeRepoTx struct {
	repository.Transaction
	tx       postgres.Tx
	beginErr error
}

func (f *fakeRepoTx) BeginTx(ctx context.Context) (postgres.Tx, error) {
	if f.beginErr != nil {
		return nil, f.beginErr
	}
	return f.tx, nil
}

type fakeResponsibleSvc struct {
	Responsible
	getBySSOIdFn     func(ctx context.Context, id string) ([]*models.Responsible, error)
	getWithChannelFn func(ctx context.Context, req *models.GetResponsibleDTO) ([]*models.ResponsibleWithChannel, error)
}

func (f *fakeResponsibleSvc) GetBySSOId(ctx context.Context, id string) ([]*models.Responsible, error) {
	if f.getBySSOIdFn != nil {
		return f.getBySSOIdFn(ctx, id)
	}
	return nil, nil
}

func (f *fakeResponsibleSvc) GetWithChannel(ctx context.Context, req *models.GetResponsibleDTO) ([]*models.ResponsibleWithChannel, error) {
	if f.getWithChannelFn != nil {
		return f.getWithChannelFn(ctx, req)
	}
	return nil, nil
}

type fakeLocationSvc struct {
	Location
	selectByDepartmentsFn func(ctx context.Context, dto *models.SelectByDepsDTO) ([]string, error)
	getByIdFn            func(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error)
	getLastFn            func(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error)
	getUsedByHolderFn    func(ctx context.Context, dto *models.GetLocationByHolderDTO) ([]*models.Location, error)
	getUsedByDeptFn      func(ctx context.Context, dto *models.GetLocationByDepartmentDTO) ([]*models.Location, error)
	setPersonFn          func(ctx context.Context, personId string) error
	setDepartmentFn      func(ctx context.Context, departmentId string) error
	createFn             func(ctx context.Context, tx postgres.Tx, dto *models.LocationDTO) error

	locationsCreated    []*models.LocationDTO
	setPersonCalls      []string
	setDepartmentCalls  []string
}

func (f *fakeLocationSvc) SelectByDepartments(ctx context.Context, dto *models.SelectByDepsDTO) ([]string, error) {
	if f.selectByDepartmentsFn != nil {
		return f.selectByDepartmentsFn(ctx, dto)
	}
	return nil, nil
}

func (f *fakeLocationSvc) GetById(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error) {
	if f.getByIdFn != nil {
		return f.getByIdFn(ctx, dto)
	}
	return nil, models.ErrNoRows
}

func (f *fakeLocationSvc) GetLast(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error) {
	if f.getLastFn != nil {
		return f.getLastFn(ctx, dto)
	}
	return nil, models.ErrNoRows
}

func (f *fakeLocationSvc) GetUsedByHolder(ctx context.Context, dto *models.GetLocationByHolderDTO) ([]*models.Location, error) {
	if f.getUsedByHolderFn != nil {
		return f.getUsedByHolderFn(ctx, dto)
	}
	return nil, nil
}

func (f *fakeLocationSvc) GetUsedByDepartment(ctx context.Context, dto *models.GetLocationByDepartmentDTO) ([]*models.Location, error) {
	if f.getUsedByDeptFn != nil {
		return f.getUsedByDeptFn(ctx, dto)
	}
	return nil, nil
}

func (f *fakeLocationSvc) SetPerson(ctx context.Context, personId string) error {
	f.setPersonCalls = append(f.setPersonCalls, personId)
	if f.setPersonFn != nil {
		return f.setPersonFn(ctx, personId)
	}
	return nil
}

func (f *fakeLocationSvc) SetDepartment(ctx context.Context, departmentId string) error {
	f.setDepartmentCalls = append(f.setDepartmentCalls, departmentId)
	if f.setDepartmentFn != nil {
		return f.setDepartmentFn(ctx, departmentId)
	}
	return nil
}

func (f *fakeLocationSvc) Create(ctx context.Context, tx postgres.Tx, dto *models.LocationDTO) error {
	f.locationsCreated = append(f.locationsCreated, dto)
	if f.createFn != nil {
		return f.createFn(ctx, tx, dto)
	}
	return nil
}

type fakeLocationRepo struct {
	repository.Location
	receivingFn        func(ctx context.Context, dto *models.ReceivingDTO) error
	forcedReceiptFn    func(ctx context.Context, dto *models.ForcedReceiptDTO) error
	forcedReceiptAllFn func(ctx context.Context) error
	createInTxFn       func(ctx context.Context, tx postgres.Tx, dto *models.LocationDTO) error
	createSeveralFn    func(ctx context.Context, dto []*models.LocationDTO) error
	selectByDeptFn     func(ctx context.Context, dto *models.SelectByDepsDTO) ([]string, error)
	receivingCalled    []*models.ReceivingDTO
	forcedCalled       []*models.ForcedReceiptDTO
	createSeveralCalled [][]*models.LocationDTO
}

func (f *fakeLocationRepo) Receiving(ctx context.Context, dto *models.ReceivingDTO) error {
	f.receivingCalled = append(f.receivingCalled, dto)
	if f.receivingFn != nil {
		return f.receivingFn(ctx, dto)
	}
	return nil
}

func (f *fakeLocationRepo) ForcedReceipt(ctx context.Context, dto *models.ForcedReceiptDTO) error {
	f.forcedCalled = append(f.forcedCalled, dto)
	if f.forcedReceiptFn != nil {
		return f.forcedReceiptFn(ctx, dto)
	}
	return nil
}

func (f *fakeLocationRepo) ForcedReceiptAll(ctx context.Context) error {
	if f.forcedReceiptAllFn != nil {
		return f.forcedReceiptAllFn(ctx)
	}
	return nil
}

func (f *fakeLocationRepo) CreateInTx(ctx context.Context, tx postgres.Tx, dto *models.LocationDTO) error {
	if f.createInTxFn != nil {
		return f.createInTxFn(ctx, tx, dto)
	}
	return nil
}

func (f *fakeLocationRepo) CreateSeveral(ctx context.Context, dto []*models.LocationDTO) error {
	f.createSeveralCalled = append(f.createSeveralCalled, dto)
	if f.createSeveralFn != nil {
		return f.createSeveralFn(ctx, dto)
	}
	return nil
}

func (f *fakeLocationRepo) SelectByDepartment(ctx context.Context, dto *models.SelectByDepsDTO) ([]string, error) {
	if f.selectByDeptFn != nil {
		return f.selectByDeptFn(ctx, dto)
	}
	return nil, nil
}

type fakeNotificationSvc struct {
	Notification
	checkUsedFn         func(ctx context.Context) error
	checkSentFn         func(ctx context.Context) error
	checkVerificationFn func(ctx context.Context) error
	checkReceivingFn    func(ctx context.Context, dto *models.DialogResponse) error
	usedCalls           int
	sentCalls           int
	verificationCalls   int
}

func (f *fakeNotificationSvc) CheckUsed(ctx context.Context) error {
	f.usedCalls++
	if f.checkUsedFn != nil {
		return f.checkUsedFn(ctx)
	}
	return nil
}

func (f *fakeNotificationSvc) CheckSent(ctx context.Context) error {
	f.sentCalls++
	if f.checkSentFn != nil {
		return f.checkSentFn(ctx)
	}
	return nil
}

func (f *fakeNotificationSvc) CheckVerification(ctx context.Context) error {
	f.verificationCalls++
	if f.checkVerificationFn != nil {
		return f.checkVerificationFn(ctx)
	}
	return nil
}

func (f *fakeNotificationSvc) CheckReceiving(ctx context.Context, dto *models.DialogResponse) error {
	if f.checkReceivingFn != nil {
		return f.checkReceivingFn(ctx, dto)
	}
	return nil
}

type fakeActivityLogSvc struct {
	ActivityLog
	logs []*models.CreateActivityLogDTO
}

func (f *fakeActivityLogSvc) LogActivity(ctx context.Context, dto *models.CreateActivityLogDTO) {
	f.logs = append(f.logs, dto)
}

func (f *fakeActivityLogSvc) Create(ctx context.Context, dto *models.CreateActivityLogDTO) {
	f.LogActivity(ctx, dto)
}

type fakeDialog struct {
	most.Dialog
	openFn func(ctx context.Context, action *models.PostAction) error
	action *models.PostAction
}

func (f *fakeDialog) Open(ctx context.Context, action *models.PostAction) error {
	f.action = action
	if f.openFn != nil {
		return f.openFn(ctx, action)
	}
	return nil
}

type fakePost struct {
	most.Post
	createFn func(ctx context.Context, dto *models.CreatePostDTO) error
	updateFn func(ctx context.Context, dto *models.UpdatePostDTO) error
	created  []*models.CreatePostDTO
	updated  []*models.UpdatePostDTO
}

func (f *fakePost) Create(ctx context.Context, dto *models.CreatePostDTO) error {
	f.created = append(f.created, dto)
	if f.createFn != nil {
		return f.createFn(ctx, dto)
	}
	return nil
}

func (f *fakePost) Update(ctx context.Context, dto *models.UpdatePostDTO) error {
	f.updated = append(f.updated, dto)
	if f.updateFn != nil {
		return f.updateFn(ctx, dto)
	}
	return nil
}

type fakeVerificationRepo struct {
	repository.Verification
	getFn                    func(ctx context.Context, req *models.GetVerificationDTO) ([]*models.Verification, error)
	getByIdFn                func(ctx context.Context, id string) (*models.Verification, error)
	getByInstrumentAndDateFn func(ctx context.Context, tx postgres.Tx, instrumentId string, date time.Time) (*models.Verification, error)
	getLastFn                func(ctx context.Context, req *models.GetVerificationDTO) (*models.Verification, error)
	createInTxFn             func(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error
	createSeveralFn          func(ctx context.Context, tx postgres.Tx, dto []*models.VerificationDTO) error
	updateFn                 func(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error
	deleteFn                 func(ctx context.Context, tx postgres.Tx, dto *models.DeleteVerificationDTO) error
}

func (f *fakeVerificationRepo) Get(ctx context.Context, req *models.GetVerificationDTO) ([]*models.Verification, error) {
	if f.getFn != nil {
		return f.getFn(ctx, req)
	}
	return nil, nil
}

func (f *fakeVerificationRepo) GetById(ctx context.Context, id string) (*models.Verification, error) {
	if f.getByIdFn != nil {
		return f.getByIdFn(ctx, id)
	}
	return nil, models.ErrNoRows
}

func (f *fakeVerificationRepo) GetByInstrumentAndDate(ctx context.Context, tx postgres.Tx, instrumentId string, date time.Time) (*models.Verification, error) {
	if f.getByInstrumentAndDateFn != nil {
		return f.getByInstrumentAndDateFn(ctx, tx, instrumentId, date)
	}
	return nil, models.ErrNoRows
}

func (f *fakeVerificationRepo) GetLast(ctx context.Context, req *models.GetVerificationDTO) (*models.Verification, error) {
	if f.getLastFn != nil {
		return f.getLastFn(ctx, req)
	}
	return nil, models.ErrNoRows
}

func (f *fakeVerificationRepo) CreateInTx(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error {
	if f.createInTxFn != nil {
		return f.createInTxFn(ctx, tx, dto)
	}
	return nil
}

func (f *fakeVerificationRepo) CreateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.VerificationDTO) error {
	if f.createSeveralFn != nil {
		return f.createSeveralFn(ctx, tx, dto)
	}
	return nil
}

func (f *fakeVerificationRepo) Update(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, tx, dto)
	}
	return nil
}

func (f *fakeVerificationRepo) Delete(ctx context.Context, tx postgres.Tx, dto *models.DeleteVerificationDTO) error {
	if f.deleteFn != nil {
		return f.deleteFn(ctx, tx, dto)
	}
	return nil
}

type fakeVerDocSvc struct {
	VerificationDoc
	getGroupedFn    func(ctx context.Context, req *models.GetGroupedVerificationDocsDTO) (*models.GroupedVerificationDocs, error)
	createSeveralFn func(ctx context.Context, tx postgres.Tx, dto []*models.VerificationDocDTO) error
	updateSeveralFn func(ctx context.Context, tx postgres.Tx, dto []*models.VerificationDocDTO) error
	deleteByDocIdFn func(ctx context.Context, tx postgres.Tx, docId string) error
	created         []*models.VerificationDocDTO
	updated         []*models.VerificationDocDTO
	deleted         []string
}

func (f *fakeVerDocSvc) GetGrouped(ctx context.Context, req *models.GetGroupedVerificationDocsDTO) (*models.GroupedVerificationDocs, error) {
	if f.getGroupedFn != nil {
		return f.getGroupedFn(ctx, req)
	}
	return &models.GroupedVerificationDocs{}, nil
}

func (f *fakeVerDocSvc) CreateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.VerificationDocDTO) error {
	f.created = append(f.created, dto...)
	if f.createSeveralFn != nil {
		return f.createSeveralFn(ctx, tx, dto)
	}
	return nil
}

func (f *fakeVerDocSvc) UpdateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.VerificationDocDTO) error {
	f.updated = append(f.updated, dto...)
	if f.updateSeveralFn != nil {
		return f.updateSeveralFn(ctx, tx, dto)
	}
	return nil
}

func (f *fakeVerDocSvc) DeleteByDocId(ctx context.Context, tx postgres.Tx, docId string) error {
	f.deleted = append(f.deleted, docId)
	if f.deleteByDocIdFn != nil {
		return f.deleteByDocIdFn(ctx, tx, docId)
	}
	return nil
}

type fakeInstrumentSvc struct {
	Instrument
	changeStatusFn func(ctx context.Context, tx postgres.Tx, dto *models.UpdateStatus) error
	statusChanges  []*models.UpdateStatus

	getByIdFn         func(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error)
	createFn          func(ctx context.Context, tx postgres.Tx, dto *models.InstrumentDTO) error
	createSeveralFn   func(ctx context.Context, dto []*models.InstrumentDTO) error
	updateFn          func(ctx context.Context, tx postgres.Tx, dto *models.InstrumentDTO) error
	deleteFn          func(ctx context.Context, dto *models.DeleteSiDTO) error
	created           []*models.InstrumentDTO
	createdSeveral    [][]*models.InstrumentDTO
	updated           []*models.InstrumentDTO
	deleted           []*models.DeleteSiDTO
}

func (f *fakeInstrumentSvc) ChangeStatus(ctx context.Context, tx postgres.Tx, dto *models.UpdateStatus) error {
	f.statusChanges = append(f.statusChanges, dto)
	if f.changeStatusFn != nil {
		return f.changeStatusFn(ctx, tx, dto)
	}
	return nil
}

func (f *fakeInstrumentSvc) GetById(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error) {
	if f.getByIdFn != nil {
		return f.getByIdFn(ctx, req)
	}
	return nil, models.ErrNoRows
}

func (f *fakeInstrumentSvc) Create(ctx context.Context, tx postgres.Tx, dto *models.InstrumentDTO) error {
	f.created = append(f.created, dto)
	if f.createFn != nil {
		return f.createFn(ctx, tx, dto)
	}
	return nil
}

func (f *fakeInstrumentSvc) CreateSeveral(ctx context.Context, dto []*models.InstrumentDTO) error {
	f.createdSeveral = append(f.createdSeveral, dto)
	if f.createSeveralFn != nil {
		return f.createSeveralFn(ctx, dto)
	}
	return nil
}

func (f *fakeInstrumentSvc) Update(ctx context.Context, tx postgres.Tx, dto *models.InstrumentDTO) error {
	f.updated = append(f.updated, dto)
	if f.updateFn != nil {
		return f.updateFn(ctx, tx, dto)
	}
	return nil
}

func (f *fakeInstrumentSvc) Delete(ctx context.Context, dto *models.DeleteSiDTO) error {
	f.deleted = append(f.deleted, dto)
	if f.deleteFn != nil {
		return f.deleteFn(ctx, dto)
	}
	return nil
}

type fakeDocumentSvc struct {
	Document
	changePathFn         func(ctx context.Context, dto *models.PathParts) error
	deleteFn             func(ctx context.Context, dto *models.DeleteDocumentDTO) error
	removeEmptyFoldersFn func(ctx context.Context) error
	pathChanges          []*models.PathParts
	deletedDocs          []*models.DeleteDocumentDTO
	folderCalls          int
}

func (f *fakeDocumentSvc) ChangePath(ctx context.Context, dto *models.PathParts) error {
	f.pathChanges = append(f.pathChanges, dto)
	if f.changePathFn != nil {
		return f.changePathFn(ctx, dto)
	}
	return nil
}

func (f *fakeDocumentSvc) Delete(ctx context.Context, dto *models.DeleteDocumentDTO) error {
	f.deletedDocs = append(f.deletedDocs, dto)
	if f.deleteFn != nil {
		return f.deleteFn(ctx, dto)
	}
	return nil
}

func (f *fakeDocumentSvc) RemoveEmptyFolders(ctx context.Context) error {
	f.folderCalls++
	if f.removeEmptyFoldersFn != nil {
		return f.removeEmptyFoldersFn(ctx)
	}
	return nil
}

type fakeRuleSvc struct {
	Rule
	getAllFn func(ctx context.Context) ([]*models.Rule, error)
}

func (f *fakeRuleSvc) GetAll(ctx context.Context) ([]*models.Rule, error) {
	if f.getAllFn != nil {
		return f.getAllFn(ctx)
	}
	return nil, nil
}

type fakeRealmSvc struct {
	Realm
	getFn     func(ctx context.Context, dto *models.GetRealmsDTO) ([]*models.Realm, error)
	getByIdFn func(ctx context.Context, req *models.GetRealmByIdDTO) (*models.Realm, error)
}

func (f *fakeRealmSvc) Get(ctx context.Context, dto *models.GetRealmsDTO) ([]*models.Realm, error) {
	if f.getFn != nil {
		return f.getFn(ctx, dto)
	}
	return nil, nil
}

func (f *fakeRealmSvc) GetById(ctx context.Context, req *models.GetRealmByIdDTO) (*models.Realm, error) {
	if f.getByIdFn != nil {
		return f.getByIdFn(ctx, req)
	}
	return &models.Realm{ID: req.ID}, nil
}

type fakeRoleSvc struct {
	Role
	getAllWithNamesFn func(ctx context.Context, dto *models.GetRolesDTO) ([]*models.RoleFull, error)
	getWithRealmFn    func(ctx context.Context, dto *models.GetRoleByRealmDTO) ([]*models.RoleWithRealm, error)
	getFn             func(ctx context.Context, name string) (*models.Role, error)
}

func (f *fakeRoleSvc) GetAllWithNames(ctx context.Context, dto *models.GetRolesDTO) ([]*models.RoleFull, error) {
	if f.getAllWithNamesFn != nil {
		return f.getAllWithNamesFn(ctx, dto)
	}
	return nil, nil
}

func (f *fakeRoleSvc) GetWithRealm(ctx context.Context, dto *models.GetRoleByRealmDTO) ([]*models.RoleWithRealm, error) {
	if f.getWithRealmFn != nil {
		return f.getWithRealmFn(ctx, dto)
	}
	return nil, nil
}

func (f *fakeRoleSvc) Get(ctx context.Context, name string) (*models.Role, error) {
	if f.getFn != nil {
		return f.getFn(ctx, name)
	}
	return nil, nil
}

type fakeAccessesSvc struct {
	Accesses
	getOriginalFn func(ctx context.Context) ([]*models.AccessesDTO, error)
}

func (f *fakeAccessesSvc) GetOriginal(ctx context.Context) ([]*models.AccessesDTO, error) {
	if f.getOriginalFn != nil {
		return f.getOriginalFn(ctx)
	}
	return nil, nil
}

type fakeReceivingSvc struct {
	Receiving
	forcedReceiptAllFn func(ctx context.Context) error
	forcedAllCalls     int
}

func (f *fakeReceivingSvc) ForcedReceiptAll(ctx context.Context) error {
	f.forcedAllCalls++
	if f.forcedReceiptAllFn != nil {
		return f.forcedReceiptAllFn(ctx)
	}
	return nil
}

type fakeUserSvc struct {
	User
	syncFn func(ctx context.Context) error
	calls  int
}

func (f *fakeUserSvc) Sync(ctx context.Context) error {
	f.calls++
	if f.syncFn != nil {
		return f.syncFn(ctx)
	}
	return nil
}

type fakeSISvc struct {
	SI
	getUsedFn func(ctx context.Context, req *models.Period) ([]*models.SiReceiving, error)
}

func (f *fakeSISvc) GetUsed(ctx context.Context, req *models.Period) ([]*models.SiReceiving, error) {
	if f.getUsedFn != nil {
		return f.getUsedFn(ctx, req)
	}
	return nil, nil
}

type fakeSectionSvc struct {
	Section
	getAllFn func(ctx context.Context, req *models.GetAllSectionsDTO) ([]*models.Section, error)
}

func (f *fakeSectionSvc) GetAll(ctx context.Context, req *models.GetAllSectionsDTO) ([]*models.Section, error) {
	if f.getAllFn != nil {
		return f.getAllFn(ctx, req)
	}
	return nil, nil
}

type fakeSIRepo struct {
	repository.SI
	getFn func(ctx context.Context, req *models.GetSiDTO) ([]*models.SI, error)
}

func (f *fakeSIRepo) Get(ctx context.Context, req *models.GetSiDTO) ([]*models.SI, error) {
	if f.getFn != nil {
		return f.getFn(ctx, req)
	}
	return nil, nil
}

type fakeVerificationSvc struct {
	Verification
	getLastFn         func(ctx context.Context, req *models.GetVerificationDTO) (*models.Verification, error)
	createFn          func(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error
	createSeveralFn   func(ctx context.Context, dto []*models.VerificationDTO) error
	updateFn          func(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error
	created           []*models.VerificationDTO
	createdSeveral    [][]*models.VerificationDTO
	updated           []*models.VerificationDTO
}

func (f *fakeVerificationSvc) GetLast(ctx context.Context, req *models.GetVerificationDTO) (*models.Verification, error) {
	if f.getLastFn != nil {
		return f.getLastFn(ctx, req)
	}
	return nil, models.ErrNoRows
}

func (f *fakeVerificationSvc) Create(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error {
	f.created = append(f.created, dto)
	if f.createFn != nil {
		return f.createFn(ctx, tx, dto)
	}
	return nil
}

func (f *fakeVerificationSvc) CreateSeveral(ctx context.Context, dto []*models.VerificationDTO) error {
	f.createdSeveral = append(f.createdSeveral, dto)
	if f.createSeveralFn != nil {
		return f.createSeveralFn(ctx, dto)
	}
	return nil
}

func (f *fakeVerificationSvc) Update(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error {
	f.updated = append(f.updated, dto)
	if f.updateFn != nil {
		return f.updateFn(ctx, tx, dto)
	}
	return nil
}

type fakeInstrumentRepo struct {
	repository.Instrument
	getByIdFn             func(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error)
	createInTxFn          func(ctx context.Context, tx postgres.Tx, dto *models.InstrumentDTO) error
	createSeveralFn       func(ctx context.Context, dto []*models.InstrumentDTO) error
	updateFn              func(ctx context.Context, tx postgres.Tx, dto *models.InstrumentDTO) error
	deleteFn              func(ctx context.Context, id string) error
	changePositionFn      func(ctx context.Context, dto *models.ChangePositionDTO) error
	changeStatusFn        func(ctx context.Context, tx postgres.Tx, dto *models.UpdateStatus) error
	changeSeveralStatusFn func(ctx context.Context, dto []*models.UpdateStatus) error
	created               []*models.InstrumentDTO
	createdSeveral        [][]*models.InstrumentDTO
	updated               []*models.InstrumentDTO
	deletedIds            []string
}

func (f *fakeInstrumentRepo) GetById(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error) {
	if f.getByIdFn != nil {
		return f.getByIdFn(ctx, req)
	}
	return nil, models.ErrNoRows
}

func (f *fakeInstrumentRepo) CreateInTx(ctx context.Context, tx postgres.Tx, dto *models.InstrumentDTO) error {
	f.created = append(f.created, dto)
	if f.createInTxFn != nil {
		return f.createInTxFn(ctx, tx, dto)
	}
	return nil
}

func (f *fakeInstrumentRepo) CreateSeveral(ctx context.Context, dto []*models.InstrumentDTO) error {
	f.createdSeveral = append(f.createdSeveral, dto)
	if f.createSeveralFn != nil {
		return f.createSeveralFn(ctx, dto)
	}
	return nil
}

func (f *fakeInstrumentRepo) Update(ctx context.Context, tx postgres.Tx, dto *models.InstrumentDTO) error {
	f.updated = append(f.updated, dto)
	if f.updateFn != nil {
		return f.updateFn(ctx, tx, dto)
	}
	return nil
}

func (f *fakeInstrumentRepo) Delete(ctx context.Context, id string) error {
	f.deletedIds = append(f.deletedIds, id)
	if f.deleteFn != nil {
		return f.deleteFn(ctx, id)
	}
	return nil
}

func (f *fakeInstrumentRepo) ChangePosition(ctx context.Context, dto *models.ChangePositionDTO) error {
	if f.changePositionFn != nil {
		return f.changePositionFn(ctx, dto)
	}
	return nil
}

func (f *fakeInstrumentRepo) ChangeStatus(ctx context.Context, tx postgres.Tx, dto *models.UpdateStatus) error {
	if f.changeStatusFn != nil {
		return f.changeStatusFn(ctx, tx, dto)
	}
	return nil
}

func (f *fakeInstrumentRepo) ChangeSeveralStatuses(ctx context.Context, dto []*models.UpdateStatus) error {
	if f.changeSeveralStatusFn != nil {
		return f.changeSeveralStatusFn(ctx, dto)
	}
	return nil
}

type fakeDepartmentRepo struct {
	repository.Department
	getAllFn func(ctx context.Context, req *models.GetDepartmentsDTO) ([]*models.Department, error)
	getByIdFn func(ctx context.Context, req *models.GetDepartmentByIdDTO) (*models.Department, error)
	createFn func(ctx context.Context, dto *models.DepartmentDTO) (string, error)
	deleteFn func(ctx context.Context, id string) error
	deleted  []string
}

func (f *fakeDepartmentRepo) GetAll(ctx context.Context, req *models.GetDepartmentsDTO) ([]*models.Department, error) {
	if f.getAllFn != nil {
		return f.getAllFn(ctx, req)
	}
	return nil, nil
}

func (f *fakeDepartmentRepo) GetById(ctx context.Context, req *models.GetDepartmentByIdDTO) (*models.Department, error) {
	if f.getByIdFn != nil {
		return f.getByIdFn(ctx, req)
	}
	return nil, models.ErrNoRows
}

func (f *fakeDepartmentRepo) Create(ctx context.Context, dto *models.DepartmentDTO) (string, error) {
	if f.createFn != nil {
		return f.createFn(ctx, dto)
	}
	return "dept-1", nil
}

func (f *fakeDepartmentRepo) Delete(ctx context.Context, id string) error {
	f.deleted = append(f.deleted, id)
	if f.deleteFn != nil {
		return f.deleteFn(ctx, id)
	}
	return nil
}

type fakeEmployeeRepo struct {
	repository.Employee
	getAllFn     func(ctx context.Context, req *models.GetEmployeesDTO) ([]*models.Employee, error)
	getUniqueFn  func(ctx context.Context, dto *models.GetUniqueEmployeeDTO) ([]*models.Employee, error)
	getByNameFn  func(ctx context.Context, req *models.GetEmployeeByNameDTO) (*models.Employee, error)
	getByIdFn    func(ctx context.Context, id string) (*models.Employee, error)
	getBySSOIdFn func(ctx context.Context, id string) (*models.Employee, error)
	getByMostFn  func(ctx context.Context, mostId string) (*models.EmployeeData, error)
	createFn     func(ctx context.Context, dto *models.WriteEmployeeDTO) error
	deleteFn     func(ctx context.Context, id string) error
	created      []*models.WriteEmployeeDTO
	deleted      []string
}

func (f *fakeEmployeeRepo) GetAll(ctx context.Context, req *models.GetEmployeesDTO) ([]*models.Employee, error) {
	if f.getAllFn != nil {
		return f.getAllFn(ctx, req)
	}
	return nil, nil
}

func (f *fakeEmployeeRepo) GetUnique(ctx context.Context, dto *models.GetUniqueEmployeeDTO) ([]*models.Employee, error) {
	if f.getUniqueFn != nil {
		return f.getUniqueFn(ctx, dto)
	}
	return nil, nil
}

func (f *fakeEmployeeRepo) GetByName(ctx context.Context, req *models.GetEmployeeByNameDTO) (*models.Employee, error) {
	if f.getByNameFn != nil {
		return f.getByNameFn(ctx, req)
	}
	return nil, models.ErrNoRows
}

func (f *fakeEmployeeRepo) GetById(ctx context.Context, id string) (*models.Employee, error) {
	if f.getByIdFn != nil {
		return f.getByIdFn(ctx, id)
	}
	return nil, models.ErrNoRows
}

func (f *fakeEmployeeRepo) GetBySSOId(ctx context.Context, id string) (*models.Employee, error) {
	if f.getBySSOIdFn != nil {
		return f.getBySSOIdFn(ctx, id)
	}
	return nil, models.ErrNoRows
}

func (f *fakeEmployeeRepo) GetByMostId(ctx context.Context, mostId string) (*models.EmployeeData, error) {
	if f.getByMostFn != nil {
		return f.getByMostFn(ctx, mostId)
	}
	return nil, models.ErrNoRows
}

func (f *fakeEmployeeRepo) Create(ctx context.Context, dto *models.WriteEmployeeDTO) error {
	f.created = append(f.created, dto)
	if f.createFn != nil {
		return f.createFn(ctx, dto)
	}
	return nil
}

func (f *fakeEmployeeRepo) Delete(ctx context.Context, id string) error {
	f.deleted = append(f.deleted, id)
	if f.deleteFn != nil {
		return f.deleteFn(ctx, id)
	}
	return nil
}

type fakeRepairSvc struct {
	Repair
	createSeveralFn func(ctx context.Context, dto []*models.RepairDTO) error
	createdSeveral  [][]*models.RepairDTO
}

func (f *fakeRepairSvc) CreateSeveral(ctx context.Context, dto []*models.RepairDTO) error {
	f.createdSeveral = append(f.createdSeveral, dto)
	if f.createSeveralFn != nil {
		return f.createSeveralFn(ctx, dto)
	}
	return nil
}

type fakePreservationSvc struct {
	Preservation
	createSeveralFn func(ctx context.Context, dto []*models.PreservationDTO) error
	createdSeveral  [][]*models.PreservationDTO
}

func (f *fakePreservationSvc) CreateSeveral(ctx context.Context, dto []*models.PreservationDTO) error {
	f.createdSeveral = append(f.createdSeveral, dto)
	if f.createSeveralFn != nil {
		return f.createSeveralFn(ctx, dto)
	}
	return nil
}

type fakeTransferToSaveSvc struct {
	TransferToSave
}

type fakeTransferToDepSvc struct {
	TransferToDepartment
	createSeveralFn func(ctx context.Context, dto []*models.TransferToDepartmentDTO) error
	createdSeveral  [][]*models.TransferToDepartmentDTO
}

func (f *fakeTransferToDepSvc) CreateSeveral(ctx context.Context, dto []*models.TransferToDepartmentDTO) error {
	f.createdSeveral = append(f.createdSeveral, dto)
	if f.createSeveralFn != nil {
		return f.createSeveralFn(ctx, dto)
	}
	return nil
}

type fakeWriteOffSvc struct {
	WriteOff
	createSeveralFn func(ctx context.Context, dto []*models.WriteOffDTO) error
	createdSeveral  [][]*models.WriteOffDTO
}

func (f *fakeWriteOffSvc) CreateSeveral(ctx context.Context, dto []*models.WriteOffDTO) error {
	f.createdSeveral = append(f.createdSeveral, dto)
	if f.createSeveralFn != nil {
		return f.createSeveralFn(ctx, dto)
	}
	return nil
}
