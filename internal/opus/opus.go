package opus

import (
	"fmt"

	"github.com/iu0jgo/gopus"
	"github.com/iu0jgo/gumble/gumble"
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

func (g *generator) getOpusApplication() gopus.Application {
	switch g.encoderMode {
	case "voip":
		fmt.Println("Use 'VoIP' encoder mode")
		return gopus.Voip
	case "lowdelay":
		fmt.Println("Use 'Restricted Low Delay' encoder mode")
		return gopus.RestrictedLowDelay
	default:
		fmt.Println("Use 'Audio' encoder mode")
		return gopus.Audio
	}
}

func (*generator) ID() int {
	return ID
}

func (g *generator) NewEncoder() gumble.AudioEncoder {
	e, _ := gopus.NewEncoder(gumble.AudioSampleRate, gumble.AudioChannels, g.getOpusApplication())
	e.SetBitrate(gopus.BitrateMaximum)
	return &Encoder{
		e,
	}
}

func (*generator) NewDecoder() gumble.AudioDecoder {
	d, _ := gopus.NewDecoder(gumble.AudioSampleRate, gumble.AudioChannels)
	return &Decoder{
		d,
	}
}

// encoder

type Encoder struct {
	*gopus.Encoder
}

func (*Encoder) ID() int {
	return ID
}

func (e *Encoder) Encode(pcm []int16, mframeSize, maxDataBytes int) ([]byte, error) {
	return e.Encoder.Encode(pcm, mframeSize, maxDataBytes)
}

func (e *Encoder) Reset() {
	e.Encoder.ResetState()
}

// decoder

type Decoder struct {
	*gopus.Decoder
}

func (*Decoder) ID() int {
	return 4
}

func (d *Decoder) Decode(data []byte, frameSize int) ([]int16, error) {
	return d.Decoder.Decode(data, frameSize, false)
}

func (d *Decoder) Reset() {
	d.Decoder.ResetState()
}
