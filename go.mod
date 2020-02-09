module src.userspace.com.au/sws

require (
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/go-chi/chi v3.3.3+incompatible
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/jackc/pgx v3.3.0+incompatible
	github.com/jmoiron/sqlx v1.2.0
	github.com/kr/pretty v0.2.0 // indirect
	github.com/lib/pq v1.3.0 // indirect
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/pkg/errors v0.8.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/shopspring/decimal v0.0.0-20200105231215-408a2507e114 // indirect
	github.com/speps/go-hashids v2.0.0+incompatible
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	src.userspace.com.au/flags v0.0.0-20200208094111-eef94fa594cc
	src.userspace.com.au/go-migrate v0.0.0-20200208102934-cf11cf76db3f
)

go 1.13

//replace src.userspace.com.au/flags => ../flags
