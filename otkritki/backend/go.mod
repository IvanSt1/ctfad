module otkritki

go 1.20

require (
    github.com/gorilla/mux v1.8.0
    github.com/gorilla/csrf v1.7.3
    github.com/gorilla/sessions v1.2.1
    github.com/gorilla/schema v1.4.1       // последняя стабильная версия :contentReference[oaicite:0]{index=0}
    golang.org/x/time v0.7.0
    golang.org/x/crypto v0.37.0           // bcrypt и прочие пакеты :contentReference[oaicite:1]{index=1}
    gorm.io/driver/mysql v1.2.3
    gorm.io/gorm v1.21.0
)