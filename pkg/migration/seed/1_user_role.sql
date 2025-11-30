INSERT INTO "public"."users" ("id", "name", "password", "email", "phone", "is_verified", "status", "avatar_url", "created_at", "updated_at") VALUES
(1, 'Super Admin', '$2a$10$QOnrE6aIpODYwng14jvImuhX2PeyizL8pDy.0TgEeWDLQXQJOaU1S', 'john.doe@gmail.com', '6281234567890', 3, 0, '', '2025-11-12 13:38:36.363524', '2025-11-12 13:38:36.363524');

INSERT INTO "public"."permissions" ("id", "name", "description", "created_at", "updated_at") VALUES
(1, 'auth', 'Basic Authentication', '2025-11-30 04:19:03.793515', '2025-11-30 04:19:03.793515'),
(2, 'user.read', 'Show User Data', '2025-11-30 04:20:11.96249', '2025-11-30 04:20:11.96249');

INSERT INTO "public"."resources" ("id", "type", "method", "path", "created_at", "updated_at") VALUES
(1, 'http', 'POST', '/logout', '2025-11-30 04:09:14.586025', '2025-11-30 04:09:14.586025'),
(2, 'http', 'GET', '/user/:id', '2025-11-30 04:09:14.586025', '2025-11-30 04:09:14.586025');

INSERT INTO "public"."resources_permissions" ("permission_id", "resource_id", "created_at", "updated_at") VALUES
(1, 1, '2025-11-30 04:19:10.457759', '2025-11-30 04:19:10.457759'),
(2, 2, '2025-11-30 04:20:27.782537', '2025-11-30 04:20:27.782537');

INSERT INTO "public"."role_permissions" ("role_id", "permission_id", "created_at", "updated_at") VALUES
(2, 1, '2025-11-30 04:21:29.327505', '2025-11-30 04:21:29.327505');

INSERT INTO "public"."roles" ("id", "name", "description", "created_at", "updated_at") VALUES
(1, 'superadmin', 'Administrator Roles', '2025-11-30 04:21:10.348797', '2025-11-30 04:21:10.348797'),
(2, 'member', 'Member Roles', '2025-11-30 04:21:10.348797', '2025-11-30 04:21:10.348797');