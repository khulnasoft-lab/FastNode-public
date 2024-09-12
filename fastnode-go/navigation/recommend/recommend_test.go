package recommend

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/navigation/git"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/navigation/ignore"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/navigation/localpath"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/navigation/metrics"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
	"github.com/stretchr/testify/require"
)

var (
	testDirString = filepath.Join(
		os.Getenv("GOPATH"), "src", "github.com", "khulnasoft-lab", "khulnasoft-lab",
		"fastnode-go", "navigation", "offline", "testdata",
	)
	testDir = localpath.Absolute(testDirString)
)

func BenchmarkKhulnasoft_labRepo(b *testing.B) {
	khulnasoft-lab := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "khulnasoft-lab", "khulnasoft-lab")
	var (
		ignoreOpts = ignore.Options{
			Root:            localpath.Absolute(khulnasoft-lab),
			IgnoreFilenames: []localpath.Relative{ignore.GitIgnoreFilename},
		}
		recOpts = Options{
			UseCommits:           true,
			ComputedCommitsLimit: git.DefaultComputedCommitsLimit,
			Root:                 localpath.Absolute(khulnasoft-lab),
			MaxFileSize:          1e6,
			MaxFiles:             1e5,
		}
		request = Request{
			MaxFileRecs:      5,
			MaxBlockRecs:     5,
			MaxFileKeywords:  -1,
			MaxBlockKeywords: 3,
			Location: Location{
				CurrentPath: filepath.Join(khulnasoft-lab, "fastnode-go", "lang", "python", "pythoncomplete", "api", "api.go"),
				CurrentLine: 50,
			},
		}
	)
	run(b, ignoreOpts, recOpts, request)
}

func BenchmarkAllValidationRepos(b *testing.B) {
	var (
		root       = filepath.Join(os.Getenv("HOME"), "nav-validation")
		ignoreOpts = ignore.Options{
			Root:           localpath.Absolute(root),
			IgnorePatterns: []string{".*", "*/*/open", "*/*/closed"},
		}
		recOpts = Options{
			UseCommits:  false,
			Root:        localpath.Absolute(root),
			MaxFileSize: 1e6,
			MaxFiles:    1e5,
		}
		request = Request{
			MaxFileRecs:      5,
			MaxBlockRecs:     5,
			MaxFileKeywords:  -1,
			MaxBlockKeywords: 3,
			Location: Location{
				CurrentPath: filepath.Join(root, "prestodb", "presto", "root", "presto-array", "src", "main", "java", "com", "facebook", "presto", "array", "Arrays.java"),
				CurrentLine: 30,
			},
		}
	)
	run(b, ignoreOpts, recOpts, request)
}

func run(b *testing.B, ignoreOpts ignore.Options, recOpts Options, request Request) {
	var initializeDuration, recommendDuration time.Duration
	for i := 0; i < b.N; i++ {
		startInitialize := time.Now()
		ignorer, err := ignore.New(ignoreOpts)
		if err != nil {
			log.Fatal(err)
		}
		s, err := git.NewStorage(git.StorageOptions{})
		if err != nil {
			log.Fatal(err)
		}
		recommender, err := NewRecommender(fastnodectx.Background(), recOpts, ignorer, s)
		if err != nil {
			log.Fatal(err)
		}
		initializeDuration += time.Since(startInitialize)

		startRecommend := time.Now()
		files, err := recommender.Recommend(fastnodectx.Background(), request)
		if err != nil {
			log.Fatal(err)
		}
		blockRequest := BlockRequest{
			Request:      request,
			InspectFiles: files,
		}
		_, err = recommender.RecommendBlocks(fastnodectx.Background(), blockRequest)
		if err != nil {
			log.Fatal(err)
		}
		recommendDuration += time.Since(startRecommend)
	}

	b.ReportMetric(float64(initializeDuration.Seconds())/float64(b.N), "s/init")
	b.ReportMetric(float64(recommendDuration.Milliseconds())/float64(b.N), "ms/rec")
	b.ReportMetric(0, "ns/op")
}

type validateRecommendTC struct {
	currentPath                  string
	inspectPath                  string
	expectedRecommendError       error
	expectedRecommendBlocksError error
}

func TestValidateRecommend(t *testing.T) {
	ignorer, err := ignore.New(ignore.Options{Root: testDir})
	require.NoError(t, err)

	recOpts := Options{
		Root:        testDir,
		MaxFileSize: 1e6,
		MaxFiles:    1e5,
	}
	s, err := git.NewStorage(git.StorageOptions{})
	require.NoError(t, err)
	r, err := NewRecommender(fastnodectx.Background(), recOpts, ignorer, s)
	require.NoError(t, err)

	absRec := filepath.Join(testDirString, "parsergo.py")
	relRec, err := filepath.Rel(testDirString, absRec)
	require.NoError(t, err)
	absVal := filepath.Join(testDirString, "astgo.py")
	relVal, err := filepath.Rel(testDirString, absVal)
	require.NoError(t, err)

	tcs := []validateRecommendTC{
		validateRecommendTC{
			currentPath:                  absRec,
			expectedRecommendBlocksError: errRelativeInspectPath,
		},
		validateRecommendTC{
			currentPath:                  relRec,
			expectedRecommendError:       errRelativeCurrentPath,
			expectedRecommendBlocksError: errRelativeInspectPath,
		},
		validateRecommendTC{
			currentPath: absRec,
			inspectPath: absVal,
		},
		validateRecommendTC{
			currentPath: absRec,
			inspectPath: absVal,
		},
		validateRecommendTC{
			currentPath:                  absRec,
			inspectPath:                  relVal,
			expectedRecommendBlocksError: errRelativeInspectPath,
		},
		validateRecommendTC{
			currentPath:                  relRec,
			inspectPath:                  absVal,
			expectedRecommendError:       errRelativeCurrentPath,
			expectedRecommendBlocksError: errRelativeCurrentPath,
		},
		validateRecommendTC{
			currentPath:                  relRec,
			inspectPath:                  relVal,
			expectedRecommendError:       errRelativeCurrentPath,
			expectedRecommendBlocksError: errRelativeInspectPath,
		},
	}

	var rankCount, batchCount int64
	for _, tc := range tcs {
		request := Request{
			MaxFileRecs:      5,
			MaxBlockRecs:     5,
			MaxFileKeywords:  -1,
			MaxBlockKeywords: 3,
			Location: Location{
				CurrentPath: tc.currentPath,
				CurrentLine: 50,
			},
		}
		blockRequest := BlockRequest{
			Request:      request,
			InspectFiles: []File{{Path: tc.inspectPath}},
		}
		_, recErr := r.Recommend(fastnodectx.Background(), request)
		require.Equal(t, tc.expectedRecommendError, recErr)
		if recErr == nil {
			rankCount++
		}
		_, recBlocksErr := r.RecommendBlocks(fastnodectx.Background(), blockRequest)
		require.Equal(t, tc.expectedRecommendBlocksError, recBlocksErr)
		if recBlocksErr == nil {
			batchCount++
		}
	}
	m := metrics.Read(true)
	require.Equal(t, int64(1), m["nav_index_count"])
	require.Equal(t, rankCount, m["nav_rank_count"])
	require.Equal(t, batchCount, m["nav_batch_count"])
}

func TestBadPatternNoError(t *testing.T) {
	ignoreOpts := ignore.Options{
		Root:           testDir,
		IgnorePatterns: []string{".*", "[a"},
	}
	ignorer, err := ignore.New(ignoreOpts)
	require.NoError(t, err)
	recOpts := Options{
		Root:        testDir,
		MaxFiles:    1e5,
		MaxFileSize: 1e6,
	}
	s, err := git.NewStorage(git.StorageOptions{})
	require.NoError(t, err)
	_, err = NewRecommender(fastnodectx.Background(), recOpts, ignorer, s)
	require.NoError(t, err)
}

type recommendTC struct {
	currentPath      string
	expectedPath     string
	expectedKeywords []string
	maxFileRecs      int
}

func TestRecommend(t *testing.T) {
	khulnasoft-lab := localpath.Absolute(os.Getenv("GOPATH")).Join("src", "github.com", "khulnasoft-lab", "khulnasoft-lab")
	ignoreOpts := ignore.Options{
		Root:           khulnasoft-lab,
		IgnorePatterns: []string{".*", "vendor", "bindata", "node_modules"},
	}
	ignorer, err := ignore.New(ignoreOpts)
	require.NoError(t, err)
	recOpts := Options{
		Root:                 khulnasoft-lab,
		MaxFileSize:          1e6,
		MaxFiles:             1e5,
		UseCommits:           true,
		ComputedCommitsLimit: git.DefaultComputedCommitsLimit,
	}
	s, err := git.NewStorage(git.StorageOptions{})
	require.NoError(t, err)
	r, err := NewRecommender(fastnodectx.Background(), recOpts, ignorer, s)
	require.NoError(t, err)

	tcs := []recommendTC{
		recommendTC{
			currentPath:      filepath.Join(testDirString, "parsergo.py"),
			expectedPath:     filepath.Join(testDirString, "astgo.py"),
			expectedKeywords: []string{"ListComprehensionExpr", "BaseComprehension"},
			maxFileRecs:      5,
		},
		recommendTC{
			currentPath:      filepath.Join(testDirString, "maingo.py"),
			expectedPath:     filepath.Join(testDirString, "datagensh.py"),
			expectedKeywords: []string{"stepsperfile", "CONTEXTSIZE"},
			maxFileRecs:      5,
		},
		recommendTC{
			currentPath:      filepath.Join(testDirString, "trainpy.py"),
			expectedPath:     filepath.Join(testDirString, "validatepy.py"),
			expectedKeywords: []string{"appraise", "mean_utility"},
			maxFileRecs:      5,
		},
		recommendTC{
			currentPath:      filepath.Join(testDirString, "Logsjs.py"),
			expectedPath:     filepath.Join(testDirString, "logscss.py"),
			expectedKeywords: []string{"logs__link", "logs__cta"},
			maxFileRecs:      5,
		},
		recommendTC{
			currentPath:      filepath.Join(testDirString, "modelpy.py"),
			expectedPath:     filepath.Join(testDirString, "modelgo.py"),
			expectedKeywords: []string{"FillUnknownIntelliJPaid", "IntelliJPaid"},
			maxFileRecs:      5,
		},

		// This test checks that all recommended files have blocks and keywords,
		// and all recommended blocks have keywords.
		recommendTC{
			currentPath:  filepath.Join(testDirString, "trainpy.py"),
			expectedPath: filepath.Join(testDirString, "validatepy.py"),
			maxFileRecs:  -1,
		},
	}
	for _, tc := range tcs {
		request := Request{
			Location: Location{
				CurrentPath: tc.currentPath,
			},
			MaxFileRecs:      tc.maxFileRecs,
			MaxBlockRecs:     5,
			MaxFileKeywords:  7,
			MaxBlockKeywords: 5,
		}
		noBlocks, err := r.Recommend(fastnodectx.Background(), request)
		require.NoError(t, err)
		blockRequest := BlockRequest{
			Request:      request,
			InspectFiles: noBlocks,
		}
		recs, err := r.RecommendBlocks(fastnodectx.Background(), blockRequest)
		if tc.maxFileRecs != -1 {
			require.Equal(t, tc.maxFileRecs, len(recs))
		}
		require.NoError(t, err)
		blocks := make(map[string]map[string]bool)
		files := make(map[string]map[string]bool)
		for _, rec := range recs {
			require.NotZero(t, len(rec.Blocks), "recommended files must have blocks")
			require.NotZero(t, len(rec.Keywords), "recommended files must have keywords")
			if request.MaxFileKeywords != -1 {
				require.True(t, len(rec.Keywords) <= request.MaxFileKeywords)
			}

			blocks[rec.Path] = make(map[string]bool)
			for _, block := range rec.Blocks {
				require.NotZero(t, len(block.Keywords), "recommended blocks must have keywords")
				for _, keyword := range block.Keywords {
					blocks[rec.Path][keyword.Word] = true
				}
			}
			files[rec.Path] = make(map[string]bool)
			for _, keyword := range rec.Keywords {
				files[rec.Path][keyword.Word] = true
			}

			require.Equal(t, len(rec.Keywords), len(files[rec.Path]))
			for i, keyword := range rec.Keywords {
				if i == 0 {
					continue
				}
				require.True(t, keyword.Score <= rec.Keywords[i-1].Score)
			}
		}

		require.Contains(t, blocks, tc.expectedPath)
		require.Contains(t, files, tc.expectedPath)
		for _, keyword := range tc.expectedKeywords {
			require.Contains(t, blocks[tc.expectedPath], keyword)
			require.Contains(t, files[tc.expectedPath], keyword)
		}
	}
}

func TestSkipRefresh(t *testing.T) {
	khulnasoft-lab := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "khulnasoft-lab", "khulnasoft-lab")
	ignoreOpts := ignore.Options{
		Root:           localpath.Absolute(khulnasoft-lab),
		IgnorePatterns: []string{".*", "vendor", "bindata", "node_modules"},
	}
	ignorer, err := ignore.New(ignoreOpts)
	require.NoError(t, err)
	recOpts := Options{
		Root:        localpath.Absolute(khulnasoft-lab),
		MaxFileSize: 1e6,
		MaxFiles:    1e5,
	}
	s, err := git.NewStorage(git.StorageOptions{})
	require.NoError(t, err)
	r, err := newRecommender(fastnodectx.Background(), recOpts, ignorer, s)
	require.NoError(t, err)

	for id, vec := range r.vectorizer.vectorSet.data {
		r.vectorizer.vectorSet.data[id] = shingleVector{
			coords:  vec.coords,
			norm:    vec.norm,
			modTime: vec.modTime.Add(-time.Second),
		}
	}

	// clear the metrics
	metrics.Read(true)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		start := time.Now()
		defer wg.Done()
		request := Request{
			Location: Location{
				CurrentPath: filepath.Join(testDirString, "parsergo.py"),
			},
			MaxFileRecs:      -1,
			MaxBlockRecs:     5,
			MaxFileKeywords:  7,
			MaxBlockKeywords: 5,
		}
		_, err := r.Recommend(fastnodectx.Background(), request)
		require.NoError(t, err)
		require.True(t, time.Since(start) > 2*time.Second)
		m := metrics.Read(true)
		require.NotZero(t, m["nav_rank_num_refreshed_files"])
	}()

	time.Sleep(time.Second)

	wg.Add(1)
	go func() {
		start := time.Now()
		defer wg.Done()
		request := Request{
			Location: Location{
				CurrentPath: filepath.Join(testDirString, "parsergo.py"),
			},
			MaxFileRecs:      -1,
			MaxBlockRecs:     5,
			MaxFileKeywords:  7,
			MaxBlockKeywords: 5,
			SkipRefresh:      true,
		}
		_, err := r.Recommend(fastnodectx.Background(), request)
		require.NoError(t, err)
		require.True(t, time.Since(start) < time.Second)
		m := metrics.Read(true)
		require.Zero(t, m["nav_rank_num_refreshed_files"])
	}()

	wg.Wait()
}
