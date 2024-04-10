start:
	@echo " > запуск docker compose..."
	@echo "______________________________"
	docker compose up --build -d
	@echo " > docker compose запущен"

stop:
	@echo " > остановка docker compose..."
	@echo "______________________________"
	docker compose down
	@echo "______________________________"
	@echo " > docker compose остановлен"

restart:
	@echo " > остановка docker compose..."
	@echo "______________________________"
	docker compose down
	@echo " > docker compose остановлен"
	@echo "______________________________"
	@echo " > запуск docker compose..."
	@echo "______________________________"
	docker compose up --build -d
	@echo "______________________________"
	@echo " > docker compose запущен"

generate_proto:
	@echo " > генерируем grpc код..."
	@echo "______________________________"
	cd ../.. && cd api-gateway/internal/tools/grpc/parsing && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative parse.proto
	cd ../.. && cd parsing-service/internal/tools/grpc/parsing && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative parsing.proto
	@echo "> нужные файлы сгенерированы"

generate_keys:
	@echo "> генерируем ключи..."
	@echo "______________________________"
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt -subj "/C=RU/ST=Some-State/L=parsing-service/O=parsing/OU=p/CN=parsing-service" -addext "subjectAltName = DNS:parsing-service"
	@echo "> ключи сгенерированы"