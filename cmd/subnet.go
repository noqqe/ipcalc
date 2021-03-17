package cmd

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var subnetCmd = &cobra.Command{
	Aliases:       []string{"sub"},
	Use:           "subnet <prefix> <n>",
	Short:         "Split the specified prefix into \"n\" sub-networks",
	SilenceErrors: true,
	Args:          cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var sub []string

		_, p := ParseCIDR(args[0])
		if p == nil {
			return fmt.Errorf("invalid prefix: %s", args[0])
		}

		n, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid number: %s", args[1])
		}

		log := math.Log2(float64(n))
		if n < 2 || log-float64(int(log)) != 0 {
			return fmt.Errorf("\"n\" must be a power of 2")
		}

		if verbose {
			for _, sp := range p.Subnets(int(log)) {
				fmt.Printf("> %s%s%s\n", Purple, sp.String(), Reset)
				_, p := ParseCIDR(sp.String())
				fmt.Printf("  Addresses:  %s%d%s\n", Pink, p.NumNodes(), Reset)
				fmt.Printf("  Netmask:    %s%s%s\n", Yellow, Explode(p.Mask), Reset)
				fmt.Printf("  First:      %s%s%s\n", Green, Explode(p.IP), Reset)
				fmt.Printf("  Last:       %s%s%s\n", Green, Explode(p.Last()), Reset)
				fmt.Printf("\n")
			}
		} else {
			for _, sp := range p.Subnets(int(log)) {
				sub = append(sub, sp.String())
			}
		  fmt.Printf("> %s%s%s\n", Purple, strings.Join(sub, ", "), Reset)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(subnetCmd)
}
