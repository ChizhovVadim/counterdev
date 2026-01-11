package arena

import (
	"context"
	"iter"
	"log"
	"sync"

	"github.com/ChizhovVadim/counterdev/internal/game"
	"golang.org/x/sync/errgroup"
)

func GenerateDataset(
	ctx context.Context,
	openingsSeq iter.Seq[game.Game],
	concurrency int,
	playerBuilder func() Player,
	outputGamePath string,
) error {
	log.Println("GenerateDataset started")
	defer log.Println("GenerateDataset finished")

	var openings = make(chan game.Game)
	var games = make(chan game.Game)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(openings)
		for opening := range openingsSeq {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case openings <- opening:
			}
		}
		return nil
	})
	g.Go(func() error {
		return saveGames(ctx, outputGamePath, games)
	})

	var wg = &sync.WaitGroup{}
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		g.Go(func() error {
			defer wg.Done()
			var player = playerBuilder()
			for g := range openings {
				var err = playGame(&g, &player, &player)
				if err != nil {
					return err
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case games <- g:
				}
			}
			return nil
		})
	}
	g.Go(func() error {
		wg.Wait()
		close(games)
		return nil
	})

	return g.Wait()
}
