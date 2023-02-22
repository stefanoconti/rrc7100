package rrc7100

import (
	"fmt"
	"strings"

	"github.com/kennygrant/sanitize"
	"layeh.com/gumble/gumble"
)

func (b *RRC7100) OnConnect(e *gumble.ConnectEvent) {
	b.Client = e.Client

	b.ConnectAttempts = 0

	b.IsConnected = true
	// turn on our online LED
	// b.LEDOn(b.OnlineLED)

	fmt.Printf("Connected to %s (%d)\n", b.Client.Conn.RemoteAddr(), b.ConnectAttempts)
	if e.WelcomeMessage != nil {
		fmt.Printf("Welcome message: %s\n", esc(*e.WelcomeMessage))
	}

	if b.ChannelName != "" {
		b.ChangeChannel(b.ChannelName)
	}
}

func (b *RRC7100) OnDisconnect(e *gumble.DisconnectEvent) {
	var reason string
	switch e.Type {
	case gumble.DisconnectError:
		reason = "connection error"
	}

	b.IsConnected = false

	// turn off our LEDs
	// b.LEDOff(b.OnlineLED)
	// b.LEDOff(b.ParticipantsLED)
	// b.LEDOff(b.TransmitLED)

	if reason == "" {
		fmt.Printf("Connection to %s disconnected, attempting again in 10 seconds...\n", b.Address)
	} else {
		fmt.Printf("Connection to %s disconnected (%s), attempting again in 10 seconds...\n", b.Address, reason)
	}

	// attempt to connect again
	b.ReConnect()
}

func (b *RRC7100) OnUserChange(e *gumble.UserChangeEvent) {
	var info string

	switch e.Type {
	case gumble.UserChangeConnected:
		info = "connected"
	case gumble.UserChangeDisconnected:
		info = "disconnected"
	case gumble.UserChangeKicked:
		info = "kicked"
	case gumble.UserChangeBanned:
		info = "banned"
	case gumble.UserChangeRegistered:
		info = "registered"
	case gumble.UserChangeUnregistered:
		info = "unregistered"
	case gumble.UserChangeName:
		info = "changed name"
	case gumble.UserChangeChannel:
		info = "changed channel"
	case gumble.UserChangeComment:
		info = "changed comment"
	case gumble.UserChangeAudio:
		info = "changed audio"
	case gumble.UserChangePrioritySpeaker:
		info = "is priority speaker"
	case gumble.UserChangeRecording:
		info = "changed recording status"
	case gumble.UserChangeStats:
		info = "changed stats"
	}

	fmt.Printf("Change event for %s: %s (%d)\n", e.User.Name, info, e.Type)

	//go b.ParticipantLEDUpdate()
}

func (b *RRC7100) OnPermissionDenied(e *gumble.PermissionDeniedEvent) {
	var info string
	switch e.Type {
	case gumble.PermissionDeniedOther:
		info = e.String
	case gumble.PermissionDeniedPermission:
		info = "insufficient permissions"
	case gumble.PermissionDeniedSuperUser:
		info = "cannot modify SuperUser"
	case gumble.PermissionDeniedInvalidChannelName:
		info = "invalid channel name"
	case gumble.PermissionDeniedTextTooLong:
		info = "text too long"
	case gumble.PermissionDeniedTemporaryChannel:
		info = "temporary channel"
	case gumble.PermissionDeniedMissingCertificate:
		info = "missing certificate"
	case gumble.PermissionDeniedInvalidUserName:
		info = "invalid user name"
	case gumble.PermissionDeniedChannelFull:
		info = "channel full"
	case gumble.PermissionDeniedNestingLimit:
		info = "nesting limit"
	}

	fmt.Printf("Permission denied: %s\n", info)
}

func (b *RRC7100) OnChannelChange(e *gumble.ChannelChangeEvent) {
	go b.ParticipantLEDUpdate()
}

func (b *RRC7100) OnTextMessage(e *gumble.TextMessageEvent) {
	fmt.Printf("Message from %s: %s\n", e.Sender.Name, strings.TrimSpace(esc(e.Message)))
}

func esc(str string) string {
	return sanitize.HTML(str)
}
