package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/luthermonson/goodhosts/pkg/hostsfile"

	"github.com/urfave/cli"
)

func Remove() cli.Command {
	return cli.Command{
		Name:      "remove",
		Aliases:   []string{"rm", "r"},
		Usage:     "Remove ip or host(s) if exists",
		Action:    remove,
		ArgsUsage: "[IP|HOST] or [IP] [HOST] ([HOST]...)",
	}
}
func remove(c *cli.Context) error {
	args := c.Args()
	hostsfile, err := loadHostsfile(c)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return cli.NewExitError("No input.", 1)
	}

	if len(args) == 1 { //could be ip or hostname
		return processSingleArg(hostsfile, args[0])
	}

	ip := args[0]
	uniqueHosts := map[string]bool{}
	var hostEntries []string

	for i := 1; i < len(args); i++ {
		uniqueHosts[args[i]] = true
	}

	for key, _ := range uniqueHosts {
		hostEntries = append(hostEntries, key)
	}

	err = hostsfile.Remove(ip, hostEntries...)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}

	err = hostsfile.Flush()
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}

	fmt.Printf("Removed: %s %s\n", ip, strings.Join(hostEntries, " "))
	return nil
}

func processSingleArg(hostsfile hostsfile.Hosts, arg string) error {
	if net.ParseIP(arg) != nil {
		fmt.Printf("Removing ip %s\n", arg)
		if err := hostsfile.RemoveByIp(arg); err != nil {
			return err
		}
		if err := hostsfile.Flush(); err != nil {
			return err
		}

		return nil
	}

	if err := hostsfile.RemoveByHostname(arg); err != nil {
		return err
	}
	if err := hostsfile.Flush(); err != nil {
		return err
	}

	return nil
}