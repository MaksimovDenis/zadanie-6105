CREATE TABLE employee (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE organization (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organization(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES employee(id) ON DELETE CASCADE
);

CREATE TYPE status_type AS ENUM (
    'Created',
    'Published',
    'Canceled',
    'Approved',
    'Rejected',
    'Closed'
);

CREATE TABLE tender (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organization(id) ON DELETE CASCADE,
    created_by INTEGER REFERENCES employee(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status status_type DEFAULT 'Created',
    service_type VARCHAR(255),
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tender_history (
    id SERIAL PRIMARY KEY,
    tender_id INTEGER REFERENCES tender(id) ON DELETE CASCADE,
    name VARCHAR(255),
    description TEXT,
    service_type VARCHAR(255),
    status status_type,
    version INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bids (
    id SERIAL PRIMARY KEY,
    tender_id INTEGER REFERENCES tender(id) ON DELETE CASCADE,
    organization_id INTEGER REFERENCES organization(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES employee(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status status_type DEFAULT 'Created', 
    author_type VARCHAR(50) DEFAULT 'User', 
    version INTEGER DEFAULT 1, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bids_decisions (
    id SERIAL PRIMARY KEY,
    bid_id INTEGER REFERENCES bids(id) ON DELETE CASCADE,
    decision_type VARCHAR(20) CHECK (decision_type IN ('Approved', 'Rejected')),
    user_id INTEGER REFERENCES employee(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE bids_history (
    id SERIAL PRIMARY KEY,
    tender_id INTEGER REFERENCES tender(id) ON DELETE CASCADE,
    bid_id INTEGER REFERENCES bids(id) ON DELETE CASCADE,
    organization_id INTEGER REFERENCES organization(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES employee(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status status_type,
    author_type VARCHAR(50),
    version INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bid_feedbacks (
    id SERIAL PRIMARY KEY,
    bid_id INTEGER REFERENCES bids(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES employee(id) ON DELETE CASCADE,
    feedback TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);