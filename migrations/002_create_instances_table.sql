-- Write your migrate up statements here

-- Create the status type
CREATE TYPE STATUS_TYPE AS ENUM ('active', 'disabled');

-- Create the instances table
CREATE TABLE instances (
    id SERIAL PRIMARY KEY,
    client_id UUID NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
    "name" VARCHAR(255) NOT NULL,
    host_address VARCHAR(255) DEFAULT NULL,
    "status" STATUS_TYPE NOT NULL DEFAULT 'active',
    "description" TEXT DEFAULT NULL,
    created_by INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Create the instance_users join table
CREATE TABLE instance_host_users (
    instance_id INT NOT NULL,
    username VARCHAR(255) NOT NULL,
    PRIMARY KEY (instance_id, username),
    FOREIGN KEY (instance_id) REFERENCES instances(id) ON DELETE CASCADE
);

-- Create the instance_users join table
CREATE TABLE instance_users (
    instance_id INT NOT NULL,
    user_id INT NOT NULL,
    instance_host_username VARCHAR(255)  NOT NULL,
    PRIMARY KEY (instance_id, user_id, instance_host_username),
    FOREIGN KEY (instance_id) REFERENCES instances(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (instance_id, instance_host_username) REFERENCES instance_host_users(instance_id, username) ON DELETE CASCADE
);

-- -- Create the instance_endpoints table
-- CREATE TABLE instance_endpoints (
--     id SERIAL PRIMARY KEY,
--     instance_id INT NOT NULL,
--     endpoint_uri VARCHAR(255) NOT NULL,
--     FOREIGN KEY (instance_id) REFERENCES instances(id) ON DELETE CASCADE
-- );

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

