// Copyright (c) 2019 Jason T. Lenz.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package pgot_test

import (
        "git.lenzplace.org/lenzj/testcli"
        "testing"
)

func TestPgot(t *testing.T) {
        testcli.RunTests(t, "../pgot")
}
