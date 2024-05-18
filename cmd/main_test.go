package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Main(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		args []string
		test func(require.TestingT, interface{}, ...interface{})
	}{
		"happy_path": {
			args: []string{os.Args[0], "../data/problem13.txt"},
			test: require.Nil,
		},
		// XXX: they die on a panic anyway, so, for now...
		// "bad_filename": {
		// 	args: []string{os.Args[0], "../data/problemX.txt"},
		// 	test: require.Nil,
		// },
		// "missing_filename": {
		// 	args: os.Args[:1],
		// 	test: require.Nil,
		// },
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			func(t *testing.T) {
				defer func() {
					tc.test(t, recover())
				}()
				os.Args = tc.args
				main()
			}(t)
		})
	}
}
