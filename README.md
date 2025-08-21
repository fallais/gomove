# GoMove

GoMove is a CLI application that prevents your computer session from locking by :

- periodically moving the mouse cursor slightly
- periodically tapping on `cmd` key (Windows key for Windows)

It's designed to work on both Windows and Linux systems.

## Why ?

You may not want your session being lock because :

- You have a long command running on terminal
- You want it

## Usage

### Configuration file

You can run `go move config create`, it will create a default configuration file for you in `~/.gomove/config.yaml`.

Or you can create it yourself :

```yaml
behavior:
  idle_timeout: 10s
  resume_after_inactivity: true
  pause_when_user_is_active: true
activities:
  - kind: mouse
    interval: 3s
debug: true
logfile: ""
```

### Start

Linux: Simply run `gomove start`

Windows: Simply run `gomove.exe start`

### Run as a service

