CREATE TABLE IF NOT EXISTS public."COMMUNITY" (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    rules TEXT,
    banner_url VARCHAR(255),
    profile_url VARCHAR(255),
    created_by VARCHAR(255) NOT NULL,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL
);
    
CREATE TABLE IF NOT EXISTS public."USER" (
    id UUID PRIMARY KEY
)
