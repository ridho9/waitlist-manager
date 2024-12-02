# Waitlist Manager

## Setup

### Dependencies

1. Node and Bun (npm should work but I developed with bun)
2. Go (I use 1.23)
3. Docker and Docker Compose

### Steps

1. Ports used: 8000, 8001, 8002, 8502. make sure this ports are unused
2. Setup db and caddy server

```
docker compose up -d
```

3. In one terminal, run the backend

```
cd backend-go
go run .

```

4. In another terminal, setup and run the frontend

```
cd frontend
bun install
bun run dev
```

5. Open [localhost:8000/book](localhost:8000/book)

6. To clean up, kill frontend and backend, and run `docker compose down`

## Docs

1. Task description [TASK.md](./TASK.md)
2. Initial Design [DESIGN.md](./DESIGN.md)

## Afterthoughts on Design

I initially designed the system to use pubsub so that the clients can listen to changes. 
But I ended up just using polling each seconds because it's easier and less error prone.

For cheking queue and serving seat, I used a goroutine worker running in the background.
These workers can be thought as a real worker monitoring the queue, and the waitress in the restaurant.
These workers can later be swapped out to a system that is able to be operated by a real waitress/system.

## Disclaimer 

Up until commit `d5d5d16` I worked with my own effort (the system is all completed), after that I used v0 to help with the design.