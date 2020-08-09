package grav

import (
	"sync"
)

// messageFilter is a series of maps that associate things about a message (its UUID, type, etc) with a boolean value to say if
// it should be allowed or denied. For each of the maps, if an entry is included, the value of the boolean is respected (true = allow, false = deny)
// Maps are either inclusive (meaning that a missing entry defaults to allow), or exclusive (meaning that a missing entry defaults to deny)
// This can be configured per map by modifiying the UUIDInclusive, TypeInclusive (etc) fields.
type messageFilter struct {
	UUIDMap       map[string]bool
	UUIDInclusive bool

	TypeMap       map[string]bool
	TypeInclusive bool

	sync.RWMutex
}

func newMessageFilter() *messageFilter {
	mf := &messageFilter{
		UUIDMap:       map[string]bool{},
		UUIDInclusive: true,
		TypeMap:       map[string]bool{},
		TypeInclusive: true,
		RWMutex:       sync.RWMutex{},
	}

	return mf
}

func (mf *messageFilter) allow(msg Message) bool {
	mf.RLock()
	defer mf.RUnlock()

	// for each map, deny the message if:
	//	- a filter entry exists and it's value is false
	//	- a filter entry doesn't exist and its inclusive rule is false

	allowType, exists := mf.TypeMap[msg.Type()]
	if exists && !allowType {
		return false
	} else if !exists && !mf.TypeInclusive {
		return false
	}

	allowUUID, exists := mf.UUIDMap[msg.UUID()]
	if exists && !allowUUID {
		return false
	} else if !exists && !mf.UUIDInclusive {
		return false
	}

	return true
}

// FilterUUID likely should not be used in normal cases, it adds a message UUID to the pod's filter.
func (mf *messageFilter) FilterUUID(uuid string, allow bool) {
	mf.Lock()
	defer mf.Unlock()

	mf.UUIDMap[uuid] = allow
}

func (mf *messageFilter) FilterType(msgType string, allow bool) {
	mf.Lock()
	defer mf.Unlock()

	mf.TypeMap[msgType] = allow
}
