-- aggregation trigger on tickets
DROP TRIGGER IF EXISTS trg_update_ticket_sales ON tickets;

-- audit triggers
DROP TRIGGER IF EXISTS trg_audit_registrations ON registrations;
DROP TRIGGER IF EXISTS trg_audit_tickets ON tickets;
DROP TRIGGER IF EXISTS trg_audit_ticket_types ON ticket_types;
DROP TRIGGER IF EXISTS trg_audit_event_schedules ON event_schedules;
DROP TRIGGER IF EXISTS trg_audit_event_categories ON event_categories;
DROP TRIGGER IF EXISTS trg_audit_events ON events;
DROP TRIGGER IF EXISTS trg_audit_categories ON categories;
DROP TRIGGER IF EXISTS trg_audit_rooms ON rooms;
DROP TRIGGER IF EXISTS trg_audit_venues ON venues;
DROP TRIGGER IF EXISTS trg_audit_user_profiles ON user_profiles;
DROP TRIGGER IF EXISTS trg_audit_users ON users;

-- updated_at triggers
DROP TRIGGER IF EXISTS trg_registrations_updated_at ON registrations;
DROP TRIGGER IF EXISTS trg_tickets_updated_at ON tickets;
DROP TRIGGER IF EXISTS trg_ticket_types_updated_at ON ticket_types;
DROP TRIGGER IF EXISTS trg_event_categories_updated_at ON event_categories;
DROP TRIGGER IF EXISTS trg_event_schedules_updated_at ON event_schedules;
DROP TRIGGER IF EXISTS trg_events_updated_at ON events;
DROP TRIGGER IF EXISTS trg_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS trg_rooms_updated_at ON rooms;
DROP TRIGGER IF EXISTS trg_venues_updated_at ON venues;
DROP TRIGGER IF EXISTS trg_user_profiles_updated_at ON user_profiles;
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;

-- functions
DROP FUNCTION IF EXISTS update_ticket_sales();
DROP FUNCTION IF EXISTS audit_trigger_func();
DROP FUNCTION IF EXISTS set_updated_at();


