-- migrate:up

ALTER TABLE users
MODIFY COLUMN id SERIAL COMMENT 'visible:"admin,user,other"',
MODIFY COLUMN username VARCHAR(255) UNIQUE NOT NULL COMMENT 'visible:"admin,user,other"',
MODIFY COLUMN password_hash VARCHAR(255) NOT NULL COMMENT 'visible:"nobody"',
MODIFY COLUMN email VARCHAR(255) UNIQUE NOT NULL COMMENT 'visible:"admin,user"',
MODIFY COLUMN email_verified BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'visible:"admin,user"',
MODIFY COLUMN email_verification_token VARCHAR(255) COMMENT 'visible:"admin"',
MODIFY COLUMN is_admin BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'visible:"admin"',
MODIFY COLUMN is_hidden BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'visible:"admin"',
MODIFY COLUMN created_at TIMESTAMP NOT NULL DEFAULT NOW() COMMENT 'visible:"admin"',
MODIFY COLUMN updated_at TIMESTAMP NOT NULL DEFAULT NOW() COMMENT 'visible:"admin"',
MODIFY COLUMN deleted_at TIMESTAMP COMMENT 'visible:"admin"';

-- migrate:down

