-- 0001_baseline.up.sql
-- Namespace ayrımı (docs/database/README.md ile hizalı):
--   auth        — hesap, oturum, cihaz, refresh token
--   character   — karakter, stat, envanter projeksiyonu
--   content     — content pack template'leri (item, race, class, dungeon, drop table)
--   world       — dünya hiyerarşisi, region, location
--   activity    — AFK contract, dungeon run
--   economy     — ledger, wallet, currency, market, vendor
--   social      — friend, block, mail, guild
--   companion   — pet, mount instance
--   progression — mastery, achievement, prestige, reputation
--   audit       — admin_action, security_event, ledger_audit, trade_log
--   api         — idempotency, rate-limit persistent bucket (Faz 01 sprint 2)
--   ops         — admin_user, admin_role, admin_user_role
--   analytics   — combat_event, market_search_event (24h hot)

CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS character;
CREATE SCHEMA IF NOT EXISTS content;
CREATE SCHEMA IF NOT EXISTS world;
CREATE SCHEMA IF NOT EXISTS activity;
CREATE SCHEMA IF NOT EXISTS economy;
CREATE SCHEMA IF NOT EXISTS social;
CREATE SCHEMA IF NOT EXISTS companion;
CREATE SCHEMA IF NOT EXISTS progression;
CREATE SCHEMA IF NOT EXISTS audit;
CREATE SCHEMA IF NOT EXISTS api;
CREATE SCHEMA IF NOT EXISTS ops;
CREATE SCHEMA IF NOT EXISTS analytics;

-- Extensions (available in stock Postgres 16 image)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
