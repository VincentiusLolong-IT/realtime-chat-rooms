# Realtime Chat Rooms

A realtime chat system built with **Go** and **WebSockets** that supports multiple chat rooms.  
Users can join different rooms and exchange messages concurrently with safe client management.

## Features
- WebSocket connection handling
- Multi-room chat logic
- Concurrent client management (Mutex / RWMutex)
- Room-based message broadcasting
- Client connect & disconnect handling
- Scalable Hub → Room → Client architecture

## Purpose
Provide a realtime communication foundation that can be extended with:
- Presence / online status
- Private messaging
- Message persistence / database storage

## Tech Stack
- Go (Golang)
- Gorilla WebSocket
- Native Go Concurrency (goroutines, channels, mutex)

## Future Improvements
- Authentication & authorization
- Message history storage
- Typing indicators
- Notifications
- Horizontal scaling with Redis / PubSub
