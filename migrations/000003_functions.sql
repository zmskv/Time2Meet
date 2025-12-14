BEGIN;

-- Scalar: total revenue by event (paid+used)
CREATE OR REPLACE FUNCTION calculate_event_revenue(p_event_id UUID)
RETURNS NUMERIC(14,2)
LANGUAGE sql
STABLE
AS $$
  SELECT COALESCE(SUM(t.amount_paid), 0)::NUMERIC(14,2)
  FROM ticket_types tt
  JOIN tickets t ON t.ticket_type_id = tt.id
  WHERE tt.event_id = p_event_id
    AND t.status IN ('paid', 'used');
$$;

-- Scalar: organizer rating (average attendance rate across their events)
CREATE OR REPLACE FUNCTION get_organizer_rating(p_user_id UUID)
RETURNS NUMERIC(6,4)
LANGUAGE sql
STABLE
AS $$
  WITH event_stats AS (
    SELECT
      e.id AS event_id,
      COUNT(t.id) FILTER (WHERE t.status IN ('paid','used'))::NUMERIC AS sold,
      COUNT(t.id) FILTER (WHERE t.status = 'used')::NUMERIC AS used
    FROM events e
    LEFT JOIN ticket_types tt ON tt.event_id = e.id
    LEFT JOIN tickets t ON t.ticket_type_id = tt.id
    WHERE e.organizer_id = p_user_id
    GROUP BY e.id
  )
  SELECT COALESCE(AVG(CASE WHEN sold > 0 THEN used / sold ELSE NULL END), 0)::NUMERIC(6,4)
  FROM event_stats;
$$;

-- Scalar: venue utilization for a given month of current year (0..1)
CREATE OR REPLACE FUNCTION get_venue_utilization(p_venue_id UUID, p_month INT)
RETURNS NUMERIC(6,4)
LANGUAGE sql
STABLE
AS $$
  WITH bounds AS (
    SELECT
      make_timestamptz(EXTRACT(YEAR FROM NOW())::INT, p_month, 1, 0, 0, 0, 'UTC') AS start_ts,
      (make_timestamptz(EXTRACT(YEAR FROM NOW())::INT, p_month, 1, 0, 0, 0, 'UTC') + INTERVAL '1 month') AS end_ts
  ),
  sched AS (
    SELECT
      GREATEST(es.start_time, b.start_ts) AS s,
      LEAST(es.end_time, b.end_ts) AS e
    FROM bounds b
    JOIN rooms r ON r.venue_id = p_venue_id
    JOIN event_schedules es ON es.room_id = r.id
    WHERE es.start_time < b.end_ts
      AND es.end_time > b.start_ts
      AND es.status IN ('planned','active','done')
  )
  SELECT
    COALESCE(
      (SELECT SUM(EXTRACT(EPOCH FROM (e - s))) FROM sched) /
      NULLIF(EXTRACT(EPOCH FROM ((SELECT end_ts FROM bounds) - (SELECT start_ts FROM bounds))), 0),
      0
    )::NUMERIC(6,4);
$$;

-- Table: sales report by event for period
CREATE OR REPLACE FUNCTION get_sales_report(p_start DATE, p_end DATE)
RETURNS TABLE(
  event_id UUID,
  event_title TEXT,
  tickets_sold BIGINT,
  revenue NUMERIC(14,2),
  unique_buyers BIGINT
)
LANGUAGE sql
STABLE
AS $$
  SELECT
    e.id,
    e.title,
    COUNT(t.id) FILTER (WHERE t.status IN ('paid','used')) AS tickets_sold,
    COALESCE(SUM(t.amount_paid) FILTER (WHERE t.status IN ('paid','used')), 0)::NUMERIC(14,2) AS revenue,
    COUNT(DISTINCT t.buyer_id) FILTER (WHERE t.status IN ('paid','used')) AS unique_buyers
  FROM events e
  LEFT JOIN ticket_types tt ON tt.event_id = e.id
  LEFT JOIN tickets t ON t.ticket_type_id = tt.id
               AND t.purchase_date >= p_start::timestamptz
               AND t.purchase_date < (p_end + 1)::timestamptz
  GROUP BY e.id, e.title
  ORDER BY revenue DESC;
$$;

-- Table: popular events within last N days
CREATE OR REPLACE FUNCTION get_popular_events(p_limit INT, p_days INT)
RETURNS TABLE(
  event_id UUID,
  title TEXT,
  registrations BIGINT,
  tickets_sold BIGINT
)
LANGUAGE sql
STABLE
AS $$
  WITH period AS (
    SELECT NOW() - (p_days::TEXT || ' days')::INTERVAL AS start_ts
  ),
  reg AS (
    SELECT r.event_id, COUNT(*) AS registrations
    FROM registrations r, period p
    WHERE r.registered_at >= p.start_ts
      AND r.status IN ('registered','attended')
    GROUP BY r.event_id
  ),
  sold AS (
    SELECT tt.event_id, COUNT(t.id) AS tickets_sold
    FROM ticket_types tt
    JOIN tickets t ON t.ticket_type_id = tt.id
    JOIN period p ON TRUE
    WHERE t.purchase_date >= p.start_ts
      AND t.status IN ('paid','used')
    GROUP BY tt.event_id
  )
  SELECT
    e.id,
    e.title,
    COALESCE(reg.registrations, 0) AS registrations,
    COALESCE(sold.tickets_sold, 0) AS tickets_sold
  FROM events e
  LEFT JOIN reg ON reg.event_id = e.id
  LEFT JOIN sold ON sold.event_id = e.id
  ORDER BY (COALESCE(reg.registrations, 0) + COALESCE(sold.tickets_sold, 0)) DESC
  LIMIT GREATEST(p_limit, 1);
$$;

-- Table: attendance stats for event by ticket type
CREATE OR REPLACE FUNCTION get_attendance_stats(p_event_id UUID)
RETURNS TABLE(
  ticket_type TEXT,
  sold BIGINT,
  used BIGINT,
  attendance_rate NUMERIC(6,4)
)
LANGUAGE sql
STABLE
AS $$
  SELECT
    tt.name AS ticket_type,
    COUNT(t.id) FILTER (WHERE t.status IN ('paid','used')) AS sold,
    COUNT(t.id) FILTER (WHERE t.status = 'used') AS used,
    COALESCE(
      (COUNT(t.id) FILTER (WHERE t.status = 'used')::NUMERIC) /
      NULLIF(COUNT(t.id) FILTER (WHERE t.status IN ('paid','used'))::NUMERIC, 0),
      0
    )::NUMERIC(6,4) AS attendance_rate
  FROM ticket_types tt
  LEFT JOIN tickets t ON t.ticket_type_id = tt.id
  WHERE tt.event_id = p_event_id
  GROUP BY tt.name
  ORDER BY tt.name;
$$;

COMMIT;

