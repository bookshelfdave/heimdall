package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

type RegionTests struct {
	Region string
	ELBs   []*AwsCertTest
}

func DoScan(region string) (*RegionTests, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Error("error:", err)
	}

	svc := elb.New(sess)
	params := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{},
	}

	resp, err := svc.DescribeLoadBalancers(params)

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	allRegionTests := RegionTests{Region: region}
	for _, elb := range resp.LoadBalancerDescriptions {
		log.Infof("%s @ %s", *elb.LoadBalancerName, *elb.DNSName)
		for _, listener := range elb.ListenerDescriptions {
			if listener.Listener.SSLCertificateId != nil {
				log.Debugf("     Cert: %s\n", *listener.Listener.SSLCertificateId)
				result, err := CheckAWSCert(*elb.LoadBalancerName,
					*elb.DNSName,
					region,
					*listener.Listener.SSLCertificateId)
				if err == nil {
					allRegionTests.ELBs = append(allRegionTests.ELBs, result)
				}
			}
		}
	}
	return &allRegionTests, nil
}
