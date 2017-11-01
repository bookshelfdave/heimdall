# Heimdall

Heimdall is a simple tool that checks AWS ELB cert expiration dates, as well as non-AWS managed certs.


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
            "Sid": "Stmt1494008280001",
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeRegions"
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

The following command will show managed EC2 certs on ELB's (classic) that are expiring in <= 30 days (I call these "managed certs"):

```shell
heimdall -ec2-region us-west-2 -warn-days 30 -skip-expired
```

You can scan all EC2 regions as well:

```
heimdall -all-regions -warn-days 30 -skip-expired
```

You can specify a file containing host:port lines to check (I call these "unmanaged certs"):

```shell
heimdall -hosts ./certhosts -warn-days 90 -skip-expired
```

With the `certhosts` file similar to:

```
www.google.com:443
www.mozilla.org:443
```

You can do "managed" and "unmanaged" certs at the same time:

```
heimdall -ec2-region us-west-2 -ec2-region us-east-1 -hosts ./foo -warn-days 30 -skip-expired
```

## Design

It's not fast, it's definitely not beautiful. It works for what I need. 

## TODO

- certs are checked serially, which I don't care about too much as I only need to run it once a month. It could be changed to run checks in parallel.

# License

http://www.apache.org/licenses/LICENSE-2.0.html

---

© 2017 Dave Parfitt
