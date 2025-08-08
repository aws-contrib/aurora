package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/aws-contrib/aurora/cmd"
	"github.com/aws-contrib/aurora/internal/database/ent"
	"github.com/aws-contrib/aurora/internal/database/ent/template"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "aurora",
		Usage: "Manage your database schema as code",
		Commands: []*cli.Command{
			{
				Name:  "migrate",
				Usage: "Manage versioned migration files",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Usage:    "select config (project) file using URL format",
						Value:    "file://aurora.hcl",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "env",
						Usage:    "set which env from the config file to use",
						Required: true,
					},
				},
				Commands: []*cli.Command{
					{
						Name:  "apply",
						Usage: "Applies pending migration files on the connected database.",
						Flags: []cli.Flag{
							&cli.DurationFlag{
								Name:  "lock-timeout",
								Usage: "set how long to wait for the database lock",
								Value: 25 * time.Minute,
							},
						},
						Action: func(ctx context.Context, command *cli.Command) error {
							repository, err := NewRepository(ctx, command)
							if err != nil {
								return err
							}

							args := &ent.LockMigrationParams{}
							args.Timeout = command.Duration("lock-timeout")
							// lock the execution
							if xerr := repository.LockMigration(ctx, args); xerr != nil {
								return xerr
							}
							// unlock the execution
							defer repository.UnlockMigration(ctx)

							migrations, err := repository.ListMigrations(ctx, &ent.ListMigrationsParams{})
							if err != nil {
								return err
							}

							state := &ent.MigrationState{}
							// prepare the status
							for _, migration := range migrations {
								if state.Next == nil {
									state.Next = migration.Revision
								}

								if err == nil {
									if state.Current == nil || state.Current.Error == nil {
										params := &ent.ApplyMigrationParams{}
										params.Migration = migration
										// apply the migration
										err = repository.ApplyMigration(ctx, params)
										// update the migration state
										migration = params.Migration
									}
								}

								if migration.Revision.ExecutedAt.IsZero() {
									state.Pending = append(state.Pending, migration.Revision)
								} else {
									state.Executed = append(state.Executed, migration.Revision)
									state.Current = migration.Revision
									state.Next = nil
								}
							}

							// print the status
							template.Execute(os.Stdout, "status", state)

							if state.Current != nil && state.Current.Error != nil {
								// return the error
								return cli.Exit("There are errors in the migration", 1)
							}
							// done!
							return err
						},
					},
					{
						Name:  "status",
						Usage: "Get information about the current migration status.",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "wait",
								Usage: "wait for the database pending migrations to be applied",
								Value: false,
							},
							&cli.DurationFlag{
								Name:  "wait-timeout",
								Usage: "set how long to wait for the database pending migrations",
								Value: 25 * time.Minute,
							},
						},
						Action: func(ctx context.Context, command *cli.Command) error {
							repository, err := NewRepository(ctx, command)
							if err != nil {
								return err
							}

							start := time.Now()
							state := &ent.MigrationState{}

							for {
								migrations, merr := repository.ListMigrations(ctx, &ent.ListMigrationsParams{})
								if merr != nil {
									return merr
								}

								state = &ent.MigrationState{}
								// prepare the status
								for _, migration := range migrations {
									if state.Next == nil {
										state.Next = migration.Revision
									}

									if migration.Revision.ExecutedAt.IsZero() {
										state.Pending = append(state.Pending, migration.Revision)
										// We should exit if there are pending migrations
										err = cli.Exit("There are pending migrations", 1)
									} else {
										state.Executed = append(state.Executed, migration.Revision)
										state.Current = migration.Revision
										state.Next = nil
									}
								}

								if state.Next != nil {
									// Wait for the migrations to be applied
									if command.Bool("wait") {
										if time.Since(start) < command.Duration("wait-timeout") {
											// Give some time for the migrations to be applied
											time.Sleep(250 * time.Millisecond)
											continue
										}
									}
								}
								// we should stop re-trying
								break
							}

							// print the status
							template.Execute(os.Stdout, "status", state)
							// done!
							return err
						},
					},
				},
			},
		},
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		cmd.Version = info.Main.Version
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func NewRepository(ctx context.Context, command *cli.Command) (*ent.MigrationRepository, error) {
	path, err := cmd.GetPath(command.String("config"))
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &cmd.Config{}
	if err := config.UnmarshalText(data); err != nil {
		return nil, err
	}

	name := command.String("env")
	// Get the environment configuration
	if config := config.GetEnvironment(name); config != nil {
		conn, err := config.GetURL()
		if err != nil {
			return nil, err
		}

		directory, err := config.Migration.GetDir()
		if err != nil {
			return nil, err
		}

		directory, err = cmd.GetPath(directory)
		if err != nil {
			return nil, err
		}

		gateway, err := ent.Open(ctx, conn)
		if err != nil {
			return nil, err
		}

		if err := gateway.CreateTableLocks(ctx); err != nil {
			return nil, err
		}

		if err := gateway.CreateTableRevisions(ctx); err != nil {
			return nil, err
		}

		repository := &ent.MigrationRepository{
			Gateway:    gateway,
			FileSystem: os.DirFS(directory),
		}

		return repository, nil
	}

	return nil, fmt.Errorf("environment %s not found in config", name)
}
