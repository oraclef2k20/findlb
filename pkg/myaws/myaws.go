package myaws

import (
	"context"
	"log"
	"reflect"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

func GetHostedZone(domain string) map[string]bool {

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := route53.NewFromConfig(cfg)
	input := &route53.ListHostedZonesInput{}
	res, _ := svc.ListHostedZones(context.TODO(), input)
	reg := regexp.MustCompile(domain)

	zones := map[string]bool{}
	for _, v := range res.HostedZones {
		if reg.MatchString(aws.ToString(v.Name)) {
			zones[aws.ToString(v.Id)[12:]] = v.Config.PrivateZone
		}
	}

	return zones
}

func GetDNSFromRecoard(hostedzone string, host string) string {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := route53.NewFromConfig(cfg)

	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedzone),
	}

	res, _ := svc.ListResourceRecordSets(context.TODO(), input)

	reg := regexp.MustCompile(`^` + host)

	for _, v := range res.ResourceRecordSets {

		if hasField(v.AliasTarget, "DNSName") {
			//			j, _ := json.Marshal(v)
			//		fmt.Println(string(j))

			if reg.MatchString(aws.ToString(v.Name)) {
				//fmt.Println(aws.ToString(v.AliasTarget.DNSName))
				//fmt.Println(aws.ToString(v.Name))
				return aws.ToString(v.AliasTarget.DNSName)
			}

		}
		//return aws.ToString(v.AliasTarget.DNSName)
		//}

	}
	return ""

}

func GetALB(DNSName string) string {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := elbv2.NewFromConfig(cfg)

	input := &elbv2.DescribeLoadBalancersInput{}
	if len(DNSName) > 0 {
		DNSName = DNSName[:len(DNSName)-1]
	}
	reg := regexp.MustCompile(DNSName)
	res, _ := svc.DescribeLoadBalancers(context.TODO(), input)

	for _, v := range res.LoadBalancers {

		if reg.MatchString(aws.ToString(v.DNSName)) {
			alb := aws.ToString(v.LoadBalancerArn)
			return alb
		}
	}

	return ""
}

func hasField(v interface{}, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByName(name).IsValid()
}
