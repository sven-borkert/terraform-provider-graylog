package sidecar

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/convert"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/util"
)

const (
	keyNodes                  = "nodes"
	keySidecars               = "sidecars"
	keyNodeID                 = "node_id"
	keyAssignments            = "assignments"
	keyAssignmentScopeVersion = "assignment_scope_version"
	systemID                  = "system"
)

func getDataFromResourceData(d *schema.ResourceData) (map[string]interface{}, error) {
	data, err := convert.GetFromResourceData(d, Resource())
	if err != nil {
		return nil, err
	}
	util.RenameKey(data, keySidecars, keyNodes)
	delete(data, keyAssignmentScopeVersion)
	return data, nil
}

func validateAssignmentScope(d *schema.ResourceData) error {
	if d.Get(keyAssignmentScopeVersion).(int) != 1 {
		return errors.New("legacy sidecar state does not identify assignment ownership safely; back up Terraform state, remove only this graylog_sidecars resource with terraform state rm, then apply its configuration again to establish ownership without clearing unrelated assignments")
	}
	return nil
}

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	if err := convert.SetResourceData(d, Resource(), data); err != nil {
		return err
	}

	d.SetId(systemID)
	return nil
}
