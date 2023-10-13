# WOPR
Can you beat WOPR at its favorite game?

WOPR is a simple TIC-TAC-TOE shell, the bot must interact with a tcp socket that hosts a game.

## Running

To run in docker, use the following:
```
docker build -t wopr .
docker run --rm -p 4000:4000 -t wopr 
```

__OR__

```
./WOPR
```

## Building
To build a new release:
```
CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath WOPR.go
```

## Solution
WOPR plays randomly so the best strategy is to simple choose 0, 1, 2 and loop until we win that way

```
go run solve/solve.go
```