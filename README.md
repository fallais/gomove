# GoMove

**GoMove** is a lightweight, cross-platform CLI application that prevents your computer session from locking by simulating minimal user activity. It intelligently moves the mouse cursor in subtle patterns and can optionally simulate keyboard activity to keep your system active.

## Features

- **Smart mouse movement**: Multiple movement patterns (square, triangle, up/down, left/right)
- **Keyboard simulation**: Optional key press simulation for additional activity
- **Intelligent detection**: Automatically pauses when user is actively working
- **Configurable timing**: Customizable intervals and idle timeouts
- **Scheduling**: Optional time-based scheduling for automatic activation
- **Easy configuration**: YAML-based configuration with sensible defaults
- **Cross-Platform**: Native support for Windows and Linux
- **Logging**: Comprehensive logging with configurable output

## Use-cases

GoMove is perfect when you need to prevent session locks during:

- **Long-running processes**: Keep your session active while builds, deployments, or data processing tasks run
- **Presentations**: Prevent screen locks during demos or presentations
- **Remote work**: Maintain connectivity during brief periods away from your desk
- **System monitoring**: Keep dashboards and monitoring tools visible

## Getting started

### Installation

1. Download the latest release for your platform from the [releases page](https://github.com/yourusername/gomove/releases)
2. Extract the binary to a directory in your PATH

### Basic Usage

1. **Create a configuration file**:

   ```bash
   gomove config create
   ```

2. **Start the service**:

   ```bash
   gomove start
   ```

That's it! GoMove will now prevent your session from locking by making subtle mouse movements.

## Configuration

GoMove uses a YAML configuration file located at `~/.gomove/config.yaml`. You can create a default configuration using:

```bash
gomove config create
```

### Configuration Options

```yaml
behavior:
  # Duration of user inactivity before starting activities
  idle_timeout: 10s
  
  # Resume activities after detecting user inactivity
  resume_after_inactivity: true
  
  # Pause activities when user is actively using the computer
  pause_when_user_is_active: true
  
  # Start activities immediately on application startup
  start_on_boot: false

activities:
  - kind: mouse                    # Activity type: mouse or keyboard
    enabled: true                  # Enable/disable this activity
    interval: 5s                   # How often to perform the activity
    pattern: square                # Movement pattern for mouse activity
    schedule:                      # Optional: schedule when this activity runs
      enabled: false
      from: "09:00"
      to: "17:00"
      days: [1, 2, 3, 4, 5]       # Monday-Friday (0=Sunday, 6=Saturday)

# Enable debug logging
debug: false

# Log file path (empty = stdout)
logfile: ""
```

## Available Commands

| Command | Description |
|---------|-------------|
| `gomove start` | Start the mouse movement service |
| `gomove config create` | Create a default configuration file |
| `gomove config show` | Display current configuration |
| `gomove --help` | Show help information |
| `gomove --version` | Show version information |

## Logging

Enable detailed logging for troubleshooting:

```yaml
debug: true
logfile: "/var/log/gomove.log"  # Linux
# logfile: "C:\\logs\\gomove.log"  # Windows
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
