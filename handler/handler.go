package handler

import (
	"myOwnRedis/resp"
	"strconv"
	"strings"
	"time"
)

type storage interface {
	Get(string) (string, bool)
	Set(string, string)
	SetTTL(string, int64)
	Delete(string) bool
	Save() error
}

type Handler struct {
	storage storage
}

func NewHandler(storage storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) Handle(v resp.Value) resp.Value {
	if v.Type != resp.TypeArray || len(v.Arr) == 0 {
		return resp.Value{Type: resp.TypeError, Str: "ERR unknown command"}
	}

	command := strings.ToUpper(v.Arr[0].Str)
	args := v.Arr[1:]

	switch command {
	case "PING":
		return h.ping(args)
	case "SET":

		return h.set(args)
	case "GET":
		return h.get(args)
	case "DEL":
		return h.delete(args)
	case "EXPIRE":
		return h.expire(args)
	case "SAVE":
		return h.save(args)
	default:
		return resp.Value{Type: resp.TypeError, Str: "ERR unknown command"}
	}
}

func (h *Handler) ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Type: resp.TypeString, Str: "PONG"}
	}
	return resp.Value{Type: resp.TypeString, Str: args[0].Str}
}

func (h *Handler) get(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Type: resp.TypeError, Str: "ERR wrong number of arguments for 'get' command"}
	}
	if val, ok := h.storage.Get(args[0].Str); ok {
		return resp.Value{Type: resp.TypeString, Str: val}
	}
	return resp.Value{Type: resp.TypeNull}
}

func (h *Handler) set(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.Value{Type: resp.TypeError, Str: "ERR wrong number of arguments for 'set' command"}
	}
	h.storage.Set(args[0].Str, args[1].Str)
	return resp.Value{Type: resp.TypeString, Str: "OK"}
}

func (h *Handler) expire(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.Value{Type: resp.TypeError, Str: "ERR wrong number of arguments for 'expire' command"}
	}

	i64, err := strconv.ParseInt(args[1].Str, 10, 64)

	if err != nil {
		return resp.Value{Type: resp.TypeError, Str: "ERR value is not an integer or out of range"}
	}

	key := args[0].Str
	expiration := time.Now().Add(time.Duration(i64) * time.Second).UnixNano()
	h.storage.SetTTL(key, expiration)

	return resp.Value{Type: resp.TypeString, Str: "OK"}
}

func (h *Handler) delete(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Type: resp.TypeError, Str: "ERR wrong number of arguments for 'delete' command"}
	}
	if ok := h.storage.Delete(args[0].Str); ok {
		return resp.Value{Type: resp.TypeString, Str: "1"}
	}
	return resp.Value{Type: resp.TypeString, Str: "0"}
}

func (h *Handler) save(args []resp.Value) resp.Value {
	if err := h.storage.Save(); err != nil {
		return resp.Value{Type: resp.TypeError, Str: "ERR while saving data"}
	}
	return resp.Value{Type: resp.TypeString, Str: "OK"}
}
