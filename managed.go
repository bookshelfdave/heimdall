package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/iam"
)

type AwsCertTest struct {
	ElbName     string
	ElbDNS      string
	Expiration  time.Time
	ExpText     string
	AWSCertType string
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

func CheckAWSCert(elbName string, dnsName string, region string, certARN string) (*AwsCertTest, error) {
	s := strings.Split(certARN, "/")
	arnPrefix, certName := s[0], s[1]
	if strings.HasPrefix(arnPrefix, "arn:aws:acm") {
		return processACMCert(elbName, dnsName, region, certARN)
	} else {
		return processIAMCert(elbName, dnsName, region, certName)
	}
}
