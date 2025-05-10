module bank-api

go 1.23.0

require (
	github.com/beevik/etree v1.3.0
	github.com/go-mail/mail v2.3.1+incompatible
	github.com/golang-jwt/jwt/v5 v5.2.1 // JWT
	github.com/gorilla/mux v1.8.1
	github.com/lib/pq v1.10.9
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/crypto v0.38.0 // для bcrypt
)

require (
	golang.org/x/sys v0.33.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/mail.v2 v2.3.1 // indirect
)
