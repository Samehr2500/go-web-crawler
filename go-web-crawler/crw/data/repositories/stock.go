package repositories

import (
	driver "crawler/data/driver"
	"crawler/data/models"
	"github.com/google/uuid"
	"time"
)

// Insert one stock in the DB
func Insert(stock models.Stock) models.Stock {

	// create the db connection
	db := driver.CreateDBConnection()

	// close the db connection
	defer db.Close()

	// Generate a new UUID
	id := uuid.New()
	// create the insert sql query
	// returning the id of the inserted stock
	sqlStatement := `INSERT INTO tbl_stocks (id, created_at, paper_name, company_name, daily_rate, market_value) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	// Set the desired timezone
	tz, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		panic(err)
	}

	// Get the current time in the desired timezone
	now := time.Now().In(tz)
	// execute the sql statement
	// Scan function will save the insert id in the id
	err = db.QueryRow(sqlStatement, id, now, stock.PaperName, stock.CompanyName, stock.DailyRate, stock.MarketValue).Scan(&stock.Id)

	if err != nil {
		panic(err)
	}
	return stock
}

// Get one stock from the DB by its id
func Get(id uuid.UUID) models.Stock {
	// create the db connection
	db := driver.CreateDBConnection()

	// close the db connection
	defer db.Close()

	// create a stock of models.Stock type
	var stock models.Stock

	// create the select sql query
	sqlStatement := `SELECT * FROM tbl_stocks WHERE id=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to stock
	err := row.Scan(
		&stock.Id, &stock.CreatedAt, &stock.PaperName, &stock.CompanyName, &stock.DailyRate, &stock.MarketValue)

	if err != nil {
		panic(err)
	}
	return stock
}

// GetAllIds stock's id from the DB
func GetAllIds() []uuid.UUID {
	// create the db connection
	db := driver.CreateDBConnection()

	// close the db connection
	defer db.Close()

	var ids []uuid.UUID

	// create the select sql query
	sqlStatement := `SELECT id FROM tbl_stocks`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		panic(err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var id uuid.UUID

		// unmarshal the row object to id
		err = rows.Scan(&id)

		if err != nil {
			panic(err)
		}

		// append the id in the ids slice
		ids = append(ids, id)

	}

	return ids
}

// GetAll stocks from the DB
func GetAll() []models.Stock {
	// create the db connection
	db := driver.CreateDBConnection()

	// close the db connection
	defer db.Close()

	var stocks []models.Stock

	// create the select sql query
	sqlStatement := `SELECT id FROM tbl_stocks`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		panic(err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var stock models.Stock

		// unmarshal the row object to stock
		err = rows.Scan(&stock.Id, &stock.CreatedAt, &stock.PaperName, &stock.CompanyName, &stock.DailyRate, &stock.MarketValue)

		if err != nil {
			panic(err)
		}

		// append the stock in the stocks slice
		stocks = append(stocks, stock)

	}

	return stocks
}

// Update stock in the DB
func Update(stock models.Stock) models.Stock {

	// create the db connection
	db := driver.CreateDBConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE tbl_stocks SET paper_name=$2, company_name=$3, daily_rate=$4, market_value=$5 WHERE id=$1`

	// execute the sql statement
	_, err := db.Exec(sqlStatement, stock.Id, stock.PaperName, stock.CompanyName, stock.DailyRate, stock.MarketValue)

	if err != nil {
		panic(err)
	}

	return stock
}

// Delete stock in the DB
func Delete(id int64) {

	// create the db connection
	db := driver.CreateDBConnection()

	// close the data connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM tbl_stocks WHERE id=$1`

	// execute the sql statement
	_, err := db.Exec(sqlStatement, id)

	if err != nil {
		panic(err)
	}
}
