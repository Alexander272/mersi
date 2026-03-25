package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres/pq_models"
	"github.com/goodsign/monday"
	"github.com/jmoiron/sqlx"
)

type SIRepo struct {
	db *sqlx.DB
	Transaction
}

func NewSIRepo(db *sqlx.DB, transaction Transaction) *SIRepo {
	return &SIRepo{
		db:          db,
		Transaction: transaction,
	}
}

type SI interface {
	Transaction
	Get(ctx context.Context, req *models.GetSiDTO) ([]*models.SI, error)
	GetVerification(ctx context.Context, req *models.Period) ([]*models.SiVerification, error)
	GetSent(ctx context.Context, req *models.GetSiDTO) ([]*models.SiReceiving, error)
	GetUsed(ctx context.Context, req *models.Period) ([]*models.SiReceiving, error)
	GetLog(ctx context.Context, req *models.Period) ([]*models.SiWithLog, error)
}

var siFieldsMap = map[string]string{
	"id":                        "i.id",
	"position":                  "i.position",
	"name":                      "i.name",
	"dateOfReceipt":             "i.date_of_receipt",
	"type":                      "i.type",
	"factoryNumber":             "i.factory_number",
	"measurementLimits":         "i.measurement_limits",
	"accuracy":                  "i.accuracy",
	"stateRegister":             "i.state_register",
	"countryOfProduce":          "i.country_of_produce",
	"manufacturer":              "i.manufacturer",
	"responsible":               "i.responsible",
	"inventory":                 "i.inventory",
	"yearOfIssue":               "i.year_of_issue",
	"interVerificationInterval": "i.inter_verification_interval",
	"actOfEntering":             "i.act_of_entering",
	"actOfEnteringId":           "i.act_of_entering_id",
	"repairInfo":                "r.period_start",
	"notes":                     "i.notes",
	"verificationDate":          "v.date",
	"nextVerificationDate":      "v.next_date",
	"department":                "l.department_id",
	"place":                     "l.place",
	"person":                    "l.person",
	"status":                    "l.status",
}

func (r *SIRepo) formatField(field string) string {
	if f, ok := siFieldsMap[field]; ok {
		return f
	}
	// Если поле не найдено, возвращаем безопасный дефолт,чтобы не упал SQL запрос из-за пустой строки
	return "i.id"
}

func (r *SIRepo) Get(ctx context.Context, req *models.GetSiDTO) ([]*models.SI, error) {
	tmp := []*pq_models.SI{}
	params := []interface{}{req.SectionId, req.Status}

	// 1. Сортировка (Безопасная)
	orderClause := r.buildOrderClause(req.Sort)
	// 2. Фильтры
	filterClause, params := r.buildFilterClause(req.Filters, params)
	// 3. Поиск
	searchClause, params := r.buildSearchClause(req.Search, params)

	// 4. Пагинация
	limitIdx := len(params) + 1
	offsetIdx := len(params) + 2
	params = append(params, req.Page.Limit, req.Page.Offset)

	query := fmt.Sprintf(`
		WITH last_verification AS (
			SELECT DISTINCT ON (instrument_id) id, instrument_id, date, next_date 
			FROM %s ORDER BY instrument_id, date DESC, created_at DESC
		),
		last_repair AS (
			SELECT DISTINCT ON (instrument_id) instrument_id, period_start, period_end, work 
			FROM %s ORDER BY instrument_id, period_start DESC
		),
		last_preservation AS (
			SELECT DISTINCT ON (instrument_id) instrument_id, date_start, date_end 
			FROM %s ORDER BY instrument_id, date_start DESC
		),
		last_transfer_save AS (
			SELECT DISTINCT ON (instrument_id) instrument_id, date_start, date_end 
			FROM %s ORDER BY instrument_id, date_start DESC
		),
		last_transfer_dep AS (
			SELECT DISTINCT ON (instrument_id) instrument_id, doc_name 
			FROM %s ORDER BY instrument_id, date DESC
		),
		last_write_off AS (
			SELECT DISTINCT ON (instrument_id) instrument_id, doc_name 
			FROM %s ORDER BY instrument_id, date DESC
		),
		last_location AS (
			SELECT DISTINCT ON (instrument_id) 
				instrument_id, l.status, l.department_id, last_place_id, person_id,
				COALESCE(lp.name,last_place) AS last_place, 
				COALESCE(e.name, NULLIF(person, ''), '') AS person,
				CASE 
					WHEN l.status = '%s' THEN COALESCE(dep.name, NULLIF(place, ''), '')
					WHEN l.status = '%s' THEN 'Резерв'
					ELSE 
						CASE 
							WHEN l.last_place != '' OR l.last_place_id IS NOT NULL 
							THEN 'Перемещение из «' || COALESCE(lp.name, l.last_place) || '»'
							ELSE 'Перемещение' 
						END
				END AS place
			FROM %s l
			LEFT JOIN %s lp ON l.last_place_id = lp.id::text
			LEFT JOIN %s e ON l.person_id = e.id
			LEFT JOIN %s dep ON l.department_id = dep.id
			ORDER BY instrument_id, date_of_issue DESC, l.created_at DESC
		)

		SELECT 
			i.id, i.position, i.name, i.date_of_receipt, i.type, i.factory_number, 
			i.measurement_limits, i.accuracy, i.state_register, i.country_of_produce, 
			i.manufacturer, i.responsible, i.inventory, i.year_of_issue, 
			i.inter_verification_interval, i.act_of_entering, i.act_of_entering_id, i.notes,
			
			COALESCE(v.date, '0001-01-01'::DATE) AS date, 
			COALESCE(v.next_date, '0001-01-01'::DATE) AS next_date,
			COALESCE(vd.name, '') AS certificate, 
			COALESCE(vd.doc_id::text, '') AS certificate_id,
			
			COALESCE(r.work, '') AS repair_work,
			COALESCE(r.period_start, '0001-01-01'::DATE) AS repair_start, 
			COALESCE(r.period_end, '0001-01-01'::DATE) AS repair_end,
			
			COALESCE(p.date_start, '0001-01-01'::DATE) AS preservation, 
			COALESCE(p.date_end, '0001-01-01'::DATE) AS de_preservation,
			
			COALESCE(ts.date_start, '0001-01-01'::DATE) AS transfer_date, 
			COALESCE(ts.date_end, '0001-01-01'::DATE) AS return_date,
			
			COALESCE(td.doc_name, '') AS transfer_to_dep, 
			COALESCE(wo.doc_name, '') AS write_off,
			
			COALESCE(l.status, 'used') AS status,
			COALESCE(l.person, '') AS person,
			COALESCE(l.place, '') AS place,
			COALESCE(l.last_place, '') AS last_place,    
			
			COUNT(*) OVER() AS total

		FROM %s AS i
		LEFT JOIN last_verification v ON v.instrument_id = i.id
		LEFT JOIN %s vd ON vd.verification_id = v.id
		LEFT JOIN last_repair r ON r.instrument_id = i.id
		LEFT JOIN last_preservation p ON p.instrument_id = i.id
		LEFT JOIN last_transfer_save ts ON ts.instrument_id = i.id
		LEFT JOIN last_transfer_dep td ON td.instrument_id = i.id
		LEFT JOIN last_write_off wo ON wo.instrument_id = i.id
		LEFT JOIN last_location l ON l.instrument_id = i.id

        WHERE i.section_id = $1 AND i.status = $2 %s %s %s
        LIMIT $%d OFFSET $%d`,
		VerificationTable,               // last_verification
		RepairTable,                     // last_repair
		PreservationTable,               // last_preservation
		TransferToSaveTable,             // last_transfer_save
		TransferToDepTable,              // last_transfer_dep
		WriteOffTable,                   // last_write_off
		constants.LocationStatusUsed,    // status check 1
		constants.LocationStatusReserve, // status check 2
		LocationTable,                   // last_location
		DepartmentTable,                 // lp
		EmployeeTable,                   // e
		DepartmentTable,                 // dep
		InstrumentsTable,                // main FROM
		VerificationDocsTable,           // vd join
		filterClause,
		searchClause,
		orderClause,
		limitIdx, offsetIdx,
	)

	// logger.Debug("get si", logger.StringAttr("query", query))
	// logger.Debug("get si", logger.AnyAttr("params", params))

	if err := r.db.SelectContext(ctx, &tmp, query, params...); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	data := r.mapToDomain(tmp)

	return data, nil
}

func (r *SIRepo) buildOrderClause(sorts []*models.Sort) string {
	if len(sorts) == 0 {
		return "ORDER BY position, i.created_at, i.id"
	}

	var parts []string
	for _, s := range sorts {
		field := r.formatField(s.Field)
		direction := "ASC"
		if strings.ToUpper(s.Type) == "DESC" {
			direction = "DESC"
		}
		parts = append(parts, fmt.Sprintf("%s %s", field, direction))
	}
	parts = append(parts, "i.created_at", "i.id")

	return "ORDER BY " + strings.Join(parts, ", ")
}
func (r *SIRepo) buildFilterClause(filters []*models.Filter, params []interface{}) (string, []interface{}) {
	if len(filters) == 0 {
		return "", params
	}

	var clauses []string

	for _, f := range filters {
		field := r.formatField(f.Field)

		// Специальная логика для департамента
		if f.Field == "department" {
			val := strings.ReplaceAll(f.Values[0].Value, ",", "|")
			params = append(params, val)
			idx := len(params)

			clauses = append(clauses, fmt.Sprintf("(%s OR (%s AND %s='moved'))",
				getFilterLine(f.Values[0].CompareType, field, idx),
				getFilterLine(f.Values[0].CompareType, "last_place_id", idx),
				r.formatField("status"),
			))
			continue
		}

		// Общая логика для остальных фильтров
		for _, sv := range f.Values {
			val := sv.Value
			if sv.CompareType == "in" {
				val = strings.ReplaceAll(val, ",", "|")
			}
			params = append(params, sv.Value)
			idx := len(params)

			clauses = append(clauses, getFilterLine(sv.CompareType, field, idx))
		}
	}

	if len(clauses) == 0 {
		return "", params
	}

	return "AND " + strings.Join(clauses, " AND "), params
}
func (r *SIRepo) buildSearchClause(search *models.Search, params []interface{}) (string, []interface{}) {
	if search == nil || len(search.Fields) == 0 || search.Value == "" {
		return "", params
	}

	params = append(params, search.Value)
	idx := len(params)

	list := []string{}
	for _, f := range search.Fields {
		list = append(list, fmt.Sprintf("%s ILIKE '%%'||$%d||'%%'", r.formatField(f), idx))
	}

	return "AND (" + strings.Join(list, " OR ") + ")", params
}

func (r *SIRepo) mapToDomain(tmp []*pq_models.SI) []*models.SI {
	data := []*models.SI{}
	for _, d := range tmp {
		repair := ""
		if !d.RepairStart.IsZero() {
			d.RepairStart = d.RepairStart.In(time.Local)
			repair = monday.Format(d.RepairStart, "Jan 2006", monday.LocaleRuRU)
		}
		if !d.RepairEnd.IsZero() {
			d.RepairEnd = d.RepairEnd.In(time.Local)
			repair += " - " + monday.Format(d.RepairEnd, "Jan 2006", monday.LocaleRuRU)
		}
		if d.RepairWork != "" {
			repair += " (" + d.RepairWork + ")"
		}

		t := d.ToModel()
		t.RepairInfo = repair
		data = append(data, t)
	}

	return data
}

func (r *SIRepo) GetVerification(ctx context.Context, req *models.Period) ([]*models.SiVerification, error) {
	query := fmt.Sprintf(`SELECT i.id, i.name, type, factory_number, year_of_issue, state_register, measurement_limits, date, next_date, 
		inter_verification_interval, manufacturer, notes, notification_channel, bid_type
		FROM %s AS i
		LEFT JOIN %s AS s ON s.id=section_id
		LEFT JOIN %s AS r ON r.id=realm_id
		LEFT JOIN LATERAL (SELECT date, next_date FROM %s WHERE instrument_id=i.id ORDER BY date DESC, created_at DESC LIMIT 1) AS v ON TRUE
		LEFT JOIN LATERAL (SELECT date AS write_off FROM %s WHERE instrument_id=i.id) AS w ON TRUE
		LEFT JOIN LATERAL (SELECT date_start AS preservation FROM %s WHERE instrument_id=i.id AND date_end='0001-01-01'::DATE) AS p ON TRUE 
		LEFT JOIN LATERAL (SELECT date_start AS transferred FROM %s WHERE instrument_id=i.id AND date_end='0001-01-01'::DATE) AS t ON TRUE
 		WHERE next_date>=$1 AND next_date<=$2 AND (CASE WHEN $3!='' THEN section_id::text=$3 ELSE true END) AND 
		is_active=true AND expiration_notice=true AND (notification_channel!='' OR $4) AND 
		deleted IS NULL AND write_off IS NULL AND preservation IS NULL AND transferred IS NULL
		ORDER BY notification_channel, i.position`,
		InstrumentsTable, SectionTable, RealmTable, VerificationTable, WriteOffTable, PreservationTable, TransferToSaveTable,
	)
	//TODO проверить

	tmp := []*pq_models.SI{}
	if err := r.db.SelectContext(ctx, &tmp, query, req.StartAt, req.FinishAt, req.SectionId, req.ChannelIsOption); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	data := []*models.SiVerification{}
	for i, s := range tmp {
		si := &models.SI{
			Id:                        s.Id,
			Name:                      s.Name,
			Type:                      s.Type,
			FactoryNumber:             s.FactoryNumber,
			YearOfIssue:               s.YearOfIssue,
			StateRegister:             s.StateRegister,
			MeasurementLimits:         s.MeasurementLimits,
			VerificationDate:          s.VerificationDate,
			NextVerificationDate:      s.NextVerificationDate,
			InterVerificationInterval: s.InterVerificationInterval,
			Manufacturer:              s.Manufacturer,
			Notes:                     s.Notes,
		}

		if i == 0 || data[len(data)-1].NotificationChannel != s.NotificationChannel {
			data = append(data, &models.SiVerification{
				NotificationChannel: s.NotificationChannel,
				BidType:             s.BidType,
				SI:                  []*models.SI{si},
			})
		} else {
			data[len(data)-1].SI = append(data[len(data)-1].SI, si)
		}
	}

	return data, nil
}

func (r *SIRepo) GetSent(ctx context.Context, req *models.GetSiDTO) ([]*models.SiReceiving, error) {
	params := []interface{}{constants.LocationStatusMoved, req.SectionId}
	filter := ""
	count := len(params) + 1

	for _, f := range req.Filters {
		for _, sv := range f.Values {
			filter += " AND " + getFilterLine(sv.CompareType, r.formatField(f.Field), count)
			if sv.CompareType == "in" {
				sv.Value = strings.ReplaceAll(sv.Value, ",", "|")
			}
			if sv.CompareType != "null" {
				params = append(params, sv.Value)
				count++
			}
		}
	}

	query := fmt.Sprintf(`SELECT i.id, i.name, factory_number, year_of_issue, state_register, measurement_limits, date, next_date,
		COALESCE(person, e.emp, '') AS person, COALESCE(place, d.dep, '') AS place, COALESCE(l.last_place, lp.name, '') AS last_place,
		COALESCE(most_channel_id, channel, '') AS notification_channel
		FROM %s AS i
		LEFT JOIN LATERAL (SELECT date, next_date FROM %s WHERE instrument_id=i.id ORDER BY date DESC, created_at DESC LIMIT 1) AS v ON TRUE
		LEFT JOIN LATERAL (SELECT status, NULLIF(person,'') AS person, NULLIF(place,'') AS place, NULLIF(last_place,'') AS last_place, 
			person_id, department_id, last_place_id FROM %s WHERE instrument_id=i.id 
			ORDER BY date_of_issue DESC, created_at DESC LIMIT 1) AS l ON TRUE
		LEFT JOIN LATERAL (SELECT NULLIF(name,'') AS emp FROM %s WHERE l.person_id::uuid=id) AS e ON TRUE
		LEFT JOIN LATERAL (SELECT NULLIF(name,'') AS dep, channel_id FROM %s WHERE l.department_id::uuid=id) AS d ON TRUE
		LEFT JOIN LATERAL (SELECT NULLIF(name,'') AS name FROM %s WHERE l.last_place_id=id::text) AS lp ON TRUE
		LEFT JOIN LATERAL (SELECT most_channel_id FROM %s WHERE id=d.channel_id) AS c ON TRUE
		LEFT JOIN LATERAL (SELECT notification_channel AS channel FROM %s AS r INNER JOIN %s AS s ON s.realm_id=r.id WHERE s.id=i.section_id) AS r ON TRUE
		WHERE l.status=$1 AND CASE WHEN $2!='' THEN section_id::text=$2 ELSE TRUE END %s
		ORDER BY channel_id, place, last_place, next_date`,
		InstrumentsTable, VerificationTable, LocationTable, EmployeeTable, DepartmentTable, DepartmentTable,
		ChannelTable, RealmTable, SectionTable, filter,
	)
	tmp := []*pq_models.SI{}

	if err := r.db.SelectContext(ctx, &tmp, query, params...); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	data := []*models.SiReceiving{}
	for i, s := range tmp {
		status := constants.LocationStatusUsed
		if s.Person == "" && s.Place == "" {
			status = constants.LocationStatusReserve
		}

		si := &models.SI{
			Id:                   s.Id,
			Name:                 s.Name,
			FactoryNumber:        s.FactoryNumber,
			YearOfIssue:          s.YearOfIssue,
			StateRegister:        s.StateRegister,
			MeasurementLimits:    s.MeasurementLimits,
			VerificationDate:     s.VerificationDate,
			NextVerificationDate: s.NextVerificationDate,
			Person:               s.Person,
			Place:                s.Place,
			LastPlace:            s.LastPlace,
		}

		notEqualDeps := len(data) > 0 && len(data[len(data)-1].SI) > 0 && data[len(data)-1].SI[0].Place != s.Place
		if i == 0 || data[len(data)-1].Channel != s.NotificationChannel || notEqualDeps {
			data = append(data, &models.SiReceiving{
				Channel: s.NotificationChannel,
				Status:  status,
				Place:   s.Place,
				SI:      []*models.SI{si},
			})
		} else {
			data[len(data)-1].SI = append(data[len(data)-1].SI, si)
		}
	}

	return data, nil
}

func (r *SIRepo) GetUsed(ctx context.Context, req *models.Period) ([]*models.SiReceiving, error) {
	query := fmt.Sprintf(`SELECT 
			i.id, i.name, i.factory_number, i.year_of_issue, i.state_register, i.measurement_limits, 
			v.date, v.next_date,
			COALESCE(l.person, e.name, '') AS person, 
			COALESCE(l.place, d.name, '') AS place, 
			COALESCE(l.last_place, lp.name, '') AS last_place,
			COALESCE(c.most_channel_id, r.channel, '') AS notification_channel
		FROM %s AS i
		-- Последняя поверка
		LEFT JOIN LATERAL (
			SELECT date, next_date 
			FROM %s 
			WHERE instrument_id = i.id 
			ORDER BY date DESC, created_at DESC 
			LIMIT 1
		) AS v ON TRUE
		-- Последнее местоположение (обязательно для фильтра status=$2)
		LEFT JOIN LATERAL (
			SELECT status, person, place, last_place, person_id, department_id, last_place_id 
			FROM %s 
			WHERE instrument_id = i.id 
			ORDER BY date_of_issue DESC, created_at DESC 
			LIMIT 1
		) AS l ON TRUE
		-- Справочники (джоиним только если ID валиден)
		LEFT JOIN %s AS e ON (l.person_id::text = e.id::text)
		LEFT JOIN %s AS d ON (l.department_id::text = d.id::text)
		LEFT JOIN %s AS lp ON (l.last_place_id::text = lp.id::text)
		-- Каналы уведомлений
		LEFT JOIN %s AS c ON (d.channel_id = c.id)
		LEFT JOIN LATERAL (
			SELECT r.notification_channel AS channel 
			FROM %s AS r 
			INNER JOIN %s AS s ON s.realm_id = r.id 
			WHERE s.id = i.section_id
		) AS r ON TRUE
		WHERE 
			($1 = '' OR i.section_id::text = $1) AND 
			l.status = $2 AND 
			v.next_date BETWEEN $3 AND $4
		ORDER BY d.channel_id, place, last_place, v.next_date`,
		InstrumentsTable, VerificationTable, LocationTable,
		EmployeeTable, DepartmentTable, DepartmentTable,
		ChannelTable, RealmTable, SectionTable,
	)

	tmp := []*pq_models.SI{}

	if err := r.db.SelectContext(ctx, &tmp, query, req.SectionId, constants.LocationStatusUsed, req.StartAt, req.FinishAt); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	data := []*models.SiReceiving{}
	for i, s := range tmp {
		si := &models.SI{
			Id:                   s.Id,
			Name:                 s.Name,
			FactoryNumber:        s.FactoryNumber,
			YearOfIssue:          s.YearOfIssue,
			StateRegister:        s.StateRegister,
			MeasurementLimits:    s.MeasurementLimits,
			VerificationDate:     s.VerificationDate,
			NextVerificationDate: s.NextVerificationDate,
			Person:               s.Person,
			Place:                s.Place,
			LastPlace:            s.LastPlace,
		}

		notEqualDeps := len(data) > 0 && len(data[len(data)-1].SI) > 0 && data[len(data)-1].SI[0].Place != s.Place
		if i == 0 || data[len(data)-1].Channel != s.NotificationChannel || notEqualDeps {
			data = append(data, &models.SiReceiving{
				Channel: s.NotificationChannel,
				Status:  constants.StatusReceiving,
				SI:      []*models.SI{si},
			})
		} else {
			data[len(data)-1].SI = append(data[len(data)-1].SI, si)
		}
	}

	return data, nil
}

func (r *SIRepo) GetLog(ctx context.Context, req *models.Period) ([]*models.SiWithLog, error) {
	query := fmt.Sprintf(`SELECT i.id, i.name, i.date_of_receipt, i.type, i.factory_number, i.responsible, 
			COALESCE(r.repair, '') AS repair, COALESCE(p.preservation, '') AS preservation,
			COALESCE(t.saving, '') AS saving, COALESCE(w.write_off, '') AS write_off
		FROM %s AS i

		LEFT JOIN LATERAL (
			SELECT string_agg(
				concat_ws(' ', 
					CASE 
						WHEN period_end < '1900-01-01'::date OR period_end IS NULL 
						THEN to_char(period_start, 'DD.MM.YYYY') 
						ELSE to_char(period_start, 'DD.MM.YYYY') || '-' || to_char(period_end, 'DD.MM.YYYY') 
					END,
					CASE WHEN work <> '' THEN '(' || work || ')' END
				), ', ' ORDER BY created_at
			) AS repair
			FROM %s 
			WHERE instrument_id = i.id
		) r ON TRUE
		
		LEFT JOIN LATERAL (
			SELECT string_agg(
				'Консервация ' || to_char(date_start, 'DD.MM.YYYY') || 
				CASE 
					WHEN date_end > '1900-01-01'::date 
					THEN ' - Расконсервация ' || to_char(date_end, 'DD.MM.YYYY') 
					ELSE '' 
				END, ', ' ORDER BY created_at
			) AS preservation
			FROM %s 
			WHERE instrument_id = i.id
		) p ON TRUE
		
		LEFT JOIN LATERAL (
			SELECT string_agg(
				'Передано ' || to_char(date_start, 'DD.MM.YYYY') || 
				CASE 
					WHEN date_end > '1900-01-01'::date 
					THEN ' - Возвращено ' || to_char(date_end, 'DD.MM.YYYY') 
					ELSE '' 
				END, ', ' ORDER BY created_at
			) AS saving
			FROM %s 
			WHERE instrument_id = i.id
		) t ON TRUE
		
		LEFT JOIN LATERAL (
			SELECT string_agg(
				'Списан ' || to_char(date, 'DD.MM.YYYY'), ', ' ORDER BY created_at
			) AS write_off
			FROM %s 
			WHERE instrument_id = i.id
		) w ON TRUE		
		WHERE i.section_id::text = $1`,
		InstrumentsTable, RepairTable, PreservationTable, TransferToSaveTable, WriteOffTable,
	)

	data := []*models.SiWithLog{}
	if err := r.db.SelectContext(ctx, &data, query, req.SectionId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}
