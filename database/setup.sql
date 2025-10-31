-- TatApps Database Setup Script
-- PostgreSQL 14+

-- Create database and user
CREATE DATABASE tatapps;
CREATE USER tatapps WITH PASSWORD 'tatapps123!';
GRANT ALL PRIVILEGES ON DATABASE tatapps TO tatapps;

-- Connect to tatapps database
\c tatapps

-- Grant schema privileges
GRANT ALL ON SCHEMA public TO tatapps;

-- Note: Tables will be automatically created by GORM AutoMigrate
-- when you first run the application

-- Optional: Create extensions if needed
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- for full-text search
