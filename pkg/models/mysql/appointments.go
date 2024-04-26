package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	"aitunews.kz/snippetbox/pkg/models"
)

type AppointmentModel struct {
	DB *sql.DB
}

func (m *AppointmentModel) Insert(user_id int, service_id, time string) (int, error) {
	stmt := `INSERT INTO appointments (user_id, service_id, time) VALUES(?, ?, ?)`
	result, err := m.DB.Exec(stmt, user_id, service_id, time)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *AppointmentModel) Update(id int, time string) error {
	stmt := `UPDATE appointments SET time = ? WHERE id = ?`
	result, err := m.DB.Exec(stmt, time, id)

	if err != nil || result == nil {
		return err
	}
	return nil
}

func (m *AppointmentModel) Delete(id int) error {
	stmt := `DELETE FROM appointments WHERE id = ?`
	result, err := m.DB.Exec(stmt, id)

	if err != nil || result == nil {
		return err
	}
	return nil
}

func (m *AppointmentModel) Get(id int) (*models.Appointment, error) {

	stmt := `SELECT * FROM appointments`

	row := m.DB.QueryRow(stmt, id)
	s := &models.Appointment{}
	err := row.Scan(&s.ID, &s.User_id, &s.Service_id, &s.Time)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// func (m *AppointmentModel) Latest(name_for string) ([]*models.Appointment, error) {
// 	stmt := `SELECT * FROM appointments`
// 	if name_for == "" {
// 		stmt = `SELECT * FROM appointments`
// 	}
// 	rows, err := m.DB.Query(stmt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	appointments := []*models.Appointment{}
// 	for rows.Next() {
// 		s := &models.Appointment{}
// 		err = rows.Scan(&s.ID, &s.User_id, &s.Service_id, &s.Time)
// 		if err != nil {
// 			return nil, err
// 		}
// 		appointments = append(appointments, s)
// 	}
// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return appointments, nil
// }

func (m *AppointmentModel) Latest(name_for, limit, offset int) ([]*models.Appointment, error) {

	stmt := fmt.Sprintf(`SELECT * FROM appointments WHERE user_id = %d`, name_for)

	if name_for == 5 {
		stmt = `SELECT * FROM appointments`
	}
	stmt += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	appointments := []*models.Appointment{}
	for rows.Next() {
		s := &models.Appointment{}
		err = rows.Scan(&s.ID, &s.User_id, &s.Service_id, &s.Time)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}
