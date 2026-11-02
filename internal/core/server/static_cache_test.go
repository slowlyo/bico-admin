package server

import "testing"

func TestIsHashedStaticAsset(t *testing.T) {
	t.Parallel()

	immutable := []string{
		"basic-D2TFQT31.js",
		"assets/js/basic-D2TFQT31.js",
		"vxe-table-CrpFQn-4.js",
		"umi.8f3a1c2d.js",
		"css/umi.abcdef0123456789.css",
		"style-a1b2c3d4.css",
	}
	for _, path := range immutable {
		path := path
		t.Run("immutable/"+path, func(t *testing.T) {
			t.Parallel()
			if !isHashedStaticAsset(path) {
				t.Fatalf("期望识别为带哈希构建产物: %q", path)
			}
		})
	}

	noCache := []string{
		"favicon.ico",
		"logo.png",
		"index.html",
		"assets/index.html",
		"vendor.js",
		"chunk-abc.js",
		"readme.txt",
	}
	for _, path := range noCache {
		path := path
		t.Run("no-cache/"+path, func(t *testing.T) {
			t.Parallel()
			if isHashedStaticAsset(path) {
				t.Fatalf("不应识别为带哈希构建产物: %q", path)
			}
		})
	}
}
