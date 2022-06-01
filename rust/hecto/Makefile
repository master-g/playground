-include .env

# TARGETS := feocore

.PHONY: lint
lint:
	@echo "  >  Linting..."
	@rustup component add clippy
	@cargo clippy -- -W clippy::pedantic


.PHONY: fix
fix:
	@echo "  >  Linting and fixing..."
	@rustup component add clippy
	@cargo clippy --fix


.PHONY: fmt
fmt:
	@echo "  >  Formatting..."
	@rustup component add rustfmt
	@cargo fmt
