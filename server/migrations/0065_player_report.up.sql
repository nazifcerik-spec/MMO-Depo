-- Player reports (DD-355): one player flags another; operators work the queue.
CREATE TABLE ops.player_report (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reporter_account   UUID NOT NULL REFERENCES identity.account(id) ON DELETE CASCADE,
    reporter_character UUID REFERENCES character.character(id) ON DELETE SET NULL,
    target_account     UUID REFERENCES identity.account(id) ON DELETE SET NULL,
    target_character   UUID REFERENCES character.character(id) ON DELETE SET NULL,
    target_name        TEXT         NOT NULL,
    category           VARCHAR(16)  NOT NULL CHECK (category IN ('harassment', 'cheating', 'name', 'spam', 'scam', 'other')),
    note               TEXT         NOT NULL DEFAULT '' CHECK (char_length(note) <= 500),
    chat_message_id    BIGINT,
    chat_excerpt       TEXT,
    state              VARCHAR(16)  NOT NULL DEFAULT 'open' CHECK (state IN ('open', 'actioned', 'dismissed')),
    resolved_by        UUID REFERENCES ops.admin_user(id) ON DELETE SET NULL,
    resolved_at        TIMESTAMPTZ,
    resolution_note    TEXT CHECK (char_length(resolution_note) <= 500),
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX ix_player_report_queue  ON ops.player_report (state, created_at);
CREATE INDEX ix_player_report_target ON ops.player_report (target_account, created_at DESC);
-- One open report per reporter, target and reason: repeats change nothing.
CREATE UNIQUE INDEX ux_player_report_open ON ops.player_report (reporter_account, target_character, category)
    WHERE state = 'open';
