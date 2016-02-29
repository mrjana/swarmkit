//go:generate protoc -I.:../../../../:../vendor --gogoswarm_out=plugins=grpc+dataaccess,import_path=github.com/docker/swarm-v2/api:. types.proto cluster.proto master.proto agent.proto

package api
