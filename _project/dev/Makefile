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