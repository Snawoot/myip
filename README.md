# myip

Reliably and quickly get your external IP address from public STUN servers. Program issues parallel queries to public STUN servers to determine public IP address and returns result as soon as quorum of matching responses reached. By default quorum is 2. Useful for scripting.

## Installation

#### Binary download

Pre-built binaries available on [releases](https://github.com/Snawoot/myip/releases/latest) page.

#### From source

Alternatively, you may install myip from source. Run within source directory

```
go install
```

#### Docker

Docker image is available as well:

```sh
docker run --rm yarmak/myip
```

#### Snap Store

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/go-myip)

```bash
sudo snap install go-myip
```

## Usage

```
$ myip
1.2.3.4
```

## Synopsis

```
  -6	use IPv6
  -q uint
    	required number of matches for success (default 2)
  -s string
    	STUN server list (default "stun.l.google.com:19302;stun.ekiga.net:3478;stun.ideasip.com:3478;stun.schlund.de:3478;stun.voiparound.com:3478;stun.voipbuster.com:3478;stun.voipstunt.com:3478")
  -t duration
    	hard timeout. Examples values: 1m, 3s, 1s500ms
```
