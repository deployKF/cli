package generate

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/v50/github"
)

type SourceHelper struct {
	GithubOwner             string // the owner of the generator source GitHub repository
	GithubRepo              string // the name of the generator source GitHub repository
	GeneratorArtifactPrefix string // the file-prefix of the generator source zip artifact
	GeneratorArtifactSuffix string // the file-suffix of the generator source zip artifact
	AssetsCacheDir          string // the sub-path under `os.UserHomeDir()` where zip artifacts will be cached
}

type SourceHelperOptions func(*SourceHelper)

func WithGithubOwner(owner string) SourceHelperOptions {
	return func(gh *SourceHelper) {
		gh.GithubOwner = owner
	}
}

func WithGithubRepo(repo string) SourceHelperOptions {
	return func(gh *SourceHelper) {
		gh.GithubRepo = repo
	}
}

func NewSourceHelper(opts ...SourceHelperOptions) *SourceHelper {
	gh := &SourceHelper{
		GithubOwner:             "deployKF",
		GithubRepo:              "deployKF",
		GeneratorArtifactPrefix: "deploykf-",
		GeneratorArtifactSuffix: "-generator.zip",
		AssetsCacheDir:          ".deploykf/assets",
	}

	for _, opt := range opts {
		opt(gh)
	}

	return gh
}

// DownloadAndUnpackSource downloads the generator source artifact for the specified version (if it's not already cached),
// unpacks it to the provided folder, then returns the local path of the artifact .zip file.
func (h *SourceHelper) DownloadAndUnpackSource(version string, unpackTargetDir string, out io.Writer) (string, error) {
	assetsCacheDir, err := h.prepareAssetsCacheDir()
	if err != nil {
		return "", err
	}

	// download the artifact, if it's not cached
	artifactName := h.GeneratorArtifactPrefix + version + h.GeneratorArtifactSuffix
	artifactIsCached, artifactPath, err := h.isArtifactCached(assetsCacheDir, artifactName)
	if err != nil {
		return "", err
	}
	if !artifactIsCached {
		fmt.Fprintf(out, "Downloading deployKF generator source version '%s' from github repo '%s/%s'\n", version, h.GithubOwner, h.GithubRepo)

		// get the GitHub release for the specified version
		githubRelease, err := h.getReleaseByVersion(version)
		if err != nil {
			return "", err
		}

		// find the artifact in the release
		var githubAsset *github.ReleaseAsset
		for _, asset := range githubRelease.Assets {
			if *asset.Name == artifactName {
				githubAsset = asset
				break
			}
		}
		if githubAsset == nil {
			return "", fmt.Errorf("generator artifact '%s' not found in release '%s'", artifactName, *githubRelease.TagName)
		}

		// download the artifact
		err = h.downloadReleaseAsset(githubAsset, artifactPath)
		if err != nil {
			return "", err
		}
	}

	// unzip the artifact
	fmt.Fprintf(out, "Using cached deployKF generator source: %s\n", artifactPath)
	err = UnzipFile(artifactPath, unpackTargetDir, "generator")
	if err != nil {
		return "", err
	}

	return artifactPath, nil
}

// prepareAssetsCacheDir creates the assets cache directory (if it doesn't exist), and returns the path.
func (h *SourceHelper) prepareAssetsCacheDir() (string, error) {
	// use an os-specific cache directory for the downloaded artifact
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	assetsCacheDir := filepath.Join(homeDir, h.AssetsCacheDir)

	// create the assets cache directory if it doesn't exist
	err = os.MkdirAll(assetsCacheDir, 0755)
	if err != nil {
		return "", err
	}

	return assetsCacheDir, nil
}

// isArtifactCached checks if a specific artifact is already cached within the assets cache directory.
func (h *SourceHelper) isArtifactCached(assetsCacheDir string, artifactName string) (bool, string, error) {
	artifactPath := filepath.Join(assetsCacheDir, artifactName)
	artifactIsCached, err := FileExists(artifactPath)
	if err != nil {
		return false, "", err
	}
	return artifactIsCached, artifactPath, nil
}

// downloadReleaseAsset downloads the specified release asset to the provided path.
func (h *SourceHelper) downloadReleaseAsset(releaseAsset *github.ReleaseAsset, downloadPath string) error {
	resp, err := http.Get(*releaseAsset.BrowserDownloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(downloadPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// getReleaseByVersion returns the `github.RepositoryRelease` corresponding to the specified version.
func (h *SourceHelper) getReleaseByVersion(version string) (*github.RepositoryRelease, error) {
	client := github.NewClient(nil)

	// the repo uses a "v" prefix for release tags
	tagName := "v" + version

	release, resp, err := client.Repositories.GetReleaseByTag(context.Background(), h.GithubOwner, h.GithubRepo, tagName)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("no github release found with tag '%s'", tagName)
		}
		return nil, err
	}

	return release, nil
}
