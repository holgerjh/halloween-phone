# 📞 Halloween-Phone 📞
Hacky program that uses baresip, ffmpeg and pulse to dial a contact and inject a wav file into the call. Used to play spooky sounds during a party 😱.

![logo](https://github.com/holgerjh/halloween-phone/blob/main/spookyphone.jpg?raw=true)


Work in progress. No tests, no dockerfile. Hacked together with very limited
time.


## Warning

**WARNING: Make sure to double-check the contact that is dialed if you run
`/dialcontact` so you do not perform calls that may cost you money. If
possible, block outgoing calls in your SIP provider. I do not take any
responsibility for costs that occure using this program. Use at your own risk!**

## How does it work

```
          ┌─────────────┐
┌─────────┤ wavs        │
│         └─────────────┴─────┐
│                             │<<loads>>
│<<consumes>>      ┌──────────▼────────────┐
│                  │                       │
│              ┌───┤    halloween-phone    │
│              │   │                       ├───────────┐
│ <<orchestrates>> └───────────────────────┘           │<<orchestrates>>
│          ┌───▼───────┐  ┌───────────────────┐   ┌────▼──────┐
└─────────►│   ffmpeg  ├──► virtual pulse mic ├───►  baresip  │
           └───────────┘  └───────────────────┘   └───────────┘
```

Halloween-phone takes a folder with wav files and registers a virtual pulse audio microphone. It then schedules the files to be played. When the time comes, it runs baresip and calls the default contact. When someone picks up, ffmpeg is utilized to stream the wav file into the microphone. As soon as the callee hangs up, the programs are terminated again. Halloween-phone loops forever until CTRL+C is pressed. The timeout between calls, as well as other things, can be configured to match your desired use-case (which is, of course, running a spooky party 🎉🎉🎉).

## Setup

### Hardware

* A physical telephone to call. I used a really old one one with pulse dialing. This had the bennefit that none of my guests could make any outgoing calls. The telephone needs to be callable. I had to solder a connector to connect the phone to my router. For german phones you should connect white and brown wires (ignore the green one if it exists), should be the same for international phones. For american phones, it should be red and green. Please double-check this against other sources before make any connections.

### Programs

* Compile with `go build`

* Install `ffmpeg`, `baresip` and `pulseaudio`

## Configuration

* Setup pulseaudio

    * Must be working. Use e.g. `pavucontrol` to check on pulse. 

* Setup baresip

    * It must register as a SIP client on startup

    * It must have a default contact it calls when `/dialcontact` is entered

    * It must use `pulse` as values for both `audio_player` and `audio_source` 

    * `ausrc_srate` and `ausrc_format` must match settings for ffmpeg and virtual pulse mic. I successfully used a rate of `8000` and format `s16`. 


* Setup halloween-phone

    * Place configuration in `~/.config/halloween-phone.cfg`

**For example configuration, see folder [example-config](https://github.com/holgerjh/halloween-phone/tree/main/example-config).**


### Prepare Tracks

* Put `wav` files into the folder defined in the configuration file. Make sure the encoding matches the settings in the configuration (rate, format etc). Audacity is a cool tool that can be used for this.
