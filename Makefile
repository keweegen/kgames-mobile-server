deploymentPath := ./deployments/local
dcPath := ${deploymentPath}/docker-compose.yml
envPath := ${deploymentPath}/.env
dcEnvPath := ${deploymentPath}/.env.docker

ifneq ("$(wildcard ${envPath})","")
	include ${envPath}
endif

ifneq ("$(wildcard ${dcEnvPath})","")
	include ${dcEnvPath}
endif

dc := docker compose \
	-f ${dcPath} \
	--env-file ${envPath} \
	--env-file ${dcEnvPath} \
	-p ${APP_NAME}

dbDSN := "dbname=${DB_NAME} user=${DB_USER} password=${DB_PASSWORD} host=${DB_HOST} port=${DB_PORT} sslmode=${DB_SSL_MODE} search_path=${DB_SCHEMA}"
dbDSNTemp := "dbname=${DB_NAME} user=${DB_USER} password=${DB_PASSWORD} host=${DB_HOST} port=${DB_PORT_TEMP} sslmode=${DB_SSL_MODE}"
migration := goose -table migrations -dir ./migrations -allow-missing postgres ${dbDSN}
migrationTemp := goose -table migrations -dir ./migrations -allow-missing postgres ${dbDSNTemp}

build:
	GOOS=linux GOARCH=arm64 go build -o ./bin/server .
	@chmod +x ./bin/server

d-build:
	@${dc} build

d-up: build
	@${dc} up -d

d-down:
	@${dc} down --remove-orphans

m-create:
	@${migration} create ${name} sql

m-up:
	@${migration} up

m-down:
	@${migration} down

m-status:
	@${migration} status

proto:
	@make proto-one name=game
	@make proto-one name=user

proto-one: proto-clean
	@protoc --go_out=api/grpc/${name} \
		--go-grpc_out=require_unimplemented_servers=false:. \
		--go-grpc_opt=paths=source_relative \
		api/grpc/${name}/*.proto

proto-clean:
	@rm -rf api/grpc/${name}/*.pb.go

model:
	@docker run --rm -d -p ${DB_PORT_TEMP}:5432 \
		--name ${APP_NAME}-migration-db \
		-e POSTGRES_USER=${DB_USER} \
		-e POSTGRES_PASSWORD=${DB_PASSWORD} \
		-e POSTGRES_DB=${DB_NAME} \
		${POSTGRES_IMAGE}
	@sleep 2
	@${migrationTemp} up
	@sed 's/= \"dbname\"/= \"$(DB_NAME)\"/g; s/= \"dbusername\"/= \"$(DB_USER)\"/g; s/= \"dbpassword\"/= \"$(DB_PASSWORD)\"/g; s/= 5432/= $(DB_PORT_TEMP)/g; s/= \"myschema\"/= \"$(DB_SCHEMA)\"/g' configs/sqlboiler.toml > ./temp_boiler_config.toml
	@sqlboiler --config ./temp_boiler_config.toml psql
	@rm ./temp_boiler_config.toml
	@docker stop ${APP_NAME}-migration-db
