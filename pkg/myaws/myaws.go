package myaws

import (
	"context"
	"reflect"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	log "github.com/sirupsen/logrus"
)

type Zone struct {
	Id      string
	Name    string
	Private bool
	Records int64
}

func GetHostedZone(domain string) []Zone {

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	log.WithFields(
		log.Fields{
			"domain": domain,
		}).Debug()

	reg := regexp.MustCompile(`^` + domain + `.$`)

	zones := []Zone{}
	svc := route53.NewFromConfig(cfg)
	input := &route53.ListHostedZonesInput{}
	res, _ := svc.ListHostedZones(context.TODO(), input)

	for _, v := range res.HostedZones {
		if reg.MatchString(aws.ToString(v.Name)) {
			zones = append(zones, Zone{
				Id:      aws.ToString(v.Id)[12:],
				Name:    aws.ToString((v.Name)),
				Private: v.Config.PrivateZone,
				Records: aws.ToInt64((v.ResourceRecordSetCount)),
			})

			log.WithFields(
				log.Fields{
					"zoneid":  aws.ToString(v.Id),
					"private": v.Config.PrivateZone,
					"name":    aws.ToString(v.Name),
					"records": aws.ToInt64(v.ResourceRecordSetCount),
				}).Debug()
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

	//	res, _ := svc.ListResourceRecordSets(context.TODO(), input)

	res, _ := ListAllResourceRecordSets(svc, input)

	reg := regexp.MustCompile(`^` + host)

	log.WithFields(
		log.Fields{
			"host": host,
		}).Debug()
	for _, v := range res {

		if hasField(v.AliasTarget, "DNSName") {
			//			j, _ := json.Marshal(v)
			//		fmt.Println(string(j))

			/*
				log.WithFields(
					log.Fields{
						"value": aws.ToString(v.Name),
						"dns":   aws.ToString(v.AliasTarget.DNSName),
					}).Debug()
			*/
			if reg.MatchString(aws.ToString(v.Name)) {

				log.WithFields(
					log.Fields{
						"revalue": aws.ToString(v.Name),
						"redns":   aws.ToString(v.AliasTarget.DNSName),
					}).Debug()
				return aws.ToString(v.AliasTarget.DNSName)
			}

		}
		//return aws.ToString(v.AliasTarget.DNSName)
		//}

	}
	return ""

}

func ListAllResourceRecordSets(svc *route53.Client, input *route53.ListResourceRecordSetsInput) (rrsets []types.ResourceRecordSet, err error) {
	//	res, _ := svc.ListResourceRecordSets(context.TODO(), input)
	//	return res

	for {
		var resp *route53.ListResourceRecordSetsOutput
		resp, err = svc.ListResourceRecordSets(context.TODO(), input)
		if err != nil {
			return
		} else {

			log.WithFields(
				log.Fields{
					"calllistrrs": "+",
				}).Debug()
			rrsets = append(rrsets, resp.ResourceRecordSets...)
			if resp.IsTruncated {
				input.StartRecordName = resp.NextRecordName
				input.StartRecordType = resp.NextRecordType
				input.StartRecordIdentifier = resp.NextRecordIdentifier
			} else {
				break
			}
		}
	}

	// unescape wildcards
	/*
		for _, rrset := range rrsets {
			rrset.Name = aws.String(unescaper.Replace(*rrset.Name))
		}
	*/

	return
}

func GetALB(DNSName string) string {

	if len(DNSName) > 0 {
		DNSName = DNSName[:len(DNSName)-1]
	} else {
		log.Warn(
			"Can not found DNSName from recoads.",
		)
		return ""
	}

	log.WithFields(
		log.Fields{
			"DNSName": DNSName,
		}).Debug()

	parts := strings.Split(DNSName, ".")
	host := parts[:2]

	// remove dualstack
	if host[0] == "dualstack" {
		host = host[1:]
	}

	// remove internal
	nparts := strings.Split(host[0], "-")
	if nparts[0] == "internal" {
		nparts = nparts[1:]
	}
	region := host[1]
	// remove uid
	nparts = nparts[:len(nparts)-1]

	lbname := strings.Join(nparts, "-")

	log.WithFields(
		log.Fields{
			"DNSName": DNSName,
			"name":    lbname,
			"host":    host,
			"region":  region,
		}).Debug()

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := elbv2.NewFromConfig(cfg)
	input := &elbv2.DescribeLoadBalancersInput{Names: []string{lbname}}
	res, _ := svc.DescribeLoadBalancers(context.TODO(), input)

	for _, v := range res.LoadBalancers {

		alb := aws.ToString(v.LoadBalancerArn)
		return alb
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
