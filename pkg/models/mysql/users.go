package mysql

import (
	"database/sql"
	"errors" // New import
	"fmt"
	"strings" // New import

	"aitunews.kz/snippetbox/pkg/models"
	"github.com/go-sql-driver/mysql" // New import
	"golang.org/x/crypto/bcrypt"     // New import
	// "gopkg.in/gomail.v2"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(full_name, email, phone, password, role string) error {

	// me := gomail.NewMessage()

	// me.SetHeader("From", "daniyar.0586@gmail.com")
	// me.SetHeader("To", "daniyar.0586@gmail.com", "daniarm146@gmail.com")
	// me.SetAddressHeader("Cc", "daniarm146@gmail.com", "Daniyar")
	// me.SetHeader("Subject", "Hello!")

	// me.SetBody("text/html", "Hello!<br>Follow this link to activate your account!")

	// d := gomail.NewDialer("smtp.gmail.com", 587, "daniyar.0586@gmail.com", "imzhanetta.121911019505")

	// if err := d.DialAndSend(me); err != nil {
	// 	panic(err)
	// }

	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (full_name, email, phone, hashed_password, created, role)
	VALUES(?, ?, ?, ?, UTC_TIMESTAMP(), ?)`

	_, err = m.DB.Exec(stmt, full_name, email, phone, string(hashedPassword), role)
	if err != nil {

		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, string, error) {
	// Retrieve the id and hashed password associated with the given email. If no
	// matching email exists, or the user is not active, we return the
	// ErrInvalidCredentials error.
	var id int
	var role string
	var hashedPassword []byte
	stmt := "SELECT id, role, hashed_password FROM users WHERE email = ? AND active = TRUE"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &role, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", models.ErrInvalidCredentials
		} else {
			return 0, "", err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, "", models.ErrInvalidCredentials
		} else {
			return 0, "", err
		}
	}
	// Otherwise, the password is correct. Return the user ID.
	return id, role, nil
}

// We'll use the Get method to fetch details for a specific user based
// on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}

func (m *UserModel) Get_confirmation_info(idd string) string {

	var c_code string
	stmt := fmt.Sprintf(`SELECT c_code FROM services WHERE id= %s`, idd)

	row := m.DB.QueryRow(stmt)
	err := row.Scan(&c_code)

	if err != nil {
		return ""
	}
	return c_code
}

func (m *UserModel) Update(idd string) error {
	stmt := `UPDATE users SET c_code = 0, confirmation = true WHERE id = ?`
	result, err := m.DB.Exec(stmt, idd)

	if err != nil || result == nil {
		return err
	}
	return nil
}
