# Stocky Backend Assignment

This repository contains the backend implementation for the Stocky assignment,
built using Golang, Gin, Logrus, and PostgreSQL.

---

## Prerequisites

Ensure the following are installed on your system:

- Go 1.20 or higher
- PostgreSQL 14 or higher
- Git
- pgAdmin or psql

---

## Clone the Repository

git clone https://github.com/<your-username>/stocky-backend-assignment.git
cd stocky-backend-assignment
---

## Install Dependencies

The project uses Go modules. Install all required dependencies using:

go mod tidy

---

## Database Setup

You can set up the database using either pgAdmin or the PostgreSQL terminal.

Create a database named `assignment`.

Using PostgreSQL terminal:

CREATE DATABASE assignment;

Alternatively, you can create the database using pgAdmin by:
- Opening pgAdmin
- Right-clicking on Databases
- Selecting Create → Database
- Naming it `assignment`

---

## Apply Database Schema

The database schema is defined in the `migrations.sql` file.

Apply the schema using the PostgreSQL terminal:

psql -U postgres -d assignment -f migrations.sql

You can also open the `migrations.sql` file in pgAdmin’s Query Tool and execute it manually
in case of any issues.

---

## Environment Configuration

Create a `.env` file in the project root with the following values:

DB_HOST=localhost  
DB_PORT=5432  
DB_USER=postgres  
DB_PASSWORD=your_password  
DB_NAME=assignment  



---

## Run the Application

Start the server using:

go run .

If the setup is correct, the application will start and listen on port 8080.
