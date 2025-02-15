# Is It Dead Yet?

A humorous, lightweight service status dashboard that answers the eternal question: "Is it dead yet?" 🤔

![Dashboard Preview](/docs/assets/look1.png)

## Features

- 🔍 Real-time service status monitoring
- ⏰ Status checks every 60 seconds
- 🚨 Visual alerts for down services
- ⌛ Downtime tracking
- 🌓 Dark/Light theme support
- 🐳 Docker support
- 🎯 Simple and lightweight

## Quick Start

### Clone the repository

```bash
git clone https://github.com/yourusername/is-it-dead-yet.git
cd is-it-dead-yet
```

### Configure your services

copy existing json config:

```bash
cp config.json.example config.json
```

or edit with editor of your choice

### Run with Docker Compose

```bash
docker-compose up -d
```

## Configuration

Edit `config.json` to add your services:

```json
{
  "services": [
    {
      "name": "My Service",
      "url": "http://service-url:port",
      "icon": "icon-name"
    }
  ]
}
```

Icons are automatically fetched from [dashboard-icons](https://github.com/walkxcode/dashboard-icons).

## Development

### Prerequisites

- Go 1.21 or higher
- Docker (optional)

### Local Development

Run locally:

```bash
go run main.go
```

# Build binary

```bash
go build -o is-it-dead-yet
```

## License

MIT License - feel free to use it as you wish! 🎉

## Contributing

Contributions are welcome! Feel free to:

1. Fork the repository
2. Create a feature branch
3. Submit a Pull Request

## Acknowledgments

- Icons from [dashboard-icons](https://github.com/walkxcode/dashboard-icons)
- Inspired by every developer asking "Is it dead yet?"

## Author

[Montynau](https://github.com/montynau)

---

Made with ❤️ and a bit of humor
