package rexUtils

import (
	"errors"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrInvalidPolicyType = errors.New("invalid policy type")
)

func BatchAddCasBinPolices(e *casbin.Enforcer, pType string, rules [][]string) error {
	var err error
	mapPolices := make(map[string][]string)
	uniquePolices := [][]string{}
	for _, rule := range rules {
		mapPolices[strings.Join(rule, "-")] = rule
	}

	if pType == "p" {
		for _, policy := range mapPolices {
			b, _ := e.HasPolicy(policy)
			if !b {
				uniquePolices = append(uniquePolices, policy)
			}
		}
		if len(uniquePolices) > 0 {
			_, err := e.AddPolicies(uniquePolices)
			if err != nil {
				logx.Errorf("BatchAddCasBinPolices p type AddPolicies failed, err :%v", err)
				return err
			}
			logx.Infof("BatchAddCasBinPolices p type AddPolicies success")
		}
	} else if pType == "g" {
		for _, policy := range mapPolices {
			// note: 判断策略是否存在
			b, _ := e.HasGroupingPolicy(policy)
			if !b {
				uniquePolices = append(uniquePolices, policy)
			}
		}
		if len(uniquePolices) > 0 {
			_, err := e.AddGroupingPolicies(uniquePolices)
			if err != nil {
				logx.Errorf("BatchAddCasBinPolices g type AddPolicies failed, err :%v", err)
				return err
			}
			logx.Infof("BatchAddCasBinPolices g type AddPolicies success")
		}
	} else {
		err = ErrInvalidPolicyType
		return err
	}
	return err
}

func BatchRemoveCasBinPolices(e *casbin.Enforcer, pType string, rules [][]string) error {
	var err error
	mapPolices := make(map[string][]string)
	uniquePolices := [][]string{}
	for _, rule := range rules {
		mapPolices[strings.Join(rule, "-")] = rule
	}

	if pType == "p" {
		for _, policy := range mapPolices {
			b, _ := e.HasPolicy(policy)
			if b {
				uniquePolices = append(uniquePolices, policy)
			}
		}
		_, err := e.RemovePolicies(uniquePolices)
		if err != nil {
			logx.Errorf("BatchAddCasBinPolices p type RemovePolicies failed, err :%v", err)
			return err
		}
		logx.Infof("BatchAddCasBinPolices p type RemovePolicies success")
	} else if pType == "g" {
		for _, policy := range mapPolices {
			// note: 判断策略是否存在
			b, _ := e.HasGroupingPolicy(policy)
			if b {
				uniquePolices = append(uniquePolices, policy)
			}
		}
		_, err := e.RemoveGroupingPolicies(uniquePolices)
		if err != nil {
			logx.Errorf("BatchAddCasBinPolices g type RemoveGroupingPolicies failed, err :%v", err)
			return err
		}
		logx.Infof("BatchAddCasBinPolices g type RemoveGroupingPolicies success")
	} else {
		err = ErrInvalidPolicyType
		return err
	}
	return err
}
