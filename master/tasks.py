from invoke import task
from docker import Client
from idflow import Utils, Docker, Flow
import os

Utils.print_system_info()

cli = Client(base_url='unix://var/run/docker.sock', timeout=600)

flow = Flow(
    repository="vjftw/orchestrate",
    prefix="master"
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
            "{0}:/go/src/github.com/vjftw/orchestrate/master".format(os.getcwd())
        ],
        working_dir="/go/src/github.com/vjftw/orchestrate/master",
        environment={}
    )
    Docker.run(cli,
        tag="{0}-dev".format(flow.get_branch_container_name()),
        command='/bin/sh -c "go test -v -cover $(glide novendor)"',
        volumes=[
            "{0}:/go/src/github.com/vjftw/orchestrate/master".format(os.getcwd())
        ],
        working_dir="/go/src/github.com/vjftw/orchestrate/master",
        environment={}
    )
