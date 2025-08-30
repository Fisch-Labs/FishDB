module github.com/Fisch-Labs/FishDB

go 1.25.0

require (
	github.com/Fisch-Labs/common v1.0.0
	github.com/Fisch-Labs/ecal v1.0.0
	github.com/gorilla/websocket v1.4.1
)

replace (
	github.com/Fisch-Labs/common v1.0.0 => ./external_deps/common
	github.com/Fisch-Labs/ecal v1.0.0 => ./external_deps/ecal
)
