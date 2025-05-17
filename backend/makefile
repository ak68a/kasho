c_m: 
	# creates a new migration
	migrate create -ext sql -dir db/migrations -seq $(name)

p_up:
	# postgres up - creates postgres server
	docker compose up -d

p_down:
	# postgres down - stops postgres server
	docker compose down

db_up:
	docker exec -it kasho_postgres createdb --username=root --owner=root kasho_db

db_down:
	docker exec -it kasho_postgres dropdb --username=root kasho_db

m_up:
	# run migrate up
	migrate -path db/migrations -database "postgres://root:root@localhost:5432/kasho_db?sslmode=disable" up

m_down:
	# run migrate down
	migrate -path db/migrations -database "postgres://root:root@localhost:5432/kasho_db?sslmode=disable" down

