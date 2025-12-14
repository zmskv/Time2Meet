BEGIN;

CREATE OR REPLACE VIEW v_event_summary AS
SELECT
  e.id AS event_id,
  e.title,
  e.status,
  e.is_public,
  e.organizer_id,
  e.created_at,
  COUNT(DISTINCT r.id) FILTER (WHERE r.status IN ('registered','attended')) AS registrations_count,
  COUNT(t.id) FILTER (WHERE t.status IN ('paid','used')) AS tickets_sold,
  COALESCE(SUM(t.amount_paid) FILTER (WHERE t.status IN ('paid','used')), 0)::NUMERIC(14,2) AS revenue
FROM events e
LEFT JOIN registrations r ON r.event_id = e.id
LEFT JOIN ticket_types tt ON tt.event_id = e.id
LEFT JOIN tickets t ON t.ticket_type_id = tt.id
GROUP BY e.id, e.title, e.status, e.is_public, e.organizer_id, e.created_at;

CREATE OR REPLACE VIEW v_venue_analytics AS
SELECT
  v.id AS venue_id,
  v.name AS venue_name,
  v.city,
  COUNT(DISTINCT r.id) AS rooms_count,
  COUNT(DISTINCT es.id) FILTER (WHERE es.status IN ('planned','active','done')) AS schedules_count,
  COALESCE(SUM(EXTRACT(EPOCH FROM (es.end_time - es.start_time))) FILTER (WHERE es.status IN ('planned','active','done')), 0) / 3600.0 AS booked_hours
FROM venues v
LEFT JOIN rooms r ON r.venue_id = v.id
LEFT JOIN event_schedules es ON es.room_id = r.id
GROUP BY v.id, v.name, v.city;

CREATE OR REPLACE VIEW v_sales_by_period AS
SELECT
  DATE_TRUNC('month', t.purchase_date) AS period_month,
  COUNT(t.id) FILTER (WHERE t.status IN ('paid','used')) AS tickets_sold,
  COALESCE(SUM(t.amount_paid) FILTER (WHERE t.status IN ('paid','used')), 0)::NUMERIC(14,2) AS revenue
FROM tickets t
GROUP BY DATE_TRUNC('month', t.purchase_date)
ORDER BY period_month DESC;

CREATE OR REPLACE VIEW v_user_activity AS
SELECT
  u.id AS user_id,
  u.full_name,
  u.role,
  COUNT(DISTINCT e.id) AS events_organized,
  COUNT(DISTINCT r.id) FILTER (WHERE r.status IN ('registered','attended')) AS events_registered,
  COUNT(DISTINCT t.id) FILTER (WHERE t.status IN ('paid','used')) AS tickets_purchased
FROM users u
LEFT JOIN events e ON e.organizer_id = u.id
LEFT JOIN registrations r ON r.user_id = u.id
LEFT JOIN tickets t ON t.buyer_id = u.id
GROUP BY u.id, u.full_name, u.role;

COMMIT;

