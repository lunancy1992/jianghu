.PHONY: dev dev-app dev-server build build-app build-server clean

# ====== 前端 ======
dev-app:
	cd app && npm run dev:h5

build-app:
	cd app && npm run build:h5

build-app-weixin:
	cd app && npm run build:mp-weixin

# ====== 后端 ======
dev-server:
	cd server && make dev

build-server:
	cd server && make build

run-server:
	cd server && make run

migrate:
	cd server && make migrate

# ====== 一键启动（两个终端分别跑） ======
dev:
	@echo "Starting server..."
	cd server && make run &
	@echo "Starting app..."
	cd app && npm run dev:h5

build:
	$(MAKE) build-server
	$(MAKE) build-app

clean:
	cd server && make clean
