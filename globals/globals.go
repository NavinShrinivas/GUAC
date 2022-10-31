package globals

import(
  // "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

var DbConn *gorm.DB;



//----------GORM MODELS----------

type Admins struct{
  // gorm.Model
  Adm_id string  `gorm:"primaryKey"`
  Adm_hash_pass string
}

type Doc struct{
  // gorm.Model
  Doc_id string `gorm:"primaryKey"`
  Def_permbit int `gorm:"type:binary(8)"`
  Adm_id string `gorm:"ForeignKey:Adm_id"`
}

type User_perms struct{
  // gorm.Model 
  Doc_id string  `gorm:"ForeignKey:Doc_id"`
  User_id string `gorm:"primaryKey"`
  Nd_permbit int `gorm:"type:binary(8)"`
}

type Auth_code struct{
  // gorm.Model
  Adm_id string `gorm:"ForeignKey:Adm_id"`
  Auth_code string
}

