start:
	@echo " > запуск docker compose..."
	docker compose -p big-earth-coursework up --build -d
	@echo " > docker compose запущен"

stop:
	@echo " > остановка docker compose..."
	docker compose -p big-earth-coursework down
	@echo " > docker compose остановлен"

restart:
	@echo " > остановка docker compose..."
	docker compose down
	@echo " > docker compose остановлен"
	@echo " > запуск docker compose..."
	docker compose up --build -d
	@echo " > docker compose запущен"

parsing_pprof_cpu:
	@echo " > загружаем профиль cpu..."
	cd ../.. && cd parsing-service/test/profiler_trace && go tool pprof -http=:6060 cpu.prof

parsing_pprof_heap:
	@echo " > загружаем профиль cpu..."
	cd ../.. && cd parsing-service/test/profiler_trace && go tool pprof -http=:6060 heap.prof

parsing_trace:
	@echo " > загружаем trace информацию..."
	cd ../.. && cd parsing-service/test/profiler_trace && go tool trace trace.out

generate_proto:
	@echo " > генерируем grpc код..."
	cd ../.. && cd api-gateway/internal/tools/grpc/parsing && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative parse.proto
	cd ../.. && cd parsing-service/internal/tools/grpc/parsing && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative parsing.proto
	@echo "> нужные файлы сгенерированы"

generate_keys:
	@echo "> генерируем сертификат..."
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem -subj "/C=RU/ST=Some-State/L=parsing-service/O=parsing/OU=p/CN=parsing-service" -addext "subjectAltName = DNS:parsing-service"
	@echo "> сертификат сгенерированы"