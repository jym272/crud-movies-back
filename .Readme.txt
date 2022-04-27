Configurar el secreto, ejecutar el comando:(antes del build hacerlo)
> export JWT_SECRET_KEY='bIFBlgmpaSx4QuvlKJyJBcBy+9wQcOFanKqaCm8MX5E='

Para empezar el server, ejecutar el comando:
> go run ./cmd/api/.


Hacer el build, generar el binario, ejecutar el comando:
> env GOOS=linux GOARCH=amd64 go build -o gomovies ./cmd/api

Para cargar la db, con datos:
->crear una db "movies"
->$ psql -d movies -f 040\ go_movies.sql

Para exportar la db movies a un archivo sql:
->$ pg_dump --no-owner movies > gm.sql

Para importar de un archivo sql a la db movies:
->$ sudo -u postgres psql -d movies -f gm.sql

Agregar password al user postgres:
Login into the psql:
->$ sudo -u postgres psql
Then in the psql console change the password and quit:
postgres=# \password postgres
Enter new password: <new-password>
postgres=# \q