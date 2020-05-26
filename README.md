# Global Transcoding Platform

-   [Global Transcoding Platform](#global-transcoding-platform)
    -   [Before You Begin](#before-you-begin)
    -   [Access Request](#access-request)         
        -   [Jira Access Request](#jira-access-request)
        -   [Confluence Access Request](#confluence-access-request)
        -   [AWS Access Request](#aws-access-request)
    -   [Development Environment Setup](#development-environment-setup)
        -   [AWS Credential Setup](#aws-credential-setup)
        -   [Local Development](#local-development)
            - [Go](#go)
            - [Docker & AWS](#docker-&-aws)
            - [Create Local PostgreSQL DB](#create-local-postgresql-db)
        -   [Build & Unit Test](#build-&-unit-test)
        -   [Pulumi setup](#pulumi-setup)
            -   [Pulumi configuration](#pulumi-configuration)
            -   [Personal Pulumi Stack Deployment](#personal-pulumi-stack-deployment)
    -   [Release Process](#release-process)
    -   [Resources](#resources)
        -   [Pulumi](#pulumi)

**Build**
[![CircleCI](https://circleci.com/gh/EurosportDigital/global-transcoding-platform.svg?style=svg&circle-token=4eb2594f171af2332fa25591eafa4344157feb75)](https://circleci.com/gh/EurosportDigital/global-transcoding-platform)

Currently, we use a monorepo to manage source code of all the micro services, source code of all the cloud infrastructure
provisioning logic, and source code of all the metrics and monitors provisioning logic.
The CI/CD pipeline will build only those projects that are changed.


## Before You Begin

It's important to note that any contribution to this repository MUST be associated with the corresponding JIRA item. Any
PR opened on behalf of the work involved in the JIRA item MUST be complete. Partial pull requests that do not cover all
cases listed in the JIRA item will not be accepted. It is ok (and preferable) to break big items into smaller tasks and
complete each one individually as long as each task brings additional value to the codebase and the project.

Additionally, every PR should be accompanied with proper unit and (where applicable) integration tests. PRs that do not
include tests will not be merged.

If you want to reference the entire GTP team, you can use @gtp-team for both Pull Requests reviews and basic communication.

## Access Request

### Jira Access Request

We use JIRA for our sprints. In order to get Jira access, create a new
[ServiceNow ticket](https://discoveryinc.service-now.com/askdiscovery?id=sc_cat_item&sys_id=e933f02d13920f84907976666144b0a0) with the following information:

Request Type: `New Access`

Business Purpose: `New member of VDP Bellevue team`

Additional Notes or Comments: `Access to eurosportdigital JIRA`

Security Roles to Add: `User`

If you don't have a Discovery email address, please contact someone from the core Bellevue team to submit a new access request on your behalf. Form to be used is available here: https://app.smartsheet.com/b/form/a55f98065ec946acad489a73ffaca933.

### Confluence Access Request

Please get create a Confluence access request to access the wiki. To make a request open a new
[ServiceNow ticket](https://discoveryinc.service-now.com/askdiscovery?id=sc_cat_item&sys_id=ffb0d9e913160f84907976666144b0f5)

Request Type: `New Access`

Business Purpose: `New member of VDP Bellevue team`

Additional Notes or Comments: `Access to scrippsnetworks confluence`

Security Roles to Add: `User`

### AWS Access Request

Before you can start using AWS you need to request access to our AWS accounts. To make a request open a new
[ServiceNow ticket](https://discoveryinc.service-now.com/askdiscovery?id=sc_cat_item&sys_id=e83233ecf47cdd009d9df1cd7143e9a1)
with the following information:

Request Type: `Grant`

Business Justification: `New member of VDP Bellevue team`

Group Membership: `CLOUD-AWS-Engineer-VDP`

Access Type: `Production Support`

Enter any additional comments if you have them.

You should receive an email when your requests are approved. This shouldn't take more than 2 business days. If it does,
ping @manca-disc.

## Development Environment Setup

### AWS Credential Setup

We use two AWS accounts for our infra deployments:

```bash
1. discovery-aws-vdp-vod-dev (580308463258)
2. discovery-aws-vdp-vod-prod (987149203146)
```

Both accounts use the role `vdp-engineer`, which you assume using your Discovery email/password.

To login to AWS Console [click here](https://discovery.okta.com/home/amazon_aws/0oa4id630pfZqa5p52p7/272).

To get AWS credentials to use AWS APIs, you need a wrapper over awscli for Okta authentication:

1. Install `gimme-aws-creds` from here: https://github.com/Nike-Inc/gimme-aws-creds

2. Run `gimme-aws-creds --config` and use the following configuration:

    - Okta Configuration Profile Name [DEFAULT]: <ENTER>

    - Okta URL for your organization: _https://discovery.okta.com_

    - URL for gimme-creds-server [appurl]: <ENTER>

    - Application url: _https://discovery.okta.com/home/amazon_aws/0oa4id630pfZqa5p52p7/272_

    - Write AWS Credentials: _y_

    - Resolve AWS alias: _y_

    - AWS Role ARN: _vdp-engineer_

    - Okta User Name: <YOUR_DISCOVERY_EMAIL>

    - AWS Default Session Duration [3600]: <ENTER>

    - Preferred MFA Device Type: _push_

    - Remember device: _y_

    - AWS Credential Profile [role]: <ENTER>

3. Run `gimme-aws-creds` and follow the prompts

4. When prompted to `Pick a role`, choose `discovery-aws-vdp-vod-dev (580308463258)`

5. `export AWS_PROFILE=vdp-engineer` (you can store this in your .bashrc/.bash_profile)

After the initial setup is done you should only need to run steps 3 and 4 from time to time.

To verify that everything works correctly, run:

`aws s3 ls`

You should see a list of current S3 buckets.

### Local Development

#### Go

We use Go as the language of choice for Global Transcoding Platform. The tooling that comes with it out of the box is utilized as much as possible.
For example, the package manager is Go Mod, unit testing framework is Go Test, build/install tools are Go Build, Go Install respectively, etc.

To get started with Go, follow the instructions here: https://golang.org/doc/install. Current version we are using is 1.14.

Also, please get familiar with the language by browsing exhaustive Golang documentation. Start here: https://golang.org/doc/code.html.

#### Docker & AWS

We use Docker to containerize our microservices, as well as AWS to deploy them. Although our deployment is cloud-agnostic we start with AWS.

To install Docker download the following package: https://download.docker.com/mac/stable/Docker.dmg

To install AWS CLI, run the following command:

```sh
brew install awscli
# or
pip3 install awscli
```

Verify AWS CLI and Docker are installed:

```sh
aws --version
docker --version
```
#### Create Local PostgreSQL DB

In order to work with a local instance of a PostgreSQL DB, you'll need to install PostgreSQL locally (in this case we're using a Windows environment)

Fist step is to download the windows installer, [here](https://www.enterprisedb.com/downloads/postgres-postgresql-downloads) you can download the installer corresponding to your OS.

![](imgs/installer-versions.jpg)

When you run the installer, you click next until the UI ask you for a user and password: ( you can modify the directory for the install, but we’ll use the default location c://Program Files/PostgreSQL ). Also, we installed the PgAdmin by default and used the default network access port. 

The username and password selected here, will be used by the default database called 'postgres' that will be created during install, so take note of them for future reference.

![Installer Passwords](imgs/installer-paswords.jpg)

After that continue, until you are requested the locale that you are going to use for the new DB cluster, here we’ll select English, United States.

![Installer Locale](imgs/installer-locale.jpg)

After finishing the install, the StackBuilder will open, this will allow you to add plugins drivers and extras to the install. You can install the psqlODBC in this step or other versions of the Database Server in case it’s needed, but you can always add them further down the line.

This will open you default browser and navigate to your server admin page, you’ll be prompted for your credentials ( the ones you set during the install) enter them and you’ll be connected to you server.
If you select servers, you’ll find the postgres default database.

![PG Admin](imgs/pg-admin-1.jpg)

You can access the postgres databa using powershell by navigation to the install folder ( C://Program Files/PostrgeSQL/12/bin ) and executing the following commands:

![commands 1](imgs/commands%201.jpg)

Complete with your password and press enter.

Next, you can create a Database with:

```sh
CREATE DATABASE dbname;
```
![CREATE DB](imgs/create%20db.jpg)

With the command `\l`  you can list the databases in your server:

![DB LIST](imgs/list-db.jpg)

You can connect to the database you created in order to work on it using \c dbname

![db connect](imgs/db-connect.jpg)

These are some useful commands to SET basic resources :
-	Create a schema in the DB :

```sh 
CREATE SCHEMA schema_name;
```
-	Create a user:

 ```sh
 CREATE USER username PASWWORD ‘password’;
 ```
-	Grant access permissions to schema:

```sh
GRANT ALL ON SCHEMA schema_name TO username;
```
-	Grant access permissions to tables: 

```sh
GRANT ALL ON ALL TABLES IN SCHEMA schema_name TO username;
```

![schema](imgs/schema%20db.jpg?raw=true)

To login to your new db with the new created user:

![login db](imgs/login%20db.jpg)

Now, let’s create a TABLE in the new schema, for that you have to declare the value name and type that will be stored on the table.

```sh
CREATE TABLE schema_name.table_name (col varchar(20));
```

![CREATE TABLE](imgs/create%20table.jpg)

Then, we can insert actual values into the database:

```sh
insert into schema_name.table_name (col) values (‘value’);
```

and you can select the values from the database to see what’s in it. 

```sh
SELECT * from schema_name.table_name;
```

![select table](imgs/select%20db.jpg)

To close the connection use `\q` and you can drop the table with 

```sh
DROP TABLE schema_name.table_name;
```

Now, we can see the changes we made in the PGadmin. 

![PG ADMIN END](imgs/pg%20admin%201.jpg)


### Build & Unit Test

Every file should have corresponding `_test.go` file which should include unit tests for that component.

Some of the handy commands to know about:

```sh
# update dependencies
go mod download
# check code style
go fmt ./...
# unit test
go test ./...
# build
go build ./...
# install -- this will put binaries in your GOPATH/bin
go install ./... 
```

### Pulumi setup

We use Pulumi to provision our Cloud infrastructure.

Install Pulumi:

```sh
curl -sSL https://get.pulumi.com | bash -s -- --version 1.14.1
echo "export PATH=${HOME}/.pulumi/bin:$PATH" >> ~/.bash_profile
```

For more information on Pulumi, check their official Getting Started Guide: https://www.pulumi.com/docs/.

### Personal Pulumi Stack Deployment

First, switch to `discovery-aws-vdp-vod-dev` AWS profile.

Our Pulumi state bucket is in `us-west-2`, for convenience you can set these environment variables to your shell
profile:

```sh
export AWS_PROFILE=vdp-engineer
export AWS_REGION=us-west-2
```

Each personal Pulumi stack should map to an already existing service stack. For example, `Sample Service`'s dev stack is located under `services/sample/`. If you want to experiment with this service using your own private stack you would create it under the same directory following the instructions below.

```sh
# login to the S3-based states store
pulumi login s3://dev-gtp-pulumi-states
# navigate to desired service which you want to test in a private stack
cd services/<SERVICE>
# create a personal stack that will contain all service's resources in it
pulumi stack init pvt-$USER-<SERVICE>
# make sure you are always working under your own stack by running
pulumi stack select pvt-$USER-<SERVICE>
```

Your personal Pulumi stack will have its environment defined in `Pulumi.pvt-$USER-<SERVICE>.yaml` file. Make sure to set initial config by running:

```sh
pulumi config set env dev
pulumi config set --secret datadog:apiKey <DATADOG_API_KEY> # found here https://app.datadoghq.com/account/settings#api
```

To preview and deploy your local Pulumi stack, follow the instructions below:

```sh
# preview your infrastructure change (use --skip-preview and --suppress-outputs flags if you dont want to see stack outout)
pulumi preview
# preview with detail:
pulumi preview --diff
# deploy!
pulumi up
# if you see errors or want more detail
pulumi up --debug --verbose 10
# clean up after you have finished your testing and no longer need deployed AWS infrastructures
pulumi destroy
# remove the stack if you don't need it anymore
pulumi stack rm
```

## Release process

TODO

## Resources

### Pulumi

Pulumi is like a training wheel for Terraform, the API is generated from Terraform schema with complete type
definitions. Whatever possible in Terraform is also possible in Pulumi, and you have type checking. To learn more about
it, you can check the generated type definitions and official Terraform documentations, or check out Pulumi examples on
github.

-   Terraform doc: https://www.terraform.io/docs/providers/aws/index.html
-   Pulumi concepts: https://www.pulumi.com/docs/reference/programming-model/
-   Storing encrypted configurations: https://www.pulumi.com/docs/reference/config/
-   Pulumi aws provider doc: https://www.pulumi.com/docs/reference/clouds/aws/
-   Insert debug info: https://www.pulumi.com/blog/unified-logs-with-pulumi-logs/
-   Pulumi examples: https://github.com/pulumi/examples

## Job API documentation - Swagger spec

The Job Service API documentation can be found [here](https://discoveryinc.atlassian.net/wiki/spaces/VDP/pages/1340342658/DRAFT+GTP-+Job+Service+API+-+Technical+Specifications).

We can also use Swagger to create the API specifications. From annotations provided by the Swagger tool, we can detail all the necessary information for the use of each of the endpoints that make up the API. Currently, you can generate a local copy of the yaml file which contains the documentation.

Note: in a later phase of the project, there will be an endpoint that will take you straight to the documentation page.

1. Install Swagger by launching the following command:
    ```sh
    go get -u github.com/go-swagger/go-swagger/cmd/swagger
    ```

2. Launch the following command to generate the .yaml or .json file  with all the docummentation collected from Swagger annotations:
    ```sh
    swagger generate spec -o ./swagger.yaml --scan-models
    ```

3. Execute the following command to raise the user interface which will show all the information of the .json or .yaml file  generated in the previous step:
    ```sh
    swagger serve -F=swagger swagger.yaml
    ```

4. You should be able to see the Swagger UI with all the API information:

![Swagger UI](imgs/swaggerUI.jpg)