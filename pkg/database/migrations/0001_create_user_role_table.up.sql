CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE user_role (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    code VARCHAR UNIQUE NOT NULL,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR,
    is_deleted BOOLEAN DEFAULT FALSE
);