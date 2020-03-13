package migrate

import "database/sql"

var migrations = []Migrations{
	{
		name:     "create_dummy_table",
		function: changeDummyTable,
	},
}

func changeDummyTable(db *sql.DB, name string) error {
	sql := "CREATE TABLE IF NOT EXISTS `dummy` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
		"`created_at` timestamp NULL DEFAULT NULL," +
		"`updated_at` timestamp NULL DEFAULT NULL," +
		"`deleted_at` timestamp NULL DEFAULT NULL," +
		"`text` varchar(255) DEFAULT NULL," +
		"PRIMARY KEY (`id`)" +
		")"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}
