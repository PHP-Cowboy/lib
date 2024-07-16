module za.game/za.lib/sqlInfo

go 1.21.0

replace (
	za.game/lib/consts => ../consts
	za.game/lib/dbconn => ../dbconn
	za.game/lib/rds => ../rds
	za.game/lib/response => ../response
	za.game/lib/sqlInfo => ../sqlInfo
	za.game/lib/tool => ../tool
	za.game/za.log => ../../za.log

)

require (
	git.dev666.cc/external/breezedup v1.0.2
	github.com/jmoiron/sqlx v1.4.0
	gorm.io/gorm v1.25.5
	za.game/lib/rds v0.0.0-00010101000000-000000000000
	za.game/lib/sqlInfo v0.0.0-00010101000000-000000000000
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.15.5 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/gomodule/redigo v1.8.9 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	gorm.io/driver/mysql v1.5.2 // indirect
	za.game/za.log v0.0.0-00010101000000-000000000000 // indirect
)
