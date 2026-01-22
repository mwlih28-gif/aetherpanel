-- Aether Panel Database Initialization Script

-- Create database if not exists
SELECT 'CREATE DATABASE aether_panel'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'aether_panel')\gexec

-- Connect to the database
\c aether_panel;

-- Create user if not exists
DO
$do$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE  rolname = 'aether') THEN

      CREATE ROLE aether LOGIN PASSWORD 'aether_secure_password';
   END IF;
END
$do$;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE aether_panel TO aether;
GRANT ALL ON SCHEMA public TO aether;

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Basic tables will be created by the application migration
