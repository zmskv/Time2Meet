CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Core tables

CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    full_name       TEXT NOT NULL,
    phone           TEXT,
    role            TEXT NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT users_role_chk CHECK (role IN ('admin', 'organizer', 'attendee')),
    CONSTRAINT users_email_chk CHECK (position('@' in email) > 1)
);

CREATE TABLE IF NOT EXISTS user_profiles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL UNIQUE,
    avatar_url      TEXT,
    bio             TEXT,
    social_links    JSONB NOT NULL DEFAULT '{}'::jsonb,
    preferences     JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT user_profiles_user_fk
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS venues (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    address         TEXT NOT NULL,
    city            TEXT NOT NULL,
    country         TEXT NOT NULL DEFAULT 'RU',
    capacity        INT NOT NULL,
    contact_phone   TEXT,
    contact_email   TEXT,
    website         TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT venues_capacity_chk CHECK (capacity >= 0),
    CONSTRAINT venues_contact_email_chk CHECK (contact_email IS NULL OR position('@' in contact_email) > 1)
);

CREATE TABLE IF NOT EXISTS rooms (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id        UUID NOT NULL,
    name            TEXT NOT NULL,
    capacity        INT NOT NULL,
    floor           INT,
    equipment       JSONB NOT NULL DEFAULT '{}'::jsonb,
    hourly_rate     NUMERIC(12,2) NOT NULL DEFAULT 0,
    is_available    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT rooms_capacity_chk CHECK (capacity >= 0),
    CONSTRAINT rooms_hourly_rate_chk CHECK (hourly_rate >= 0),
    CONSTRAINT rooms_venue_fk
        FOREIGN KEY (venue_id) REFERENCES venues(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT rooms_venue_name_uniq UNIQUE (venue_id, name)
);

CREATE TABLE IF NOT EXISTS categories (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL UNIQUE,
    description     TEXT,
    icon            TEXT,
    slug            TEXT NOT NULL UNIQUE,
    parent_id       UUID,
    sort_order      INT NOT NULL DEFAULT 0,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT categories_parent_fk
        FOREIGN KEY (parent_id) REFERENCES categories(id)
        ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT categories_slug_chk CHECK (length(slug) > 0)
);

CREATE TABLE IF NOT EXISTS events (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organizer_id        UUID NOT NULL,
    title               TEXT NOT NULL,
    description         TEXT,
    status              TEXT NOT NULL,
    is_public           BOOLEAN NOT NULL DEFAULT TRUE,
    max_participants    INT,
    cover_image         TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT events_status_chk CHECK (status IN ('draft', 'published', 'cancelled', 'completed')),
    CONSTRAINT events_max_participants_chk CHECK (max_participants IS NULL OR max_participants >= 0),
    CONSTRAINT events_organizer_fk
        FOREIGN KEY (organizer_id) REFERENCES users(id)
        ON UPDATE CASCADE ON DELETE RESTRICT
);

-- N:M between events and categories
CREATE TABLE IF NOT EXISTS event_categories (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id        UUID NOT NULL,
    category_id     UUID NOT NULL,
    is_primary      BOOLEAN NOT NULL DEFAULT FALSE,
    source          TEXT NOT NULL DEFAULT 'manual',
    assigned_by     UUID,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT event_categories_event_category_uniq UNIQUE (event_id, category_id),
    CONSTRAINT event_categories_event_fk
        FOREIGN KEY (event_id) REFERENCES events(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT event_categories_category_fk
        FOREIGN KEY (category_id) REFERENCES categories(id)
        ON UPDATE CASCADE ON DELETE RESTRICT
    ,
    CONSTRAINT event_categories_assigned_by_fk
        FOREIGN KEY (assigned_by) REFERENCES users(id)
        ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS event_schedules (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id        UUID NOT NULL,
    room_id         UUID NOT NULL,
    start_time      TIMESTAMPTZ NOT NULL,
    end_time        TIMESTAMPTZ NOT NULL,
    status          TEXT NOT NULL,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT event_schedules_time_chk CHECK (end_time > start_time),
    CONSTRAINT event_schedules_status_chk CHECK (status IN ('planned', 'active', 'cancelled', 'done')),
    CONSTRAINT event_schedules_event_fk
        FOREIGN KEY (event_id) REFERENCES events(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT event_schedules_room_fk
        FOREIGN KEY (room_id) REFERENCES rooms(id)
        ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS ticket_types (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id         UUID NOT NULL,
    name             TEXT NOT NULL,
    price            NUMERIC(12,2) NOT NULL,
    currency         CHAR(3) NOT NULL DEFAULT 'RUB',
    quantity_total   INT NOT NULL,
    quantity_sold    INT NOT NULL DEFAULT 0,
    sale_start       TIMESTAMPTZ,
    sale_end         TIMESTAMPTZ,
    description      TEXT,
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ticket_types_price_chk CHECK (price >= 0),
    CONSTRAINT ticket_types_qty_chk CHECK (quantity_total >= 0 AND quantity_sold >= 0 AND quantity_sold <= quantity_total),
    CONSTRAINT ticket_types_sale_range_chk CHECK (
        sale_start IS NULL OR sale_end IS NULL OR sale_end > sale_start
    ),
    CONSTRAINT ticket_types_event_fk
        FOREIGN KEY (event_id) REFERENCES events(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT ticket_types_event_name_uniq UNIQUE (event_id, name)
);

CREATE TABLE IF NOT EXISTS tickets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_type_id  UUID NOT NULL,
    buyer_id        UUID NOT NULL,
    purchase_date   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status          TEXT NOT NULL,
    qr_code         TEXT NOT NULL UNIQUE,
    amount_paid     NUMERIC(12,2) NOT NULL DEFAULT 0,
    used_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT tickets_status_chk CHECK (status IN ('paid', 'refunded', 'void', 'used')),
    CONSTRAINT tickets_amount_paid_chk CHECK (amount_paid >= 0),
    CONSTRAINT tickets_used_at_chk CHECK (used_at IS NULL OR used_at >= purchase_date),
    CONSTRAINT tickets_type_fk
        FOREIGN KEY (ticket_type_id) REFERENCES ticket_types(id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT tickets_buyer_fk
        FOREIGN KEY (buyer_id) REFERENCES users(id)
        ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS registrations (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id              UUID NOT NULL,
    event_id             UUID NOT NULL,
    status               TEXT NOT NULL,
    registered_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    attendance_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    notes                TEXT,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT registrations_status_chk CHECK (status IN ('registered', 'cancelled', 'attended', 'no_show')),
    CONSTRAINT registrations_user_event_uniq UNIQUE (user_id, event_id),
    CONSTRAINT registrations_user_fk
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT registrations_event_fk
        FOREIGN KEY (event_id) REFERENCES events(id)
        ON UPDATE CASCADE ON DELETE CASCADE
);

-- Audit table
CREATE TABLE IF NOT EXISTS audit_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_name      TEXT NOT NULL,
    record_id       TEXT NOT NULL,
    action          TEXT NOT NULL,
    old_data        JSONB,
    new_data        JSONB,
    changed_by      UUID,
    ip_address      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT audit_action_chk CHECK (action IN ('INSERT', 'UPDATE', 'DELETE')),
    CONSTRAINT audit_changed_by_fk
        FOREIGN KEY (changed_by) REFERENCES users(id)
        ON UPDATE CASCADE ON DELETE SET NULL
);


