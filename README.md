# Heimdall

Heimdall is a (crude) tool that checks AWS ELB IAM/ACM cert expiration dates. Output is in JSON via stdout, logging goes to stderr. I'm sure this will change, as it's not very convenient. It's a work-in-progress. 


## Building

This is a side project that I'm not putting a ton of effort into, so fetching deps/building is a bit yolo at the moment:

```
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
./heimdall --region us-east-1 --region us-west-2 | jq .
```
Specify as many regions as you want with additional `--region foo` flags.

Output is in JSON via stdout, logging goes to stderr. I prefer to use [jq](https://stedolan.github.io/jq/) to pretty-print the json output.

## Design

It's not fast, it's not beautiful. It works for what I need. 

#License

http://www.apache.org/licenses/LICENSE-2.0.html

---

© 2017 Dave Parfitt
