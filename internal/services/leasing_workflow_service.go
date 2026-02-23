package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	errs "github.com/HendraaaIrwn/honda-leasing-api/internal/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// LeasingWorkflowService contains end-to-end business process from request until delivery.
type LeasingWorkflowService interface {
	SubmitApplication(ctx context.Context, input SubmitApplicationInput) (*models.LeasingContract, error)
	ProcessAutoScoring(ctx context.Context, input AutoScoringDecisionInput) error
	ProcessSurveyResult(ctx context.Context, input SurveyDecisionInput) error
	ProcessFinalApproval(ctx context.Context, input FinalApprovalInput) error
	ExecuteAkad(ctx context.Context, input AkadInput) error
	RecordInitialPayment(ctx context.Context, input InitialPaymentInput) error
	ProcessDealerFulfillment(ctx context.Context, input DealerFulfillmentInput) error
	CompleteDelivery(ctx context.Context, input DeliveryCompletionInput) error
}

type leasingWorkflowService struct {
	db *gorm.DB
}

func NewLeasingWorkflowService(db *gorm.DB) LeasingWorkflowService {
	return &leasingWorkflowService{db: db}
}

func (s *leasingWorkflowService) SubmitApplication(ctx context.Context, input SubmitApplicationInput) (*models.LeasingContract, error) {
	if input.CustomerID < 1 || input.MotorID < 1 || input.ProductID < 1 || input.DPDibayar <= 0 {
		return nil, errs.ErrInvalidInput
	}

	requestDate := time.Now()
	if input.RequestDate != nil && !input.RequestDate.IsZero() {
		requestDate = *input.RequestDate
	}
	requestDate = requestDate.UTC()

	var created models.LeasingContract

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var motor models.Motor
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&motor, "motor_id = ?", input.MotorID).Error; err != nil {
			return err
		}
		if strings.ToLower(motor.StatusUnit) != MotorStatusReady {
			return errs.ErrMotorUnitNotReady
		}

		var product models.LeasingProduct
		if err := tx.First(&product, "product_id = ?", input.ProductID).Error; err != nil {
			return err
		}

		tenor := product.TenorBulan
		if input.TenorBulan > 0 {
			tenor = input.TenorBulan
		}
		if tenor <= 0 {
			return errs.ErrInvalidInput
		}

		nilaiKendaraan := motor.HargaOTR
		minDP := nilaiKendaraan * product.DPPersenMin / 100
		maxDP := nilaiKendaraan * product.DPPersenMax / 100
		if input.DPDibayar < minDP || input.DPDibayar > maxDP {
			return errs.ErrDPOutOfRange
		}

		pokokPinjaman := nilaiKendaraan - input.DPDibayar
		totalPinjaman := calculateTotalPinjaman(pokokPinjaman, product.BungaFlat, tenor)
		cicilan := totalPinjaman / float64(tenor)
		mulaiCicil := requestDate.AddDate(0, 1, 0)

		contract := models.LeasingContract{
			RequestDate:       requestDate,
			TanggalMulaiCicil: mulaiCicil,
			TenorBulan:        tenor,
			NilaiKendaraan:    nilaiKendaraan,
			DPDibayar:         input.DPDibayar,
			PokokPinjaman:     pokokPinjaman,
			TotalPinjaman:     totalPinjaman,
			CicilanPerBulan:   cicilan,
			Status:            ContractStatusDraft,
			CustomerID:        input.CustomerID,
			MotorID:           input.MotorID,
			ProductID:         input.ProductID,
		}
		if err := tx.Create(&contract).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Motor{}).
			Where("motor_id = ?", motor.MotorID).
			Update("status_unit", MotorStatusBooked).Error; err != nil {
			return err
		}

		for _, doc := range input.Documents {
			if strings.TrimSpace(doc.FileName) == "" || strings.TrimSpace(doc.FileURL) == "" {
				continue
			}
			document := models.LeasingContractDocument{
				FileName:   doc.FileName,
				FileSize:   doc.FileSize,
				FileType:   doc.FileType,
				FileURL:    doc.FileURL,
				ContractID: contract.ContractID,
			}
			if err := tx.Create(&document).Error; err != nil {
				return err
			}
		}

		if err := s.bootstrapTasksFromTemplate(tx, contract.ContractID, requestDate); err != nil {
			return err
		}

		created = contract
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (s *leasingWorkflowService) ProcessAutoScoring(ctx context.Context, input AutoScoringDecisionInput) error {
	if input.ContractID < 1 {
		return errs.ErrInvalidInput
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		contract, err := s.lockContract(tx, input.ContractID)
		if err != nil {
			return err
		}
		if contract.Status != ContractStatusDraft {
			return errs.ErrContractNotDraft
		}

		if err := s.updateTaskStatusByKeyword(tx, contract.ContractID, "scoring", TaskStatusCompleted); err != nil {
			return err
		}

		if input.AutoApproved {
			if err := s.transitionContractStatus(tx, contract, ContractStatusApproved); err != nil {
				return err
			}
			if err := s.updateTaskStatusByKeyword(tx, contract.ContractID, "pre-approval", TaskStatusCompleted); err != nil {
				return err
			}
			return s.appendTaskNoteByKeyword(tx, contract.ContractID, "scoring", "auto_scoring_note", input.Note, TaskAttrStatusCompleted)
		}

		if err := s.updateTaskStatusByKeyword(tx, contract.ContractID, "review", TaskStatusInProgress); err != nil {
			return err
		}
		if err := s.appendTaskNoteByKeyword(tx, contract.ContractID, "review", "manual_review_note", input.Note, TaskAttrStatusPending); err != nil {
			return err
		}

		if !input.ManualReviewReady {
			return nil
		}

		if input.ManualApproved {
			if err := s.transitionContractStatus(tx, contract, ContractStatusApproved); err != nil {
				return err
			}
			return s.updateTaskStatusByKeyword(tx, contract.ContractID, "review", TaskStatusCompleted)
		}

		if err := s.transitionContractStatus(tx, contract, ContractStatusCanceled); err != nil {
			return err
		}
		if err := s.releaseMotorIfBooked(tx, contract.MotorID); err != nil {
			return err
		}
		return s.updateTaskStatusByKeyword(tx, contract.ContractID, "review", TaskStatusCancelled)
	})
}

func (s *leasingWorkflowService) ProcessSurveyResult(ctx context.Context, input SurveyDecisionInput) error {
	if input.ContractID < 1 {
		return errs.ErrInvalidInput
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		contract, err := s.lockContract(tx, input.ContractID)
		if err != nil {
			return err
		}
		if contract.Status != ContractStatusApproved {
			return errs.ErrContractNotApproved
		}

		switch input.Decision {
		case SurveyDecisionApprove:
			if err := s.updateTaskStatusByKeyword(tx, contract.ContractID, "survei", TaskStatusCompleted); err != nil {
				return err
			}
			return s.appendTaskNoteByKeyword(tx, contract.ContractID, "survei", "survey_note", input.Note, TaskAttrStatusCompleted)

		case SurveyDecisionReject:
			if err := s.transitionContractStatus(tx, contract, ContractStatusCanceled); err != nil {
				return err
			}
			if err := s.releaseMotorIfBooked(tx, contract.MotorID); err != nil {
				return err
			}
			if err := s.updateTaskStatusByKeyword(tx, contract.ContractID, "survei", TaskStatusCancelled); err != nil {
				return err
			}
			return s.appendTaskNoteByKeyword(tx, contract.ContractID, "survei", "survey_reject_reason", input.Note, TaskAttrStatusCancelled)

		case SurveyDecisionRequestAdditionalDP:
			if input.AdditionalDP <= contract.DPDibayar {
				return errs.ErrDPOutOfRange
			}

			var product models.LeasingProduct
			if err := tx.First(&product, "product_id = ?", contract.ProductID).Error; err != nil {
				return err
			}
			minDP := contract.NilaiKendaraan * product.DPPersenMin / 100
			maxDP := contract.NilaiKendaraan * product.DPPersenMax / 100
			if input.AdditionalDP < minDP || input.AdditionalDP > maxDP {
				return errs.ErrDPOutOfRange
			}

			pokokPinjaman := contract.NilaiKendaraan - input.AdditionalDP
			totalPinjaman := calculateTotalPinjaman(pokokPinjaman, product.BungaFlat, contract.TenorBulan)
			cicilan := totalPinjaman / float64(contract.TenorBulan)

			updates := map[string]interface{}{
				"dp_dibayar":        input.AdditionalDP,
				"pokok_pinjaman":    pokokPinjaman,
				"total_pinjaman":    totalPinjaman,
				"cicilan_per_bulan": cicilan,
				"status":            ContractStatusDraft,
			}
			if err := tx.Model(&models.LeasingContract{}).
				Where("contract_id = ?", contract.ContractID).
				Updates(updates).Error; err != nil {
				return err
			}
			if err := s.updateTaskStatusByKeyword(tx, contract.ContractID, "survei", TaskStatusCompleted); err != nil {
				return err
			}
			return s.appendTaskNoteByKeyword(tx, contract.ContractID, "survei", "additional_dp_request", input.Note, TaskAttrStatusPending)

		default:
			return errs.ErrInvalidDecision
		}
	})
}

func (s *leasingWorkflowService) ProcessFinalApproval(ctx context.Context, input FinalApprovalInput) error {
	if input.ContractID < 1 {
		return errs.ErrInvalidInput
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		contract, err := s.lockContract(tx, input.ContractID)
		if err != nil {
			return err
		}
		if contract.Status != ContractStatusApproved {
			return errs.ErrContractNotApproved
		}

		if input.Approved {
			if err := s.updateTaskStatusByKeyword(tx, contract.ContractID, "approval", TaskStatusCompleted); err != nil {
				return err
			}
			return s.appendTaskNoteByKeyword(tx, contract.ContractID, "approval", "final_approval_note", input.Note, TaskAttrStatusCompleted)
		}

		if err := s.transitionContractStatus(tx, contract, ContractStatusCanceled); err != nil {
			return err
		}
		if err := s.releaseMotorIfBooked(tx, contract.MotorID); err != nil {
			return err
		}
		if err := s.updateTaskStatusByKeyword(tx, contract.ContractID, "approval", TaskStatusCancelled); err != nil {
			return err
		}
		return s.appendTaskNoteByKeyword(tx, contract.ContractID, "approval", "final_approval_reject_reason", input.Note, TaskAttrStatusCancelled)
	})
}

func (s *leasingWorkflowService) ExecuteAkad(ctx context.Context, input AkadInput) error {
	if input.ContractID < 1 {
		return errs.ErrInvalidInput
	}

	akadDate := input.AkadDate
	if akadDate.IsZero() {
		akadDate = time.Now()
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		contract, err := s.lockContract(tx, input.ContractID)
		if err != nil {
			return err
		}
		if contract.Status != ContractStatusApproved {
			return errs.ErrContractNotApproved
		}

		contractNumber := strings.TrimSpace(input.ContractNumber)
		if contractNumber == "" && (contract.ContractNumber == nil || strings.TrimSpace(*contract.ContractNumber) == "") {
			if input.GenerateContractCode {
				contractNumber = fmt.Sprintf("KTR-%d-%06d", time.Now().Year(), contract.ContractID)
			}
		}

		mulaiCicil := akadDate.AddDate(0, 1, 0)
		if input.TanggalMulaiCicil != nil && !input.TanggalMulaiCicil.IsZero() {
			mulaiCicil = *input.TanggalMulaiCicil
		}

		updates := map[string]interface{}{
			"tanggal_akad":        akadDate,
			"tanggal_mulai_cicil": mulaiCicil,
		}
		if contractNumber != "" {
			updates["contract_number"] = contractNumber
		}

		if err := tx.Model(&models.LeasingContract{}).
			Where("contract_id = ?", contract.ContractID).
			Updates(updates).Error; err != nil {
			return err
		}

		return s.updateTaskStatusByKeyword(tx, contract.ContractID, "akad", TaskStatusCompleted)
	})
}

func (s *leasingWorkflowService) RecordInitialPayment(ctx context.Context, input InitialPaymentInput) error {
	if input.ContractID < 1 || strings.TrimSpace(input.NomorBukti) == "" || input.JumlahBayar <= 0 {
		return errs.ErrInvalidInput
	}
	if input.JumlahBayar <= 0 {
		return errs.ErrInvalidPaymentAmount
	}

	tanggalBayar := input.TanggalBayar
	if tanggalBayar.IsZero() {
		tanggalBayar = time.Now()
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		contract, err := s.lockContract(tx, input.ContractID)
		if err != nil {
			return err
		}
		if contract.Status != ContractStatusApproved && contract.Status != ContractStatusActive {
			return errs.ErrInvalidStatusTransition
		}

		payment := models.Payment{
			NomorBukti:       strings.TrimSpace(input.NomorBukti),
			JumlahBayar:      input.JumlahBayar,
			TanggalBayar:     tanggalBayar,
			MetodePembayaran: strings.TrimSpace(input.MetodePembayaran),
			Provider:         strings.TrimSpace(input.Provider),
			ContractID:       contract.ContractID,
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		return s.updateTaskStatusByKeyword(tx, contract.ContractID, "pembayaran dp", TaskStatusCompleted)
	})
}

func (s *leasingWorkflowService) ProcessDealerFulfillment(ctx context.Context, input DealerFulfillmentInput) error {
	if input.ContractID < 1 {
		return errs.ErrInvalidInput
	}
	if !input.UnitReadyStock && input.EstimatedIndentWeek < 0 {
		return errs.ErrInvalidInput
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		contract, err := s.lockContract(tx, input.ContractID)
		if err != nil {
			return err
		}
		if contract.Status != ContractStatusApproved {
			return errs.ErrContractNotApproved
		}

		if input.UnitReadyStock {
			if err := s.updateTaskStatusByKeyword(tx, contract.ContractID, "po", TaskStatusCompleted); err != nil {
				return err
			}
			return s.appendTaskNoteByKeyword(tx, contract.ContractID, "po", "unit_stock_note", "unit ready stock", TaskAttrStatusCompleted)
		}

		note := strings.TrimSpace(input.Note)
		if note == "" {
			note = fmt.Sprintf("inden unit %d minggu", input.EstimatedIndentWeek)
		}
		if err := s.updateTaskStatusByKeyword(tx, contract.ContractID, "po", TaskStatusInProgress); err != nil {
			return err
		}
		return s.appendTaskNoteByKeyword(tx, contract.ContractID, "po", "indent_info", note, TaskAttrStatusPending)
	})
}

func (s *leasingWorkflowService) CompleteDelivery(ctx context.Context, input DeliveryCompletionInput) error {
	if input.ContractID < 1 || !input.CustomerReceived {
		return errs.ErrInvalidInput
	}

	deliveryDate := input.DeliveryDate
	if deliveryDate.IsZero() {
		deliveryDate = time.Now()
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		contract, err := s.lockContract(tx, input.ContractID)
		if err != nil {
			return err
		}
		if contract.Status != ContractStatusApproved {
			return errs.ErrContractNotApproved
		}

		if err := tx.Model(&models.Motor{}).
			Where("motor_id = ?", contract.MotorID).
			Update("status_unit", MotorStatusLeased).Error; err != nil {
			return err
		}

		updates := map[string]interface{}{
			"status": ContractStatusActive,
		}
		if input.TanggalMulaiCicil != nil && !input.TanggalMulaiCicil.IsZero() {
			updates["tanggal_mulai_cicil"] = *input.TanggalMulaiCicil
		}

		if err := tx.Model(&models.LeasingContract{}).
			Where("contract_id = ?", contract.ContractID).
			Updates(updates).Error; err != nil {
			return err
		}

		for _, doc := range input.ContractDocUploads {
			if strings.TrimSpace(doc.FileName) == "" || strings.TrimSpace(doc.FileURL) == "" {
				continue
			}
			document := models.LeasingContractDocument{
				FileName:   doc.FileName,
				FileSize:   doc.FileSize,
				FileType:   doc.FileType,
				FileURL:    doc.FileURL,
				ContractID: contract.ContractID,
			}
			if err := tx.Create(&document).Error; err != nil {
				return err
			}
		}

		if err := s.updateTaskStatusByKeyword(tx, contract.ContractID, "delivery", TaskStatusCompleted); err != nil {
			return err
		}
		if input.DocumentHandover {
			if err := s.appendTaskNoteByKeyword(tx, contract.ContractID, "delivery", "document_handover", input.HandoverNote, TaskAttrStatusCompleted); err != nil {
				return err
			}
		}

		return s.appendTaskNoteByKeyword(tx, contract.ContractID, "delivery", "delivery_date", deliveryDate.Format("2006-01-02"), TaskAttrStatusCompleted)
	})
}

func (s *leasingWorkflowService) bootstrapTasksFromTemplate(tx *gorm.DB, contractID int64, startDate time.Time) error {
	var templates []models.TemplateTask
	if err := tx.Order("teta_id ASC").Preload("Attributes").Find(&templates).Error; err != nil {
		return err
	}
	if len(templates) == 0 {
		return nil
	}

	for idx, tmpl := range templates {
		task := models.LeasingTask{
			TaskName:   tmpl.TetaName,
			StartDate:  startDate,
			EndDate:    startDate.AddDate(0, 0, 14),
			SequenceNo: idx + 1,
			Status:     TaskStatusInProgress,
			ContractID: contractID,
			RoleID:     tmpl.TetaRoleID,
		}
		if err := tx.Create(&task).Error; err != nil {
			return err
		}

		for _, attr := range tmpl.Attributes {
			taskAttr := models.LeasingTaskAttribute{
				TasaName:   attr.TetatName,
				TasaValue:  "",
				TasaStatus: TaskAttrStatusPending,
				TasaLetaID: task.TaskID,
			}
			if err := tx.Create(&taskAttr).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *leasingWorkflowService) lockContract(tx *gorm.DB, contractID int64) (*models.LeasingContract, error) {
	var contract models.LeasingContract
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&contract, "contract_id = ?", contractID).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}

func (s *leasingWorkflowService) transitionContractStatus(tx *gorm.DB, contract *models.LeasingContract, nextStatus string) error {
	if contract.Status == nextStatus {
		return nil
	}
	if !isAllowedContractTransition(contract.Status, nextStatus) {
		return errs.ErrInvalidStatusTransition
	}
	if err := tx.Model(&models.LeasingContract{}).
		Where("contract_id = ?", contract.ContractID).
		Update("status", nextStatus).Error; err != nil {
		return err
	}
	contract.Status = nextStatus
	return nil
}

func isAllowedContractTransition(current, next string) bool {
	allowed := map[string]map[string]bool{
		ContractStatusDraft: {
			ContractStatusApproved: true,
			ContractStatusCanceled: true,
		},
		ContractStatusApproved: {
			ContractStatusDraft:    true, // additional DP / rework setelah survey
			ContractStatusActive:   true,
			ContractStatusCanceled: true,
		},
		ContractStatusActive: {
			ContractStatusCanceled: true,
		},
	}

	return allowed[current][next]
}

func (s *leasingWorkflowService) releaseMotorIfBooked(tx *gorm.DB, motorID int64) error {
	return tx.Model(&models.Motor{}).
		Where("motor_id = ? AND status_unit = ?", motorID, MotorStatusBooked).
		Update("status_unit", MotorStatusReady).Error
}

func (s *leasingWorkflowService) updateTaskStatusByKeyword(tx *gorm.DB, contractID int64, keyword, status string) error {
	keyword = strings.TrimSpace(strings.ToLower(keyword))
	if keyword == "" {
		return nil
	}

	now := time.Now().UTC()
	updates := map[string]interface{}{
		"status": status,
	}
	if status == TaskStatusCompleted {
		updates["actual_startdate"] = now
		updates["actual_enddate"] = now
	}
	if status == TaskStatusCancelled {
		updates["actual_enddate"] = now
	}

	return tx.Model(&models.LeasingTask{}).
		Where("contract_id = ? AND LOWER(task_name) LIKE ?", contractID, "%"+keyword+"%").
		Updates(updates).Error
}

func (s *leasingWorkflowService) appendTaskNoteByKeyword(tx *gorm.DB, contractID int64, keyword, name, note, status string) error {
	note = strings.TrimSpace(note)
	if note == "" {
		return nil
	}

	var task models.LeasingTask
	err := tx.Where("contract_id = ? AND LOWER(task_name) LIKE ?", contractID, "%"+strings.ToLower(strings.TrimSpace(keyword))+"%").
		Order("sequence_no ASC").
		First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	attr := models.LeasingTaskAttribute{
		TasaName:   strings.TrimSpace(name),
		TasaValue:  note,
		TasaStatus: status,
		TasaLetaID: task.TaskID,
	}
	return tx.Create(&attr).Error
}

func calculateTotalPinjaman(pokokPinjaman, bungaFlat float64, tenorBulan int16) float64 {
	if tenorBulan <= 0 {
		return pokokPinjaman
	}
	margin := pokokPinjaman * (bungaFlat / 100) * (float64(tenorBulan) / 12)
	return pokokPinjaman + margin
}
