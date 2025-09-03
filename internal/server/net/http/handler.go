package http

import (
	"log"
	"net/http"

	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/otie173/odinbit/internal/server/game/texture"
	"github.com/otie173/odinbit/internal/server/game/world"
	"github.com/vmihailenco/msgpack/v5"
)

type Handler struct {
	router         *Router
	textureStorage *texture.Storage
	world          *world.World
}

func NewHandler(router *Router, storage *texture.Storage, world *world.World) *Handler {
	h := &Handler{
		router:         router,
		textureStorage: storage,
		world:          world,
	}

	h.router.setupRoutes(h)
	return h
}

func (h *Handler) Run(addr string) error {
	if err := http.ListenAndServe(addr, h.router.mux); err != nil {
		return err
	}
	return nil
}

func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ping"))
}

func createPacket(pktType packet.PacketType, pktData []byte) packet.Packet {
	return packet.Packet{
		Type:    pktType,
		Payload: pktData,
	}
}

func createBPacket(pkt packet.Packet) ([]byte, error) {
	binaryPkt, err := msgpack.Marshal(&pkt)
	if err != nil {
		return nil, err
	}
	return binaryPkt, nil
}

func (h *Handler) getTextures(w http.ResponseWriter, r *http.Request) {
	data, err := h.textureStorage.GetTextures()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Printf("Error! Cant get textures: %v\n", err)
		return
	}

	pkt := createPacket(packet.GetTexturesType, data)
	binaryPkt, err := msgpack.Marshal(&pkt)
	if err != nil {
		log.Printf("Error! Cant marshal packet: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/msgpack")
	if _, err := w.Write(binaryPkt); err != nil {
		log.Printf("Error! Cant send binary packet: %v\n", err)
	}
}
