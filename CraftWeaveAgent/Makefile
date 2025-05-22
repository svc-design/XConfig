BINARY=cw-agent

all: build

init:
	@echo "🔍 Checking for Rust toolchain..."
	@if ! command -v cargo >/dev/null 2>&1; then \
		echo "🚀 Rust not found. Installing via rustup..."; \
		curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y; \
		. $$HOME/.cargo/env; \
	fi

	@echo "🛠️ Initializing Rust project..."
	@test -f Cargo.toml || cargo init --vcs none

	@echo "📦 Ensuring required dependencies..."
	@cargo add tokio --features full
	@cargo add serde --features derive
	@cargo add serde_json
	@cargo add reqwest --features json
	@cargo add git2
	@cargo add chrono
	@cargo add clap --features derive
	@cargo add log
	@cargo add thiserror
	@cargo add anyhow

	@echo "📁 Creating module files..."
	@mkdir -p src
	@touch src/{main.rs,scheduler.rs,config.rs,executor.rs,result_store.rs,models.rs}

	@echo "✅ cw-agent init complete."

build:
	cargo build --release

run:
	cargo run -- --oneshot

install:
	install -m 755 target/release/$(BINARY) /usr/local/bin/$(BINARY)

test:
	cargo test

fmt:
	cargo fmt

lint:
	cargo clippy -- -D warnings

clean:
	cargo clean

