// Package notification implements the in-game inbox, push preferences and
// the push queue with FCM/APNs delivery (Sprint 23, docs/notification/*).
package notification

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Mode string

const (
	ModePushAndInbox Mode = "push_and_inbox"
	ModeInboxOnly    Mode = "inbox_only"
	ModeSilentPush   Mode = "silent_push"
	ModeOff          Mode = "off"
)

func (m Mode) Valid() bool {
	switch m {
	case ModePushAndInbox, ModeInboxOnly, ModeSilentPush, ModeOff:
		return true
	}
	return false
}

// Kinds and their default modes (docs/notification/preferences.md §3).
const KindSystemBroadcast = "system_broadcast"

var DefaultModes = map[string]Mode{
	"activity_completed":   ModePushAndInbox,
	"dungeon_completed":    ModePushAndInbox,
	"craft_completed":      ModePushAndInbox,
	"travel_completed":     ModePushAndInbox,
	"energy_restored":      ModeOff,
	"world_event_starting": ModePushAndInbox,
	"world_event_active":   ModePushAndInbox,
	"guild_event":          ModePushAndInbox,
	"mail_received":        ModePushAndInbox,
	"whisper_received":     ModePushAndInbox,
	"friend_online":        ModeOff,
	"market_sold":          ModePushAndInbox,
	"market_expiring":      ModeOff,
	"party_invite":         ModePushAndInbox,
	"raid_matched":         ModePushAndInbox, // the raid queue formed a party (DD-299)
	"guild_invite":         ModePushAndInbox,
	"achievement_unlocked": ModeSilentPush,
	"support_reply":        ModePushAndInbox, // DD-234
	KindReward:             ModeInboxOnly,    // DD-242
	KindSystemBroadcast:    ModePushAndInbox,
}

// RateRule caps pushes per account per kind (push-notification.md §9).
type RateRule struct {
	MaxPerHour int
	Cooldown   time.Duration
}

var RateRules = map[string]RateRule{
	"activity_completed":   {6, 5 * time.Minute},
	"dungeon_completed":    {6, 5 * time.Minute},
	"craft_completed":      {12, 2 * time.Minute},
	"travel_completed":     {8, 2 * time.Minute},
	"energy_restored":      {2, 30 * time.Minute},
	"world_event_starting": {2, 10 * time.Minute},
	"world_event_active":   {2, 10 * time.Minute},
	"guild_event":          {4, 5 * time.Minute},
	"mail_received":        {10, time.Minute},
	"whisper_received":     {20, 30 * time.Second},
	"friend_online":        {5, time.Minute},
	"market_sold":          {20, time.Minute},
	"party_invite":         {10, 30 * time.Second},
	"raid_matched":         {6, time.Minute},
	"guild_invite":         {10, 30 * time.Second},
	"achievement_unlocked": {10, time.Minute},
}

// RateAllowed reports whether one more push fits the kind's budget given
// pushes already sent in the last hour and the time of the latest one.
// system_broadcast is never rate limited.
func RateAllowed(kind string, sentLastHour int, last *time.Time, now time.Time) bool {
	rule, ok := RateRules[kind]
	if !ok {
		return true
	}
	if sentLastHour >= rule.MaxPerHour {
		return false
	}
	return last == nil || now.Sub(*last) >= rule.Cooldown
}

// ParseClock parses "HH:MM" into minutes after midnight.
func ParseClock(s string) (int, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 || len(parts[0]) != 2 || len(parts[1]) != 2 {
		return 0, fmt.Errorf("notification: bad clock %q", s)
	}
	h, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, fmt.Errorf("notification: bad clock %q", s)
	}
	return h*60 + m, nil
}

// InQuietHours reports whether now falls inside [start, end) in tz. The
// window may wrap midnight ("22:30"–"08:00"). Unset or invalid → false.
func InQuietHours(now time.Time, start, end, tz string) bool {
	if start == "" || end == "" {
		return false
	}
	s, err1 := ParseClock(start)
	e, err2 := ParseClock(end)
	if err1 != nil || err2 != nil || s == e {
		return false
	}
	loc := time.UTC
	if tz != "" {
		if l, err := time.LoadLocation(tz); err == nil {
			loc = l
		}
	}
	t := now.In(loc)
	cur := t.Hour()*60 + t.Minute()
	if s < e {
		return cur >= s && cur < e
	}
	return cur >= s || cur < e
}

// Decision is the outcome of applying preferences to one notice.
type Decision struct {
	Inbox    bool
	Push     bool
	Priority string // high | normal | silent
}

// Decide applies the account's mode, global mute and quiet hours.
// system_broadcast ignores all of them (preferences.md §4-5).
func Decide(kind string, mode Mode, muted, quiet bool, highPriority bool) Decision {
	prio := "normal"
	if highPriority {
		prio = "high"
	}
	if kind == KindSystemBroadcast {
		return Decision{Inbox: true, Push: true, Priority: "high"}
	}
	d := Decision{}
	switch mode {
	case ModeOff:
		return d
	case ModeInboxOnly:
		d.Inbox = true
		return d
	case ModeSilentPush:
		d.Push, prio = true, "silent"
	case ModePushAndInbox:
		d.Inbox, d.Push = true, true
	}
	if muted {
		d.Push = false
	}
	if quiet {
		prio = "silent"
	}
	d.Priority = prio
	return d
}
