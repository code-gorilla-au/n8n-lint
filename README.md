# n8n-lint
Simple linter tool for n8n workflows JSON files.

## Motivation

n8n workflows can become complex and hard to maintain. `n8n-lint` helps you maintain high-quality workflows by identifying common issues and anti-patterns.

## Features 

- Run lint on all n8n workflows within a directory.
- Configurable GLOB include / exclude patterns.
- Multiple built-in rules for workflow validation.
- Configurable reporting levels (error, warn, off).

## CLI Usage

### Installation

You can build the CLI using `task`:

```bash
task build
```

This will create a binary named `n8n-lint-<os>-<arch>`.

### Commands

#### `check`

The `check` command is used to lint n8n workflow files.

```bash
./n8n-lint-darwin-arm64 check [flags]
```

**Flags:**
- `--config`, `-c`: Path to the configuration file (default: `.n8n-lint.yaml`).
- `--verbose`, `-v`: Enable verbose logging.
- `--help`, `-h`: Show help.

## Configuration

`n8n-lint` uses a YAML configuration file to define which rules to run and which files to include or ignore.

### Example Configuration (`.n8n-lint.yaml`)

```yaml
include:
  - "**/*.json"
ignore:
  - "node_modules/**"
  - "backups/**"

rules:
  no_infinite_loop:
    report: error
  no_dead_ends:
    report: warn
    allowed_names:
      - "End Node"
      - "Respond to Webhook"
  no_dangling_ifs:
    report: error
  no_disabled_nodes:
    report: warn
    allowed_names:
      - "Debug Node"
```

### Ruleset Configuration

The `rules` section allows you to configure individual rules:

| Rule Name | Description | Configuration Options |
|-----------|-------------|-----------------------|
| `no_infinite_loop` | Detects potential infinite loops in the workflow. | `report` |
| `no_dead_ends` | Detects nodes that do not lead to any subsequent action. | `report`, `allowed_names` (list of strings) |
| `no_dangling_ifs` | Detects IF nodes where one or both branches are not connected. | `report` |
| `no_disabled_nodes` | Detects disabled nodes in the workflow. | `report`, `allowed_names` (list of strings) |

#### Reporting Levels
- `error`: Fails the lint check (exit code 1).
- `warn`: Prints a warning but does not fail the check.
- `off`: Disables the rule.

## Getting started

Small guide on getting started with local dev.

### Install tools

```bash
task install-local
```

### Install dependencies

```bash
task install
```

### Run CI

```bash
task ci
```

## Add new rule

Use the template to add a new rule.

```bash
# Generate new rule scaffold 
task gen-rule -- --name <rule_name>
```

## License

MIT License. See [LICENSE](LICENSE) for details.