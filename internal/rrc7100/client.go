package rrc7100

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/iu0jgo/gumble/gumble"
	"github.com/iu0jgo/gumble/gumbleopenal"
	"github.com/iu0jgo/gumble/gumbleutil"
	"github.com/kennygrant/sanitize"
)

func (b *RRC7100) Init() {
	b.Config.Attach(gumbleutil.AutoBitrate)
	b.Config.Attach(b)

	// b.initGPIO()

	b.Connect()

	// our main run loop here... keep things alive
	keepAlive := make(chan bool)
	exitStatus := 0

	<-keepAlive
	os.Exit(exitStatus)
}

func (b *RRC7100) CleanUp() {
	b.Client.Disconnect()
	b.LEDOffAll()
}

func (b *RRC7100) Connect() {
	var err error
	b.ConnectAttempts++

	_, err = gumble.DialWithDialer(new(net.Dialer), b.Address, b.Config, &b.TLSConfig)
	if err != nil {
		fmt.Printf("Connection to %s failed (%s), attempting again in 10 seconds...\n", b.Address, err)
		b.ReConnect()
	} else {
		b.OpenStream()
	}
}

func (b *RRC7100) ReConnect() {
	if b.Client != nil {
		b.Client.Disconnect()
	}

	if b.ConnectAttempts < 100 {
		go func() {
			time.Sleep(10 * time.Second)
			b.Connect()
		}()
		return
	} else {
		fmt.Fprintf(os.Stderr, "Unable to connect, giving up\n")
		os.Exit(1)
	}
}

func (b *RRC7100) OpenStream() {
	// Audio
	if os.Getenv("ALSOFT_LOGLEVEL") == "" {
		os.Setenv("ALSOFT_LOGLEVEL", "0")
	}

	if stream, err := gumbleopenal.New(b.Client); err != nil {
		fmt.Fprintf(os.Stderr, "Stream open error (%s)\n", err)
		os.Exit(1)
	} else {
		b.Stream = stream
		b.TransmitStart() //Instantly start Transmitting
	}
}

func (b *RRC7100) ResetStream() {
	b.Stream.Destroy()

	// Sleep a bit and re-open
	time.Sleep(50 * time.Millisecond)

	b.OpenStream()
}

func (b *RRC7100) TransmitStart() {
	if b.IsConnected == false {
		return
	}

	b.IsTransmitting = true

	// turn on our transmit LED
	b.LEDOn(b.TransmitLED)

	b.Stream.StartSource()
}

func (b *RRC7100) TransmitStop() {
	if b.IsConnected == false {
		return
	}

	b.Stream.StopSource()

	b.LEDOff(b.TransmitLED)

	b.IsTransmitting = false
}

func (b *RRC7100) OnConnect(e *gumble.ConnectEvent) {
	b.Client = e.Client

	b.ConnectAttempts = 0

	b.IsConnected = true
	// turn on our online LED
	b.LEDOn(b.OnlineLED)

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
	b.LEDOff(b.OnlineLED)
	b.LEDOff(b.ParticipantsLED)
	b.LEDOff(b.TransmitLED)

	if reason == "" {
		fmt.Printf("Connection to %s disconnected, attempting again in 10 seconds...\n", b.Address)
	} else {
		fmt.Printf("Connection to %s disconnected (%s), attempting again in 10 seconds...\n", b.Address, reason)
	}

	// attempt to connect again
	b.ReConnect()
}

func (b *RRC7100) ChangeChannel(ChannelName string) {
	channel := b.Client.Channels.Find(ChannelName)
	if channel != nil {
		b.Client.Self.Move(channel)
	} else {
		fmt.Printf("Unable to find channel: %s\n", ChannelName)
	}
}

func (b *RRC7100) ParticipantLEDUpdate() {
	time.Sleep(100 * time.Millisecond)

	// If we have more than just ourselves in the channel, turn on the participants LED, otherwise, turn it off

	var participantCount = len(b.Client.Self.Channel.Users)

	if participantCount > 1 {
		fmt.Printf("Channel '%s' has %d participants\n", b.Client.Self.Channel.Name, participantCount)
		b.LEDOn(b.ParticipantsLED)
	} else {
		fmt.Printf("Channel '%s' has no other participants\n", b.Client.Self.Channel.Name)
		b.LEDOff(b.ParticipantsLED)
	}
}

func (b *RRC7100) OnTextMessage(e *gumble.TextMessageEvent) {
	fmt.Printf("Message from %s: %s\n", e.Sender.Name, strings.TrimSpace(esc(e.Message)))
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

	go b.ParticipantLEDUpdate()
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

func (b *RRC7100) OnUserList(e *gumble.UserListEvent) {
}

func (b *RRC7100) OnACL(e *gumble.ACLEvent) {
}

func (b *RRC7100) OnBanList(e *gumble.BanListEvent) {
}

func (b *RRC7100) OnContextActionChange(e *gumble.ContextActionChangeEvent) {
}

func (b *RRC7100) OnServerConfig(e *gumble.ServerConfigEvent) {
}

func esc(str string) string {
	return sanitize.HTML(str)
}
