# 📞 Halloween-Phone 📞
Hacky program that uses baresip, ffmpeg and pulse to dial a contact and inject a wav file into the call. Used to play spooky sounds during a party 😱.

```
((/////////////////////////////*///*********///////////////////////(((((((((((((
(///////////////////////////////************/***////////////////////////((((((((
////////////////////////********************************///////////////////////(
///////////**************,,,,,/(,,,,#&@@&&&&&&&&&@@%/**********/////////////////
////////******,,,,,,,,,*%@&%#((//(//(/((///((##%%%%&*@#&&@@@@@@%////////////////
///////*********#&&&&&%%#((((#%&@@@%%%%&&&&&&&&@&@&@@@@&&%%#%&&&@&@@////////////
//////*******@@@@@&&&&%%%%#%%,&%@@@@@@@@@@@@@@@@@@&@@@,,&@%#//(#%&@@@///////////
////*****/**@@@@&%(**/#%&@@@,,,,@@@@@&##%#%&&@%@(.@%,,(@%@@%%##%%@@@@@//////////
///*****@&@@@@@@&&%%%%%&@@%#@,,,,,,,,,@@@%&&%%@@&%,,,,,@@&&&&&@&&@@@@@//////////
//*****/@@@&#@&@&&##@%#@@&@@@#,,,,,.@....&%(#(/%@,,&#,,,,**%@&@@@@@@@@//////////
//****/@@#@*@@@@@@@@&@@@@&@,,,,,,,@%.@,&@/,,,,,,,#& .(@*,****(&(*******/////////
//*****@/&/*********,(@/,,,.&%%#/*@%(*#(#@/....,,(@&&&&@&@/******#%#/*//////////
///**/@&@@************@###@%%%##(@%&,&@*%(%(@@@@&%&&&&&@(**********//@&/////////
//////@@@@*********/@(@@##/@@@@%@@@@@%/@&&,..(/,/@#%#@&&(%#**////////////(@@&@&&
/////@@%@*/****@@/@@@@@@@##(@@@@@@@@@@@@,&&*,*#,./@@@@@@@@/%&*//////////////////
////@@#@(@@@&&@@@@###@@@@@&##@@@@@#@@#.#&&&&@@@@@@@@@@@@@@@@#@&**////////////(((
///@@@@%@@@@&%%#//***/#((@@@%#@@&&@@&&&%#../%&&&@@@&&&&%&&&&%@@@&**/**//////////
//@@@@#%(//*********//&##&@@@((/@(/,*#&&&%%&&%####&&&@&@&&@@@@@@&@@*//*/////////
((#@@%(//************//&#&%@@@@@@@@@&&&%%%%#(*/#&&@@@@@@@@@@@@@@@@@/////////////
((##((/***************/(&&@#@@@@&%@@@@&&@@&&&&&&&&&&&@&@@@@@@@@@@@@/////////////
((((((/********,,******/(%@#&@@@@@@@&&&&&&@@&&&@@@@@@@@@@@@@@@@@@@@((///////////
(//((((/***************/(#&%@@@&@@@&%%%((#@&@&@@@@@@@&@@@@@@@@@@@@%((///////((((
(((/((((/***************/((&@@@@@@@&@&@&@@&@&&@@@@@@@@@@&#(((/((////////(/((((((
((///((((//*////(((((///////(&@@@@&%#(/////////////////////////((((///(/((((((((
((////((((/*/((((((((((((((////////////////////////////////////(((/////(((((((((
((/////((((/(((((((((((((((((/////////////////////////////////((///////(((((((((
```

Work in progress. No tests, no dockerfile. Hacked together with very limited
time.

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

## Setup

### Programs

* Compile with `go build`

* Install `ffmpeg`, `baresip` and `pulseaudio`

## Configuration

* Setup pulseaudio

    * Must be working. Use e.g. `pavucontrol` to check on pulse. 

* Setup baresip

    * It must register as a SIP client on startup

    * It must have has a default contact it calls when `/dialcontact` is entered

    * It must use `pulse` as values for both `audio_player` and `audio_source` 

    * `ausrc_srate` and `ausrc_format` must match settings for ffmpeg and virtual pulse mic. I successfully used a rate of `8000` and format `s16`. 


* Setup halloween-phone

    * Place configuration in `~/.config/halloween-phone.cfg`

**For example configuration, see folder exmaple-config.**


### Prepare Tracks

* Put `wav` files into the folder defined in the configuration file. Make sure the encoding matches the settings in the configuration (rate, format etc). Audacity is a cool tool that can be used for this.
