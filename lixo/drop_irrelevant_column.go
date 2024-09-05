// package migrations

// import (
//      "gorm.io/gorm"
// )

// func DropIrrelevantColumn(db *gorm.DB) error {
//     // Defina o nome da tabela e da coluna
//     tableName := "companies"
//     columnName := "name"

//     // Remova a coluna diretamente
//     return db.Migrator().DropColumn(tableName, columnName)
// }