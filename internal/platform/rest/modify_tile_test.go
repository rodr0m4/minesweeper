package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"minesweeper/internal"
	"minesweeper/internal/operation"
	"minesweeper/internal/platform/game"
	"minesweeper/internal/platform/rest/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_Tap_Should_Fail_When_Passed_Invalid_JSON(t *testing.T) {
	handler := ModifyTileHandler{}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	registerModifyTile(r, handler)

	req := newTapRequestFromBytes([]byte("<title/>"))

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func Test_Tap_Should_Fail_When_Tapper_Fails(t *testing.T) {
	row := 1
	column := 1

	handler := ModifyTileHandler{
		Tapper: tapperThatExpectsAndReturns(t, row, column, func() (internal.TapResult, error) {
			return internal.NothingResult, errors.New("oh no")
		}),
	}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	registerModifyTile(r, handler)

	r.ServeHTTP(rr, newTapRequest(row, column))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func Test_Mark_Should_Fail_With_Bad_Request_When_Invalid_Mark(t *testing.T) {
	handler := ModifyTileHandler{}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	registerModifyTile(r, handler)

	r.ServeHTTP(rr, newMarkRequest(1, 1, "other"))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func Test_Mark_Should_Fail_When_Marker_Fails(t *testing.T) {
	row := 1
	column := 1

	handler := ModifyTileHandler{
		Marker: markerThatExpectsAndReturns(t, row, column, internal.FlagMark, errors.New("oh no")),
	}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	registerModifyTile(r, handler)

	r.ServeHTTP(rr, newMarkRequest(row, column, "flag"))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func Test_Mark_Should_Call_ShowBoard_When_Marker_Does_Not_Fail(t *testing.T) {
	row := 1
	column := 1

	handler := ModifyTileHandler{
		Game:   gameWhoseBoardSucceedsWith(internal.Board{}),
		Marker: markerThatExpectsAndReturns(t, row, column, internal.QuestionMark, nil),
		BoardDrawer: boardDrawerMock(func(board internal.Board) operation.ShowedGame {
			return operation.ShowedGame{}
		}),
	}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	registerModifyTile(r, handler)

	r.ServeHTTP(rr, newMarkRequest(row, column, "question"))

	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_Tap_Should_Convert_TapResult_To_String(t *testing.T) {
	type Case struct {
		given  internal.TapResult
		expect string
	}

	cases := []Case{
		{
			given:  internal.NothingResult,
			expect: "Nothing Happened",
		},
		{
			given:  internal.ExplosionResult,
			expect: "BOOM!",
		},
	}

	for _, tt := range cases {
		name := fmt.Sprintf("should serialize %d into %s correctly", tt.given, tt.expect)

		t.Run(name, func(t *testing.T) {
			row := 1
			column := 1

			tapper := tapperThatExpectsAndReturns(t, row, column, func() (internal.TapResult, error) {
				return tt.given, nil
			})

			drawer := boardDrawerMock(func(board internal.Board) operation.ShowedGame {
				return operation.ShowedGame{}
			})

			handler := ModifyTileHandler{
				Game:        gameWhoseBoardSucceedsWith(internal.Board{}),
				Tapper:      tapper,
				BoardDrawer: drawer,
			}

			rr := httptest.NewRecorder()
			_, r := gin.CreateTestContext(rr)

			registerModifyTile(r, handler)

			r.ServeHTTP(rr, newTapRequest(row, column))

			assert.Equal(t, http.StatusOK, rr.Code)

			var response modifyTileResponse
			_ = json.Unmarshal(rr.Body.Bytes(), &response)

			assert.Equal(t, tt.expect, response.Result)
		})
	}
}

// Helpers

func newTapRequestFromBytes(buf []byte) *http.Request {
	return httptest.NewRequest(http.MethodPost, "/game/tap", bytes.NewBuffer(buf))
}

func newMarkRequestFromBytes(buf []byte) *http.Request {
	return httptest.NewRequest(http.MethodPost, "/game/mark", bytes.NewBuffer(buf))
}

func newTapRequest(row, column int) *http.Request {
	body := gin.H{
		"row":    row,
		"column": column,
	}
	buf, _ := json.Marshal(body)
	return newTapRequestFromBytes(buf)
}

func newMarkRequest(row, column int, mark string) *http.Request {
	body := gin.H{
		"row":    row,
		"column": column,
		"mark":   mark,
	}
	buf, _ := json.Marshal(body)
	return newMarkRequestFromBytes(buf)
}

func registerModifyTile(r *gin.Engine, handler ModifyTileHandler) {
	r.Use(middleware.ErrorLogger())
	r.POST("/game/tap", handler.Tap)
	r.POST("/game/mark", handler.Mark)
}

// Mocks

type TapperMock func(game game.Game, row, column int) (internal.TapResult, error)

func (t TapperMock) Tap(game game.Game, row, column int) (internal.TapResult, error) {
	return t(game, row, column)
}

type markerMock func(game game.Game, row, column int, mark internal.TileMark) error

func (m markerMock) Mark(game game.Game, row, column int, mark internal.TileMark) error {
	return m(game, row, column, mark)
}

func tapperThatExpectsAndReturns(t *testing.T, row, column int, f func() (internal.TapResult, error)) TapperMock {
	return func(game game.Game, actualRow, actualColumn int) (internal.TapResult, error) {
		assert.Equal(t, row, actualRow)
		assert.Equal(t, column, actualColumn)
		return f()
	}
}

func markerThatExpectsAndReturns(t *testing.T, row, column int, mark internal.TileMark, err error) markerMock {
	return func(game game.Game, actualRow, actualColumn int, actualMark internal.TileMark) error {
		assert.Equal(t, row, actualRow)
		assert.Equal(t, column, actualColumn)
		assert.Equal(t, mark, actualMark)

		return err
	}
}
