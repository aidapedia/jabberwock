CREATE TABLE IF NOT EXISTS users(
   "id" BIGSERIAL PRIMARY KEY,
   "type" SMALLINT NOT NULL,
   "name" VARCHAR (50) NOT NULL,
   "password" VARCHAR NOT NULL,
   "email" VARCHAR (50) NOT NULL,
   "phone" VARCHAR (20) NOT NULL,
   "is_verified" BIGINT DEFAULT 0,
   "status" SMALLINT DEFAULT 0,
   "avatar_url" VARCHAR ,
   "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_key ON users (email);
CREATE UNIQUE INDEX IF NOT EXISTS users_phone_key ON users (phone);