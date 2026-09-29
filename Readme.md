# Multiplayer 2D Game (Godot + Go)

A small real-time multiplayer prototype. The **client** is a 2D platformer built with Godot 4, and the **server** is a Go WebSocket server that assigns player IDs and relays player events and positions between clients.

## Features

- Players connect over WebSocket and receive a unique ID
- Every client sees all other players spawn when they join and disappear when they leave
- Live position syncing (20 updates/second) with smoothing for remote players
- Only the local player reads keyboard input; remote players follow network updates

## Tech stack

| Part   | Technology                                                    |
| ------ | ------------------------------------------------------------- |
| Client | Godot 4.3+ (GDScript, `CharacterBody2D`, `WebSocketPeer`)     |
| Server | Go, [`gorilla/websocket`](https://github.com/gorilla/websocket) |

## Project structure

```
.
├── server/
│   └── main.go            # Go WebSocket server
└── client/                # Godot project
    ├── scenes/
    │   ├── game.tscn      # Main scene: tilemap, camera, "Players" container
    │   └── player.tscn    # Player: CharacterBody2D + animated sprite + collision
    └── scripts/
        ├── Network.gd     # Autoload: WebSocket connection + message handling
        ├── Game.gd        # Spawns / removes / updates players
        └── Movement.gd    # Local input + position sending; remote smoothing
```

> Adjust the paths above if your folders are laid out differently.

## How it works

```
 Godot client A ──┐                          ┌── Godot client B
                  ├── WebSocket  /game ──►  Go server
 Godot client C ──┘   (JSON messages)        └── ...
```

1. A client connects to `ws://<host>:8080/game`.
2. The server assigns an ID and replies with `welcome`, then sends one `player_joined` for every player already connected.
3. The server tells everyone else that the new player joined.
4. While playing, each client sends its position; the server relays it to all other clients.
5. When a client disconnects, the server broadcasts `player_left`.

### Message protocol

All messages are JSON text frames with a `type` field.

**Server → client**

| Type           | Fields                    | Meaning                                  |
| -------------- | ------------------------- | ---------------------------------------- |
| `welcome`      | `player_id`               | Your assigned ID                         |
| `player_joined`| `player_id`               | Another player is in the game            |
| `player_left`  | `player_id`               | A player disconnected                    |
| `player_moved` | `player_id`, `x`, `y`     | A player's new position                  |

**Client → server**

| Type   | Fields | Meaning                     |
| ------ | ------ | --------------------------- |
| `move` | `x`, `y` | The local player's position |

## Getting started

### Prerequisites

- [Go](https://go.dev/dl/) 1.18+ (the server uses `any`)
- [Godot](https://godotengine.org/download) 4.3 or newer

### 1. Run the server

```bash
cd server
go mod init game-server        # first time only
go get github.com/gorilla/websocket
go run main.go
```

The server listens on `0.0.0.0:8080`, endpoint `/game`.

### 2. Point the client at the server

In `scripts/Network.gd`, set the server address:

```gdscript
const SERVER_URL = "ws://192.168.100.30:8080/game"
```

Use `ws://localhost:8080/game` if the server runs on the same machine, or the server's LAN IP if you play across devices.

### 3. Register the autoload

In Godot: **Project → Project Settings → Globals → Autoload**, add `scripts/Network.gd` with the name `Network`.

### 4. Run the game

1. Open `game.tscn`, select the root node, and make sure `Game.gd` is attached with `player_scene` set to `player.tscn`.
2. Make sure `game.tscn` has a `Players` (`Node2D`) child and no hardcoded `Player` node.
3. Set `game.tscn` as the main scene.
4. To test multiplayer on one machine: **Debug → Run Multiple Instances → Run 2 Instances**, then press play.

### Controls

| Action | Key                 |
| ------ | ------------------- |
| Move   | Left / Right arrows |
| Jump   | Space / Enter (`ui_accept`) |

## Troubleshooting

- **"Failed to connect" / stuck on "Connecting to server..."**: check `SERVER_URL`, that the server is running, and that port 8080 is allowed through your firewall.
- **Players don't appear**: check that `Network` is registered as an autoload and that `player_scene` is set on the `Game` node.
- **All players move with my keyboard**: `is_local` must be set only for the local player (done in `Game.gd`).

## Known limitations / roadmap

- [ ] Send existing players' last known positions to newly joined players
- [ ] Server-side validation of positions (currently trusts clients)
- [ ] Sync animations and facing direction
- [ ] Player names and a lobby / rooms
- [ ] Reconnect handling
- [ ] TLS (`wss://`) for internet play
