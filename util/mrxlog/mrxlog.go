package mrxlog

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/cespare/xxhash/v2"
)

const (
	logKey      = "MRXPath"
	chainIDKey  = "chainID"
	parentIDKey = "parentID"
)

// MRXHistory points to the parents of all previous
// mrxIds in a register chain.
type MRXHistory struct {
	MrxID  string
	action string
	parent *MRXHistory
	extra  map[string]any
	origin string
	// include internal depth?

}

/*
Prospective action list

- search
- transform but how
- schema match?
- API at certain location
- run custom action?

*/

// String is a Custom string function to ensure recursive pointers are found
// instead of writing the memory address.
func (mh MRXHistory) String() string {

	return fmt.Sprintf("{%v %v %v %v %v}", mh.MrxID, mh.action, mh.parent, mh.extra, mh.origin)
}

// MarshalJSON keeps the MRXhistory unexported
// updating with this function allows all the data to be exported.
func (mh MRXHistory) MarshalJSON() ([]byte, error) {

	m := mrxHistory{
		MrxID:  mh.MrxID,
		Action: mh.action,
		Parent: mh.parent,
		Extra:  mh.extra,
		Origin: mh.origin,
	}

	return json.Marshal(m)
}

// convert before the json bytes stage
type mrxHistory struct {
	MrxID  string
	Action string
	Parent *MRXHistory    `json:"Parent,omitempty"`
	Extra  map[string]any `json:"Extra,omitempty"`
	Origin string
}

// NewMRX generates a new mrxHistory object, with any extra options
func NewMRX(mrxId string, options ...func(*MRXHistory)) *MRXHistory {
	mh := &MRXHistory{MrxID: mrxId}
	for _, opt := range options {
		opt(mh)
	}
	return mh

}

// WithID sets the metarexID
func WithID(id string) func(t *MRXHistory) {

	return func(m *MRXHistory) {
		m.MrxID = id
	}
}

// WithAction sets the action forming this metarexID
func WithAction(action string) func(t *MRXHistory) {

	return func(m *MRXHistory) {
		m.action = action
	}
}

// WithOrigin sets the Origin of this metadata
func WithOrigin(origin string) func(t *MRXHistory) {

	return func(m *MRXHistory) {
		m.origin = origin
	}
}

// WithExtra sets the any extra fields
func WithExtra(extra map[string]any) func(t *MRXHistory) {

	return func(m *MRXHistory) {
		m.extra = extra
	}
}

// Depth returns the depth (how many parents) of an mrx object.
func (m MRXHistory) Depth() (n int) {
	middleM := m
	n = 1 // start the count at 1
	for middleM.parent != nil {
		n++
		middleM = *middleM.parent
	}

	return n
}

// PushChild adds an MRX child to the current error message chain
// the pointer that is returned is to the child.
func (parent *MRXHistory) PushChild(child MRXHistory) *MRXHistory {
	child.parent = parent
	return &child
}

// PopChild takes the children of depth nCount from an mrxHistroy object
// the pointer that is returned is to the parents of n count.
func (m MRXHistory) PopChild(nCount int) (*MRXHistory, error) {

	// child.parent = Parent
	d := m.Depth()
	switch {
	case d < nCount:
		return nil, fmt.Errorf("more ncount of %v greater than the total number of children %v", nCount, d)
	case d == nCount:
		return &m, nil
	default:
		// loop through
		pos := m
		children := make([]MRXHistory, nCount)
		for i := 0; i < nCount; i++ {
			children[i] = pos
			pos = *pos.parent
		}

		child := children[0]
		ref := &child
		for _, parent := range children[1:] {
			ref.parent = &parent
			ref = &parent
		}

		// cut the parent pointer off
		ref.parent = nil

		return &child, nil
	}
	/*
		return &m, nil

		/*
			go up N times, as long as depth allows it

		return &MRXHistory{}, nil*/
}

// CutChild cuts nCount of children out of the object
func (m MRXHistory) CutChild(child MRXHistory, nCount int) (*MRXHistory, error) {

	// child.parent = Parent
	d := m.Depth()
	switch {
	case d < nCount:
		return nil, fmt.Errorf("more ncount of %v greater than the total number of children %v", nCount, d)
	default:
		// loop through
		for i := 0; i < nCount; i++ {
			m = *m.parent
		}
	}
	return &m, nil
}

// LogDebug logs to slog, adding the error message as an argument so its not
// manually updated fo reach instance
func (m *MRXHistory) LogDebug(msg string, args ...any) {
	ckey, cID, pkey, pid := m.GetID()
	baseArgs := append([]any{logKey, *m, ckey, cID, pkey, pid}, args...)
	slog.Debug(msg, baseArgs...)
}

// GetID generates the hash of the Chain of parents and the parents ID.
// This can be used to trace objects through the register more easily.
func (m MRXHistory) GetID() (chainID, chain, parentID, parent string) {
	middleM := m
	// start the count at 1
	path := make([]string, 0)
	for middleM.parent != nil {
		path = append(path, middleM.MrxID)
		middleM = *middleM.parent
	}
	// set the origin of the parent to help set it to be more unique
	path = append(path, []string{middleM.MrxID, middleM.origin}...)
	// revers the order so parent IDs can be calculared
	slices.Reverse(path)

	// length greater than 2 as you don't want to include the parent
	if len(path) > 2 {
		parentPaths := strings.Join(path[:len(path)-1], "")
		parent = fmt.Sprintf("%08x", xxhash.Sum64([]byte(parentPaths)))

	}

	paths := strings.Join(path, "")
	chain = fmt.Sprintf("%08x", xxhash.Sum64([]byte(paths)))
	chainID = chainIDKey
	parentID = parentIDKey
	return
}

// LogInfo logs to slog, adding the error message as an argument so its not
// manually updated fo reach instance
func (m *MRXHistory) LogInfo(msg string, args ...any) {

	// utilise this later
	// m, _ = m.PopChild(1)

	ckey, cID, pkey, pid := m.GetID()
	baseArgs := append([]any{logKey, *m, ckey, cID, pkey, pid}, args...)
	slog.Info(msg, baseArgs...)
}

// LogWarn logs to slog, adding the error message as an argument so its not
// manually updated fo reach instance
func (m *MRXHistory) LogWarn(msg string, args ...any) {

	ckey, cID, pkey, pid := m.GetID()
	baseArgs := append([]any{logKey, *m, ckey, cID, pkey, pid}, args...)
	slog.Warn(msg, baseArgs...)
}

// LogError logs to slog, adding the error message as an argument so its not
// manually updated fo reach instance
func (m *MRXHistory) LogError(msg string, args ...any) {
	ckey, cID, pkey, pid := m.GetID()
	baseArgs := append([]any{logKey, *m, ckey, cID, pkey, pid}, args...)
	slog.Error(msg, baseArgs...)
}
