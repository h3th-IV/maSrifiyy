package db

var (
	createSellerTableStmt = `
		CREATE TABLE IF NOT EXISTS sellers (
			id SERIAL PRIMARY KEY,
			user_id VARCHAR(255) NOT NULL,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

	createGoodsTableStmt = `
	CREATE TABLE IF NOT EXISTS goods (
		id SERIAL PRIMARY KEY,
		product_id VARCHAR(255) NOT NULL,
		name VARCHAR(100) NOT NULL,
		quantity INT NOT NULL,
		max_threshold INT NOT NULL,
		min_threshold INT NOT NULL,
		created_by INT REFERENCES sellers(id) ON DELETE CASCADE
	);`
	dropSellerTableStmt = `DROP TABLE IF EXISTS sellers;`
	dropGoodsTableStmt  = `DROP TABLE IF EXISTS goods;`
	createSellerAcct    = `INSERT INTO sellers (user_id, first_name, last_name, email, password, created_at)
	VALUES ($1, $2, $3, $4, $5, NOW())`
	updateSellerAcct = `UPDATE sellers SET first_name = $1, last_name = $2, email = $3, password = $4
	WHERE id = $5`
	getSellerAcctbyEmail  = `SELECT id, user_id, first_name, last_name, email, password, created_at FROM sellers WHERE email = $1`
	getSellerAcctbyID     = `SELECT id, user_id, first_name, last_name, email, password, created_at FROM sellers WHERE id = $1`
	getSellerAcctbyUserId = `SELECT id, user_id, first_name, last_name, email, password, created_at FROM sellers WHERE user_id = $1`
	createItemInInventory = `
	INSERT INTO goods (product_id, name, quantity, max_threshold, min_threshold, created_by)
	VALUES ($1, $2, $3, $4, $5, $6)`
	getItembyProductID    = `SELECT id, product_id, name, quantity, max_threshold, min_threshold, created_by FROM goods WHERE product_id = $1`
	getAllItem            = `SELECT * from goods`
	updateItemInInventory = `UPDATE goods SET name = $1, quantity = $2, max_threshold = $3
	WHERE id = $4 AND created_by = $5`
	setMaxThreshold     = `UPDATE goods SET max_threshold = $1 WHERE product_id = $2 AND created_by = $3`
	getItembyID         = `SELECT id, product_id, name, quantity, max_threshold, min_threshold, created_by FROM goods WHERE id = $1`
	addItem             = `UPDATE goods SET quantity = quantity + $1 WHERE product_id = $2 AND quantity + $1 <= max_threshold`
	incrementItem       = `UPDATE goods SET quantity = quantity + $1 WHERE product_id = $2 AND quantity + $1 <= max_threshold`
	decrementItem       = `UPDATE goods SET quantity = quantity - $1 WHERE product_id = $2 AND quantity - $1 >= min_threshold`
	getLowStockPorducts = `
	SELECT g.id, g.product_id, g.name, g.quantity, g.max_threshold, g.min_threshold, g.created_by, s.first_name, s.email FROM goods g INNER JOIN sellers s ON g.created_by = s.id WHERE g.quantity <= g.min_threshold OR g.quantity <= g.min_threshold + 0.1 * g.min_threshold;`
)
