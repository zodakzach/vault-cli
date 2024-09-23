# CLI Password Manager

## Commands Overview
**Root command:** `vault`

1. **`add`**:
   - **Description**: Add a new password entry to the vault.
   - **Usage**:
     ```bash
     <root> add --service "Service Name" --username "user@example.com" --password "password123"
     <root> add --service "Service Name" --username "user@example.com" --generate
     ```

2. **`delete`**:
   - **Description**: Delete a specific password entry from the vault.
   - **Usage**:
     ```bash
     <root> delete --service "Service Name"
     ```

3. **`list`**:
   - **Description**: List all stored services in the vault.
   - **Usage**:
     ```bash
     <root> list
     ```

4. **`help`**:
   - **Description**: Display help for all available commands.
   - **Usage**:
     ```bash
     <root> help
     ```

5. **`get`**:
   - **Description**: Retrieve a username and password for a specific service, using fuzzy search if necessary.
   - **Usage**:
     ```bash
     <root> get --service "Service Name"
     ```

6. **`update`**:
   - **Description**: Update the username or password for an existing service.
   - **Usage**:
     ```bash
     <root> update --service "Service Name" --username "newuser@example.com" --password "newpassword123"
     ```

7. **`generate`**:
   - **Description**: Generate a strong random password based on specified criteria.
   - **Usage**:
     ```bash
     <root> generate --length 16 --special --numbers --uppercase
     ```

8. **`set-master`**:
   - **Description**: Set or update the master password for accessing the vault.
   - **Usage**:
     ```bash
     <root> set-master --password "newmasterpassword"
     ```

9. **`search`**:
   - **Description**: Search for a service or username in the vault using fuzzy search.
   - **Usage**:
     ```bash
     <root> search --query "searchterm"
     ```

10. **`import`**:
    - **Description**: Import password entries from a file (CSV or JSON).
    - **Usage**:
      ```bash
      <root> import --file "path/to/file.csv" --format csv
      ```

11. **`export`**:
    - **Description**: Export all stored password entries to a file (CSV or JSON).
    - **Usage**:
      ```bash
      <root> export --file "path/to/file.csv" --format csv
      ```

12. **`backup`**:
    - **Description**: Create an encrypted backup of the password database.
    - **Usage**:
      ```bash
      <root> backup --file "path/to/backup.enc"
      ```

13. **`restore`**:
    - **Description**: Restore the vault from an encrypted backup file.
    - **Usage**:
      ```bash
      <root> restore --file "path/to/backup.enc"
      ```

14. **`unlock`**:
    - **Description**: Unlock the vault by providing the master password.
    - **Usage**:
      ```bash
      <root> unlock --password "masterpassword"
      ```

15. **`lock`**:
    - **Description**: Manually lock the vault to require the master password for future operations.
    - **Usage**:
      ```bash
      <root> lock
      ```

16. **`config`**:
    - **Description**: Configure settings such as the auto-lock timeout period.
    - **Usage**:
      ```bash
      <root> config --lock-timeout 10  # Set auto-lock to 10 minutes of inactivity
      ```
