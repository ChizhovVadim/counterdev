package arena

import (
	"context"
	"log"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/ChizhovVadim/counterdev/internal/game"
	"golang.org/x/sync/errgroup"
)

func GenerateDataset(
	ctx context.Context,
	gamesCount int,
	openingRandomness int,
	openingSize int,
	openingPlayerBuilder func() Player,
	playerBuilder func() Player,
	outputGamePath string,
) error {
	log.Println("GenerateDataset started")
	defer log.Println("GenerateDataset finished")

	var concurrency = runtime.GOMAXPROCS(0)

	var games = make(chan game.Game)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return saveGames(ctx, outputGamePath, games)
	})

	var gameIndex = new(int32)
	var wg = &sync.WaitGroup{}
	for range concurrency {
		wg.Add(1)
		g.Go(func() error {
			defer wg.Done()
			return playDatasetGames(ctx, gameIndex, int32(gamesCount), openingRandomness, openingSize, openingPlayerBuilder, playerBuilder, games)
		})
	}
	g.Go(func() error {
		wg.Wait()
		close(games)
		return nil
	})

	return g.Wait()
}

func playDatasetGames(
	ctx context.Context,
	gameIndex *int32,
	gamesCount int32,
	openingRandomness int,
	openingSize int,
	openingPlayerBuilder func() Player, // должен поддерживать Randomness!
	playerBuilder func() Player,
	games chan<- game.Game,
) error {
	var openingPlayer = openingPlayerBuilder()
	var player = playerBuilder()
	var repeats = make(map[uint64]struct{})

	for {
		var g, _ = game.NewGame("")
		var ok = playRandomOpening(&g, openingPlayer.Engine, openingRandomness, openingSize, openingPlayer.TimeLimit)
		if !ok {
			continue
		}
		if _, found := repeats[g.Position.Key]; found {
			continue
		}
		repeats[g.Position.Key] = struct{}{}

		var err = playGame(&g, &player, &player)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case games <- g:
			if atomic.AddInt32(gameIndex, 1) >= gamesCount {
				return nil
			}
		}
	}
}
