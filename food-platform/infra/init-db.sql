-- Create databases for each microservice (database-per-service pattern)
-- Reference: docs/ARCHITECTURE.md Section 6.1

CREATE DATABASE auth_db;
CREATE DATABASE restaurants_db;
CREATE DATABASE menus_db;
CREATE DATABASE orders_db;
CREATE DATABASE payments_db;
CREATE DATABASE dispatch_db;
CREATE DATABASE drivers_db;
CREATE DATABASE notifications_db;
CREATE DATABASE fraud_db;
CREATE DATABASE promos_db;
CREATE DATABASE audit_db;
CREATE DATABASE analytics_db;

-- Enable required extensions on each database
\connect auth_db
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\connect restaurants_db
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";

\connect menus_db
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\connect orders_db
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";

\connect payments_db
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\connect dispatch_db
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";

\connect drivers_db
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\connect notifications_db
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\connect fraud_db
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\connect promos_db
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\connect audit_db
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto"; -- for digest functions

\connect analytics_db
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
