package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/iam"
)

type AwsCertTest struct {
	ElbName     string
	ElbDNS      string
	Expiration  time.Time
	ExpText     string
	AWSCertType string
}

type RegionTests struct {
	Region string
	ELBs   []*AwsCertTest
}

func listEC2Regions() ([]string, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		log.Error("error:", err)
	}

	var allRegions []string

	svc := ec2.New(sess)
	params := &ec2.DescribeRegionsInput{}
	resp, err := svc.DescribeRegions(params)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	for _, region := range resp.Regions {
		log.Debug(*region.RegionName)
		allRegions = append(allRegions, *region.RegionName)
	}
	return allRegions, nil
}

func processRegionELBs(region string) (*RegionTests, error) {
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
				result, err := checkAWSCert(*elb.LoadBalancerName,
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

func checkAWSCert(elbName string, dnsName string, region string, certARN string) (*AwsCertTest, error) {
	s := strings.Split(certARN, "/")
	arnPrefix, certName := s[0], s[1]
	if strings.HasPrefix(arnPrefix, "arn:aws:acm") {
		return processACMCert(elbName, dnsName, region, certARN)
	} else {
		return processIAMCert(elbName, dnsName, region, certName)
	}
}

func processACMCert(elbName string, dnsName string, region string, certARN string) (*AwsCertTest, error) {
	log.Debug("Detected ACM cert")
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))

	svc := acm.New(sess)
	params := &acm.DescribeCertificateInput{
		CertificateArn: aws.String(certARN),
	}
	resp, err := svc.DescribeCertificate(params)

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	exp := resp.Certificate.NotAfter

	log.Infof("Expiration: %s\n", exp)
	year, month, day, hour := DiffDate(time.Now(), *exp)
	expText := fmt.Sprintf("%d years, %d months, %d days, %d hours", year, month, day, hour)
	log.Debug(expText)
	result := AwsCertTest{elbName, dnsName, *exp, expText, "ACM"}
	return &result, nil
}

func processIAMCert(elbName string, dnsName string, region string, certName string) (*AwsCertTest, error) {
	log.Debugf("Cert name = %s\n", certName)
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))

	svc := iam.New(sess)

	params := &iam.GetServerCertificateInput{
		ServerCertificateName: aws.String(certName), // Required
	}
	resp, err := svc.GetServerCertificate(params)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	exp := resp.ServerCertificate.ServerCertificateMetadata.Expiration
	log.Infof("Expiration: %s\n", exp)
	year, month, day, hour := DiffDate(time.Now(), *exp)
	log.Infof("%d years, %d months, %d days, %d hours\n", year, month, day, hour)
	expText := fmt.Sprintf("%d years, %d months, %d days, %d hours", year, month, day, hour)
	log.Debug(expText)
	result := AwsCertTest{elbName, dnsName, *exp, expText, "ACM"}
	return &result, nil
}

func showManagedExpirations(allRegions []*RegionTests, warnDays int, skipExpired bool) {
	for _, region := range allRegions {
		log.Debugf("Region = %s\n", region.Region)
		if region.ELBs == nil {
			continue
		}
		for _, certTest := range region.ELBs {
			log.Debugf(" ELB = %s\n", certTest.ElbName)
			exp := time.Until(certTest.Expiration)
			daysUntil := int(exp.Hours() / 24)
			if daysUntil < 0 {
				if !skipExpired {
					log.Errorf("[%s] %s (%s) cert has expired: %s",
						certTest.AWSCertType, certTest.ElbName, certTest.ElbDNS, certTest.ExpText)
				}
			} else if daysUntil <= warnDays {
				log.Errorf("[%s] %s (%s) cert is expiring soon: %s (%s)",
					certTest.AWSCertType, certTest.ElbName, certTest.ElbDNS, certTest.Expiration, certTest.ExpText)
			}
		}
	}
}
