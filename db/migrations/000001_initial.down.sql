TRUNCATE users CASCADE ;
TRUNCATE user_role_permissions CASCADE;
TRUNCATE user_roles CASCADE ;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS user_role_permissions;
DROP TABLE IF EXISTS user_apitoken;
DROP TABLE IF EXISTS user_roles;


DROP EXTENSION IF EXISTS pgcrypto;