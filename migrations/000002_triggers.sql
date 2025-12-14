BEGIN;

-- Updated_at helper
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Audit trigger function
CREATE OR REPLACE FUNCTION audit_trigger_func()
RETURNS TRIGGER AS $$
DECLARE
  v_old JSONB;
  v_new JSONB;
  v_record_id TEXT;
  v_user_id UUID;
  v_ip TEXT;
BEGIN
  IF (TG_OP = 'DELETE') THEN
    v_old := to_jsonb(OLD);
    v_new := NULL;
  ELSIF (TG_OP = 'UPDATE') THEN
    v_old := to_jsonb(OLD);
    v_new := to_jsonb(NEW);
  ELSE
    v_old := NULL;
    v_new := to_jsonb(NEW);
  END IF;

  -- app.user_id and app.ip are optional; backend can set them per-connection:
  --   SET LOCAL app.user_id = '123';
  --   SET LOCAL app.ip = '1.2.3.4';
  v_user_id := NULLIF(current_setting('app.user_id', true), '')::UUID;
  v_ip := NULLIF(current_setting('app.ip', true), '');

  v_record_id := COALESCE(
    v_new->>'id',
    v_old->>'id',
    CASE
      WHEN v_new ? 'event_id' AND v_new ? 'category_id' THEN (v_new->>'event_id') || ':' || (v_new->>'category_id')
      WHEN v_old ? 'event_id' AND v_old ? 'category_id' THEN (v_old->>'event_id') || ':' || (v_old->>'category_id')
      ELSE NULL
    END,
    'unknown'
  );

  INSERT INTO audit_log(table_name, record_id, action, old_data, new_data, changed_by, ip_address, created_at)
  VALUES (TG_TABLE_NAME, v_record_id, TG_OP, v_old, v_new, v_user_id, v_ip, NOW());

  IF (TG_OP = 'DELETE') THEN
    RETURN OLD;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Ticket sales aggregation trigger
CREATE OR REPLACE FUNCTION update_ticket_sales()
RETURNS TRIGGER AS $$
BEGIN
  IF (TG_OP = 'INSERT') THEN
    IF NEW.status IN ('paid', 'used') THEN
      UPDATE ticket_types
      SET quantity_sold = quantity_sold + 1,
          updated_at = NOW()
      WHERE id = NEW.ticket_type_id;
    END IF;
    RETURN NEW;
  ELSIF (TG_OP = 'DELETE') THEN
    IF OLD.status IN ('paid', 'used') THEN
      UPDATE ticket_types
      SET quantity_sold = GREATEST(quantity_sold - 1, 0),
          updated_at = NOW()
      WHERE id = OLD.ticket_type_id;
    END IF;
    RETURN OLD;
  ELSE
    -- UPDATE
    IF OLD.ticket_type_id <> NEW.ticket_type_id THEN
      IF OLD.status IN ('paid', 'used') THEN
        UPDATE ticket_types
        SET quantity_sold = GREATEST(quantity_sold - 1, 0),
            updated_at = NOW()
        WHERE id = OLD.ticket_type_id;
      END IF;
      IF NEW.status IN ('paid', 'used') THEN
        UPDATE ticket_types
        SET quantity_sold = quantity_sold + 1,
            updated_at = NOW()
        WHERE id = NEW.ticket_type_id;
      END IF;
      RETURN NEW;
    END IF;

    IF OLD.status NOT IN ('paid', 'used') AND NEW.status IN ('paid', 'used') THEN
      UPDATE ticket_types
      SET quantity_sold = quantity_sold + 1,
          updated_at = NOW()
      WHERE id = NEW.ticket_type_id;
    ELSIF OLD.status IN ('paid', 'used') AND NEW.status NOT IN ('paid', 'used') THEN
      UPDATE ticket_types
      SET quantity_sold = GREATEST(quantity_sold - 1, 0),
          updated_at = NOW()
      WHERE id = NEW.ticket_type_id;
    END IF;
    RETURN NEW;
  END IF;
END;
$$ LANGUAGE plpgsql;

-- updated_at triggers
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_user_profiles_updated_at ON user_profiles;
CREATE TRIGGER trg_user_profiles_updated_at
BEFORE UPDATE ON user_profiles
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_venues_updated_at ON venues;
CREATE TRIGGER trg_venues_updated_at
BEFORE UPDATE ON venues
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_rooms_updated_at ON rooms;
CREATE TRIGGER trg_rooms_updated_at
BEFORE UPDATE ON rooms
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_categories_updated_at ON categories;
CREATE TRIGGER trg_categories_updated_at
BEFORE UPDATE ON categories
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_events_updated_at ON events;
CREATE TRIGGER trg_events_updated_at
BEFORE UPDATE ON events
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_event_schedules_updated_at ON event_schedules;
CREATE TRIGGER trg_event_schedules_updated_at
BEFORE UPDATE ON event_schedules
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_event_categories_updated_at ON event_categories;
CREATE TRIGGER trg_event_categories_updated_at
BEFORE UPDATE ON event_categories
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_ticket_types_updated_at ON ticket_types;
CREATE TRIGGER trg_ticket_types_updated_at
BEFORE UPDATE ON ticket_types
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_tickets_updated_at ON tickets;
CREATE TRIGGER trg_tickets_updated_at
BEFORE UPDATE ON tickets
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_registrations_updated_at ON registrations;
CREATE TRIGGER trg_registrations_updated_at
BEFORE UPDATE ON registrations
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- audit triggers (for all core tables, including event_categories)
DROP TRIGGER IF EXISTS trg_audit_users ON users;
CREATE TRIGGER trg_audit_users
AFTER INSERT OR UPDATE OR DELETE ON users
FOR EACH ROW EXECUTE FUNCTION audit_trigger_func();

DROP TRIGGER IF EXISTS trg_audit_user_profiles ON user_profiles;
CREATE TRIGGER trg_audit_user_profiles
AFTER INSERT OR UPDATE OR DELETE ON user_profiles
FOR EACH ROW EXECUTE FUNCTION audit_trigger_func();

DROP TRIGGER IF EXISTS trg_audit_venues ON venues;
CREATE TRIGGER trg_audit_venues
AFTER INSERT OR UPDATE OR DELETE ON venues
FOR EACH ROW EXECUTE FUNCTION audit_trigger_func();

DROP TRIGGER IF EXISTS trg_audit_rooms ON rooms;
CREATE TRIGGER trg_audit_rooms
AFTER INSERT OR UPDATE OR DELETE ON rooms
FOR EACH ROW EXECUTE FUNCTION audit_trigger_func();

DROP TRIGGER IF EXISTS trg_audit_categories ON categories;
CREATE TRIGGER trg_audit_categories
AFTER INSERT OR UPDATE OR DELETE ON categories
FOR EACH ROW EXECUTE FUNCTION audit_trigger_func();

DROP TRIGGER IF EXISTS trg_audit_events ON events;
CREATE TRIGGER trg_audit_events
AFTER INSERT OR UPDATE OR DELETE ON events
FOR EACH ROW EXECUTE FUNCTION audit_trigger_func();

DROP TRIGGER IF EXISTS trg_audit_event_categories ON event_categories;
CREATE TRIGGER trg_audit_event_categories
AFTER INSERT OR UPDATE OR DELETE ON event_categories
FOR EACH ROW EXECUTE FUNCTION audit_trigger_func();

DROP TRIGGER IF EXISTS trg_audit_event_schedules ON event_schedules;
CREATE TRIGGER trg_audit_event_schedules
AFTER INSERT OR UPDATE OR DELETE ON event_schedules
FOR EACH ROW EXECUTE FUNCTION audit_trigger_func();

DROP TRIGGER IF EXISTS trg_audit_ticket_types ON ticket_types;
CREATE TRIGGER trg_audit_ticket_types
AFTER INSERT OR UPDATE OR DELETE ON ticket_types
FOR EACH ROW EXECUTE FUNCTION audit_trigger_func();

DROP TRIGGER IF EXISTS trg_audit_tickets ON tickets;
CREATE TRIGGER trg_audit_tickets
AFTER INSERT OR UPDATE OR DELETE ON tickets
FOR EACH ROW EXECUTE FUNCTION audit_trigger_func();

DROP TRIGGER IF EXISTS trg_audit_registrations ON registrations;
CREATE TRIGGER trg_audit_registrations
AFTER INSERT OR UPDATE OR DELETE ON registrations
FOR EACH ROW EXECUTE FUNCTION audit_trigger_func();

-- aggregation trigger on tickets
DROP TRIGGER IF EXISTS trg_update_ticket_sales ON tickets;
CREATE TRIGGER trg_update_ticket_sales
AFTER INSERT OR UPDATE OR DELETE ON tickets
FOR EACH ROW EXECUTE FUNCTION update_ticket_sales();

COMMIT;

