# Minesweeper

REST implementation of minesweeper

## Existing Resources

### Creating a game

`POST /game`

It expects a body of the form:

```typescript
{
  "rows": int
  "columns": int
  "bombs": int
}
```

* `rows`, `columns` and `bombs` must be bigger than zero.
* Game must not be already started.

In any of those cases a HTTP 400 error with the cause will be returned.

If given a valid payload, a HTTP 201 will be returned with a JSON representing the state of the board.

### Querying the state of the game

`GET /game`

It will return a JSON representing the state of the board.

It will fail if game has not started yet.

### "Tapping" (left click) a Tile

`POST /game/tap`

It receives a body of the form:

```typescript
{
  "row": int,
  "column": int,
}
```

* The game must be already started.
* `row` and `column` must be withing the board's limits.
* The Tile at that position must not be revealed already.

In any of those cases a HTTP 400 error with the cause will be returned.

It will Tap (left click) that tile, if it was a bomb the game will be considered finished.

It returns a JSON that contains a `result` field that signals if the game is lost or not, and the current state of the game (if the game is lost, every tile will be revealed).

## Missing Resources

### Delete a Game

### Mark (right click) a Tile

## Missing Features

This code is very WIP, several features are missing:

* Client Library
* A persistant implementation of Game
* Time Tracking
* Multiple accounts
* Multiple games

## Design Decisions

The structure of the API is as follows:

```
minesweeper ->
  cmd ->
    rest -> REST API, it calls providers and eventually runs a Server
      main.go
  internal ->
    operation -> Use cases, it only knows about symbols from `internal` and "orchestates" them (every rest handler should mimic 1:1 one of those)
    platform -> 
      game -> Game interface and implementations
        game.go Represents an existing game, with operations for accessing and syncing a Game
        fake_game.go Game used for tests
        in_memory_game.go Current implementation, it uses a volatile in memory Board.
      provide -> Dependency injection
      random -> math/rand façade
      rest -> gin handlers
        middleware -> gin middlewares
        routes.go gin routing
        server.go Server type, it has singleton dependencies as internal state
    board.go Board type
    draw_board.go Utility for drawing a Board into a string
    tile.go Tile type
```

### Domain

* A Tile represents a square in the game, it can be hidden, revealed or marked (with a Flag or a Question Mark) and can have a bomb or not.
* A Board is a matrix of randomly assigned tiles
* A Game is a Board handler, used for abstracting over a persistent/volatile store.
* A Game Holder stores multiple game sessions.