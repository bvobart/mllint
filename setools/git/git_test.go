package git_test

import (
	"io/ioutil"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/setools/git"
	"github.com/bvobart/mllint/utils/exec"
)

func TestDetect(t *testing.T) {
	dir := "."
	require.True(t, git.Detect(dir))

	dir = ".."
	require.True(t, git.Detect(dir))

	dir = os.TempDir()
	require.False(t, git.Detect(dir))
}

func TestIsTracking(t *testing.T) {
	dir := "."
	require.True(t, git.IsTracking(dir, "git_test.go"))
	require.True(t, git.IsTracking(dir, "git*.go"))
	require.False(t, git.IsTracking(dir, "non-existant-file"))

	file, err := ioutil.TempFile(dir, "git.is-tracking.test-resource.*.txt")
	require.NoError(t, err)
	require.False(t, git.IsTracking(dir, file.Name()))

	require.NoError(t, os.Remove(file.Name())) // cleanup
}

func TestFindLargeFiles(t *testing.T) {
	dir := "."

	threshold := uint64(1)
	largeFiles, err := git.FindLargeFiles(dir, threshold)
	require.NoError(t, err)
	require.Len(t, largeFiles, 2)

	// test that largeFiles is sorted by filesize in descending order (i.e. largest files first)
	prevSize := uint64(math.MaxUint64)
	for _, file := range largeFiles {
		require.Truef(t, file.Size < prevSize, "Should be sorted by filesize in descending order: %+v", largeFiles)
		prevSize = file.Size
	}

	threshold = uint64(1000000000)
	largeFiles, err = git.FindLargeFiles(dir, threshold)
	require.NoError(t, err)
	require.Len(t, largeFiles, 0)
}

func TestFindLargeFilesInHistory(t *testing.T) {
	dir := "."
	threshold := uint64(4000)
	exec.PipelineOutput = func(execdir string, commands ...[]string) ([]byte, error) {
		require.Equal(t, dir, execdir)
		require.Equal(t, []string{"git", "rev-list", "--objects", "--all"}, commands[0])
		require.Equal(t, []string{"git", "cat-file", "--batch-check=%(objecttype) %(objectname) %(objectsize) %(rest)"}, commands[1])
		return []byte(mockPipelineOutput), nil
	}

	expected := []git.FileSize{
		{Path: "go.sum", CommitHash: "c57ea8a8f6301694a7c5e066ef47622c19eec971", Size: 29524},
		{Path: "go.sum", CommitHash: "a61a4e7af7463f270c4222c964865af5c65cfd8c", Size: 29274},
		{Path: "projectlinters/dependencies.go", CommitHash: "070e0b0dc26f83093da32d58c77f49babee5fd39", Size: 8570},
		{Path: "build/setup.py", CommitHash: "d7c8acfa5c917844da0d99eceb596c2c3cd9a94f", Size: 4454},
	}
	largeFiles, err := git.FindLargeFilesInHistory(dir, threshold)
	require.NoError(t, err)
	require.Equal(t, expected, largeFiles)
	// test that largeFiles is sorted by filesize in descending order (i.e. largest files first)
	prevSize := uint64(math.MaxUint64)
	for _, file := range largeFiles {
		require.Truef(t, file.Size < prevSize, "Should be sorted by filesize in descending order: %+v", largeFiles)
		prevSize = file.Size
	}
}

const mockPipelineOutput = `blob 7aba7e93fc6f8da4e63a683f2fda4a657d7835b2 3484 .github/workflows/build-publish.yml
tree 89417d9bb896948c6e612a9acd99f2f1161a8df0 460 
tree 44401018fc8e5d6848e11002c2d58069f5f86d7d 148 build
blob d7c8acfa5c917844da0d99eceb596c2c3cd9a94f 4454 build/setup.py
tree 97b671f32309f50286c511c54854c94b5481332a 460 
tree 28f2cf7ebadc50bfb8e2b36fe39ee4c22014a69d 148 build
blob f87a14883e720b2af12494f4267609a4d149c273 3782 build/setup.py
tree bdb4c319ba66c8d7006a218842b9b26a5fb497b9 460 
tree 89541d8890e272594697dff05e3f3bdae45caee3 148 build
blob 609bfb46865f88448215798f5be1c62fcadf576b 3736 build/setup.py
tree cb69cfc5fd8a9a35c5112db2cc05a0b42f924ffd 460 
tree 320fae3ca2435842b140017bac23cbab4e62b46d 36 .github
tree 0530886ac4d375cbb3c7eeac65818e837ef65257 45 .github/workflows
blob 8be29be1847fd360d5bf20763e078a4cb333a032 3566 .github/workflows/build-publish.yml
tree 7ab6ad8e91d403ccbcc49bae1331192a2138f001 148 build
blob ac8ca53c87a4226021f251c03f554b0d079a6176 3737 build/setup.py
tree a21b3b22066823c261f9deb25b58eb3b06242965 460 
tree 6872cc56a2fbe07a7abd2b8c28c41952f6cf0f27 111 build
tree 9f163712801b314d3c469a47696bf5cfadca8597 460 
tree 20d7af8672b0034824ec04865847713abbe032b9 36 .github
tree a4728930d086adc09fc8a7555f7c7786fb84d23c 45 .github/workflows
blob f1998f6f05174d35c5b89670fb8f73e8c20429b3 3648 .github/workflows/build-publish.yml
tree 0819f626b9574465d743ef53e8982c7935e5f905 460 
tree 978cd67dd2a05733b3480989fab49ac0d583339a 36 .github
tree ac95dc13947fbe33bcdfca43d5552c37e64bc0cf 45 .github/workflows
blob 131373c99151e34468712b003c8f4ddcdfec4f5d 3510 .github/workflows/build-publish.yml
tree 3d9b19465d876d89f9b0a28f6cd474b472526462 460 
tree 1ba3b649f01e0643307c87800a159d3e57f0b14e 36 .github
tree 85e1af4bc06d989a94d389522e4bc1f494ffcdad 45 .github/workflows
blob 76faca8829411caa320e76f0dbc114f6dc9ea650 3083 .github/workflows/build-publish.yml
tree 71c0743c0de6dbf6fb95e46e2e0a38d7810432d2 460 
tree 42d714838bec0fd05a8bebd3a1e65734edf3675a 36 .github
tree 9f9295bed487adc05aad278a73930b49d4b189c1 45 .github/workflows
blob 1cafb7253d00b74414b11220c1c956cf04320189 3018 .github/workflows/build-publish.yml
blob 56e93e17e1679b8b7f7b793d9e2b2f926bc68201 155 .gitignore
blob 032bdd1f31db2f87d3b72b1e7cc607e8219f4817 1478 Makefile
tree b1947f2ce4e26b6a5664563aaa59035fa348f393 111 build
blob f4dc3eb65e4c34f57a4e801cae8cceb2fd1eb027 3740 build/setup.py
tree 9ea9556ea087686cfdfc24b3c732815c6af89caf 426 
blob f49c7e1a4ff7420e709a9f4feea890882282f419 605 .gitlab-ci.yml
tree d83c8a0cbdaaefcdd4fec725994efbe309d4db79 390 
blob 74314ca9e58759eade3cd29cc4f0c0f904065152 116 .gitignore
tree 3cb716c195b0543ee3471981b53584d90b70d718 111 build
blob 32682c6106ade6fd24c2c7e7d40f28df9db33b75 1730 build/setup.py
tree f23aaf89218f5010ae81c26e38e1766e4bcb2bad 320 
tree bb56cd7a7d2b22c66c4d4ea25a5a9e27415ab6cc 320 
tree 6c7356f0e7fb11b8d33230401990e6404ce1f263 35 commands
blob 0826d5a76ef14c7584ac717e5dee1241e84a2f7a 2274 commands/root.go
tree 5a9620bd5f75989ee1aeebcd88367ad9b6d840fc 239 projectlinters
blob c4b7c1d307a6b349af7aa4fbf7df05b89dd72986 1906 projectlinters/git.go
tree 972129e7d8a0fce726bc8c5a474347239099b696 320 
tree fca12202b7a43f55bf497752b88325e274eb5dc8 75 api
blob 8a522c1ad3bfbd46d9d9300d2d30c9713c2618c8 719 api/results.go
tree a476f8fc0c6ef6291ddda2461e173a0b6a4b3909 35 commands
blob 3971eb2fda73494aaf18d7ad8fcf8af414f10f48 2078 commands/root.go
tree 29759bc50ceec801cb6bc1e55a83f55dd31fb0dd 239 projectlinters
blob 070e0b0dc26f83093da32d58c77f49babee5fd39 8570 projectlinters/dependencies.go
tree 90b601cc8938505e84a15c28753f1841711ad93f 320 
tree 4f8e1b1b4c21ffde4a7543bffd984bb946c6bea4 239 projectlinters
blob 856985bcb919a45bf8dfa694945fa0a155889641 2726 projectlinters/dependencies.go
blob ba37e0192d3e1a24b8f0e13c6fa250aca0a3d42d 2017 projectlinters/dependencies_test.go
tree b8426a008ae6aa3e18699137eb813b7916de7f12 39 projectlinters/test-resources
tree 13bbdb8d27f7519750728cca99ab07ee44500bec 229 projectlinters/test-resources/dependencies
tree 061d9999e3071a6dd50928ae06f9eaab7c1a2aca 138 utils
tree b96027d694822c13029940eb0cd3ebd0c271aa78 147 utils/depsmgmt
blob 41511fb3d7fed539337e8e7e0f243a645b5dbbc3 865 utils/depsmgmt/managers.go
blob 213baff231a15ede8653e5aa080409d13e4efa19 272 utils/depsmgmt/pip.go
tree 0d1e4d3af108c908eaf4282cac942b23941124f5 320 
tree f4e98e4f32404a9925c3892af73e6fe2234e4072 111 projectlinters
blob 6d40a8550be778fb21a52f3bdc9ca223d493ea0a 2716 projectlinters/dependencies.go
tree a85588db60d735a5205a4eedde4fc3acdf864c67 320 
tree cafb5c8a5c10b970cf96c6195785e16cfd76ea5d 68 projectlinters
blob aec876cfc270905b1c0ced9cbe1daad4705e9122 273 projectlinters/all.go
tree 9f47190ce941502d5c0acea6cdfb259bfe2e2dd5 103 utils
tree 67802af8871a179b0bb78fa93a68bca19f978252 320 
blob e7cbf45c93702c0ea505e56b769e214f5d99058b 614 .gitlab-ci.yml
tree 886fe75abc1231a5713b4dfd6997ea2c2838e69c 320 
tree 01c638cc528d39fe9b9e1176369198935520943b 75 api
blob 45ae5a9693b52a4d361f42f43c55e6ad7b10acb1 648 api/results.go
tree e5469e6b2dc62418a0d933188fc615573239e512 35 commands
blob 4ab0aef6dbbf6c73de13d7f9d38acb98106907da 2056 commands/root.go
blob f4672304e7046e01076e2c0e214958c82917f4bf 154 go.mod
blob c57ea8a8f6301694a7c5e066ef47622c19eec971 29524 go.sum
tree 914dc327d7cc1545595584543599e1a52ffdd8a4 68 projectlinters
blob 0e7e171d6bd738a8a21c685477ad7e9845a558ad 142 projectlinters/all.go
blob bd844ce823d3ff64e0891c300b90f71b7b63d570 487 projectlinters/git.go
tree 8ac5628092973b86ac92ab44f37e6df2ffb14494 66 utils
tree 78f78409aa065041c9815c37fc9817cf7ab94032 73 utils/git
blob 1f335a5070b045f86ff8e2e7de6b493cbf4a92ea 244 utils/git/git.go
blob d458f1ca1828d9e5e8638d9008071d8b1583cebe 309 utils/git/git_test.go
tree 05e8f325f4aafb4b7629d9ee6ec139eb4698b29f 278 
blob 6458083944d9b663f0382c856c6690f8363b42bb 118 go.mod
blob a61a4e7af7463f270c4222c964865af5c65cfd8c 29274 go.sum
tree a7a43724c0e9871f7168d22c37fc7720df16750a 66 utils
tree 540d7e93ce31afc7bb7c044a3d29bdb5d74512d9 34 utils/git
tree 211bbf25aa576baa5e7a95c2fc20cee7099d6942 241 
tree 2db84144581daae349c85a4f33aebc9ebe3b904e 241 
tree c70fc5b3ee01299c745468d45d932c785864fb00 75 api
blob 5b5e1cc192919e08ba523ae01a52cda6d9a8995f 213 api/results.go
tree c3a132798c688377bae8928364453d546c8ce761 35 commands
blob 60624f0b34415d84962ace5e5891dd357d5b46cf 1786 commands/root.go
tree 9defc4250bc9f5ad53091e4c793aa0337b6d8e51 68 projectlinters
blob d852c9fdb99450c523717e39791598ff71ed1341 257 projectlinters/git.go
tree 5a439f85965f0ca9a0c320b0c0c1d14779b116e8 36 utils
`
