-- Create databases directly (will fail if they exist, but that's handled by Docker init)
-- CREATE DATABASE hypertube;
CREATE DATABASE keycloak;

-- Create schema (this supports IF NOT EXISTS)
\c keycloak
CREATE SCHEMA  keycloak_schema;
