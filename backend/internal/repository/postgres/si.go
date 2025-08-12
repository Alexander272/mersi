package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres/pq_models"
	"github.com/jmoiron/sqlx"
)

type SIRepo struct {
	db *sqlx.DB
}

func NewSIRepo(db *sqlx.DB) *SIRepo {
	return &SIRepo{
		db: db,
	}
}

type SI interface {
	Get(ctx context.Context, req *models.GetSiDTO) ([]*models.SI, error)
	GetVerification(ctx context.Context, req *models.Period) ([]*models.SiVerification, error)
	GetSent(ctx context.Context, req *models.GetSiDTO) ([]*models.SiReceiving, error)
	GetUsed(ctx context.Context, req *models.Period) ([]*models.SiReceiving, error)
}

func (r *SIRepo) formatField(field string) string {
	format := make(map[string]string)

	format["position"] = "position"
	format["name"] = "name"
	format["dateOfReceipt"] = "date_of_receipt"
	format["type"] = "type"
	format["factoryNumber"] = "factory_number"
	format["measurementLimits"] = "measurement_limits"
	format["accuracy"] = "accuracy"
	format["stateRegister"] = "state_register"
	format["countryOfProduce"] = "country_of_produce"
	format["manufacturer"] = "manufacturer"
	format["responsible"] = "responsible"
	format["inventory"] = "inventory"
	format["yearOfIssue"] = "year_of_issue"
	format["interVerificationInterval"] = "inter_verification_interval"
	format["actOfEntering"] = "act_of_entering"
	format["actOfEnteringId"] = "act_of_entering_id"
	format["notes"] = "notes"
	format["verificationDate"] = "date"
	format["nextVerificationDate"] = "next_date"
	format["department"] = "department_id"
	format["place"] = "place"
	format["person"] = "person"
	format["status"] = "l.status"

	return format[field]
}

func (r *SIRepo) Get(ctx context.Context, req *models.GetSiDTO) ([]*models.SI, error) {
	data := []*models.SI{}
	params := []interface{}{req.SectionId, req.Status}
	count := len(params) + 1

	order := " ORDER BY "
	for _, s := range req.Sort {
		order += fmt.Sprintf("%s %s, ", r.formatField(s.Field), s.Type)
	}
	if len(req.Sort) == 0 {
		order += "position, "
	}
	order += "created_at, id"

	filter := ""
	if len(req.Filters) > 0 {
		filter += " AND "
		filters := []string{}

		for _, ns := range req.Filters {
			if ns.Field == "department" {
				filters = append(filters, "("+
					getFilterLine(ns.Values[0].CompareType, r.formatField(ns.Field), count)+" OR ("+
					getFilterLine(ns.Values[0].CompareType, "last_place_id", count)+
					" AND "+r.formatField("status")+"='moved'))",
				)
				// filter += " AND (" + getFilterLine(ns.Values[0].CompareType, r.formatField(ns.Field), count) + " OR (" +
				// 	getFilterLine(ns.Values[0].CompareType, "last_place_id", count) + "AND m.status='moved'))"

				ns.Values[0].Value = strings.ReplaceAll(ns.Values[0].Value, ",", "|")
				params = append(params, ns.Values[0].Value)
				count++
				continue
			}
			for _, sv := range ns.Values {
				filters = append(filters, getFilterLine(sv.CompareType, r.formatField(ns.Field), count))
				if sv.CompareType == "in" {
					sv.Value = strings.ReplaceAll(sv.Value, ",", "|")
				}
				params = append(params, sv.Value)
				count++
			}
		}
		filter += strings.Join(filters, " AND ")
	}

	search := ""
	if req.Search != nil {
		search = " AND ("

		list := []string{}
		for _, f := range req.Search.Fields {
			list = append(list, fmt.Sprintf("%s ILIKE '%%'||$%d||'%%'", r.formatField(f), count))
		}
		params = append(params, req.Search.Value)
		count++
		search += strings.Join(list, " OR ") + ")"
	}

	params = append(params, req.Page.Limit, req.Page.Offset)

	query := fmt.Sprintf(`SELECT i.id, position, name, date_of_receipt, type, factory_number, measurement_limits, accuracy, state_register, 
		COALESCE(l.status, 'used') AS status, country_of_produce, manufacturer, responsible, inventory, year_of_issue, inter_verification_interval, 
		act_of_entering, act_of_entering_id, notes,
		v.date, v.next_date, COALESCE(cert, '') AS certificate, COALESCE(cert_id, '') AS certificate_id, COALESCE(repair, '') AS repair,
		COALESCE(p.date_start, 0) AS preservation, COALESCE(p.date_end, 0) AS de_preservation,
		COALESCE(ts.date_start, 0) AS transfer_date, COALESCE(ts.date_end, 0) AS return_date, 
		COALESCE(td.doc_name, '') AS transfer_to_dep, COALESCE(wo.doc_name, '') AS write_off,
		COALESCE(person, e.emp, '') AS person, COALESCE(place, dep.dep, '') AS place, COALESCE(l.last_place, '') AS last_place,
		COUNT(*) OVER() AS total
		FROM %s AS i
		LEFT JOIN LATERAL (SELECT id, date, next_date FROM %s WHERE instrument_id=i.id ORDER BY date DESC, created_at DESC LIMIT 1) AS v ON TRUE
		LEFT JOIN LATERAL (SELECT name AS cert, doc_id::text AS cert_id FROM %s WHERE verification_id=v.id) AS d ON TRUE
		LEFT JOIN LATERAL (SELECT date_part('year', to_timestamp(period_end)) || ' (' || work || ')' AS repair FROM %s 
			WHERE instrument_id=i.id ORDER BY period_end DESC LIMIT 1) AS r ON TRUE
		LEFT JOIN LATERAL (SELECT date_start, date_end FROM %s WHERE instrument_id=i.id ORDER BY date_start DESC LIMIT 1) AS p ON TRUE
		LEFT JOIN LATERAL (SELECT date_start, date_end FROM %s WHERE instrument_id=i.id ORDER BY date_start DESC LIMIT 1) AS ts ON TRUE
		LEFT JOIN LATERAL (SELECT doc_name FROM %s WHERE instrument_id=i.id ORDER BY date DESC LIMIT 1) AS td ON TRUE
		LEFT JOIN LATERAL (SELECT doc_name FROM %s WHERE instrument_id=i.id ORDER BY date DESC LIMIT 1) AS wo ON TRUE
		LEFT JOIN LATERAL (SELECT (CASE WHEN status='%s' THEN place WHEN status='%s' THEN 'Резерв' ELSE
			(CASE WHEN last_place!='' OR last_place_id IS NOT NULL 
				THEN 'Перемещение из «'||COALESCE(lp.name,last_place)||'»' ELSE 'Перемещение' END) END) AS place, last_place, status,
			person, person_id, department_id, last_place_id FROM %s AS l
			LEFT JOIN LATERAL (SELECT name FROM %s WHERE l.last_place_id::uuid=id) AS lp ON true
			WHERE instrument_id=i.id ORDER BY date_of_issue DESC, created_at DESC LIMIT 1) AS l ON TRUE
		LEFT JOIN LATERAL (SELECT name as emp FROM %s WHERE l.person_id::uuid=id) AS e ON true
		LEFT JOIN LATERAL (SELECT name as dep FROM %s WHERE l.department_id::uuid=id) AS dep ON true
		
		WHERE section_id=$1 AND i.status=$2 %s%s%s LIMIT $%d OFFSET $%d`,
		InstrumentsTable, VerificationTable, VerificationDocsTable, RepairTable, PreservationTable,
		TransferToSaveTable, TransferToDepTable, WriteOffTable,
		constants.LocationStatusUsed, constants.LocationStatusReserve, LocationTable, DepartmentTable, EmployeeTable, DepartmentTable,
		filter, search, order, count, count+1,
	)
	// logger.Debug("get si", logger.StringAttr("query", query))

	if err := r.db.SelectContext(ctx, &data, query, params...); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *SIRepo) GetVerification(ctx context.Context, req *models.Period) ([]*models.SiVerification, error) {
	query := fmt.Sprintf(`SELECT i.id, i.name, type, factory_number, year_of_issue, state_register, measurement_limits, date, next_date, 
		inter_verification_interval, manufacturer, notes, notification_channel, bid_type
		FROM %s AS i
		LEFT JOIN %s AS s ON s.id=section_id
		LEFT JOIN %s AS r ON r.id=realm_id
		LEFT JOIN LATERAL (SELECT date, next_date FROM %s WHERE instrument_id=i.id ORDER BY date DESC, created_at DESC LIMIT 1) AS v ON TRUE
		LEFT JOIN LATERAL (SELECT date AS write_off FROM %s WHERE instrument_id=i.id) AS w ON TRUE
		LEFT JOIN LATERAL (SELECT date_start AS preservation FROM %s WHERE instrument_id=i.id AND date_end=0) AS p ON TRUE 
		LEFT JOIN LATERAL (SELECT date_start AS transferred FROM %s WHERE instrument_id=i.id AND date_end=0) AS t ON TRUE
 		WHERE next_date>=$1 AND next_date<=$2 AND (CASE WHEN $3!='' THEN section_id::text=$3 ELSE true END) AND 
		is_active=true AND expiration_notice=true AND notification_channel!='' AND 
		deleted IS NULL AND write_off IS NULL AND preservation IS NULL AND transferred IS NULL
		ORDER BY notification_channel, i.position`,
		InstrumentsTable, SectionTable, RealmTable, VerificationTable, WriteOffTable, PreservationTable, TransferToSaveTable,
	)

	tmp := []*pq_models.SI{}
	if err := r.db.SelectContext(ctx, &tmp, query, req.StartAt, req.FinishAt, req.SectionId); err != nil {
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
	params := []interface{}{req.SectionId, constants.LocationStatusMoved}
	filter := ""
	count := 2

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
		LEFT JOIN LATERAL (SELECT status, person, place, last_place, person_id, department_id, last_place_id FROM %s WHERE instrument_id=i.id 
			ORDER BY date_of_issue DESC, created_at DESC LIMIT 1) AS l ON TRUE
		LEFT JOIN LATERAL (SELECT name AS emp FROM %s WHERE l.person_id::uuid=id) AS e ON true
		LEFT JOIN LATERAL (SELECT name AS dep, channel_id FROM %s WHERE l.department_id::uuid=id) AS d ON true
		LEFT JOIN LATERAL (SELECT name FROM %s WHERE l.last_place_id::uuid=id) AS lp ON true
		LEFT JOIN LATERAL (SELECT most_channel_id FROM %s WHERE id=d.channel_id) AS c ON TRUE
		LEFT JOIN LATERAL (SELECT notification_channel AS channel FROM %s AS r INNER JOIN %s AS s ON s.realm_id=r.id WHERE s.id=$1) AS r ON TRUE
		WHERE section_id=$1 AND l.status=$2 %s
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

func (r *SIRepo) GetUsed(ctx context.Context, req *models.Period) ([]*models.SiReceiving, error) {
	query := fmt.Sprintf(`SELECT i.id, i.name, factory_number, year_of_issue, state_register, measurement_limits, date, next_date,
		COALESCE(person, e.emp, '') AS person, COALESCE(place, d.dep, '') AS place, COALESCE(l.last_place, lp.name, '') AS last_place,
		COALESCE(most_channel_id, channel, '') AS notification_channel
		FROM %s AS i
		LEFT JOIN LATERAL (SELECT date, next_date FROM %s WHERE instrument_id=i.id ORDER BY date DESC, created_at DESC LIMIT 1) AS v ON TRUE
		LEFT JOIN LATERAL (SELECT status, person, place, last_place, person_id, department_id, last_place_id FROM %s WHERE instrument_id=i.id 
			ORDER BY date_of_issue DESC, created_at DESC LIMIT 1) AS l ON TRUE
		LEFT JOIN LATERAL (SELECT name AS emp FROM %s WHERE l.person_id::uuid=id) AS e ON true
		LEFT JOIN LATERAL (SELECT name AS dep, channel_id FROM %s WHERE l.department_id::uuid=id) AS d ON true
		LEFT JOIN LATERAL (SELECT name FROM %s WHERE l.last_place_id::uuid=id) AS lp ON true
		LEFT JOIN LATERAL (SELECT most_channel_id FROM %s WHERE id=d.channel_id) AS c ON TRUE
		LEFT JOIN LATERAL (SELECT notification_channel AS channel FROM %s AS r INNER JOIN %s AS s ON s.realm_id=r.id WHERE s.id=$1) AS r ON TRUE
		WHERE section_id=$1 AND l.status=$2 AND next_date>=$3 AND next_date<=$4
		ORDER BY channel_id, place, last_place, next_date`,
		InstrumentsTable, VerificationTable, LocationTable, EmployeeTable, DepartmentTable, DepartmentTable,
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
