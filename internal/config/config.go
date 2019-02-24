// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package config

import (
	"context"
	"strings"
	"time"
)

// GetContext gets the context and cancellation for a gRPC call.
func GetContext(timeout int) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		return context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	} else {
		return context.WithCancel(context.Background())
	}
}

// ExtractKVs extracts a series of semicolon-separated key:value pairs
// into a string:string key-value mapping.
// If the same key appears more than once, it will be overwritten in
// the map by the latter pair.
func ExtractKVs(cfgValue string) map[string]string {
	cfgs := map[string]string{}

	// first, split by semicolon
	splitPairs := strings.Split(cfgValue, ";")
	for _, p := range splitPairs {
		// now, split pair by FIRST colon
		kv := strings.SplitN(p, ":", 2)
		if len(kv) == 2 {
			cfgs[kv[0]] = kv[1]
		}
	}

	return cfgs
}
