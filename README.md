[![CodeQL Advanced](https://github.com/ngmisl/dotclaude/actions/workflows/codeql.yml/badge.svg)](https://github.com/ngmisl/dotclaude/actions/workflows/codeql.yml)

# dotclaude

A curated collection of AI development agents and prompts to enhance your AI-assisted development workflow, supporting both Claude Code and Gemini CLI.

## Purpose

This repository provides:
- **Development guidelines** (CLAUDE.md) with SOLID principles and best practices
- **Specialized agents** for code quality, debugging, validation, and more
- **Prompt sanitizer** to detect malicious content in prompts/agents
- **Multi-AI support** for both Claude Code and Gemini CLI agents

## Usage

### Setup Your Project

#### For Claude Code

Copy these files to your project root:

1. **CLAUDE.md** - Development guidelines that Claude will follow
2. **.claude/agents/** - Pre-configured specialized agents

```bash
cp CLAUDE.md /path/to/your/project/
cp -r .claude /path/to/your/project/
```

#### For Gemini CLI

Copy these files to your project root:

1. **CLAUDE.md** - Development guidelines (referenced in AGENTS.md)
2. **.gemini/** - Gemini agent configuration

```bash
cp CLAUDE.md /path/to/your/project/
cp -r .gemini /path/to/your/project/
```

The Gemini setup includes:
- `.gemini/AGENTS.md` - Unified agent definitions converted from Claude agents
- `.gemini/settings.json` - Configuration to load agents from AGENTS.md

### Scan Before Using (Important!)

**Always run the prompt sanitizer first** before using any prompts or agents from external sources:

```bash
# Basic scan
go run main.go

# Build once, run multiple times
go build -o prompt-sanitizer main.go
./prompt-sanitizer

# JSON output for CI/CD integration
./prompt-sanitizer --output=json

# Scan specific directory
./prompt-sanitizer /path/to/directory
```

#### Why Scan First?

The prompt sanitizer detects hidden malicious content that could compromise your AI interactions:

- **Hidden characters**: Zero-width spaces, direction overrides that hide instructions
- **Encoded payloads**: Base64-encoded malicious commands (exec, eval, override)
- **Homoglyphs**: Cyrillic/Greek characters masquerading as Latin letters
- **Control characters**: Non-printable characters used for obfuscation

**Example threats**:
- `ignore all previous instructions` hidden in zero-width characters
- Base64-encoded `exec('malicious command')`
- Cyrillic 'e' in `еxecute` that looks identical to Latin 'e'

#### Scanner Features

- **Comprehensive Testing**: Unit tests ensure reliability and prevent regressions
- **Configurable Scanning**: Customize ignored directories and file extensions via `config.yaml`
- **JSON Output**: Machine-readable results with `--output=json` for automation
- **CI/CD Ready**: Exits with non-zero status on high-severity issues for build integration
- **Improved Detection**: Enhanced Base64 pattern matching catches more potential threats

Running the scanner ensures your prompts and agents are safe before integrating them into your workflow.

## Configuration

The scanner uses `config.yaml` to customize its behavior:

```yaml
ignored_directories:
  - "node_modules"
  - "vendor"
  - ".git"
  - "dist"
  - "build"
  - ".next"

supported_extensions:
  - ".xml"
  - ".md"
  - ".yaml"
  - ".yml"
  - ".txt"
```

Customize this file to match your project's structure and the file types you want to scan.

## Testing

The project includes comprehensive unit tests to ensure scanner reliability:

```bash
# Run all tests
go test -v

# Run tests with coverage
go test -v -cover

# View detailed coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Tests validate:
- Detection of zero-width characters, homoglyphs, and control characters
- Base64 encoded content identification
- JSON output formatting
- Configuration loading

## Recent Improvements

**v1.1.0** - Gemini Support & Enhanced Scanner
- Added Gemini CLI agent support with unified AGENTS.md configuration
- Implemented comprehensive unit test suite (main_test.go)
- Fixed Base64 detection regex for improved pattern matching
- Added config.yaml for customizable scanning behavior
- Implemented JSON output format for CI/CD integration
- Added proper exit codes for automated build pipelines

## Agent Sources

- [ClaudeCodeAgents](https://github.com/darcyegb/ClaudeCodeAgents)

## Support

- [Donate](https://fourzerofour.fkey.id)
- [$MTDV Developer Token](https://swap.cow.fi/#/8453/swap/ETH/MTDV)
- [Farcaster Social Media](https://farcaster.xyz/~/code/I1TPIM) (ref)
- [Private Wallet](https://i.fluidkey.com/SZC2U1) (ref)
- [AI Trading Bot FREE](https://senpi.ai/?r=M1977) (ref)
