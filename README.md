# Minesweeper

REST implementation of minesweeper

## How to Run Locally

You can run the code locally with `go run cmd/rest/main.go`.

## How to browse through the code

You can first check in `routes.go` for the Handler method that represents that resource (if there are dependencies you can check for them in the `provide` package).

## Configuration

* If the `PORT` env var is set it will use that port, otherwise it will default to `8080`
* If the `REVEAL_EVERYTHING` env var is set to `true`, every time that the board would be shown in the hidden internal state will be revealed (`X` for bomb and `O` for no bomb. `F` for Flag, `?` for Question Mark, `H` for Hidden and `R` for Revealed). Otherwise, the game would be played like the desktop version.
* If the `SHOW_TILES` env var is set to `true`, every time that the API would show information about a game, it will display the board as a JSON array of Tile entities (optimized for consumption as a REST resource). It is `false` by default.
* If the `SHOW_LINES` env var is set to `false` every time that the API would show information about a game, it will display the board as a JSON array of strings that simulate the desktop game (optimized for playing directly by making requests to the API). It is `true` by default.

### Why env vars?

`REVEAL_EVERYTING`, `SHOW_TILES` and `SHOW_LINES` are not env vars because I think it is a very good decision, but because it was easier to iterate through the design using them as such. The "best" implementation would have:

* `REVEAL_EVERYTHING` as a permission for a super user .
* `SHOW_TILES` and `SHOW_LINES` as preferences stored in a persistent store or passed by header.


## Existing Resources

The API url is:
* https://minesweeper-rodroma.herokuapp.com/ for the Heroku Dyno with the last changes in master
* http://localhost:$PORT/ for the local version

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

#### An example

`curl --data '{"rows": 10, "columns": 10, "bombs": 20}' --request POST 'https://minesweeper-rodroma.herokuapp.com/games' | jq .`

### Querying the state of the game

`GET /game`

It will return a JSON representing the state of the board.

It will fail if game has not started yet.

#### An example

`curl 'https://minesweeper-rodroma.herokuapp.com/games/0' | jq .`

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

In any of those cases an HTTP 400 error with the cause will be returned.

It will Tap (left click) that tile, if it was a bomb the game will be considered finished.

It returns a JSON that contains a `result` field that signals if the game is lost or not, and the current state of the game (if the game is lost, every tile will be revealed).

#### An example

`curl --data '{"row": 0, "column": 6}' --request POST 'https://minesweeper-rodroma.herokuapp.com/games/0/tap' | jq .`

### Mark (right click) a Tile

`POST /game/mark`

It receives a body of the form:

```typescript
{
  "row": int,
  "column": int,
  "mark": "flag" | "question"
}
```

* The game must be already started.
* `row` and `column` must be withing the board's limits.
* `mark` should be a valid mark.
* The Tile at that position must not be revealed already.

In any of those cases an HTTP 400 error with the cause will be returned.

It will Mark (right click) that tile with the given mark

It returns a JSON with the current state of the game.

#### An example

`curl --data '{"row": 0, "column": 6, "mark": "question"}' --request POST 'http://localhost:8080/games/0/mark' | jq .`

## Missing Resources

* Delete a Game
* Sign up
* Log in

## Missing Features

This code is very WIP, several features are missing:

* Client Library
* A persistent implementation of Game
* Time Tracking
* Multiple accounts

## API Structure

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

## The Domain

* A Tile represents a square in the game, it can be hidden, revealed or marked (with a Flag or a Question Mark) and can have a bomb or not.
* A Board is a matrix of randomly assigned tiles
* A Game is a Board handler, used for abstracting over a persistent/volatile store.
* A Game Holder stores multiple game sessions.