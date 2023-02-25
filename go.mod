module src.userspace.com.au/sws

require (
	github.com/blend/go-sdk v2.0.0+incompatible // indirect
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-chi/chi v4.0.3+incompatible
	github.com/go-chi/jwtauth v4.0.4+incompatible
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/hashicorp/golang-lru v0.5.4
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.3.0 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/mssola/user_agent v0.5.1
	github.com/oschwald/maxminddb-golang v1.6.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/shopspring/decimal v0.0.0-20200105231215-408a2507e114 // indirect
	github.com/speps/go-hashids v2.0.0+incompatible
	github.com/wcharczuk/go-chart v2.0.1+incompatible
	golang.org/x/crypto v0.1.0
	golang.org/x/image v0.0.0-20200119044424-58c23975cae1 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	src.userspace.com.au/go-migrate v0.0.0-20200208102934-cf11cf76db3f
	src.userspace.com.au/templates v0.0.0-20200320100635-3f0c3fbfaf31
)

//replace src.userspace.com.au/templates => ../templates

go 1.14
