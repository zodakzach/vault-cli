# CLI Sensitive Data Manager
`vault-cli` is a command-line interface (CLI) tool built with Go and Cobra for managing data securely. The tool supports various operations like adding, deleting, listing, importing, and exporting data, while ensuring that all sensitive data values are encrypted in a SQLite database.

## Commands Overview

**Root Command:** `vault-cli`

The root command `vault-cli` is the entry point for all operations. Below is an overview of the available subcommands.


1. **`add`** - Add a new data entry

The `add` command allows users to securely add new sensitive data entries to the vault. These entries can include usernames, email addresses, API keys, or other secret values associated with a service.

```bash
vault-cli add --service <service_name>
```

2. **`unlock`** - Unlock the vault

The `unlock` command is used to unlock the vault by providing the correct master password. This allows you to access or modify the sensitive data stored within the vault.

```bash
vault-cli unlock --password <master_password>
```

3. **`lock`** - Lock the vault

The `lock` command is used to secure the vault, preventing access to sensitive data until it is unlocked again.

```bash
vault-cli lock
```

4. **`delete`** - Delete a stored entry from the vault

The `delete` command allows users to remove a stored sensitive data entry from the vault using the specified service and identifier.

```bash
vault-cli delete --service <service_name> --identifier <identifier_value>
```

5. **`get`** - Retrieve a sensitive data entry from the vault

The `get` command allows users to retrieve a specific sensitive data entry from the vault by providing the associated service and identifier.

```bash
vault-cli get --service <service_name> --identifier <identifier_value>
```

6. **`set-master`** - Set or update the master password

The `set-master` command allows users to set a new master password for accessing the vault or update the existing one.

```bash
vault-cli set-master --password <new_master_password> [--old-password <old_master_password>]
```

7. **`update`** - Update a sensitive data entry in the vault

The `update` command allows users to modify the value or identifier for a specific service stored in the vault.

```bash
vault-cli update --service <service_name> --identifier <identifier_value>
```

8. **`list`** - List all stored services and identifiers

The `list` command provides a way for users to view all services stored in the vault along with their associated identifiers.

```bash
vault-cli list [--id-type <identifier_type>]
```

9. **`generate`** - Generate a random password

The `generate` command allows users to create a secure random password of a specified length.

```bash
vault-cli generate [--length <length>]
```

10. **`export`** - Export all sensitive data entries to a file (CSV or JSON)

The `export` command allows users to export all stored sensitive data entries from the vault into a specified file format, either CSV or JSON.

```bash
vault-cli export --file <file_path> --format <format>
```

11. **`import`** - Import password entries from a file (CSV or JSON)

The `import` command allows users to import password entries into the vault from a specified file in either CSV or JSON format.

```bash
vault-cli import --file <file_path>
```