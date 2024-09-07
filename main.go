package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type Package struct {
	Name       string     `json:"name"`
	Owner      Owner      `json:"owner"`
	Repository Repository `json:"repository"`
}

type Owner struct {
	Login string `json:"login"`
}

type Repository struct {
	FullName string `json:"full_name"`
}

type Version struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

type OutputVersion struct {
	Package string `json:"package"`
	Name    string `json:"name"`
	ID      int64  `json:"id"`
}

func main() {
	var rootCmd = &cobra.Command{Use: "ghcr-cli"}

	var listCmd = &cobra.Command{
		Use:   "list [owner/repo]",
		Short: "List Docker images in GitHub Container Registry",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := checkGHInstalled(); err != nil {
				log.Fatal(err)
			}

			owner, repo, err := parseOwnerRepo(args[0])
			if err != nil {
				log.Fatal(err)
			}

			packageName, err := getPackageName(owner, repo)
			if err != nil {
				log.Fatal(err)
			}

			versions, err := getPackageVersions(packageName)
			if err != nil {
				log.Fatal(err)
			}

			outputVersions := make([]OutputVersion, len(versions))
			for i, version := range versions {
				outputVersions[i] = OutputVersion{
					Package: packageName,
					Name:    version.Name,
					ID:      version.ID,
				}
			}

			jsonOutput, err := json.MarshalIndent(outputVersions, "", "  ")
			if err != nil {
				log.Fatalf("Error creating JSON output: %v", err)
			}

			fmt.Println(string(jsonOutput))
		},
	}

	var deleteCmd = &cobra.Command{
		Use:   "delete [owner/repo] [version-id]",
		Short: "Delete a Docker image version from GitHub Container Registry",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := checkGHInstalled(); err != nil {
				log.Fatal(err)
			}

			owner, repo, err := parseOwnerRepo(args[0])
			if err != nil {
				log.Fatal(err)
			}

			packageName, err := getPackageName(owner, repo)
			if err != nil {
				log.Fatal(err)
			}

			versionID := args[1]

			if !confirmDeletion(packageName, versionID) {
				fmt.Println("Deletion cancelled.")
				return
			}

			if err := deletePackageVersion(packageName, versionID); err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Successfully deleted image %s version %s\n", packageName, versionID)
		},
	}

	rootCmd.AddCommand(listCmd, deleteCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getPackageName(owner, repo string) (string, error) {
	// https://docs.github.com/ja/rest/packages/packages?apiVersion=2022-11-28#list-packages-for-the-authenticated-users-namespace
	cmd := exec.Command("gh", "api",
		"-H", "Accept: application/vnd.github+json",
		"-H", "X-GitHub-Api-Version: 2022-11-28",
		"/user/packages?package_type=container")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error listing packages: %v\nOutput: %s", err, string(output))
	}

	var packages []Package
	if err := json.Unmarshal(output, &packages); err != nil {
		return "", fmt.Errorf("error parsing JSON: %v", err)
	}

	for _, pkg := range packages {
		if pkg.Repository.FullName == fmt.Sprintf("%s/%s", owner, repo) {
			return pkg.Name, nil
		}
	}

	return "", fmt.Errorf("package not found for %s/%s", owner, repo)
}

func getPackageVersions(packageName string) ([]Version, error) {
	// https://docs.github.com/ja/rest/packages/packages?apiVersion=2022-11-28#list-package-versions-for-a-package-owned-by-the-authenticated-user
	cmd := exec.Command("gh", "api",
		"-H", "Accept: application/vnd.github+json",
		"-H", "X-GitHub-Api-Version: 2022-11-28",
		fmt.Sprintf("/user/packages/container/%s/versions", packageName))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error getting versions: %v\nOutput: %s", err, string(output))
	}

	var versions []Version
	if err := json.Unmarshal(output, &versions); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return versions, nil
}

func deletePackageVersion(packageName, versionID string) error {
	cmd := exec.Command("gh", "api",
		"-H", "Accept: application/vnd.github+json",
		"-H", "X-GitHub-Api-Version: 2022-11-28",
		"-X", "DELETE",
		fmt.Sprintf("/user/packages/container/%s/versions/%s", packageName, versionID))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error deleting version: %v\nOutput: %s", err, string(output))
	}

	return nil
}

func parseOwnerRepo(input string) (string, string, error) {
	parts := strings.Split(input, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid owner/repo format. Expected 'owner/repo'")
	}
	return parts[0], parts[1], nil
}

func checkGHInstalled() error {
	_, err := exec.LookPath("gh")
	if err != nil {
		return fmt.Errorf("GitHub CLI (gh) is not installed or not in PATH. Please install it and try again")
	}
	return nil
}

func confirmDeletion(packageName, versionID string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Are you sure you want to delete image %s version %s? (y/N): ", packageName, versionID)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
