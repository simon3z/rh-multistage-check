package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"gopkg.in/yaml.v2"
)

type ArrayFlag []string

func (i *ArrayFlag) String() string {
	return strings.Join(*i, ", ")
}

func (i *ArrayFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var cmdFlags = struct {
	Branch       string
	MetadataFile string
	Repositories ArrayFlag
}{}

func branchReferenceFormat(branch string) string {
	return fmt.Sprintf("refs/heads/%s", branch)
}

func init() {
	flag.StringVar(&cmdFlags.Branch, "branch", "", "Branch to check")
	flag.StringVar(&cmdFlags.MetadataFile, "metadata", "container.yaml", "Metadata file to check")
	flag.Var(&cmdFlags.Repositories, "repository", "Repositories to checkout and check")
}

type ContainerMetadata struct {
	Go struct {
		Modules []map[string]string
	}
	Compose struct {
		Packages []string
	}
}

func main() {
	flag.Parse()

	for _, repository := range cmdFlags.Repositories {
		fmt.Printf("\u25A0 Verifying repository: %s (branch: %s)\n", repository, cmdFlags.Branch)

		fs := memfs.New()
		storer := memory.NewStorage()

		_, err := git.Clone(storer, fs, &git.CloneOptions{
			URL:           repository,
			ReferenceName: plumbing.ReferenceName(branchReferenceFormat(cmdFlags.Branch)),
			SingleBranch:  true,
		})

		if err != nil {
			panic(err)
		}

		fmt.Printf("  \u2713 Checkout completed\n")

		f, err := fs.Open(cmdFlags.MetadataFile)

		if err != nil {
			panic(err)
		}

		defer f.Close()

		d, err := ioutil.ReadAll(f)

		if err != nil {
			panic(err)
		}

		metadata := ContainerMetadata{}

		err = yaml.Unmarshal([]byte(d), &metadata)

		if err != nil {
			panic(err)
		}

		if len(metadata.Go.Modules) == 0 {
			fmt.Printf("  \u2717 Error: GO Modules not found!\n")
		} else {
			fmt.Printf("  \u2713 GO Modules found\n")
		}
	}
}
