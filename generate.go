// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

//go:generate mkdir -p ./avro
//go:generate $GOPATH/bin/gogen-avro ./avro application_version_installed.avsc
