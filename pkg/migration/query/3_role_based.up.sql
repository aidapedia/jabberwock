CREATE TABLE IF NOT EXISTS "resources"(
   "id" BIGSERIAL PRIMARY KEY,
   "type" VARCHAR(20) NOT NULL,   -- http, rpc
   "method" VARCHAR(20),          -- GET, POST (null for RPC)
   "path" VARCHAR(255) NOT NULL,   -- /orders, /orders/:id, MethodName
   "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "permissions"(
   "id" BIGSERIAL PRIMARY KEY,
   "name" VARCHAR NOT NULL,
   "description" TEXT,
   "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "resources_permissions"(
   "permission_id" BIGINT NOT NULL REFERENCES permissions(id),
   "resource_id" BIGINT NOT NULL REFERENCES resources(id),
   "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   CONSTRAINT "resources_permissions_pkey" PRIMARY KEY (resource_id, permission_id)
);

CREATE TABLE IF NOT EXISTS "roles"(
   "id" BIGSERIAL PRIMARY KEY,
   "name" VARCHAR NOT NULL,
   "description" TEXT,
   "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS roles_name_key ON roles ("name");

CREATE TABLE IF NOT EXISTS "role_permissions"(
   "role_id" BIGINT NOT NULL,
   "permission_id" BIGINT NOT NULL,
   "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   CONSTRAINT "role_permissions_pkey" PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS "user_roles"(
   "user_id" BIGINT NOT NULL REFERENCES users(id),
   "role_id" BIGINT NOT NULL REFERENCES roles(id),
   "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   CONSTRAINT "user_roles_pkey" PRIMARY KEY (user_id, role_id)
);

ALTER TABLE "users" DROP COLUMN "type";