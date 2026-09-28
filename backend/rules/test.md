Client
   |
   | HTTP request
   | "I want a WebSocket connection"
   ↓
Go server
   |
   | Upgrade HTTP → WebSocket
   ↓
WebSocket connection


Go functions commonly return: result + error 

defer means:
Execute this when gameHandler finishes.

                 Multiplayer 3D Game
                        │
          ┌─────────────┴─────────────┐
          │                           │
      Unity Client              Go Server
          │                           │
     ┌────┴────┐                 ┌────┴────┐
     │         │                 │         │
   3D world  UI/HUD          WebSocket   Game logic
     │                           │
 Player movement              Players
 Camera                       Position
 Animations                   Health
                              Combat