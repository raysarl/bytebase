package api

import "fmt"

// PlanType is the type for a plan.
type PlanType int

const (
	// FREE is the plan type for FREE.
	FREE PlanType = iota
	// TEAM is the plan type for TEAM.
	TEAM
	// ENTERPRISE is the plan type for ENTERPRISE.
	ENTERPRISE
)

// String returns the string format of plan type.
func (p PlanType) String() string {
	switch p {
	case FREE:
		return "FREE"
	case TEAM:
		return "TEAM"
	case ENTERPRISE:
		return "ENTERPRISE"
	}
	return ""
}

// FeatureType is the type of a feature.
type FeatureType string

const (
	// Change Workflow

	// FeatureBackwardCompatibilty checks if a DDL change is backward compatible.
	// See https://docs.bytebase.com/features/sql-advisor/backward-compatibility-migration-check
	//
	// Currently, we only support MySQL/TiDB thanks to github.com/pingcap/parser
	FeatureBackwardCompatibilty FeatureType = "bb.feature.backward-compatibility"
	// FeatureSchemaDrift detects if there occurs schema drift.
	// See https://docs.bytebase.com/features/drift-detection
	FeatureSchemaDrift FeatureType = "bb.feature.schema-drift"
	// FeatureTaskScheduleTime allows user to run task at a scheduled time
	FeatureTaskScheduleTime FeatureType = "bb.feature.task-schedule-time"
	// FeatureMultiTenancy allows user to enable tenant mode for the project.
	//
	// Tenant mode allows user to track a group of homogeneous database changes together.
	// e.g. A game studio may deploy many servers, each server is fully isolated with its
	// own database. When a new game version is released, it may require to upgrade the
	// underlying database schema, then tenant mode will help the studio to track the
	// schema change across all databases.
	FeatureMultiTenancy FeatureType = "bb.feature.multi-tenancy"
	// FeatureDBAWorkflow enforces the DBA workflow.
	//
	// - Developers can't create and view instances since they are exclusively by DBA, they can
	//   only access database.
	// - Developers can submit troubleshooting issue.
	FeatureDBAWorkflow FeatureType = "bb.feature.dba-workflow"
	// FeatureDataSource exposes the data source concept.
	//
	// Currently, we DO NOT expose this feature.
	//
	// Internally Bytebase stores instance username/password in a seperate data source model.
	// This allows a single instance to have multiple data sources (e.g. one RW and one RO).
	// And from the user's perspective, the username/password
	// look like the property of the instance, which are not. They are the property of data source which
	// in turns belongs to the instance.
	// - Support defining extra data source for a database and exposing the related data source UI.
	FeatureDataSource FeatureType = "bb.feature.data-source"

	// Policy Control

	// FeatureApprovalPolicy allows user to specify approval policy for the environment
	//
	// e.g. One can configure to NOT require approval for dev environment while require
	//      manual approval for production.
	FeatureApprovalPolicy FeatureType = "bb.feature.approval-policy"
	// FeatureBackupPolicy allows user to specify backup policy for the environment
	//
	// e.g. One can configure to NOT require backup for dev environment while require
	//      weekly backup for staging and daily backup for production.
	FeatureBackupPolicy FeatureType = "bb.feature.backup-policy"

	// Admin & Security

	// FeatureRBAC enables RBAC.
	//
	// - Workspace level RBAC
	// - Project level RBAC
	FeatureRBAC FeatureType = "bb.feature.rbac"

	// Feature3rdPartyLogin allows user to login using 3rd party account.
	//
	// Currently, we only support GitLab EE/CE OAuth login.
	Feature3rdPartyLogin FeatureType = "bb.feature.3rd-party-login"
)

func (e FeatureType) String() string {
	switch e {
	case FeatureBackwardCompatibilty:
		return "bb.feature.backward-compatibility"
	case FeatureSchemaDrift:
		return "bb.feature.schema-drift"
	case FeatureTaskScheduleTime:
		return "bb.feature.task-schedule-time"
	case FeatureMultiTenancy:
		return "bb.feature.multi-tenancy"
	case FeatureDBAWorkflow:
		return "bb.feature.dba-workflow"
	case FeatureDataSource:
		return "bb.feature.data-source"
	case FeatureApprovalPolicy:
		return "bb.feature.approval-policy"
	case FeatureBackupPolicy:
		return "bb.feature.backup-policy"
	case FeatureRBAC:
		return "bb.feature.rbac"
	case Feature3rdPartyLogin:
		return "bb.feature.3rd-party-login"
	}
	return ""
}

// Name returns a readable name of the feature
func (e FeatureType) Name() string {
	switch e {
	case FeatureBackwardCompatibilty:
		return "Backward compatibility"
	case FeatureSchemaDrift:
		return "Schema drift"
	case FeatureTaskScheduleTime:
		return "Task schedule time"
	case FeatureMultiTenancy:
		return "Multi-tenancy"
	case FeatureDBAWorkflow:
		return "DBA workflow"
	case FeatureDataSource:
		return "Data source"
	case FeatureApprovalPolicy:
		return "Approval policy"
	case FeatureBackupPolicy:
		return "Backup policy"
	case FeatureRBAC:
		return "RBAC"
	case Feature3rdPartyLogin:
		return "3rd party login"
	}
	return ""
}

// AccessErrorMessage returns a error message with feature name and minimum supported plan.
func (e FeatureType) AccessErrorMessage() string {
	plan := e.minimumSupportedPlan()
	return fmt.Sprintf("%s is a %s feature, please upgrade to access it.", e.Name(), plan.String())
}

// minimumSupportedPlan will find the minimum plan which support the target feature.
func (e FeatureType) minimumSupportedPlan() PlanType {
	for i, enabled := range FeatureMatrix[e] {
		if enabled {
			return PlanType(i)
		}
	}

	return ENTERPRISE
}

// FeatureMatrix is a map from the a particular feature to the respective enablement of a particular plan
var FeatureMatrix = map[FeatureType][3]bool{
	"bb.feature.backward-compatibility": {false, true, true},
	"bb.feature.schema-drift":           {false, true, true},
	"bb.feature.task-schedule-time":     {false, true, true},
	"bb.feature.multi-tenancy":          {false, true, true},
	"bb.feature.dba-workflow":           {false, false, true},
	"bb.feature.data-source":            {false, false, false},
	"bb.feature.approval-policy":        {false, true, true},
	"bb.feature.backup-policy":          {false, true, true},
	"bb.feature.rbac":                   {false, true, true},
	"bb.feature.3rd-party-login":        {false, true, true},
}

// Plan is the API message for a plan.
type Plan struct {
	Type PlanType `jsonapi:"attr,type"`
}

// PlanPatch is the API message for patching a plan.
type PlanPatch struct {
	Type PlanType `jsonapi:"attr,type"`
}
