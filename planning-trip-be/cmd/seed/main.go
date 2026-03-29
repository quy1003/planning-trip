package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"planning-trip-be/internal/config"
	"planning-trip-be/internal/database"
	"planning-trip-be/internal/model"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type seedState struct {
	users      []model.User
	trips      []model.Trip
	places     []model.Place
	tripPlaces []model.TripPlace
	albums     []model.Album
	photos     []model.Photo
	comments   []model.Comment
}

func main() {
	env, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	db, err := database.OpenPostgres(env.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("migration failed before seed: %v", err)
	}

	seedID := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seedID))

	if err := db.Transaction(func(tx *gorm.DB) error {
		state := &seedState{}

		if err := seedUsers(tx, state, seedID); err != nil {
			return err
		}
		if err := seedTripsAndMembers(tx, state, rng); err != nil {
			return err
		}
		if err := seedPlacesAndTripPlaces(tx, state, rng); err != nil {
			return err
		}
		if err := seedSchedule(tx, state, rng); err != nil {
			return err
		}
		if err := seedAlbumsAndPhotos(tx, state, rng); err != nil {
			return err
		}
		if err := seedComments(tx, state, rng); err != nil {
			return err
		}
		if err := seedReactions(tx, state, rng); err != nil {
			return err
		}
		if err := seedPhotoTags(tx, state, rng); err != nil {
			return err
		}
		if err := seedNotifications(tx, state, rng); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("seed failed: %v", err)
	}

	log.Printf("seed completed successfully (seed_id=%d)", seedID)
}

func seedUsers(tx *gorm.DB, state *seedState, seedID int64) error {
	users := []model.User{
		{FullName: "Quy Nguyen", Email: fmt.Sprintf("quy.%d.1@example.com", seedID), PasswordHash: "hashed_password_1", AvatarURL: "https://picsum.photos/seed/quy/300", Bio: "Team planner"},
		{FullName: "Linh Tran", Email: fmt.Sprintf("linh.%d.2@example.com", seedID), PasswordHash: "hashed_password_2", AvatarURL: "https://picsum.photos/seed/linh/300", Bio: "Food hunter"},
		{FullName: "Minh Le", Email: fmt.Sprintf("minh.%d.3@example.com", seedID), PasswordHash: "hashed_password_3", AvatarURL: "https://picsum.photos/seed/minh/300", Bio: "Photographer"},
		{FullName: "Anh Pham", Email: fmt.Sprintf("anh.%d.4@example.com", seedID), PasswordHash: "hashed_password_4", AvatarURL: "https://picsum.photos/seed/anh/300", Bio: "Logistics"},
		{FullName: "Thao Vu", Email: fmt.Sprintf("thao.%d.5@example.com", seedID), PasswordHash: "hashed_password_5", AvatarURL: "https://picsum.photos/seed/thao/300", Bio: "Nature lover"},
	}

	for i := range users {
		if err := tx.Create(&users[i]).Error; err != nil {
			return err
		}
	}
	state.users = users
	return nil
}

func seedTripsAndMembers(tx *gorm.DB, state *seedState, rng *rand.Rand) error {
	start := time.Now().UTC().AddDate(0, 1, 0)
	end := start.AddDate(0, 0, 3)
	tripTitles := []string{"Nha Trang 4N3D", "Da Lat Weekend", "Phu Quoc Chill Tour"}

	for i, title := range tripTitles {
		creator := state.users[i%len(state.users)]
		trip := model.Trip{
			CreatorID:   creator.ID,
			Title:       title,
			Description: "Trip planned by group with map, schedule, comments, and album.",
			StartDate:   &start,
			EndDate:     &end,
			Visibility:  "group",
			Status:      "published",
		}
		if err := tx.Create(&trip).Error; err != nil {
			return err
		}
		state.trips = append(state.trips, trip)

		for _, u := range state.users {
			role := "viewer"
			if u.ID == creator.ID {
				role = "owner"
			} else if rng.Intn(3) == 0 {
				role = "editor"
			}

			member := model.TripMember{
				TripID:   trip.ID,
				UserID:   u.ID,
				Role:     role,
				JoinedAt: time.Now().UTC(),
			}
			if err := tx.Create(&member).Error; err != nil {
				return err
			}
		}

		start = start.AddDate(0, 0, 5)
		end = end.AddDate(0, 0, 5)
	}

	return nil
}

func seedPlacesAndTripPlaces(tx *gorm.DB, state *seedState, rng *rand.Rand) error {
	basePlaces := []struct {
		Name    string
		Address string
		Lat     float64
		Lng     float64
	}{
		{"Thap Ba Ponagar", "Nha Trang", 12.2651, 109.1955},
		{"VinWonders Nha Trang", "Hon Tre, Nha Trang", 12.2270, 109.2436},
		{"Hon Chong", "Nha Trang", 12.2796, 109.2013},
		{"Cho Dam", "Nha Trang", 12.2548, 109.1897},
		{"Ba Ho Waterfall", "Ninh Ich", 12.4978, 109.1954},
	}

	for _, p := range basePlaces {
		createdBy := state.users[rng.Intn(len(state.users))].ID
		place := model.Place{
			Name:      p.Name,
			Address:   p.Address,
			Lat:       p.Lat,
			Lng:       p.Lng,
			CreatedBy: &createdBy,
		}
		if err := tx.Create(&place).Error; err != nil {
			return err
		}
		state.places = append(state.places, place)
	}

	for _, trip := range state.trips {
		for i := 0; i < 3; i++ {
			place := state.places[rng.Intn(len(state.places))]
			createdBy := state.users[rng.Intn(len(state.users))].ID
			tripPlace := model.TripPlace{
				TripID:        trip.ID,
				PlaceID:       place.ID,
				Title:         place.Name,
				Note:          "Suggested stop in itinerary",
				DayIndex:      i + 1,
				OrderIndex:    i,
				EstimatedCost: float64(100000 + rng.Intn(400000)),
				CreatedBy:     createdBy,
			}
			if err := tx.Create(&tripPlace).Error; err != nil {
				return err
			}
			state.tripPlaces = append(state.tripPlaces, tripPlace)
		}
	}

	return nil
}

func seedSchedule(tx *gorm.DB, state *seedState, _ *rand.Rand) error {
	for _, tp := range state.tripPlaces {
		start := time.Now().UTC().Add(time.Duration(tp.DayIndex) * 24 * time.Hour)
		start = time.Date(start.Year(), start.Month(), start.Day(), 8+tp.OrderIndex*2, 0, 0, 0, time.UTC)
		end := start.Add(90 * time.Minute)

		item := model.TripScheduleItem{
			TripID:      tp.TripID,
			TripPlaceID: &tp.ID,
			Title:       "Visit " + tp.Title,
			Description: "Group activity and check-in",
			StartTime:   &start,
			EndTime:     &end,
			DayIndex:    tp.DayIndex,
			OrderIndex:  tp.OrderIndex,
			CreatedBy:   tp.CreatedBy,
		}
		if err := tx.Create(&item).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedAlbumsAndPhotos(tx *gorm.DB, state *seedState, rng *rand.Rand) error {
	for _, trip := range state.trips {
		creator := state.users[rng.Intn(len(state.users))]
		album := model.Album{
			TripID:      trip.ID,
			Title:       "Album " + trip.Title,
			Description: "Shared memories",
			CreatedBy:   creator.ID,
		}
		if err := tx.Create(&album).Error; err != nil {
			return err
		}
		state.albums = append(state.albums, album)

		for i := 0; i < 4; i++ {
			uploader := state.users[rng.Intn(len(state.users))]
			tripPlace := state.tripPlaces[rng.Intn(len(state.tripPlaces))]
			if tripPlace.TripID != trip.ID {
				continue
			}

			takenAt := time.Now().UTC().Add(-time.Duration(rng.Intn(360)) * time.Hour)
			photo := model.Photo{
				AlbumID:      album.ID,
				TripID:       trip.ID,
				UploaderID:   uploader.ID,
				URL:          fmt.Sprintf("https://picsum.photos/seed/%s-%d/1200/800", trip.ID, i),
				ThumbnailURL: fmt.Sprintf("https://picsum.photos/seed/%s-thumb-%d/400/300", trip.ID, i),
				Caption:      "Great moment with the team",
				TakenAt:      &takenAt,
				TripPlaceID:  &tripPlace.ID,
			}
			if err := tx.Create(&photo).Error; err != nil {
				return err
			}
			state.photos = append(state.photos, photo)
		}
	}
	return nil
}

func seedComments(tx *gorm.DB, state *seedState, rng *rand.Rand) error {
	messages := []string{
		"Looks great!",
		"Can we add a coffee stop here?",
		"I love this place",
		"Schedule seems tight but doable",
	}

	for _, trip := range state.trips {
		for i := 0; i < 3; i++ {
			author := state.users[rng.Intn(len(state.users))]
			comment := model.Comment{
				AuthorID:   author.ID,
				TargetType: "trip",
				TargetID:   trip.ID,
				Content:    messages[rng.Intn(len(messages))],
			}
			if err := tx.Create(&comment).Error; err != nil {
				return err
			}
			state.comments = append(state.comments, comment)
		}
	}

	for _, photo := range state.photos {
		author := state.users[rng.Intn(len(state.users))]
		comment := model.Comment{
			AuthorID:   author.ID,
			TargetType: "photo",
			TargetID:   photo.ID,
			Content:    "Awesome photo!",
		}
		if err := tx.Create(&comment).Error; err != nil {
			return err
		}
		state.comments = append(state.comments, comment)
	}

	return nil
}

func seedReactions(tx *gorm.DB, state *seedState, rng *rand.Rand) error {
	reactionTypes := []string{"like", "love", "wow", "sad", "angry"}

	for _, trip := range state.trips {
		for _, u := range state.users {
			reaction := model.Reaction{
				UserID:       u.ID,
				TargetType:   "trip",
				TargetID:     trip.ID,
				ReactionType: reactionTypes[rng.Intn(len(reactionTypes))],
			}
			if err := tx.Create(&reaction).Error; err != nil {
				return err
			}
		}
	}

	for _, comment := range state.comments {
		u := state.users[rng.Intn(len(state.users))]
		reaction := model.Reaction{
			UserID:       u.ID,
			TargetType:   "comment",
			TargetID:     comment.ID,
			ReactionType: reactionTypes[rng.Intn(len(reactionTypes))],
		}
		if err := tx.Create(&reaction).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedPhotoTags(tx *gorm.DB, state *seedState, rng *rand.Rand) error {
	for _, photo := range state.photos {
		tagged := state.users[rng.Intn(len(state.users))]
		taggedBy := state.users[rng.Intn(len(state.users))]
		tag := model.PhotoTag{
			PhotoID:   photo.ID,
			UserID:    tagged.ID,
			TaggedBy:  taggedBy.ID,
			CreatedAt: time.Now().UTC(),
		}
		if err := tx.Create(&tag).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedNotifications(tx *gorm.DB, state *seedState, rng *rand.Rand) error {
	for _, u := range state.users {
		payload, err := json.Marshal(map[string]interface{}{
			"type":      "trip_comment",
			"trip_id":   state.trips[rng.Intn(len(state.trips))].ID,
			"message":   "New comment on your trip",
			"createdBy": state.users[rng.Intn(len(state.users))].ID,
		})
		if err != nil {
			return err
		}

		notification := model.Notification{
			UserID: u.ID,
			Type:   "trip_comment",
			Data:   datatypes.JSON(payload),
			IsRead: rng.Intn(2) == 0,
		}
		if err := tx.Create(&notification).Error; err != nil {
			return err
		}
	}
	return nil
}
