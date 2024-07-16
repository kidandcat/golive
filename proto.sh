#!/bin/bash
curl --output profile.proto "https://raw.githubusercontent.com/google/pprof/main/proto/profile.proto"
protoc --go_opt=Mprofile.proto=github.com/kidandcat/golive/frontend --go_out=frontend --go_opt=paths=source_relative profile.proto