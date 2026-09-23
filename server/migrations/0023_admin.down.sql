-- 0023_admin.down.sql
DROP TABLE IF EXISTS ops.account_sanction;
DROP TABLE IF EXISTS audit.admin_action;
DROP TABLE IF EXISTS ops.admin_session;
DROP TABLE IF EXISTS ops.admin_user_role;
DROP TABLE IF EXISTS ops.admin_user;
