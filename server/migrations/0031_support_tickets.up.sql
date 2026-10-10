-- Support tickets (DD-234): in-app reports answered from the admin panel.
CREATE TABLE ops.support_ticket (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id   UUID NOT NULL REFERENCES auth.account(id) ON DELETE CASCADE,
    character_id UUID REFERENCES character.character(id) ON DELETE SET NULL,
    category     VARCHAR(16)  NOT NULL CHECK (category IN ('bug', 'account', 'abuse', 'feedback', 'other')),
    subject      VARCHAR(120) NOT NULL,
    body         TEXT         NOT NULL CHECK (length(body) <= 4000),
    client_info  JSONB        NOT NULL DEFAULT '{}'::jsonb,
    state        VARCHAR(16)  NOT NULL DEFAULT 'open' CHECK (state IN ('open', 'answered', 'closed')),
    reply        TEXT CHECK (length(reply) <= 4000),
    replied_by   UUID REFERENCES ops.admin_user(id) ON DELETE SET NULL,
    replied_at   TIMESTAMPTZ,
    closed_at    TIMESTAMPTZ,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX ix_support_ticket_queue   ON ops.support_ticket (state, created_at);
CREATE INDEX ix_support_ticket_account ON ops.support_ticket (account_id, created_at DESC);
