package step

import (
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/v2/env"
)

type fakeEnvRepo struct {
	env.Repository
	envs map[string]string
}

func (r fakeEnvRepo) Get(key string) string { return r.envs[key] }

func TestCacheKeys(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     []string
	}{
		{name: "default layout", envValue: "", want: keys},
		{name: "relocated by Build Cache for Xcode", envValue: "/Users/vagrant/.bitrise/cache/xcode-spm", want: xcelerateKeys},
		{name: "blank env var falls back to the default layout", envValue: "   ", want: keys},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cacheKeys(fakeEnvRepo{envs: map[string]string{EnvSwiftPackagesPath: tt.envValue}})

			if strings.Join(got, ",") != strings.Join(tt.want, ",") {
				t.Errorf("cacheKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Neither namespace may sit under the other's prefix fallback, or one layout restores over the other.
func TestNamespacesAreDisjoint(t *testing.T) {
	defaultPrefix, xceleratePrefix := keys[1], xcelerateKeys[1]

	if strings.HasPrefix(xcelerateKeys[0], defaultPrefix) {
		t.Errorf("xcelerate key %q is matched by the default prefix fallback %q", xcelerateKeys[0], defaultPrefix)
	}
	if strings.HasPrefix(keys[0], xceleratePrefix) {
		t.Errorf("default key %q is matched by the xcelerate prefix fallback %q", keys[0], xceleratePrefix)
	}
}
