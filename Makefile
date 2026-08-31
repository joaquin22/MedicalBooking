BACKEND_DIR := backend
COMPOSE := cd $(BACKEND_DIR) && docker compose

.PHONY: help build up down start stop logs ps restart clean migrate

help:
	@echo "Comandos disponibles:"
	@echo "  make build     Construir las imágenes de Docker"
	@echo "  make up        Levantar los contenedores en segundo plano"
	@echo "  make down      Detener y eliminar los contenedores"
	@echo "  make start     Levantar los contenedores"
	@echo "  make stop      Detener los contenedores"
	@echo "  make logs      Ver los logs del backend"
	@echo "  make ps        Ver el estado de los contenedores"
	@echo "  make restart   Reiniciar los contenedores"
	@echo "  make clean     Detener y eliminar contenedores y volúmenes"

build:
	$(COMPOSE) build

up: build
	$(COMPOSE) up

down:
	$(COMPOSE) down

start:
	$(COMPOSE) start

stop:
	$(COMPOSE) stop

logs:
	$(COMPOSE) logs -f backend

ps:
	$(COMPOSE) ps

restart:
	$(COMPOSE) restart

clean:
	$(COMPOSE) down -v
