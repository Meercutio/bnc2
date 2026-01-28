# 🎮 Bulls & Cows PvP (Realtime MVP)

Real-time browser-based PvP game inspired by **Bulls and Cows**, built with **Go + WebSocket** and a lightweight SPA frontend.

The project focuses on **correct game logic, deterministic server-side state, testability, and production-ready infrastructure**, while keeping the MVP simple and extensible.

---

## ✨ Features

- Real-time PvP gameplay via WebSocket
- Server-authoritative game state
- Deterministic Bulls & Cows scoring (with repeated digits)
- Round-based gameplay with optional timer
- Full round history visible to both players
- Automatic reconnect support (client-side)
- CI/CD-ready setup
- Dockerized and PaaS-friendly (Render/Fly/VPS)

---

## 🎯 Game Rules

- Two players participate in a match.
- Each player sets a **4-digit secret number**:
- leading zeros allowed (`0007`)
- repeated digits allowed (`1122`)
- Game proceeds in **rounds**:
- both players submit guesses simultaneously
- if both submit early — the round ends immediately
- if the timer expires and a player did not submit — the guess is marked as **missed**
- After each round, the server calculates:
- **Bulls** — correct digit in correct position
- **Cows** — correct digit in wrong position (multiset-based)
- **Win condition**:
- first player to guess opponent’s secret wins
- if both guess correctly in the same round → **draw**
- Full round history is shared with both players.

All rules and results are enforced **server-side**.

---

## 🏗 Architecture

### Backend
- Go (`net/http`)
- WebSocket for real-time communication

### Frontend
- Single-page application (HTML/CSS/JS)
- Same-origin HTTP + WebSocket
- UI is **state-driven** (renders only from server `state`)

---

Run
make up
localhost:8080

## 📁 Project Structure

