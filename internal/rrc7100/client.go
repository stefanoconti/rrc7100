package rrc7100

import (
	"fmt"
	"net"
	"os"
	"time"

	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleopenal"
	"layeh.com/gumble/gumbleutil"
)

func (b *RRC7100) Init() {
	b.Config.Attach(gumbleutil.AutoBitrate)
	//b.Config.Attach(b)
	b.Config.Attach(gumbleutil.Listener{
		Connect:          b.OnConnect,
		Disconnect:       b.OnDisconnect,
		ChannelChange:    b.OnChannelChange,
		TextMessage:      b.OnTextMessage,
		UserChange:       b.OnUserChange,
		PermissionDenied: b.OnPermissionDenied,
	})

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
	// b.LEDOffAll()
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
	if !b.IsConnected {
		return
	}

	b.IsTransmitting = true

	// turn on our transmit LED
	// b.LEDOn(b.TransmitLED)

	b.Stream.StartSource()
}

func (b *RRC7100) TransmitStop() {
	if b.IsConnected {
		return
	}

	b.Stream.StopSource()

	// b.LEDOff(b.TransmitLED)

	b.IsTransmitting = false
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
		// b.LEDOn(b.ParticipantsLED)
	} else {
		fmt.Printf("Channel '%s' has no other participants\n", b.Client.Self.Channel.Name)
		// b.LEDOff(b.ParticipantsLED)
	}
}
