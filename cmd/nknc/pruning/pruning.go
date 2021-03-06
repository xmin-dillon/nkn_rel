package pruning

import (
	"fmt"

	"github.com/nknorg/nkn/v2/chain/store"
	"github.com/urfave/cli"
)

func pruningAction(c *cli.Context) error {
	if c.NumFlags() == 0 {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	switch {
	case c.Bool("currentheight"):
		cs, err := store.NewLedgerStore()
		if err != nil {
			return err
		}
		_, h, err := cs.GetCurrentBlockHashFromDB()
		if err != nil {
			return err
		}
		fmt.Println(h)
	case c.Bool("startheights"):
		cs, err := store.NewLedgerStore()
		if err != nil {
			return err
		}
		refCountStartHeight, pruningStartHeight := cs.GetPruningStartHeight()
		fmt.Println(refCountStartHeight, pruningStartHeight)
	case c.Bool("pruning"):
		cs, err := store.NewLedgerStore()
		if err != nil {
			return err
		}

		if c.Bool("seq") {
			err := cs.SequentialPrune()
			if err != nil {
				return err
			}
		} else if c.Bool("lowmem") {
			err := cs.PruneStatesLowMemory(true)
			if err != nil {
				return err
			}
		} else {
			err := cs.PruneStates()
			if err != nil {
				return err
			}
		}
	case c.Bool("dumpnodes"):
		cs, err := store.NewLedgerStore()
		if err != nil {
			return err
		}
		err = cs.TrieTraverse(true)
		if err != nil {
			return err
		}
	case c.Bool("verifystate"):
		cs, err := store.NewLedgerStore()
		if err != nil {
			return err
		}
		err = cs.VerifyState()
		if err != nil {
			return err
		}
	default:
		cli.ShowSubcommandHelp(c)
		return nil
	}

	return nil
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:        "pruning",
		Usage:       "state trie pruning for nknd",
		Description: "state trie pruning for nknd.",
		ArgsUsage:   "[args]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "currentheight",
				Usage: "current block height of offline db",
			},
			cli.BoolFlag{
				Name:  "startheights",
				Usage: "start heights of refcount and pruning",
			},
			cli.BoolFlag{
				Name:  "pruning",
				Usage: "prune state trie",
			},
			cli.BoolFlag{
				Name:  "seq",
				Usage: "prune state trie sequential mode",
			},
			cli.BoolFlag{
				Name:  "lowmem",
				Usage: "prune state trie low memory mode",
			},
			cli.BoolFlag{
				Name:  "dumpnodes",
				Usage: "dump nodes of trie",
			},
			cli.BoolFlag{
				Name:  "verifystate",
				Usage: "verify state of ledger",
			},
		},
		Action: pruningAction,
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			return cli.NewExitError("", 1)
		},
	}
}
