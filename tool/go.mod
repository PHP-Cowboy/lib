module za.game/za.lib/tool

go 1.21.0

require gorm.io/gorm v1.25.5

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
)

replace (
	za.game/lib/rds => ../rds
	za.game/lib/sqlInfo => ../sqlInfo
	za.game/za.log => ../../za.log
)
