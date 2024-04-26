package mysql

import (
	"database/sql"
	"errors"

	"aitunews.kz/snippetbox/pkg/models"
)

type ServiceModel struct {
	DB *sql.DB
}

func (m *ServiceModel) Insert(title, content, master string, price int) (int, error) {
	stmt := `INSERT INTO services (title, content, master, price) VALUES(?, ?, ?, ?)`
	result, err := m.DB.Exec(stmt, title, content, master, price)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *ServiceModel) Update(title string, price int) error {
	stmt := `UPDATE services SET price = ? WHERE title = ?`
	result, err := m.DB.Exec(stmt, price, title)

	if err != nil || result == nil {
		return err
	}
	return nil
}

func (m *ServiceModel) Delete(title string) error {
	stmt := `DELETE FROM services WHERE title = ?`
	result, err := m.DB.Exec(stmt, title)

	if err != nil || result == nil {
		return err
	}
	return nil
}

func (m *ServiceModel) Get(id int) (*models.Service, error) {

	stmt := `SELECT * FROM services`

	row := m.DB.QueryRow(stmt, id)
	s := &models.Service{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Master, &s.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *ServiceModel) Latest(name_for, sort, sort_type string) ([]*models.Service, error) {
	stmt := `SELECT * FROM services`
	if name_for == "" {
		stmt = `SELECT * FROM services`
	}
	if sort != "" {
		stmt += " ORDER BY " + sort
		if sort_type != "" {
			stmt += " " + sort_type
		}
	}

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	services := []*models.Service{}
	for rows.Next() {
		s := &models.Service{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Master, &s.Price)
		if err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return services, nil
}
