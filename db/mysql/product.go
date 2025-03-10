package mysql

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func (db *MysqlDB) GetProducts(limit, offset int, name string) ([]*Product, error) {
	query := "SELECT id, name, price, quantity FROM products WHERE name LIKE ? LIMIT ? OFFSET ?"
	rows, err := db.DB.Query(query, "%"+name+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}
