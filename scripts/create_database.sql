-- Script to create the PostgreSQL database for D&D Combat API
-- Run as postgres superuser

-- Create database if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'dnd_combat') THEN
        CREATE DATABASE dnd_combat
            WITH 
            OWNER = postgres
            ENCODING = 'UTF8'
            LC_COLLATE = 'en_US.utf8'
            LC_CTYPE = 'en_US.utf8'
            TEMPLATE = template0
            CONNECTION LIMIT = -1;
        
        COMMENT ON DATABASE dnd_combat IS 'Database for D&D Combat Virtual Tabletop Application';
    END IF;
END
$$;

-- Connect to the database
\connect dnd_combat;

-- Create application role if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'dnd_app_user') THEN
        CREATE ROLE dnd_app_user WITH
            LOGIN
            NOSUPERUSER
            NOCREATEDB
            NOCREATEROLE
            INHERIT
            NOREPLICATION
            CONNECTION LIMIT -1
            PASSWORD 'dnd_password';
            
        COMMENT ON ROLE dnd_app_user IS 'Application role for D&D Combat API';
    END IF;
END
$$;

-- Grant privileges to the application user
GRANT ALL PRIVILEGES ON DATABASE dnd_combat TO dnd_app_user; 