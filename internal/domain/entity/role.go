package entity

type Role string

const (
	RoleUser          Role = "user"
	RoleAdmin         Role = "admin"
	RolePremium       Role = "premium"
	RoleSupport       Role = "support"     // Destek ekibi
	RoleModerator     Role = "moderator"   // İçerik moderatörü
	RoleAnalyst       Role = "analyst"     // Kullanıcı davranış analisti
	RoleSecurityAdmin Role = "sec_admin"   // Güvenlik yöneticisi
	RoleSystemMonitor Role = "sys_monitor" // Sistem monitör
)

// RolePermissions her rolün izinlerini tanımlar
var RolePermissions = map[Role][]Permission{
	RoleAdmin: {
		PermissionAll,
	},
	RoleSecurityAdmin: {
		PermissionUserBlock,
		PermissionUserUnblock,
		PermissionViewAuditLogs,
		PermissionViewSecurityLogs,
		PermissionManageRoles,
	},
	RoleModerator: {
		PermissionUserBlock,
		PermissionUserUnblock,
		PermissionViewAuditLogs,
	},
	RoleSupport: {
		PermissionViewUserDetails,
		PermissionResetUserPassword,
		PermissionViewAuditLogs,
	},
	RoleAnalyst: {
		PermissionViewAnalytics,
		PermissionViewAuditLogs,
		PermissionExportData,
	},
	RoleSystemMonitor: {
		PermissionViewMetrics,
		PermissionViewLogs,
		PermissionViewAlerts,
	},
}
