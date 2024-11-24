# Database

This project uses **dbmate** and **sqlboiler** for managing database migrations and ORM-based data handling, respectively. The environment variables for database configuration are automatically loaded with **devenv**, simplifying database operations like migrations and generating Go models.

## Requirements

- devenv
- docker

## Setup

Before using dbmate or sqlboiler, ensure that you have the appropriate environment variables set. These are automatically loaded through **devenv**, so you can simply navigate to the `back` directory to begin working with the database.

### Environment Variables

Environment variables are automatically managed using **devenv**. This makes it easy to load the necessary configuration and start working with the database without having to manually set up variables each time.

---

## Migrations with dbmate

**dbmate** is used for managing database migrations. It simplifies the process of creating, applying, and rolling back database migrations.

### Commands

- **Create a new migration**:  
  To create a new migration file, run the following command from the `back` directory:

  ```bash
  dbmate new <migration_name>
  ```

- **Apply migrations**:  
  To apply any pending migrations to the database, use:

  ```bash
  dbmate up
  ```

- **Rollback migrations**:  
  To rollback the most recent migration, use:

  ```bash
  dbmate down
  ```

- **Reset the database**:  
  To reset the database, i.e., drop all tables and reapply all migrations, run:

  ```bash
  dbmate reset
  ```

### Documentation for dbmate

For more detailed information about how to use **dbmate** for database migrations, refer to the official documentation:

- [dbmate GitHub Repository](https://github.com/amacneil/dbmate)

---

## ORM with sqlboiler

**sqlboiler** is used to generate Go models from the database schema. It simplifies working with the database by providing strongly typed Go structs for database tables.

### Generate Models

After updating the database schema (e.g., by running migrations), you can regenerate the Go models by running:

```bash
go generate .
```

This will automatically generate or update the models based on the current state of the database.

### Documentation for sqlboiler

For more information on how to use **sqlboiler** to work with Go models and interact with the database, refer to the official documentation:

- [sqlboiler GitHub Repository](https://github.com/volatiletech/sqlboiler)
- [sqlboiler Documentation](https://github.com/volatiletech/sqlboiler#readme)

---

## Best Practices

- **Version Control**: Keep migration files in version control to ensure consistency across environments.
- **Generate Models**: After each migration, run `go generate .` to regenerate the Go models so they stay in sync with the database schema.
- **Environment Consistency**: Use **devenv** to ensure that the environment variables are always loaded correctly before interacting with the database.
