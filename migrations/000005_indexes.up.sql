-- WHERE/ORDER BY
CREATE INDEX IF NOT EXISTS idx_events_status_created ON events(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_tickets_purchase_date ON tickets(purchase_date);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);

-- JOIN helpers
CREATE INDEX IF NOT EXISTS idx_event_schedules_event_room ON event_schedules(event_id, room_id);
CREATE INDEX IF NOT EXISTS idx_tickets_buyer ON tickets(buyer_id);
CREATE INDEX IF NOT EXISTS idx_tickets_type ON tickets(ticket_type_id);
CREATE INDEX IF NOT EXISTS idx_ticket_types_event ON ticket_types(event_id);
CREATE INDEX IF NOT EXISTS idx_rooms_venue ON rooms(venue_id);
CREATE INDEX IF NOT EXISTS idx_registrations_event ON registrations(event_id);

-- Audit queries
CREATE INDEX IF NOT EXISTS idx_audit_log_table_time ON audit_log(table_name, created_at DESC);


