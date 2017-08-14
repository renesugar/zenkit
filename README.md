zenkit
======

_zenkit_ is a Go microservice toolkit. Its purpose is to pull together
best-of-breed technologies in a known good configuration. With its companion
repo, [zenkit-template](https://github.com/zenoss/zenkit-template), you can be
up and running with an immediately deployable microservice in under a minute.

## Quick Reference
* _`make`_ to regenerate code after modifying `design/*.go`
* _`make test`_ to run tests
* _`make api-test`_ to run api tests
* _`make run`_ to run the service using docker-compose
* _`make build`_ to build the image
* _`make vendor`_ to update dependencies based on `glide.yaml`
* _`make local-dev`_ to install dev tools locally

## Prerequisites
To develop and run a zenkit microservice, you will need:
* make
* docker-ce >= 17.05 (Official installation instructions for
  [Ubuntu](https://docs.docker.com/engine/installation/linux/ubuntu/)
  | [CentOS](https://docs.docker.com/engine/installation/linux/centos/)
  | [macOS](https://docs.docker.com/docker-for-mac/install/))
* docker-compose (Install from here via `sudo make docker-compose`, or see the
  [official instructions](https://docs.docker.com/compose/install/))
* A Go environment. [gvm](https://github.com/moovweb/gvm) is an easy way to get
  one. Use the most recent release.

Additional helpful utilities include:
* [cobra](https://github.com/spf13/cobra) (`go get github.com/spf13/cobra/cobra`)
* [httpie](https://httpie.org/) (`apt install httpie` on Ubuntu)
* [jq](https://stedolan.github.io/jq/) (`apt install jq` on Ubuntu)
* Also, run `make local-dev` to install build and test tools in your local
  environment. Everything will still work if you don't do this, using a Docker
  container, but you may find it convenient.

## Quickstart
Just run this to create a microservice named `examplesvc`:

    bash <(curl -sSL https://git.io/vQB98) examplesvc

This will ask you a series of questions. You can always change the answers
later, except the first one, which is prefilled for you. This quickstart
assumes you chose the default port.

Once it's generated, go into your new directory and run `make` to pull in
dependencies and get everything set up:

    cd examplesvc
    make

Now you can run tests, if you want:

    make test

And you can start the thing, too. It doesn't do much, but it will run:

    make run

You can make requests against the included example endpoint:

    http :8080/hello/dolly
    http :8080/hello/newman

Then you can check metrics:

    http :8080/_admin/metrics

And browse the currently-trivial Swagger spec:

    http :8080/_admin/swagger

You'll need to [add your source](#create-github-repo) to Github and 
[create Jenkins jobs](create-jenkins-jobs) for continuous integration.

When you're ready to add business logic to your new service,
[proceed](#microservice-development).

## Zenkit Components
### The zenkit library
The _zenkit_ library provides:
* A standard base microservice with middleware preconfigured to support
  authentication, logging, metrics, and request tracing
* Generation of controller scaffolding, test code and
  [Swagger](http://swagger.io/) from common definitions
* Application instrumentation helpers
* Registration of common configuration arguments
* External service client creation
* ...and much, much more!

The zenkit library is the proper place for functionality that, when changed,
may affect all microservices.

### zenkit-template
[zenkit-template](https://github.com/zenoss/zenkit-template) is
a [boilr](https://github.com/tmrts/boilr) template that generates a fully
functional microservice scaffold (with a dummy endpoint) ready for you to fill
with business logic. Once generated, the original repository is no longer
referenced, and the microservice may be customized as you please.

Besides generating boilerplate, the zenkit template provides a Makefile that's
very convenient for development. Wrapping [Docker](https://docker.com) and
[docker-compose](https://docs.docker.com/compose/), it allows you to build,
test, manage vendored dependencies, and run in a local environment without
cumbersome setup.

Changes to the zenkit template will, for obvious reasons, only affect new
microservices created using zenkit.

### zenoss/zenkit-build Image
The [zenoss/zenkit-build](https://hub.docker.com/r/zenoss/zenkit-build/)
([GitHub](https://github.com/zenoss/zenkit-build)) Docker image is used to run
tests, build the service binaries, generate coverage reports, etc. Its purpose
is to remove steps a developer must perform to get started. The Makefile in
zenkit-template uses this image automatically.  _The image version of 
zenkit-build is specified in `.env` and may be updated there when new versions 
of zenkit-build are released._

You may find it convenient, however, to install the testing tools locally
rather than running them in a container. You may do this by running `make
local-dev`. This will install the Go testing and coverage tools into your
existing Go environment. The Makefile is smart enough to use local tools if you
installed them, so you can keep running `make test`.

### Technologies Used
* [Cobra](https://github.com/spf13/cobra) for CLI. Cobra files live under the
  `cmd` directory and are created using the `cobra` command line application.
* [Viper](https://github.com/spf13/viper) for configuration. All configuration
  is able to be specified via environment variables and config file, and live
  reloading of configuration is supported.
* [Goa](https://goa.design/) for APIs, service boilerplate, security, and
  [Swagger](http://swagger.io/) generation. Much of the development process
  involves modifying the resources defined in `design/resources.go` and using
  the `goagen` tool (encapsulated fully by the Makefile) to regenerate
  scaffolding code and boilerplate, then adding business logic.
* [go-metrics](https://github.com/rcrowley/go-metrics) for metrics.
* [Logrus](https://github.com/sirupsen/logrus) for structured logging.
* [Ginkgo](https://onsi.github.io/ginkgo/) and
  [Gomega](https://onsi.github.io/gomega) for testing.
* [Dredd](http://dredd.readthedocs.io/en/latest/) for api tests.

## Microservice Development
1. Add or modify resources and actions in `design/resources.go`, using [Goa's
   DSL](https://goa.design/reference/goa/design/apidsl/). The
   [goa-cellar](https://github.com/goadesign/goa-cellar) example implementation
   may also be a useful reference.

2. Define examples for your request and response objects in the design
   specification.  This provides richer swagger documentation and allows dredd
   tests to work automatically.  Documentation on how to implement examples is 
   found in the Goa documentation for 
   [Example](https://goa.design/reference/goa/design/apidsl/#func-example-a-name-apidsl-example-a) 
   and 
   [Metadata](https://goa.design/reference/goa/design/apidsl/#func-metadata-a-name-apidsl-metadata-a)
   functions.  See [the provided example
   resource](https://github.com/zenoss/zenkit-template/blob/master/template/design/resources.go)
   for a functional example.

3. `make`. This will generate scaffolding code in the `resources`
   directory, or modify existing scaffolding.

   Note: Goa generates _all_ the generated code under `resources/app`. Don't
   bother modifying it if you wish to avoid needless frustration.

4. Implement the resource action you've just defined. You'll find commented
   body in the boilerplate methods:

        // ControllerName_Action: start_implement

        // Put your logic here

        // ControllerName_Action: end_implement

   Like it says, put your logic in between the two outer comments. Leave the
   stuff on the outside of those comments alone. This allows `goagen` to
   regenerate the scaffolding around your logic as needed.

5. Add tests for your new code. There may already be a `CONTROLLER_test.go`
   defined. If not, run `ginkgo generate CONTROLLER`, where `CONTROLLER` is, of
   course, the name of the Go file containing your controller implementations.
   Goa generates test helpers for all resources to validate the contract, so
   that the DSL matches the implementation matches the Swagger output. You can
   lean on these in your tests to write them much faster, simply passing in the
   arguments that you expect to trigger each response. See the [tests for the
   example resource](https://github.com/zenoss/zenkit-template/blob/master/template/resources/example_test.go) for a functional example.

7. `make test`.  You may also run tests automatically on save by running
   `ginkgo watch resources` or `ginkgo watch -r`.

8. Add hooks for api tests that will handle environment setup and teardown for
   each api test. (TODO: need an example service that demonstrates this
   implementation).

9. `make api-test`. Starts service and dependencies (as defined in
   docker-compose.yml) and runs dredd tests in a container within a private
   network.

10. `make run` to rebuild the image and redeploy the service locally. This will
   bring it up on port {{Port}}, allowing you to use `curl` or `httpie`.  You
   may also simply use `go build {{Name}}`, then run the resulting binary
   manually, although if supporting services are required, the `docker-compose`
   functionality the Makefile implements is very convenient.

## Build and Repo Setup

### Create GitHub Repo
The zenkit creation script initialized your microservice as a git repo and
committed the initial code already. Now you need to add it to GitHub.

To create a new repo on Github, navigate to the [new repository
page](https://github.com/organizations/zenoss/repositories/new).  Enter the
same name you used to create the repository, since it's a Go package.  Also
a helpful description would be good. Make sure the repository will be private,
disable options to automatically add files to the repo, and click the *Create
repository* button.  Follow the instructions to _Push an existing repository
from the command line_.

There are several GitHub configuration items that are required for the new repo.
Browse to the _Settings_ tab in the new repo.  In the _Collaborations & teams_ 
settings, add the following teams and permissions:

| Team | Permission |
| ---- | ---------- |
| Administrators | Admin |
| Developers | Write |
| Automation | Write |
| Employees | Read |

In the _Branches_ settings, choose the _master_ branch for protection and
ensure that the options for "Require pull request reviews before merging"
and "Require status checks to pass before merging" are enabled.

### Create Jenkins Jobs
To create the suite of Jenkins Jobs, browse to the
[Microservice Job Builder](http://jenkins.zing.zenoss.eng/job/job_create/build?delay=0sec),
enter the name of the service repo, and click "Build".  This will create a
new suite of jobs for your service in the
[Microservices folder](http://jenkins.zing.zenoss.eng/job/micro-services/).

## Microservice Release
Once your microservice is ready for release, you will need to ensure that it is 
configured to deploy to the staging environment.  This is managed in
[zing-deploy](https://github.com/zenoss/zing-deploy).

1. Create a branch off of master of
   [zing-deploy](https://github.com/zenoss/zing-deploy).

2. Update the docker-compose files with your service's configuration.

3. Add your service to ci/release.groovy. In that file, you will see a list of
   objects, each with the same two fields, `path` and `services`.  Path is the
   path to the service's folder in jenkins, usually
   "micro-services/{{SERVICE}}".  Services is the list of service names that
   will use the image.  Note that ci/release.groovy only pertains to services
   that are released to staging/production.

4. Open a pull request against the master branch.  Jenkins will automatically
   deploy the last successful version of your service image into staging and
   run acceptance tests.

5. Once the pull request has been merged into master, your service will
   automatically get added to the stack.

6. Finally, you will need to enable production releases for the service so that
   pull requests service will automatically run tests in staging and merges
   into master will automatically update the stack in staging.  To do this,
   delete the disable-release config file from service's folder in Jenkins.
   This is typically located at
   `http://jenkins.zing.zenoss.eng/job/micro-services/job/{{SERVICE}}/configfiles`

## Issues?
[Zenoss Jira](https://jira.zenoss.com). Open an issue, ZING project, Zenkit
component.

## Troubleshooting

* When I run `make`, I see a message similar to
```
    [ERROR]    Error scanning github.com/dgrijalva/jwt-go: cannot find package "."
        in:     /home/user/.glide/cache/src/https-github.com-dgrijalva-jwt-go
```
Your package is missing from the glide.yaml.  Add the package and rerun `make`.
