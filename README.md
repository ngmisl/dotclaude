# dotclaude

A curated collection of Claude Code agents and development prompts to enhance your AI-assisted development workflow.

## Purpose

This repository provides:
- **Development guidelines** (CLAUDE.md) with SOLID principles and best practices
- **Specialized agents** for code quality, debugging, validation, and more
- **Prompt sanitizer** to detect malicious content in prompts/agents

## Usage

### Setup Your Project

Copy these files to your project root:

1. **CLAUDE.md** - Development guidelines that Claude will follow
2. **.claude/agents/** - Pre-configured specialized agents

```bash
cp CLAUDE.md /path/to/your/project/
cp -r .claude /path/to/your/project/
```

### Scan Before Using (Important!)

**Always run the prompt sanitizer first** before using any prompts or agents from external sources:

```bash
go run main.go
# or build it first
go build -o prompt-sanitizer main.go
./prompt-sanitizer
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

Running the scanner ensures your prompts and agents are safe before integrating them into your workflow.

## Agent Sources

- [ClaudeCodeAgents](https://github.com/darcyegb/ClaudeCodeAgents)

## Support

- [Donate](https://fourzerofour.fkey.id)
- [Farcaster Social Media](https://farcaster.xyz/~/code/I1TPIM) (ref)
- [Private Wallet](https://i.fluidkey.com/SZC2U1) (ref)
- [AI Trading Bot FREE](https://senpi.ai/?r=M1977) (ref)
