# Vite Makefile

# Path to package manager (npm/yarn/pnpm)
PKG_MANAGER = npm

# Commands
build: lint format
	$(PKG_MANAGER) run build

dev:
	$(PKG_MANAGER) run dev

lint:
	$(PKG_MANAGER) run lint

format:
	$(PKG_MANAGER) run format

check:
	$(PKG_MANAGER) run lint
	$(PKG_MANAGER) run format -- --check

clean:
	rm -rf dist
