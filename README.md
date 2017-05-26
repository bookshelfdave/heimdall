# Heimdall

Heimdall is a simple tool that checks AWS ELB IAM/ACM cert expiration dates. Output is in JSON via stdout, logging goes to stderr. I'm sure this will change, as it's not very convenient. It's a work-in-progress. 


## Building

**Go 1.8 is required to build this project.**

This is a side project that I'm not putting a ton of effort into, so fetching deps/building is a bit yolo at the moment:

```
cd $GOPATH
mkdir -p ./src/github.com/metadave
cd ./src/github.com/metadave
git clone https://github.com/metadave/heimdall.git
cd heimdall
make deps
make build
```

## IAM Setup

You'll need a user with API access using the following IAM policy:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Stmt1494008280000",
            "Effect": "Allow",
            "Action": [
                "elasticloadbalancing:DescribeLoadBalancers"
            ],
            "Resource": [
                "*"
            ]
        },
        {
            "Sid": "Stmt1494008362000",
            "Effect": "Allow",
            "Action": [
                "acm:DescribeCertificate"
            ],
            "Resource": [
                "*"
            ]
        },
        {
            "Sid": "Stmt1494008388000",
            "Effect": "Allow",
            "Action": [
                "iam:GetServerCertificate"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
```

## Usage

```shell
# The following command will show managed EC2 certs on ELB's (classic) that are expiring in <= 30 days
./heimdall -ec2-region us-west-2 -warn-days 30 -skip-expired

# You can specify a file containing host:port lines to check:
./heimdall -hosts ./certhosts -warn-days 90 -skip-expired

With the `certhosts` file similar to:

```
www.google.com:443
www.mozilla.org:443
```


# for json output:
./heimdall -ec2-region eu-west-1 -ec2-region us-east-1 -ec2-region us-west-2 -ec2-region ap-northeast-1 -json | jq .
```
Specify as many regions as you want with additional `-ec2-region foo` flags.

## Design

It's not fast, it's not beautiful. It works for what I need. 

## TODO

- General cert expiration checking via OpenSSL

# License

http://www.apache.org/licenses/LICENSE-2.0.html

---

© 2017 Dave Parfitt
