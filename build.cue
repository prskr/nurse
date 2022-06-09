package nurse

import (
	"dagger.io/dagger"
	"universe.dagger.io/go"
)

dagger.#Plan & {
	client: {
		filesystem: {
			"./": read: {
				contents: dagger.#FS
				exclude: [
					"README.md",
					"build.cue",
				]
			}
		}
		network: "unix:///var/run/docker.sock": connect: dagger.#Socket // Docker daemon socket
	}
	actions: {
		test: go.#Test & {
			package: "./..."
			source:  client.filesystem."./".read.contents
			mounts: docker: {
				dest:     "/var/run/docker.sock"
				contents: client.network."unix:///var/run/docker.sock".connect
			}
		}
	}
}
