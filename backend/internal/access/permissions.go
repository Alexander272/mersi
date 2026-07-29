package access

const (
	ResourceSi                 ResourceSlug = "si"
	ResourceDepartment         ResourceSlug = "department"
	ResourceEmployee           ResourceSlug = "employee"
	ResourceLocation           ResourceSlug = "location"
	ResourceChannel            ResourceSlug = "channel"
	ResourceVerification       ResourceSlug = "verification"
	ResourceReserve            ResourceSlug = "reserve"
	ResourceDocuments          ResourceSlug = "documents"
	ResourceRoles              ResourceSlug = "roles"
	ResourceUsers              ResourceSlug = "users"
	ResourceRealms             ResourceSlug = "realms"
	ResourceSections           ResourceSlug = "sections"
	ResourceColumns            ResourceSlug = "columns"
	ResourceVerificationFields ResourceSlug = "verification-fields"
	ResourceCreatingForm       ResourceSlug = "creating-form"
	ResourceContextMenu        ResourceSlug = "context-menu"
	ResourceToolsMenu          ResourceSlug = "tools-menu"
	ResourceRepair             ResourceSlug = "repair"
	ResourcePreservation       ResourceSlug = "preservation"
	ResourceTransferToSave     ResourceSlug = "transfer-to-save"
	ResourceTransferToDep      ResourceSlug = "transfer-to-department"
	ResourceWriteOff           ResourceSlug = "write-off"
	ResourceHistoryTypes       ResourceSlug = "history-types"
)

var OrderOfResources = map[ResourceSlug]int{
	ResourceSi:                 1,
	ResourceDepartment:         2,
	ResourceEmployee:           3,
	ResourceLocation:           4,
	ResourceChannel:            5,
	ResourceVerification:       6,
	ResourceReserve:            7,
	ResourceDocuments:          8,
	ResourceSections:           9,
	ResourceColumns:            10,
	ResourceVerificationFields: 11,
	ResourceCreatingForm:       12,
	ResourceContextMenu:        13,
	ResourceToolsMenu:          14,
	ResourceRepair:             15,
	ResourcePreservation:       16,
	ResourceTransferToSave:     17,
	ResourceTransferToDep:      18,
	ResourceWriteOff:           19,
	ResourceHistoryTypes:       20,
	ResourceRealms:             21,
	ResourceUsers:              22,
	ResourceRoles:              23,
}

var Reg = NewRegistry(
	Resource{
		Slug:           ResourceSi,
		Name:           "СИ",
		Group:          "Операции",
		Description:    "Управление средствами измерений",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceDepartment,
		Name:           "Подразделения",
		Group:          "Справочники",
		Description:    "Управление подразделениями",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceEmployee,
		Name:           "Сотрудники",
		Group:          "Справочники",
		Description:    "Управление сотрудниками",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceLocation,
		Name:           "Размещение",
		Group:          "Операции",
		Description:    "Управление размещением",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceChannel,
		Name:           "Каналы уведомлений",
		Group:          "Уведомления",
		Description:    "Управление каналами уведомлений",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceVerification,
		Name:           "Поверка",
		Group:          "Операции",
		Description:    "Управление поверками",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceReserve,
		Name:           "Резерв",
		Group:          "Операции",
		Description:    "Управление резервом",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceDocuments,
		Name:           "Документы",
		Group:          "Операции",
		Description:    "Управление документами",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceSections,
		Name:           "Разделы",
		Group:          "Справочники",
		Description:    "Управление разделами",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceColumns,
		Name:           "Колонки",
		Group:          "Справочники",
		Description:    "Управление колонками",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceVerificationFields,
		Name:           "Поля поверки",
		Group:          "Справочники",
		Description:    "Управление полями поверки",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceCreatingForm,
		Name:           "Форма создания",
		Group:          "Справочники",
		Description:    "Управление формой создания",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceContextMenu,
		Name:           "Контекстное меню",
		Group:          "Настройки",
		Description:    "Управление контекстным меню",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceToolsMenu,
		Name:           "Меню инструментов",
		Group:          "Настройки",
		Description:    "Управление меню инструментов",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceRepair,
		Name:           "Ремонт",
		Group:          "Операции",
		Description:    "Управление ремонтами",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourcePreservation,
		Name:           "Консервация",
		Group:          "Операции",
		Description:    "Управление консервациями",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceTransferToSave,
		Name:           "Передача на хранение",
		Group:          "Операции",
		Description:    "Управление передачами на хранение",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceTransferToDep,
		Name:           "Передача в подразделение",
		Group:          "Операции",
		Description:    "Управление передачами в подразделение",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceWriteOff,
		Name:           "Списание",
		Group:          "Операции",
		Description:    "Управление списаниями",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceHistoryTypes,
		Name:           "Типы истории",
		Group:          "Справочники",
		Description:    "Управление типами истории",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceRealms,
		Name:           "Области",
		Group:          "Администрирование",
		Description:    "Управление областями",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceUsers,
		Name:           "Пользователи",
		Group:          "Администрирование",
		Description:    "Управление пользователями",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceRoles,
		Name:           "Роли",
		Group:          "Администрирование",
		Description:    "Управление ролями пользователей",
		AllowedActions: actions(Read, Write, Delete),
	},
)
