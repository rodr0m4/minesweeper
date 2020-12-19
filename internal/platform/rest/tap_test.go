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

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func Test_Tap_Should_Fail_When_Passed_Invalid_JSON(t *testing.T) {
	handler := TapHandler{}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	registerTap(r, handler)

	req := newTapRequestFromBytes([]byte("<title/>"))

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func Test_Tap_Should_Fail_When_Tapper_Fails(t *testing.T) {
	row := 1
	column := 1

	tapper := TapperThatExpectsAndReturns(t, row, column, func() (internal.TapResult, error) {
		return internal.NothingResult, errors.New("oh no")
	})

	handler := TapHandler{
		Tapper: tapper,
	}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	registerTap(r, handler)

	r.ServeHTTP(rr, newTapRequest(row, column))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func Test_Tap_Should_Fail_If_Tap_Passes_And_Shower_Fails(t *testing.T) {
	row := 1
	column := 1

	tapper := TapperThatExpectsAndReturns(t, row, column, func() (internal.TapResult, error) {
		return internal.NothingResult, nil
	})

	shower := gameShowerMock(func(g game.Game) (operation.ShowedGame, error) {
		return operation.ShowedGame{}, errors.New("oh no")
	})

	handler := TapHandler{
		Tapper:     tapper,
		GameShower: shower,
	}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	registerTap(r, handler)

	r.ServeHTTP(rr, newTapRequest(1, 1))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func Test_Tap_Should_Convert_TapResult_To_String(t *testing.T) {
	type Case struct {
		given  internal.TapResult
		expect string
	}

	cases := []Case{
		{
			given:  internal.NothingResult,
			expect: "Nothing",
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

			tapper := TapperThatExpectsAndReturns(t, row, column, func() (internal.TapResult, error) {
				return tt.given, nil
			})

			shower := gameShowerMock(func(g game.Game) (operation.ShowedGame, error) {
				return operation.ShowedGame{}, nil
			})

			handler := TapHandler{
				Tapper:     tapper,
				GameShower: shower,
			}

			rr := httptest.NewRecorder()
			_, r := gin.CreateTestContext(rr)

			registerTap(r, handler)

			r.ServeHTTP(rr, newTapRequest(row, column))

			assert.Equal(t, http.StatusOK, rr.Code)

			var response tapHandlerResponse
			_ = json.Unmarshal(rr.Body.Bytes(), &response)

			assert.Equal(t, tt.expect, response.Result)
		})
	}
}

// Helpers

func newTapRequestFromBytes(buf []byte) *http.Request {
	return httptest.NewRequest(http.MethodPatch, "/games", bytes.NewBuffer(buf))
}

func newTapRequest(row, column int) *http.Request {
	body := gin.H{
		"row":    row,
		"column": column,
	}
	buf, _ := json.Marshal(body)
	return newTapRequestFromBytes(buf)
}

func registerTap(r *gin.Engine, handler TapHandler) {
	r.Use(middleware.ErrorLogger())
	r.PATCH("/games", handler.Tap)
}

// Mocks

type TapperMock func(game game.Game, row, column int) (internal.TapResult, error)

func (t TapperMock) Tap(game game.Game, row, column int) (internal.TapResult, error) {
	return t(game, row, column)
}

func TapperThatExpectsAndReturns(t *testing.T, row, column int, f func() (internal.TapResult, error)) TapperMock {
	return func(game game.Game, actualRow, actualColumn int) (internal.TapResult, error) {
		assert.Equal(t, row, actualRow)
		assert.Equal(t, column, actualColumn)
		return f()
	}
}
