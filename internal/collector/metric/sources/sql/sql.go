package sql

import "database/sql"

func RowsToMaps(rows *sql.Rows) ([]map[string]any, error) {
	var results []map[string]any

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		data := make(map[string]any)
		columns := make([]string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		for i, colName := range cols {
			data[colName] = columns[i]
		}

		results = append(results, data)
	}

	return results, nil
}
