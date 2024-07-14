-- Write your migrate up statements here

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(15) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user', -- or use an ENUM type if preferred
    admin BOOLEAN NOT NULL DEFAULT FALSE
);


---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

