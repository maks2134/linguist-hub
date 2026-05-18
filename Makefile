.PHONY: up down dev-backend dev-frontend logs clean swagger

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

swagger:
	cd backend && make swagger

dev-backend:
	cd backend && make dev

dev-frontend:
	cd frontend && npm run dev

clean:
	docker compose down -v
	rm -rf backend/bin
	rm -rf backend/docs
	rm -rf frontend/.next
	rm -rf frontend/node_modules

install-frontend:
	cd frontend && npm install

build:
	docker compose build

restart: down up
