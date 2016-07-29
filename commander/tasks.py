from invoke import task
from docker import Client
from idflow import Utils, Docker, Flow
import os

Utils.print_system_info()

cli = Client(base_url='unix://var/run/docker.sock', timeout=600)

flow = Flow(
    repository="vjftw/orchestrate",
    prefix="commander"
)

@task
def test(ctx):
    Docker.build(cli,
        dockerfile='Dockerfile.dev',
        tag="{0}-dev".format(flow.get_branch_container_name())
    )

    Docker.run(cli,
        tag="{0}-dev".format(flow.get_branch_container_name()),
        command='glide install',
        volumes=[
            "{0}:/go/src/github.com/vjftw/orchestrate/commander".format(os.getcwd())
        ],
        working_dir="/go/src/github.com/vjftw/orchestrate/commander",
        environment={}
    )
    Docker.run(cli,
        tag="{0}-dev".format(flow.get_branch_container_name()),
        command='/bin/sh -c "go test -v -cover $(glide novendor)"',
        volumes=[
            "{0}:/go/src/github.com/vjftw/orchestrate/commander".format(os.getcwd())
        ],
        working_dir="/go/src/github.com/vjftw/orchestrate/commander",
        environment={
            "TERM": "xterm-color"
        }
    )

@task
def build(ctx):
    Docker.build(cli,
        dockerfile='Dockerfile.dev',
        tag="{0}-dev".format(flow.get_branch_container_name())
    )

    Docker.run(cli,
        tag="{0}-dev".format(flow.get_branch_container_name()),
        command='glide install',
        volumes=[
            "{0}:/go/src/github.com/vjftw/orchestrate/commander".format(os.getcwd())
        ],
        working_dir="/go/src/github.com/vjftw/orchestrate/commander",
        environment={}
    )
    Docker.run(cli,
        tag="{0}-dev".format(flow.get_branch_container_name()),
        command='/bin/sh -c "go build -a -installsuffix cgo -o dist/commander"',
        volumes=[
            "{0}:/go/src/github.com/vjftw/orchestrate/commander".format(os.getcwd())
        ],
        working_dir="/go/src/github.com/vjftw/orchestrate/commander",
        environment={
            "CGO_ENABLED": 0,
            "GOOS": "linux"
        }
    )

    Docker.build(cli,
        dockerfile='Dockerfile.app',
        tag=flow.get_build_container_name()
    )
