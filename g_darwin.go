// Copyright 2016 aletheia7. All rights reserved. Use of this source code is
// governed by a BSD-2-Clause license that can be found in the LICENSE file.
package main

import (
	"github.com/aletheia7/ul"
	"log"
)

func setup_log() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(ul.New_object(`godoc-proxy`, ``))
}
