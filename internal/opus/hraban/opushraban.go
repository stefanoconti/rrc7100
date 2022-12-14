package opushraban

import (
	"fmt"

	"github.com/iu0jgo/gumble/gumble"
	hopus "gopkg.in/hraban/opus.v2"
)

var Codec gumble.AudioCodec

const ID = 4

func Register(encoderMode string) {
	Codec = &generator{encoderMode: encoderMode}
	gumble.RegisterAudioCodec(4, Codec)
}

// generator

type generator struct {
	encoderMode string
}

func (g *generator) getOpusApplication() hopus.Application {
	switch g.encoderMode {
	case "voip":
		fmt.Println("Use 'VoIP' encoder mode")
		return hopus.AppVoIP
	case "audio":
		fmt.Println("Use 'Audio' encoder mode")
		return hopus.AppAudio
	case "lowdelay":
		fmt.Println("Use 'Restricted Low Delay' encoder mode")
		return hopus.AppRestrictedLowdelay
	default:
		fmt.Println("Use 'Audio' encoder mode")
		return hopus.AppAudio
	}
}

func (*generator) ID() int {
	return ID
}

func (g *generator) NewEncoder() gumble.AudioEncoder {
	e, _ := hopus.NewEncoder(gumble.AudioSampleRate, gumble.AudioChannels, g.getOpusApplication())
	e.SetBitrate(24000)
	e.SetComplexity(3)
	e.SetMaxBandwidth(hopus.Wideband)
	return &Encoder{
		e,
	}
}

func (*generator) NewDecoder() gumble.AudioDecoder {
	d, _ := hopus.NewDecoder(gumble.AudioSampleRate, gumble.AudioChannels)
	return &Decoder{
		d,
		gumble.AudioChannels,
	}
}

// encoder

type Encoder struct {
	*hopus.Encoder
}

func (*Encoder) ID() int {
	return ID
}

func (e *Encoder) Encode(pcm []int16, mframeSize, maxDataBytes int) ([]byte, error) {
	buf := make([]byte, maxDataBytes)
	n, err := e.Encoder.Encode(pcm, buf)
	return buf[:n], err

}

func (e *Encoder) Reset() {
	// e.Encoder.ResetState()
}

// decoder

type Decoder struct {
	*hopus.Decoder
	channels int
}

func (*Decoder) ID() int {
	return 4
}

func (d *Decoder) Decode(data []byte, frameSize int) ([]int16, error) {
	pcmBuf := make([]int16, d.channels*frameSize)
	n, err := d.Decoder.Decode(data, pcmBuf)
	return pcmBuf[:n], err
}

func (d *Decoder) Reset() {
	// d.Decoder.ResetState()
}
