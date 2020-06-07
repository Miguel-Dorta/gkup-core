package io

import (
	"github.com/Miguel-Dorta/gkup-core/pkg/syncUtils"
	"sync"
)

type interactiveRequest struct {
	ID            int      `json:"id"`
	InteractionID uint64   `json:"interaction_id"`
	Prompt        string   `json:"prompt"`
	ValidReplies  []string `json:"valid_replies"`
	reply         string
	wg            *sync.WaitGroup
}

var (
	activeRequests = make(map[uint64]*interactiveRequest, 10)
	idCounter      syncUtils.AtomicCounter
)

func init() {
	callback[IdInteractiveErr] = handleInteractiveErrors
}

func handleInteractiveErrors(j JSON) JSON {
	req := activeRequests[j["interaction_id"].(uint64)]
	req.reply = j["response"].(string)
	if containsString(req.reply, req.ValidReplies) {
		return JSON{
			"id":    IdErr,
			"error": "invalid reply",
		}
	}
	req.wg.Done()
	return nil
}

func AskInteractiveError(prompt string, validResponses []string) string {
	r := &interactiveRequest{
		ID:            IdInteractiveErr,
		InteractionID: idCounter.Add(),
		Prompt:        prompt,
		ValidReplies:  validResponses,
		wg:            new(sync.WaitGroup),
	}
	r.wg.Add(1)
	activeRequests[r.InteractionID] = r
	if err := stdout.Encode(r); err != nil {
		panic("error encoding message: " + err.Error())
	}
	r.wg.Wait()
	return r.reply
}

func containsString(s string, list []string) bool {
	for i := range list {
		if list[i] == s {
			return true
		}
	}
	return false
}
