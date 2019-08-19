package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	gitlab "github.com/xanzy/go-gitlab"
	lab "github.com/zaquestion/lab/internal/gitlab"
)

var (
	mrLabels       []string
	mrState        string
	mrTargetBranch string
	mrNumRet       int
	mrAll          bool
	ciStatus       bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list [remote]",
	Aliases: []string{"ls"},
	Short:   "List merge requests",
	Long:    ``,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rn, _, err := parseArgs(args)
		if err != nil {
			log.Fatal(err)
		}

		num := mrNumRet
		if mrAll {
			num = -1
		}
		mrs, err := lab.MRList(rn, gitlab.ListProjectMergeRequestsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: mrNumRet,
			},
			Labels:       mrLabels,
			State:        &mrState,
			TargetBranch: &mrTargetBranch,
			OrderBy:      gitlab.String("updated_at"),
		}, num)
		if err != nil {
			log.Fatal(err)
		}
		for _, mr := range mrs {
			fmt.Printf("#%d %s", mr.IID, mr.Title)
			if ciStatus {
				fmt.Printf("\t(%s)", mr.Pipeline.Status)
			}
			fmt.Println("")
		}
	},
}

func init() {
	listCmd.Flags().StringSliceVarP(
		&mrLabels, "label", "l", []string{}, "filter merge requests by label")
	listCmd.Flags().StringVarP(
		&mrState, "state", "s", "opened",
		"filter merge requests by state (opened/closed/merged)")
	listCmd.Flags().IntVarP(
		&mrNumRet, "number", "n", 10,
		"number of merge requests to return")
	listCmd.Flags().StringVarP(
		&mrTargetBranch, "target-branch", "t", "",
		"filter merge requests by target branch")
	listCmd.Flags().BoolVarP(&mrAll, "all", "a", false, "List all MRs on the project")
	listCmd.Flags().BoolVarP(&ciStatus, "ci-status", "c", false, "Include CI Status in the results")

	listCmd.MarkZshCompPositionalArgumentCustom(1, "__lab_completion_remote")
	listCmd.MarkFlagCustom("state", "(opened closed merged)")
	mrCmd.AddCommand(listCmd)
}
