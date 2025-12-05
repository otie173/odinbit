package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/otie173/odinbit/internal/pkg/server"
	"github.com/otie173/odinbit/internal/server/core/manager"
	"github.com/otie173/odinbit/internal/server/core/ticker"
	"github.com/otie173/odinbit/internal/server/game/player"
	"github.com/otie173/odinbit/internal/server/game/texture"
	"github.com/otie173/odinbit/internal/server/game/world"
	"github.com/otie173/odinbit/internal/server/net/http"
	"github.com/otie173/odinbit/internal/server/net/tcp"
)

func runMigrations(db *sqlx.DB) error {
	files, err := filepath.Glob("migrations/*.up.sql")
	if err != nil {
		return err
	}

	sort.Strings(files)
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		if _, err := db.Exec(string(content)); err != nil {
			return err
		}
		log.Printf("Info! Database migration applied: %s\n", file)
	}

	return nil
}

func main() {
	db, err := sqlx.Connect("sqlite3", "database.db")
	if err != nil {
		log.Println(db, err)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		log.Fatalf("Error! Cant apply data migration: %v\n", err)
	}
	log.Println("Info! Database migrations applied successfully")

	textures := texture.NewPack()
	overworld := world.New(textures)
	playerStorage := player.NewStorage(db, 16)

	textureHandler := texture.NewHandler(textures)
	playerHandler := player.NewHandler(playerStorage)
	worldHandler := world.NewHandler(overworld, playerStorage)

	router := http.NewRouter(chi.NewRouter())
	handler := http.NewHandler(router, textures, overworld)

	dispatcher := tcp.NewDispatcher(playerStorage, textureHandler)
	listener := tcp.NewListener(dispatcher)

	tps := 10
	ticker := ticker.New(tps)

	components := manager.Components{
		Textures:      textures,
		Overworld:     overworld,
		Players:       playerStorage,
		WorldHandler:  worldHandler,
		PlayerHandler: playerHandler,
		Handler:       handler,
		Listener:      listener,
		Ticker:        ticker,
	}
	manager := manager.New(components)

	server := server.New("0.0.0.0:8080", manager)
	server.Load()
	server.Run()
}
