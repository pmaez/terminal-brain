package app

import (
	"fmt"

	"terminal-brain/internal/core"
	"terminal-brain/internal/modules/cmdmem"
	"terminal-brain/internal/modules/duck"
	"terminal-brain/internal/modules/dump"
	"terminal-brain/internal/modules/stack"
	"terminal-brain/internal/modules/timer"
	"terminal-brain/internal/ui"
)

func Run(dbPath string) error {
	db, err := core.OpenDB(dbPath)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	if err := core.Migrate(db); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	state := &core.AppState{
		DB:      db,
		Stores:  core.NewSQLiteStores(db),
		Modules: map[string]core.Module{},
	}

	if err := registerModules(state); err != nil {
		return err
	}

	model := ui.NewShellModel(state)
	return ui.Run(model)
}

func registerModules(state *core.AppState) error {
	mods := []core.Module{
		stack.New(),
		dump.New(),
		duck.New(),
		cmdmem.New(),
		timer.New(),
	}

	state.Order = make([]string, 0, len(mods))
	for _, mod := range mods {
		name := mod.Name()
		state.Modules[name] = mod
		state.Order = append(state.Order, name)
		if err := mod.Init(state); err != nil {
			return err
		}
	}

	if len(state.Order) > 0 {
		state.Active = state.Order[0]
	}

	return nil
}
