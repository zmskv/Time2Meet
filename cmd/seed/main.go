package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"time2meet/internal/infrastructure/config"
	"time2meet/internal/infrastructure/persistence/postgres"
	"time2meet/pkg/logger"

	"github.com/brianvoe/gofakeit/v6"
	"go.uber.org/zap"
)

func main() {
	var (
		seed       = flag.Int64("seed", time.Now().UnixNano(), "rng seed")
		venuesN    = flag.Int("venues", 50, "venues count")
		roomsN     = flag.Int("rooms", 300, "rooms count")
		catsN      = flag.Int("categories", 50, "categories count")
		usersN     = flag.Int("users", 500, "users count")
		eventsN    = flag.Int("events", 1000, "events count")
		schedulesN = flag.Int("schedules", 1500, "event schedules count")
		typesN     = flag.Int("ticket-types", 2000, "ticket types count")
		ticketsN   = flag.Int("tickets", 5000, "tickets count")
		regsN      = flag.Int("registrations", 3000, "registrations count")
	)
	flag.Parse()

	log := logger.New()
	defer func() { _ = log.Sync() }()
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Error("config load failed", zap.Error(err))
		os.Exit(1)
	}

	db, err := postgres.NewDB(cfg.Database, log)
	if err != nil {
		log.Error("db connect failed", zap.Error(err))
		os.Exit(1)
	}
	defer db.Close()

	rnd := rand.New(rand.NewSource(*seed))
	gofakeit.Seed(*seed)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	log.Info("seeding start",
		zap.Int("venues", *venuesN),
		zap.Int("rooms", *roomsN),
		zap.Int("categories", *catsN),
		zap.Int("users", *usersN),
		zap.Int("events", *eventsN),
		zap.Int("schedules", *schedulesN),
		zap.Int("ticket_types", *typesN),
		zap.Int("tickets", *ticketsN),
		zap.Int("registrations", *regsN),
	)

	catIDs := make([]string, 0, *catsN)
	for i := 0; i < *catsN; i++ {
		name := fmt.Sprintf("%s %s", gofakeit.HackerNoun(), gofakeit.HackerVerb())
		slug := fmt.Sprintf("%s-%d", gofakeit.Username(), i+1)
		var id string
		err := db.QueryRowxContext(ctx, `
			INSERT INTO categories (name, description, icon, slug, sort_order, is_active)
			VALUES ($1, $2, $3, $4, $5, true)
			RETURNING id
		`, name, gofakeit.Sentence(10), "tag", slug, i).Scan(&id)
		if err != nil {
			log.Warn("insert category failed", zap.Int("i", i), zap.Error(err))
			continue
		}
		catIDs = append(catIDs, id)
	}

	venueIDs := make([]string, 0, *venuesN)
	for i := 0; i < *venuesN; i++ {
		addr := gofakeit.Address()
		var id string
		err := db.QueryRowxContext(ctx, `
			INSERT INTO venues (name, address, city, country, capacity, contact_phone, contact_email, website, is_active)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,true)
			RETURNING id
		`,
			fmt.Sprintf("%s Center", gofakeit.Company()),
			addr.Address,
			addr.City,
			"RU",
			rnd.Intn(5000)+100,
			gofakeit.Phone(),
			gofakeit.Email(),
			gofakeit.URL(),
		).Scan(&id)
		if err != nil {
			log.Warn("insert venue failed", zap.Int("i", i), zap.Error(err))
			continue
		}
		venueIDs = append(venueIDs, id)
	}

	roomIDs := make([]string, 0, *roomsN)
	for i := 0; i < *roomsN; i++ {
		venueID := venueIDs[rnd.Intn(len(venueIDs))]
		name := fmt.Sprintf("Room-%d", i+1)
		capacity := rnd.Intn(800) + 20
		floor := rnd.Intn(8)
		equipment := fmt.Sprintf(`{"projector":%t,"sound":%t,"wifi":true}`, rnd.Intn(2) == 0, rnd.Intn(2) == 0)
		hourly := fmt.Sprintf("%.2f", float64(rnd.Intn(20000)+500)/100.0)
		var id string
		err := db.QueryRowxContext(ctx, `
			INSERT INTO rooms (venue_id, name, capacity, floor, equipment, hourly_rate, is_available)
			VALUES ($1,$2,$3,$4,$5::jsonb,$6,true)
			RETURNING id
		`, venueID, name, capacity, floor, equipment, hourly).Scan(&id)
		if err != nil {
			log.Warn("insert room failed", zap.Int("i", i), zap.Error(err))
			continue
		}
		roomIDs = append(roomIDs, id)
	}

	userIDs := make([]string, 0, *usersN)
	organizerIDs := make([]string, 0, *usersN/5)
	for i := 0; i < *usersN; i++ {
		role := "attendee"
		if i == 0 {
			role = "admin"
		} else if rnd.Intn(5) == 0 {
			role = "organizer"
		}
		email := gofakeit.Email()
		var id string
		err := db.QueryRowxContext(ctx, `
			INSERT INTO users (email, password_hash, full_name, phone, role, is_active)
			VALUES ($1,$2,$3,$4,$5,true)
			RETURNING id
		`, email, "seed_hash", gofakeit.Name(), gofakeit.Phone(), role).Scan(&id)
		if err != nil {
			log.Warn("insert user failed", zap.Int("i", i), zap.Error(err))
			continue
		}
		userIDs = append(userIDs, id)
		if role == "organizer" {
			organizerIDs = append(organizerIDs, id)
		}
		_, _ = db.ExecContext(ctx, `
			INSERT INTO user_profiles (user_id, avatar_url, bio, social_links, preferences)
			VALUES ($1,$2,$3,$4::jsonb,$5::jsonb)
			ON CONFLICT (user_id) DO NOTHING
		`, id, gofakeit.ImageURL(128, 128), gofakeit.Sentence(12), `{"tg":"@`+gofakeit.Username()+`"}`, `{"lang":"ru"}`)
	}

	eventIDs := make([]string, 0, *eventsN)
	for i := 0; i < *eventsN; i++ {
		org := organizerIDs[rnd.Intn(len(organizerIDs))]
		status := "published"
		if rnd.Intn(20) == 0 {
			status = "draft"
		}
		if rnd.Intn(50) == 0 {
			status = "cancelled"
		}
		var maxp any
		if rnd.Intn(3) != 0 {
			v := rnd.Intn(2000) + 20
			maxp = v
		}
		var id string
		err := db.QueryRowxContext(ctx, `
			INSERT INTO events (organizer_id, title, description, status, is_public, max_participants, cover_image)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
			RETURNING id
		`, org, fmt.Sprintf("%s Meetup %d", gofakeit.BuzzWord(), i+1), gofakeit.Paragraph(1, 4, 12, " "), status, true, maxp, gofakeit.ImageURL(640, 360)).Scan(&id)
		if err != nil {
			log.Warn("insert event failed", zap.Int("i", i), zap.Error(err))
			continue
		}
		eventIDs = append(eventIDs, id)

		for k := 0; k < 1+rnd.Intn(3); k++ {
			cid := catIDs[rnd.Intn(len(catIDs))]
			_, _ = db.ExecContext(ctx, `
				INSERT INTO event_categories (event_id, category_id, is_primary)
				VALUES ($1,$2,$3)
				ON CONFLICT DO NOTHING
			`, id, cid, k == 0)
		}
	}

	for i := 0; i < *schedulesN; i++ {
		eid := eventIDs[rnd.Intn(len(eventIDs))]
		rid := roomIDs[rnd.Intn(len(roomIDs))]
		start := time.Now().Add(time.Duration(rnd.Intn(120*24)) * time.Hour).UTC()
		end := start.Add(time.Duration(1+rnd.Intn(6)) * time.Hour)
		status := "planned"
		_, err := db.ExecContext(ctx, `
			INSERT INTO event_schedules (event_id, room_id, start_time, end_time, status, notes)
			VALUES ($1,$2,$3,$4,$5,$6)
		`, eid, rid, start, end, status, gofakeit.Sentence(8))
		if err != nil {
			log.Warn("insert schedule failed", zap.Int("i", i), zap.Error(err))
		}
	}

	type typeInfo struct {
		ID        string
		Remaining int
	}
	types := make([]typeInfo, 0, *typesN)
	for i := 0; i < *typesN; i++ {
		eid := eventIDs[rnd.Intn(len(eventIDs))]
		qty := rnd.Intn(300) + 20
		price := fmt.Sprintf("%.2f", float64(rnd.Intn(500000)+10000)/100.0)
		name := []string{"Standard", "VIP", "Student", "Early Bird"}[rnd.Intn(4)]
		var id string
		err := db.QueryRowxContext(ctx, `
			INSERT INTO ticket_types (event_id, name, price, currency, quantity_total, quantity_sold, sale_start, sale_end, description, is_active)
			VALUES ($1,$2,$3,'RUB',$4,0,$5,$6,$7,true)
			RETURNING id
		`, eid, fmt.Sprintf("%s-%d", name, i+1), price, qty, time.Now().Add(-24*time.Hour).UTC(), time.Now().Add(90*24*time.Hour).UTC(), gofakeit.Sentence(10)).Scan(&id)
		if err != nil {
			log.Warn("insert ticket_type failed", zap.Int("i", i), zap.Error(err))
			continue
		}
		types = append(types, typeInfo{ID: id, Remaining: qty})
	}

	for i := 0; i < *ticketsN; i++ {
		var idx int
		for tries := 0; tries < 50; tries++ {
			idx = rnd.Intn(len(types))
			if types[idx].Remaining > 0 {
				break
			}
		}
		if types[idx].Remaining <= 0 {
			break
		}
		tt := &types[idx]
		buyer := userIDs[rnd.Intn(len(userIDs))]
		qr := fmt.Sprintf("QR-%d-%d", time.Now().UnixNano(), i)
		amount := fmt.Sprintf("%.2f", float64(rnd.Intn(500000)+10000)/100.0)
		_, err := db.ExecContext(ctx, `
			INSERT INTO tickets (ticket_type_id, buyer_id, purchase_date, status, qr_code, amount_paid)
			VALUES ($1,$2,$3,'paid',$4,$5)
		`, tt.ID, buyer, time.Now().Add(-time.Duration(rnd.Intn(60*24))*time.Hour).UTC(), qr, amount)
		if err != nil {
			log.Warn("insert ticket failed", zap.Int("i", i), zap.Error(err))
			continue
		}
		tt.Remaining--
	}

	createdRegs := 0
	for i := 0; i < *regsN && createdRegs < *regsN; i++ {
		uid := userIDs[rnd.Intn(len(userIDs))]
		eid := eventIDs[rnd.Intn(len(eventIDs))]
		status := "registered"
		_, err := db.ExecContext(ctx, `
			INSERT INTO registrations (user_id, event_id, status, attendance_confirmed, notes)
			VALUES ($1,$2,$3,false,$4)
			ON CONFLICT (user_id, event_id) DO NOTHING
		`, uid, eid, status, gofakeit.Sentence(6))
		if err != nil {
			log.Warn("insert registration failed", zap.Int("i", i), zap.Error(err))
			continue
		}
		createdRegs++
	}

	log.Info("seeding done")
}
