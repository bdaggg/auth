package entity

type Permission string

const (
	PermissionAll               Permission = "*"
	PermissionUserBlock         Permission = "user:block"
	PermissionUserUnblock       Permission = "user:unblock"
	PermissionViewAuditLogs     Permission = "audit:view"
	PermissionViewSecurityLogs  Permission = "security:view"
	PermissionManageRoles       Permission = "roles:manage"
	PermissionViewUserDetails   Permission = "user:view"
	PermissionResetUserPassword Permission = "user:reset_password"
	PermissionViewAnalytics     Permission = "analytics:view"
	PermissionExportData        Permission = "data:export"
	PermissionViewMetrics       Permission = "metrics:view"
	PermissionViewLogs          Permission = "logs:view"
	PermissionViewAlerts        Permission = "alerts:view"
)
