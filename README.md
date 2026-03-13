# crud_golang_testing_example_private

Este projeto é um guia para construir uma **API profissional em Golang**, utilizando boas práticas usadas em empresas.

## Stack utilizada

* **Gin** → Framework HTTP
* **database/sql** → Acesso ao banco
* **PostgreSQL**
* **golang-migrate** → Migrations
* **Docker + Docker Compose** → Ambiente de execução
* **.env** → Configurações da aplicação

A arquitetura segue um modelo **layered / clean**, separando responsabilidades entre HTTP, regras de negócio e acesso ao banco.

---

# 1. Criar o Projeto

```bash
mkdir students-api
cd students-api
go mod init github.com/seuuser/students-api
```

---

# 2. Estrutura Inicial de Pastas

```
students-api
│
├── cmd
│   └── api
│       └── main.go
│
├── internal
│   ├── config
│   ├── database
│   ├── handler
│   ├── repository
│   ├── service
│   └── model
│
├── migrations
│
├── docker
│
├── .env
├── docker-compose.yml
├── go.mod
└── go.sum
```

## Função de cada pasta

| Pasta               | Função                               |
| ------------------- | ------------------------------------ |
| cmd/api             | Ponto de entrada da aplicação        |
| internal/config     | Carregamento das variáveis do `.env` |
| internal/database   | Conexão com o banco                  |
| internal/model      | Structs da aplicação                 |
| internal/repository | Acesso ao banco de dados             |
| internal/service    | Regras de negócio                    |
| internal/handler    | Camada HTTP                          |
| migrations          | Arquivos de migration                |
| docker              | Dockerfile                           |

---

# 3. Arquivo `.env`

Criar na raiz do projeto:

```
.env
```

Variáveis típicas:

```
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_SSLMODE=
API_PORT=
```

---

# 4. Configuração da Aplicação

Criar:

```
internal/config/config.go
```

Responsabilidades:

* Carregar o `.env`
* Criar uma struct `Config`
* Centralizar configurações da aplicação

Pacote útil:

```bash
go get github.com/joho/godotenv
```

---

# 5. Conexão com Banco de Dados

Criar pasta:

```
internal/database
```

Arquivo:

```
connection.go
```

Responsabilidades:

* Criar conexão com PostgreSQL
* Retornar `*sql.DB`
* Configurar pool de conexões

Instalar driver.

Opção tradicional:

```bash
go get github.com/lib/pq
```

Opção moderna:

```bash
go get github.com/jackc/pgx/v5/stdlib
```

---

# 6. Migrations

Instalar CLI do migrate.

Mac:

```bash
brew install golang-migrate
```

Ou via Go:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

---

# 7. Criar Migration

```bash
migrate create -ext sql -dir migrations -seq create_students_table
```

Arquivos gerados:

```
migrations/
  000001_create_students_table.up.sql
  000001_create_students_table.down.sql
```

---

# 8. Rodar Migrations

Executar migrations:

```bash
migrate -path migrations -database "postgres://user:pass@localhost:5432/db?sslmode=disable" up
```

Rollback:

```bash
migrate -path migrations -database "postgres://..." down
```

Ver versão:

```bash
migrate -path migrations -database "postgres://..." version
```

---

# 9. Models

Criar pasta:

```
internal/model
```

Arquivo exemplo:

```
student.go
```

Models representam:

* tabelas do banco
* estruturas utilizadas na aplicação

Não deve existir lógica aqui.

---

# 10. Repository (Acesso ao Banco)

Criar pasta:

```
internal/repository
```

Arquivo exemplo:

```
student_repository.go
```

Responsabilidades:

* CRUD
* Queries SQL
* Comunicação com `database/sql`

---

## CreateUser

```go
func (r *Repository) CreateUser(req CreateUserRequest, hashedPassword string) (*User, error) {

	query := `
        INSERT INTO users (name, email, password)
        VALUES ($1, $2, $3)
        RETURNING id, name, email, created_at, updated_at
    `

	user := &User{}

	err := r.db.QueryRow(query, req.Name, req.Email, hashedPassword).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	return user, nil
}
```

---

## FindAll

```go
func (r *Repository) FindAll() ([]User, error){

	query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        ORDER BY id
    `

	rows, err := r.db.Query(query)
	if err != nil{
		return nil, fmt.Errorf("erro ao buscar usuarios: %w", err)
	}
	defer rows.Close()

	var users []User

	for rows.Next(){

		var u User

		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil{
			return nil, fmt.Errorf("Erro ao ler o usuario: %w", err)
		}

		users = append(users, u)
	}

	return users, nil
}
```

---

## FindByID

```go
func (r *Repository) FindByID(id int) (*User, error){

	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &User{}

	err := r.db.QueryRow(query, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil{
		return nil, fmt.Errorf("Erro ao buscar o usuario: %w", err)
	}

	return user, nil
}
```

---

## Update

```go
func (r *Repository) Update(id int, req UpdateUserRequest)(*User, error){

	query := `
		UPDATE users
		SET name = $1, email = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, email, created_at, updated_at
	`

	user := &User{}

	err := r.db.QueryRow(query, req.Name, req.Email, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows{
		return nil, nil
	}

	if err != nil{
		return nil, fmt.Errorf("Erro ao atualizar o usuario: %w", err)
	}

	return user, nil
}
```

---

## Delete

```go
func (r *Repository) Delete(id int) error {

	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
```

---

# 11. Service (Regras de Negócio)

Criar pasta:

```
internal/service
```

Arquivo:

```
student_service.go
```

Responsabilidades:

* Validações
* Regras de negócio
* Orquestração entre camadas

Fluxo:

```
handler → service → repository
```

A camada **service não deve conhecer HTTP**.

---

# 12. Handler (Camada HTTP)

Criar pasta:

```
internal/handler
```

Arquivo exemplo:

```
student_handler.go
```

Responsabilidades:

* Receber requests
* Validar JSON
* Chamar services
* Retornar responses

Aqui é onde o **Gin** é utilizado.

---

# 13. Router

Criar arquivo:

```
internal/handler/router.go
```

Responsável por registrar rotas da API.

Exemplo de endpoints:

```
POST   /students
GET    /students
GET    /students/:id
PUT    /students/:id
DELETE /students/:id
```

---

# 14. Main da Aplicação

Arquivo:

```
cmd/api/main.go
```

Responsabilidades:

1. Carregar configuração
2. Conectar no banco
3. Inicializar repositories
4. Inicializar services
5. Inicializar handlers
6. Configurar router
7. Subir servidor HTTP

Fluxo de inicialização:

```
config
 ↓
database
 ↓
repository
 ↓
service
 ↓
handler
 ↓
router
 ↓
start server
```

---

# 15. Dockerfile

Criar pasta:

```
docker
```

Arquivo:

```
Dockerfile
```

Responsabilidades:

* Buildar binário Go
* Executar aplicação dentro do container

---

# 16. Docker Compose

Arquivo:

```
docker-compose.yml
```

Serviços:

```
api
postgres
```

Usar volumes para persistência do banco.

---

# 17. Rodar o Projeto

Fluxo comum em ambiente de desenvolvimento.

### Subir banco

```bash
docker compose up -d postgres
```

### Rodar migrations

```bash
migrate up
```

### Subir API

```bash
docker compose up api
```

---

# Estrutura Final do Projeto

```
students-api
│
├── cmd
│   └── api
│       └── main.go
│
├── internal
│   ├── config
│   │   └── config.go
│   │
│   ├── database
│   │   └── connection.go
│   │
│   ├── model
│   │   └── student.go
│   │
│   ├── repository
│   │   └── student_repository.go
│   │
│   ├── service
│   │   └── student_service.go
│   │
│   └── handler
│       ├── student_handler.go
│       └── router.go
│
├── migrations
├── docker
│   └── Dockerfile
├── .env
├── docker-compose.yml
└── go.mod
```

---

# Princípio Importante

Nunca misturar responsabilidades.

Separação correta:

```
Handler → Service → Repository
```

* **Handler** → HTTP
* **Service** → regras de negócio
* **Repository** → banco de dados

Seguindo essa estrutura, a API fica **organizada, escalável e alinhada com projetos profissionais em Go**.
