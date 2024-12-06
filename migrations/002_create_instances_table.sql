-- Write your migrate up statements here

-- Create the instances table
CREATE TABLE instances (
    id SERIAL PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL UNIQUE,
    created_by INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Create the instance_users join table
CREATE TABLE instance_users (
    instance_id INT NOT NULL,
    username VARCHAR(255) NOT NULL,
    PRIMARY KEY (instance_id, username),
    FOREIGN KEY (instance_id) REFERENCES instances(id) ON DELETE CASCADE
);

-- Create the instance_access_users join table
CREATE TABLE instance_access_users (
    instance_id INT NOT NULL,
    user_id INT NOT NULL,
    instance_user_username VARCHAR(255)  NOT NULL,
    PRIMARY KEY (instance_id, user_id, instance_user_username),
    FOREIGN KEY (instance_id) REFERENCES instances(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (instance_id, instance_user_username) REFERENCES instance_users(instance_id, username) ON DELETE CASCADE
);
-- Create the instance_endpoints table
CREATE TABLE instance_endpoints (
    id SERIAL PRIMARY KEY,
    instance_id INT NOT NULL,
    endpoint_uri VARCHAR(255) NOT NULL,
    FOREIGN KEY (instance_id) REFERENCES instances(id) ON DELETE CASCADE
);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

