package aws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-aws/aws/internal/tfawsresource"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/wafv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAwsWafv2WebACL_basic(t *testing.T) {
	var v wafv2.WebACL
	webACLName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_wafv2_web_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsWafv2WebACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsWafv2WebACLConfig_Basic(webACLName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsWafv2WebACLExists(resourceName, &v),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "wafv2", regexp.MustCompile(`regional/webacl/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", webACLName),
					resource.TestCheckResourceAttr(resourceName, "description", webACLName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "scope", wafv2.ScopeRegional),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.allow.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.block.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.0.cloudwatch_metrics_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.0.metric_name", "friendly-metric-name"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.0.sampled_requests_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccAwsWafv2WebACLConfig_BasicUpdate(webACLName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsWafv2WebACLExists(resourceName, &v),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "wafv2", regexp.MustCompile(`regional/webacl/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", webACLName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated"),
					resource.TestCheckResourceAttr(resourceName, "scope", wafv2.ScopeRegional),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.allow.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.block.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.0.cloudwatch_metrics_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.0.metric_name", "friendly-metric-name"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.0.sampled_requests_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "2"),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"name":             "rule-2",
						"priority":         "10",
						"action.#":         "1",
						"action.0.allow.#": "0",
						"action.0.block.#": "0",
						"action.0.count.#": "1",
						"statement.#":      "1",
						"statement.0.size_constraint_statement.#":                                 "1",
						"statement.0.size_constraint_statement.0.comparison_operator":             "LT",
						"statement.0.size_constraint_statement.0.field_to_match.#":                "1",
						"statement.0.size_constraint_statement.0.field_to_match.0.query_string.#": "1",
						"statement.0.size_constraint_statement.0.size":                            "50",
						"statement.0.size_constraint_statement.0.text_transformation.#":           "2",
					}),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*.statement.0.size_constraint_statement.0.text_transformation.*", map[string]string{
						"priority": "2",
						"type":     "CMD_LINE",
					}),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*.statement.0.size_constraint_statement.0.text_transformation.*", map[string]string{
						"priority": "5",
						"type":     "NONE",
					}),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"name":             "rule-1",
						"priority":         "1",
						"action.#":         "1",
						"action.0.allow.#": "1",
						"action.0.block.#": "0",
						"action.0.count.#": "0",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAwsWafv2WebACLImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccAwsWafv2WebACL_ChangeNameForceNew(t *testing.T) {
	var before, after wafv2.WebACL
	webACLName := acctest.RandomWithPrefix("tf-acc-test")
	ruleGroupNewName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_wafv2_web_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsWafv2WebACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsWafv2WebACLConfig_Basic(webACLName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsWafv2WebACLExists(resourceName, &before),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "wafv2", regexp.MustCompile(`regional/webacl/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", webACLName),
					resource.TestCheckResourceAttr(resourceName, "description", webACLName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "scope", wafv2.ScopeRegional),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.allow.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.block.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.0.cloudwatch_metrics_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.0.metric_name", "friendly-metric-name"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.0.sampled_requests_enabled", "false"),
				),
			},
			{
				Config: testAccAwsWafv2WebACLConfig_Basic(ruleGroupNewName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsWafv2WebACLExists(resourceName, &after),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "wafv2", regexp.MustCompile(`regional/webacl/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", ruleGroupNewName),
					resource.TestCheckResourceAttr(resourceName, "description", ruleGroupNewName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "scope", wafv2.ScopeRegional),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.allow.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.block.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.0.cloudwatch_metrics_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.0.metric_name", "friendly-metric-name"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.0.sampled_requests_enabled", "false"),
				),
			},
		},
	})
}

func TestAccAwsWafv2WebACL_Disappears(t *testing.T) {
	var v wafv2.WebACL
	webACLName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_wafv2_web_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsWafv2WebACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsWafv2WebACLConfig_Minimal(webACLName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsWafv2WebACLExists(resourceName, &v),
					testAccCheckResourceDisappears(testAccProvider, resourceAwsWafv2WebACL(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAwsWafv2WebACL_ManagedRuleGroupStatement(t *testing.T) {
	var v wafv2.WebACL
	webACLName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_wafv2_web_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsWafv2WebACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsWafv2WebACLConfig_ManagedRuleGroupStatement(webACLName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsWafv2WebACLExists(resourceName, &v),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "wafv2", regexp.MustCompile(`regional/webacl/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", webACLName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"name":                      "rule-1",
						"action.#":                  "0",
						"override_action.#":         "1",
						"override_action.0.count.#": "0",
						"override_action.0.none.#":  "1",
						"statement.#":               "1",
						"statement.0.managed_rule_group_statement.#":                 "1",
						"statement.0.managed_rule_group_statement.0.name":            "AWSManagedRulesCommonRuleSet",
						"statement.0.managed_rule_group_statement.0.vendor_name":     "AWS",
						"statement.0.managed_rule_group_statement.0.excluded_rule.#": "0",
					}),
				),
			},
			{
				Config: testAccAwsWafv2WebACLConfig_ManagedRuleGroupStatement_Update(webACLName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsWafv2WebACLExists(resourceName, &v),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "wafv2", regexp.MustCompile(`regional/webacl/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", webACLName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"name":                      "rule-1",
						"action.#":                  "0",
						"override_action.#":         "1",
						"override_action.0.count.#": "1",
						"override_action.0.none.#":  "0",
						"statement.#":               "1",
						"statement.0.managed_rule_group_statement.#":                      "1",
						"statement.0.managed_rule_group_statement.0.name":                 "AWSManagedRulesCommonRuleSet",
						"statement.0.managed_rule_group_statement.0.vendor_name":          "AWS",
						"statement.0.managed_rule_group_statement.0.excluded_rule.#":      "2",
						"statement.0.managed_rule_group_statement.0.excluded_rule.0.name": "SizeRestrictions_QUERYSTRING",
						"statement.0.managed_rule_group_statement.0.excluded_rule.1.name": "NoUserAgent_HEADER",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAwsWafv2WebACLImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccAwsWafv2WebACL_Minimal(t *testing.T) {
	var v wafv2.WebACL
	webACLName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_wafv2_web_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsWafv2WebACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsWafv2WebACLConfig_Minimal(webACLName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsWafv2WebACLExists(resourceName, &v),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "wafv2", regexp.MustCompile(`regional/webacl/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", webACLName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "scope", wafv2.ScopeRegional),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.allow.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.block.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.0.cloudwatch_metrics_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.0.metric_name", "friendly-metric-name"),
					resource.TestCheckResourceAttr(resourceName, "visibility_config.0.sampled_requests_enabled", "false"),
				),
			},
		},
	})
}

func TestAccAwsWafv2WebACL_RateBasedStatement(t *testing.T) {
	var v wafv2.WebACL
	webACLName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_wafv2_web_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsWafv2WebACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsWafv2WebACLConfig_RateBasedStatement(webACLName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsWafv2WebACLExists(resourceName, &v),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "wafv2", regexp.MustCompile(`regional/webacl/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", webACLName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"name":                               "rule-1",
						"action.#":                           "1",
						"action.0.allow.#":                   "0",
						"action.0.block.#":                   "0",
						"action.0.count.#":                   "1",
						"statement.#":                        "1",
						"statement.0.rate_based_statement.#": "1",
						"statement.0.rate_based_statement.0.aggregate_key_type":     "IP",
						"statement.0.rate_based_statement.0.limit":                  "50000",
						"statement.0.rate_based_statement.0.scope_down_statement.#": "0",
					}),
				),
			},
			{
				Config: testAccAwsWafv2WebACLConfig_RateBasedStatement_Update(webACLName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsWafv2WebACLExists(resourceName, &v),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "wafv2", regexp.MustCompile(`regional/webacl/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", webACLName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"name":                               "rule-1",
						"action.#":                           "1",
						"action.0.allow.#":                   "0",
						"action.0.block.#":                   "0",
						"action.0.count.#":                   "1",
						"statement.#":                        "1",
						"statement.0.rate_based_statement.#": "1",
						"statement.0.rate_based_statement.0.aggregate_key_type":                                           "IP",
						"statement.0.rate_based_statement.0.limit":                                                        "10000",
						"statement.0.rate_based_statement.0.scope_down_statement.#":                                       "1",
						"statement.0.rate_based_statement.0.scope_down_statement.0.geo_match_statement.#":                 "1",
						"statement.0.rate_based_statement.0.scope_down_statement.0.geo_match_statement.0.country_codes.#": "2",
						"statement.0.rate_based_statement.0.scope_down_statement.0.geo_match_statement.0.country_codes.0": "US",
						"statement.0.rate_based_statement.0.scope_down_statement.0.geo_match_statement.0.country_codes.1": "NL",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAwsWafv2WebACLImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccAwsWafv2WebACL_RuleGroupReferenceStatement(t *testing.T) {
	var v wafv2.WebACL
	var idx int
	webACLName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_wafv2_web_acl.test"
	excludedRules := []interface{}{
		map[string]interface{}{
			"name": "rule-to-exclude-b",
		},
		map[string]interface{}{
			"name": "rule-to-exclude-a",
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsWafv2WebACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsWafv2WebACLConfig_RuleGroupReferenceStatement(webACLName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsWafv2WebACLExists(resourceName, &v),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "wafv2", regexp.MustCompile(`regional/webacl/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", webACLName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					computeWafv2RuleGroupRefStatementIndex(&v, &idx, []interface{}{}),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"name":                      "rule-1",
						"override_action.#":         "1",
						"override_action.0.count.#": "1",
						"override_action.0.none.#":  "0",
						"statement.#":               "1",
						"statement.0.rule_group_reference_statement.#":                 "1",
						"statement.0.rule_group_reference_statement.0.excluded_rule.#": "0",
					}),
					// TODO: TypeSet check need helper for regex
					testAccMatchResourceAttrArnWithIndexesAddr(resourceName, "rule.%d.statement.0.rule_group_reference_statement.0.arn", &idx, regexp.MustCompile(`regional/rulegroup/.+$`)),
				),
			},
			{
				Config: testAccAwsWafv2WebACLConfig_RuleGroupReferenceStatement_Update(webACLName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsWafv2WebACLExists(resourceName, &v),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "wafv2", regexp.MustCompile(`regional/webacl/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", webACLName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					computeWafv2RuleGroupRefStatementIndex(&v, &idx, excludedRules),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"name":                      "rule-1",
						"override_action.#":         "1",
						"override_action.0.count.#": "1",
						"override_action.0.none.#":  "0",
						"statement.#":               "1",
						"statement.0.rule_group_reference_statement.#":                      "1",
						"statement.0.rule_group_reference_statement.0.excluded_rule.#":      "2",
						"statement.0.rule_group_reference_statement.0.excluded_rule.0.name": "rule-to-exclude-b",
						"statement.0.rule_group_reference_statement.0.excluded_rule.1.name": "rule-to-exclude-a",
					}),
					// TODO: TypeSet check need helper for regex
					testAccMatchResourceAttrArnWithIndexesAddr(resourceName, "rule.%d.statement.0.rule_group_reference_statement.0.arn", &idx, regexp.MustCompile(`regional/rulegroup/.+$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAwsWafv2WebACLImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccAwsWafv2WebACL_Tags(t *testing.T) {
	var v wafv2.WebACL
	webACLName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_wafv2_web_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsWafv2WebACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsWafv2WebACLConfig_OneTag(webACLName, "Tag1", "Value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsWafv2WebACLExists(resourceName, &v),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "wafv2", regexp.MustCompile(`regional/webacl/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Tag1", "Value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAwsWafv2WebACLImportStateIdFunc(resourceName),
			},
			{
				Config: testAccAwsWafv2WebACLConfig_TwoTags(webACLName, "Tag1", "Value1Updated", "Tag2", "Value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsWafv2WebACLExists(resourceName, &v),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "wafv2", regexp.MustCompile(`regional/webacl/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Tag1", "Value1Updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.Tag2", "Value2"),
				),
			},
			{
				Config: testAccAwsWafv2WebACLConfig_OneTag(webACLName, "Tag2", "Value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsWafv2WebACLExists(resourceName, &v),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "wafv2", regexp.MustCompile(`regional/webacl/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Tag2", "Value2"),
				),
			},
		},
	})
}

// Reference: https://github.com/terraform-providers/terraform-provider-aws/issues/13862
func TestAccAwsWafv2WebACL_MaxNestedRateBasedStatements(t *testing.T) {
	var v wafv2.WebACL
	webACLName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_wafv2_web_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsWafv2WebACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsWafv2WebACLConfig_multipleNestedRateBasedStatements(webACLName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsWafv2WebACLExists(resourceName, &v),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "wafv2", regexp.MustCompile(`regional/webacl/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", webACLName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"statement.#":                                                                                                      "1",
						"statement.0.rate_based_statement.#":                                                                               "1",
						"statement.0.rate_based_statement.0.limit":                                                                         "300",
						"statement.0.rate_based_statement.0.aggregate_key_type":                                                            "IP",
						"statement.0.rate_based_statement.0.scope_down_statement.#":                                                        "1",
						"statement.0.rate_based_statement.0.scope_down_statement.0.not_statement.#":                                        "1",
						"statement.0.rate_based_statement.0.scope_down_statement.0.not_statement.0.statement.#":                            "1",
						"statement.0.rate_based_statement.0.scope_down_statement.0.not_statement.0.statement.0.or_statement.#":             "1",
						"statement.0.rate_based_statement.0.scope_down_statement.0.not_statement.0.statement.0.or_statement.0.statement.#": "2",
						"statement.0.rate_based_statement.0.scope_down_statement.0.not_statement.0.statement.0.or_statement.0.statement.0.regex_pattern_set_reference_statement.#": "1",
						"statement.0.rate_based_statement.0.scope_down_statement.0.not_statement.0.statement.0.or_statement.0.statement.1.ip_set_reference_statement.#":            "1",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAwsWafv2WebACLImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccAwsWafv2WebACL_MaxNestedOperatorStatements(t *testing.T) {
	var v wafv2.WebACL
	webACLName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_wafv2_web_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsWafv2WebACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsWafv2WebACLConfig_multipleNestedOperatorStatements(webACLName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsWafv2WebACLExists(resourceName, &v),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "wafv2", regexp.MustCompile(`regional/webacl/.+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", webACLName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"statement.#":                                                                                    "1",
						"statement.0.and_statement.#":                                                                    "1",
						"statement.0.and_statement.0.statement.#":                                                        "2",
						"statement.0.and_statement.0.statement.0.not_statement.#":                                        "1",
						"statement.0.and_statement.0.statement.0.not_statement.0.statement.#":                            "1",
						"statement.0.and_statement.0.statement.0.not_statement.0.statement.0.or_statement.#":             "1",
						"statement.0.and_statement.0.statement.0.not_statement.0.statement.0.or_statement.0.statement.#": "2",
						"statement.0.and_statement.0.statement.0.not_statement.0.statement.0.or_statement.0.statement.0.regex_pattern_set_reference_statement.#": "1",
						"statement.0.and_statement.0.statement.0.not_statement.0.statement.0.or_statement.0.statement.1.ip_set_reference_statement.#":            "1",
						"statement.0.and_statement.0.statement.1.geo_match_statement.#":                                                                          "1",
						"statement.0.and_statement.0.statement.1.geo_match_statement.0.country_codes.0":                                                          "NL",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAwsWafv2WebACLImportStateIdFunc(resourceName),
			},
		},
	})
}

// Calculates the index which isn't static because ARN is generated as part of the test
func computeWafv2RuleGroupRefStatementIndex(r *wafv2.WebACL, idx *int, e []interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ruleResource := resourceAwsWafv2WebACL().Schema["rule"].Elem.(*schema.Resource)
		rule := map[string]interface{}{
			"name":     "rule-1",
			"priority": 1,
			"action":   []interface{}{},
			"override_action": []interface{}{
				map[string]interface{}{
					"none":  []interface{}{},
					"count": make([]interface{}, 1),
				},
			},
			"statement": []interface{}{
				map[string]interface{}{
					"and_statement":                         []interface{}{},
					"byte_match_statement":                  []interface{}{},
					"geo_match_statement":                   []interface{}{},
					"ip_set_reference_statement":            []interface{}{},
					"managed_rule_group_statement":          []interface{}{},
					"not_statement":                         []interface{}{},
					"or_statement":                          []interface{}{},
					"rate_based_statement":                  []interface{}{},
					"regex_pattern_set_reference_statement": []interface{}{},
					"rule_group_reference_statement": []interface{}{
						map[string]interface{}{
							"arn":           aws.StringValue(r.Rules[0].Statement.RuleGroupReferenceStatement.ARN),
							"excluded_rule": e,
						},
					},
					"size_constraint_statement": []interface{}{},
					"sqli_match_statement":      []interface{}{},
					"xss_match_statement":       []interface{}{},
				},
			},
			"visibility_config": []interface{}{
				map[string]interface{}{
					"cloudwatch_metrics_enabled": false,
					"metric_name":                "friendly-rule-metric-name",
					"sampled_requests_enabled":   false,
				},
			},
		}

		f := schema.HashResource(ruleResource)
		*idx = f(rule)

		return nil
	}
}

func testAccCheckAwsWafv2WebACLDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_wafv2_web_acl" {
			continue
		}

		conn := testAccProvider.Meta().(*AWSClient).wafv2conn
		resp, err := conn.GetWebACL(
			&wafv2.GetWebACLInput{
				Id:    aws.String(rs.Primary.ID),
				Name:  aws.String(rs.Primary.Attributes["name"]),
				Scope: aws.String(rs.Primary.Attributes["scope"]),
			})

		if err == nil {
			if resp == nil || resp.WebACL == nil {
				return fmt.Errorf("Error getting WAFv2 WebACL")
			}
			if aws.StringValue(resp.WebACL.Id) == rs.Primary.ID {
				return fmt.Errorf("WAFv2 WebACL %s still exists", rs.Primary.ID)
			}
		}

		// Return nil if the WebACL is already destroyed
		if isAWSErr(err, wafv2.ErrCodeWAFNonexistentItemException, "") {
			return nil
		}

		return err
	}

	return nil
}

func testAccCheckAwsWafv2WebACLExists(n string, v *wafv2.WebACL) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No WAFv2 WebACL ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).wafv2conn
		resp, err := conn.GetWebACL(&wafv2.GetWebACLInput{
			Id:    aws.String(rs.Primary.ID),
			Name:  aws.String(rs.Primary.Attributes["name"]),
			Scope: aws.String(rs.Primary.Attributes["scope"]),
		})

		if err != nil {
			return err
		}

		if resp == nil || resp.WebACL == nil {
			return fmt.Errorf("Error getting WAFv2 WebACL")
		}

		if aws.StringValue(resp.WebACL.Id) == rs.Primary.ID {
			*v = *resp.WebACL
			return nil
		}

		return fmt.Errorf("WAFv2 WebACL (%s) not found", rs.Primary.ID)
	}
}

func testAccAwsWafv2WebACLConfig_Basic(name string) string {
	return fmt.Sprintf(`
resource "aws_wafv2_web_acl" "test" {
  name        = "%[1]s"
  description = "%[1]s"
  scope       = "REGIONAL"

  default_action {
    allow {}
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "friendly-metric-name"
    sampled_requests_enabled   = false
  }
}
`, name)
}

func testAccAwsWafv2WebACLConfig_BasicUpdate(name string) string {
	return fmt.Sprintf(`
resource "aws_wafv2_web_acl" "test" {
  name        = "%s"
  description = "Updated"
  scope       = "REGIONAL"

  default_action {
    block {}
  }

  rule {
    name     = "rule-2"
    priority = 10

    action {
      count {}
    }

    statement {
      size_constraint_statement {
        comparison_operator = "LT"
        size                = 50

        field_to_match {
          query_string {}
        }

        text_transformation {
          priority = 5
          type     = "NONE"
        }

        text_transformation {
          priority = 2
          type     = "CMD_LINE"
        }
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = false
      metric_name                = "friendly-rule-metric-name"
      sampled_requests_enabled   = false
    }
  }

  rule {
    name     = "rule-1"
    priority = 1

    action {
      allow {}
    }

    statement {
      geo_match_statement {
        country_codes = ["US", "NL"]
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = false
      metric_name                = "friendly-rule-metric-name"
      sampled_requests_enabled   = false
    }
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "friendly-metric-name"
    sampled_requests_enabled   = false
  }
}
`, name)
}

func testAccAwsWafv2WebACLConfig_ManagedRuleGroupStatement(name string) string {
	return fmt.Sprintf(`
resource "aws_wafv2_web_acl" "test" {
  name        = "%[1]s"
  description = "%[1]s"
  scope       = "REGIONAL"

  default_action {
    allow {}
  }

  rule {
    name     = "rule-1"
    priority = 1

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesCommonRuleSet"
        vendor_name = "AWS"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = false
      metric_name                = "friendly-rule-metric-name"
      sampled_requests_enabled   = false
    }
  }

  tags = {
    Tag1 = "Value1"
    Tag2 = "Value2"
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "friendly-metric-name"
    sampled_requests_enabled   = false
  }
}
`, name)
}

func testAccAwsWafv2WebACLConfig_ManagedRuleGroupStatement_Update(name string) string {
	return fmt.Sprintf(`
resource "aws_wafv2_web_acl" "test" {
  name        = "%[1]s"
  description = "%[1]s"
  scope       = "REGIONAL"

  default_action {
    allow {}
  }

  rule {
    name     = "rule-1"
    priority = 1

    override_action {
      count {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesCommonRuleSet"
        vendor_name = "AWS"

        excluded_rule {
          name = "SizeRestrictions_QUERYSTRING"
        }

        excluded_rule {
          name = "NoUserAgent_HEADER"
        }
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = false
      metric_name                = "friendly-rule-metric-name"
      sampled_requests_enabled   = false
    }
  }

  tags = {
    Tag1 = "Value1"
    Tag2 = "Value2"
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "friendly-metric-name"
    sampled_requests_enabled   = false
  }
}
`, name)
}

func testAccAwsWafv2WebACLConfig_RateBasedStatement(name string) string {
	return fmt.Sprintf(`
resource "aws_wafv2_web_acl" "test" {
  name        = "%[1]s"
  description = "%[1]s"
  scope       = "REGIONAL"

  default_action {
    block {}
  }

  rule {
    name     = "rule-1"
    priority = 1

    action {
      count {}
    }

    statement {
      rate_based_statement {
        limit = 50000
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = false
      metric_name                = "friendly-rule-metric-name"
      sampled_requests_enabled   = false
    }
  }

  tags = {
    Tag1 = "Value1"
    Tag2 = "Value2"
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "friendly-metric-name"
    sampled_requests_enabled   = false
  }
}
`, name)
}

func testAccAwsWafv2WebACLConfig_RateBasedStatement_Update(name string) string {
	return fmt.Sprintf(`
resource "aws_wafv2_web_acl" "test" {
  name        = "%[1]s"
  description = "%[1]s"
  scope       = "REGIONAL"

  default_action {
    block {}
  }

  rule {
    name     = "rule-1"
    priority = 1

    action {
      count {}
    }

    statement {
      rate_based_statement {
        limit              = 10000
        aggregate_key_type = "IP"

        scope_down_statement {
          geo_match_statement {
            country_codes = ["US", "NL"]
          }
        }
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = false
      metric_name                = "friendly-rule-metric-name"
      sampled_requests_enabled   = false
    }
  }

  tags = {
    Tag1 = "Value1"
    Tag2 = "Value2"
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "friendly-metric-name"
    sampled_requests_enabled   = false
  }
}
`, name)
}

func testAccAwsWafv2WebACLConfig_RuleGroupReferenceStatement(name string) string {
	return fmt.Sprintf(`
resource "aws_wafv2_rule_group" "test" {
  capacity = 10
  name     = "rule-group-%[1]s"
  scope    = "REGIONAL"

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "friendly-metric-name"
    sampled_requests_enabled   = false
  }
}

resource "aws_wafv2_web_acl" "test" {
  name  = "%[1]s"
  scope = "REGIONAL"

  default_action {
    block {}
  }

  rule {
    name     = "rule-1"
    priority = 1

    override_action {
      count {}
    }

    statement {
      rule_group_reference_statement {
        arn = aws_wafv2_rule_group.test.arn
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = false
      metric_name                = "friendly-rule-metric-name"
      sampled_requests_enabled   = false
    }
  }

  tags = {
    Tag1 = "Value1"
    Tag2 = "Value2"
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "friendly-metric-name"
    sampled_requests_enabled   = false
  }
}
`, name)
}

func testAccAwsWafv2WebACLConfig_RuleGroupReferenceStatement_Update(name string) string {
	return fmt.Sprintf(`
resource "aws_wafv2_rule_group" "test" {
  capacity = 10
  name     = "rule-group-%[1]s"
  scope    = "REGIONAL"

  rule {
    name     = "rule-1"
    priority = 1

    action {
      count {}
    }

    statement {
      geo_match_statement {
        country_codes = ["NL"]
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = false
      metric_name                = "friendly-rule-metric-name"
      sampled_requests_enabled   = false
    }
  }

  rule {
    name     = "rule-to-exclude-a"
    priority = 10

    action {
      allow {}
    }

    statement {
      geo_match_statement {
        country_codes = ["US"]
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = false
      metric_name                = "friendly-rule-metric-name"
      sampled_requests_enabled   = false
    }
  }

  rule {
    name     = "rule-to-exclude-b"
    priority = 15

    action {
      allow {}
    }

    statement {
      geo_match_statement {
        country_codes = ["GB"]
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = false
      metric_name                = "friendly-rule-metric-name"
      sampled_requests_enabled   = false
    }
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "friendly-metric-name"
    sampled_requests_enabled   = false
  }
}

resource "aws_wafv2_web_acl" "test" {
  name  = "%[1]s"
  scope = "REGIONAL"

  default_action {
    block {}
  }

  rule {
    name     = "rule-1"
    priority = 1

    override_action {
      count {}
    }

    statement {
      rule_group_reference_statement {
        arn = aws_wafv2_rule_group.test.arn

        excluded_rule {
          name = "rule-to-exclude-b"
        }

        excluded_rule {
          name = "rule-to-exclude-a"
        }
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = false
      metric_name                = "friendly-rule-metric-name"
      sampled_requests_enabled   = false
    }
  }

  tags = {
    Tag1 = "Value1"
    Tag2 = "Value2"
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "friendly-metric-name"
    sampled_requests_enabled   = false
  }
}
`, name)
}

func testAccAwsWafv2WebACLConfig_Minimal(name string) string {
	return fmt.Sprintf(`
resource "aws_wafv2_web_acl" "test" {
  name  = "%s"
  scope = "REGIONAL"

  default_action {
    allow {}
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "friendly-metric-name"
    sampled_requests_enabled   = false
  }
}
`, name)
}

func testAccAwsWafv2WebACLConfig_OneTag(name, tagKey, tagValue string) string {
	return fmt.Sprintf(`
resource "aws_wafv2_web_acl" "test" {
  name        = "%[1]s"
  description = "%[1]s"
  scope       = "REGIONAL"

  default_action {
    allow {}
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "friendly-metric-name"
    sampled_requests_enabled   = false
  }

  tags = {
    "%s" = "%s"
  }
}
`, name, tagKey, tagValue)
}

func testAccAwsWafv2WebACLConfig_TwoTags(name, tag1Key, tag1Value, tag2Key, tag2Value string) string {
	return fmt.Sprintf(`
resource "aws_wafv2_web_acl" "test" {
  name        = "%[1]s"
  description = "%[1]s"
  scope       = "REGIONAL"

  default_action {
    allow {}
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "friendly-metric-name"
    sampled_requests_enabled   = false
  }

  tags = {
    "%s" = "%s"
    "%s" = "%s"
  }
}
`, name, tag1Key, tag1Value, tag2Key, tag2Value)
}

func testAccAwsWafv2WebACLConfig_multipleNestedRateBasedStatements(name string) string {
	return fmt.Sprintf(`
resource "aws_wafv2_regex_pattern_set" "test" {
  name  = "%[1]s"
  scope = "REGIONAL"
  regular_expression {
    regex_string = "one"
  }
}

resource "aws_wafv2_ip_set" "test" {
  name               = "%[1]s"
  scope              = "REGIONAL"
  ip_address_version = "IPV4"
  addresses          = ["1.2.3.4/32", "5.6.7.8/32"]
}

resource "aws_wafv2_web_acl" "test" {
  name        = "%[1]s"
  description = "%[1]s"
  scope       = "REGIONAL"

  default_action {
    allow {}
  }

  rule {
    name     = "rule"
    priority = 0

    action {
      block {}
    }

    statement {
      rate_based_statement {
        limit              = 300
        aggregate_key_type = "IP"

        scope_down_statement {
          not_statement {
            statement {
              or_statement {
                statement {
                  regex_pattern_set_reference_statement {
                    arn = aws_wafv2_regex_pattern_set.test.arn
                    field_to_match {
                      uri_path {}
                    }
                    text_transformation {
                      type     = "LOWERCASE"
                      priority = 1
                    }
                  }
                }
                statement {
                  ip_set_reference_statement {
                    arn = aws_wafv2_ip_set.test.arn
                  }
                }
              }
            }
          }
        }
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = false
      metric_name                = "rule"
      sampled_requests_enabled   = false
    }
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "waf"
    sampled_requests_enabled   = false
  }
}
`, name)
}

func testAccAwsWafv2WebACLConfig_multipleNestedOperatorStatements(name string) string {
	return fmt.Sprintf(`
resource "aws_wafv2_regex_pattern_set" "test" {
  name  = "%[1]s"
  scope = "REGIONAL"
  regular_expression {
    regex_string = "one"
  }
}

resource "aws_wafv2_ip_set" "test" {
  name               = "%[1]s"
  scope              = "REGIONAL"
  ip_address_version = "IPV4"
  addresses          = ["1.2.3.4/32", "5.6.7.8/32"]
}

resource "aws_wafv2_web_acl" "test" {
  name        = "%[1]s"
  description = "%[1]s"
  scope       = "REGIONAL"

  default_action {
    allow {}
  }

  rule {
    name     = "rule"
    priority = 0

    action {
      block {}
    }

    statement {
      and_statement {
        statement {
          not_statement {
            statement {
              or_statement {
                statement {
                  regex_pattern_set_reference_statement {
                    arn = aws_wafv2_regex_pattern_set.test.arn
                    field_to_match {
                      uri_path {}
                    }
                    text_transformation {
                      type     = "LOWERCASE"
                      priority = 1
                    }
                  }
                }
                statement {
                  ip_set_reference_statement {
                    arn = aws_wafv2_ip_set.test.arn
                  }
                }
              }
            }
          }
        }
        statement {
          geo_match_statement {
            country_codes = ["NL"]
          }
        }
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = false
      metric_name                = "rule"
      sampled_requests_enabled   = false
    }
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "waf"
    sampled_requests_enabled   = false
  }
}
`, name)
}

func testAccAwsWafv2WebACLImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s/%s/%s", rs.Primary.ID, rs.Primary.Attributes["name"], rs.Primary.Attributes["scope"]), nil
	}
}
