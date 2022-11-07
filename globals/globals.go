package globals

import (
	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DbConn *gorm.DB

//----------GORM MODELS----------

//[TODO] Foriegn keys dont seem to be applied in SQL tables.

type Admins struct {
	// gorm.Model
	Adm_id        string `gorm:"primaryKey"`
	Adm_hash_pass string
}

// We have a constraint that a given Doc can only have one admin.
type Doc struct {
	// gorm.Model
	Doc_id      string `gorm:"primaryKey"`
	Def_permbit int8   `gorm:"type:int(1)"`
	Adm_id      string `gorm:"foreignKey:Adm_id"`
}

type User_perms struct {
	// gorm.Model
	Doc_id     string `gorm:"foreignKey:Doc_id"`
	User_id    string
	Nd_permbit int8 `gorm:"type:int(1)"`
}

type Auth_code struct {
	// gorm.Model
	Adm_id    string `gorm:"ForeignKey:Adm_id"`
	Auth_code string `gorm:"primaryKey"`
}
