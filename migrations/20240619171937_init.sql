-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE SCHEMA IF NOT EXISTS "kgames";

CREATE TABLE IF NOT EXISTS "kgames"."users"
(
    "id"         uuid PRIMARY KEY     DEFAULT "uuid_generate_v4"(),
    "name"       VARCHAR     NOT NULL,
    "email"      VARCHAR     NOT NULL UNIQUE,
    "locale"     VARCHAR     NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "kgames"."games"
(
    "id"                 uuid PRIMARY KEY        DEFAULT "uuid_generate_v4"(),
    "type_code"          VARCHAR        NOT NULL,
    "state_code"         VARCHAR        NOT NULL,
    "bid"                DECIMAL(10, 2) NOT NULL,
    "max_players"        SMALLINT       NOT NULL,
    "finish_reason_code" VARCHAR        NULL,
    "created_at"         timestamptz    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by"         VARCHAR        NOT NULL,
    "started_at"         timestamptz    NULL,
    "finished_at"        timestamptz    NULL
);

CREATE INDEX IF NOT EXISTS "idx_games_state_code" ON "kgames"."games" USING HASH ("state_code");

CREATE TABLE IF NOT EXISTS "kgames"."game_players"
(
    "game_id"   uuid           NOT NULL REFERENCES "kgames"."games" ("id"),
    "player_id" uuid           NOT NULL REFERENCES "kgames"."users" ("id"),
    "position"  SMALLINT       NOT NULL,
    "profit"    DECIMAL(10, 2) NOT NULL DEFAULT 0.0,
    "fee"       DECIMAL(10, 2) NOT NULL DEFAULT 0.0,
    "gave_up"   BOOLEAN        NOT NULL DEFAULT FALSE,
    "draw"      BOOLEAN        NOT NULL DEFAULT FALSE,
    PRIMARY KEY ("game_id", "player_id")
);

CREATE TABLE IF NOT EXISTS "kgames"."game_moves"
(
    "id"         uuid PRIMARY KEY     DEFAULT "uuid_generate_v4"(),
    "game_id"    uuid        NOT NULL REFERENCES "kgames"."games" ("id"),
    "striker_id" uuid        NOT NULL REFERENCES "kgames"."users" ("id"),
    "batter_id"  uuid        NOT NULL REFERENCES "kgames"."users" ("id"),
    "unit_code"  VARCHAR     NOT NULL,
    "position"   jsonb       NOT NULL DEFAULT '{}'::jsonb,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS "kgames"."game_moves";
DROP TABLE IF EXISTS "kgames"."game_players";
DROP INDEX IF EXISTS "idx_games_state_code";
DROP TABLE IF EXISTS "kgames"."games";
DROP TABLE IF EXISTS "kgames"."users";
