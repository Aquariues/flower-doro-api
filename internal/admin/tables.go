package admin

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
)

func GetFlowersTable(ctx *context.Context) table.Table {
	flowers := newTable(ctx, db.Int, "id")

	info := flowers.GetInfo()
	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("Kind", "kind", db.Varchar).FieldSortable().FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("English Name", "english_name", db.Varchar).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Vietnamese Name", "vietnamese_name", db.Varchar).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Rarity", "rarity", db.Varchar).FieldFilterable()
	info.AddField("Asset Name", "asset_name", db.Varchar)
	info.AddField("Updated At", "updated_at", db.Timestamp).FieldSortable()
	info.SetTable("flowers").SetTitle("Flowers").SetDescription("Flower catalog")

	formList := flowers.GetForm()
	formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("Kind", "kind", db.Varchar, form.Text)
	formList.AddField("English Name", "english_name", db.Varchar, form.Text)
	formList.AddField("Vietnamese Name", "vietnamese_name", db.Varchar, form.Text)
	formList.AddField("English Description", "english_description", db.Text, form.TextArea)
	formList.AddField("Vietnamese Description", "vietnamese_description", db.Text, form.TextArea)
	formList.AddField("Rarity", "rarity", db.Varchar, form.SelectSingle).FieldOptions(rarityOptions()).FieldDefault("common")
	formList.AddField("Asset Name", "asset_name", db.Varchar, form.Text)
	formList.AddField("Created At", "created_at", db.Timestamp, form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("Updated At", "updated_at", db.Timestamp, form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.SetTable("flowers").SetTitle("Flowers").SetDescription("Flower catalog")

	return flowers
}

func GetUsersTable(ctx *context.Context) table.Table {
	users := newTable(ctx, db.Varchar, "id")

	info := users.GetInfo()
	info.AddField("ID", "id", db.Varchar).FieldSortable()
	info.AddField("Device ID", "device_id", db.Varchar).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Display Name", "display_name", db.Varchar).FieldEditAble(editType.Text)
	info.AddField("Locale", "locale", db.Varchar).FieldFilterable()
	info.AddField("Created At", "created_at", db.Timestamp).FieldSortable()
	info.SetTable("users").SetTitle("Users").SetDescription("FlowerDoro users")

	formList := users.GetForm()
	formList.AddField("ID", "id", db.Varchar, form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("Device ID", "device_id", db.Varchar, form.Text)
	formList.AddField("Display Name", "display_name", db.Varchar, form.Text)
	formList.AddField("Locale", "locale", db.Varchar, form.SelectSingle).FieldOptions(localeOptions()).FieldDefault("vi")
	formList.AddField("Created At", "created_at", db.Timestamp, form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("Updated At", "updated_at", db.Timestamp, form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.SetTable("users").SetTitle("Users").SetDescription("FlowerDoro users")

	return users
}

func GetGardenFlowersTable(ctx *context.Context) table.Table {
	gardenFlowers := newTable(ctx, db.Varchar, "id")

	info := gardenFlowers.GetInfo()
	info.AddField("ID", "id", db.Varchar).FieldSortable()
	info.AddField("User ID", "user_id", db.Varchar).FieldFilterable()
	info.AddField("Flower Kind", "flower_kind", db.Varchar).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Focus Minutes", "focus_minutes", db.Int).FieldSortable()
	info.AddField("Earned At", "earned_at", db.Timestamp).FieldSortable().FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.SetTable("garden_flowers").SetTitle("Garden Flowers").SetDescription("Earned flower rewards")

	formList := gardenFlowers.GetForm()
	formList.AddField("ID", "id", db.Varchar, form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("User ID", "user_id", db.Varchar, form.Text)
	formList.AddField("Flower Kind", "flower_kind", db.Varchar, form.Text)
	formList.AddField("Focus Session ID", "focus_session_id", db.Varchar, form.Text)
	formList.AddField("Focus Minutes", "focus_minutes", db.Int, form.Number).FieldDefault("30")
	formList.AddField("Earned At", "earned_at", db.Timestamp, form.Datetime)
	formList.AddField("Created At", "created_at", db.Timestamp, form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.SetTable("garden_flowers").SetTitle("Garden Flowers").SetDescription("Earned flower rewards")

	return gardenFlowers
}

func GetFocusSessionsTable(ctx *context.Context) table.Table {
	sessions := newTable(ctx, db.Varchar, "id")

	info := sessions.GetInfo()
	info.AddField("ID", "id", db.Varchar).FieldSortable()
	info.AddField("User ID", "user_id", db.Varchar).FieldFilterable()
	info.AddField("Focus Minutes", "focus_minutes", db.Int).FieldSortable()
	info.AddField("Rewarded", "rewarded", db.Bool).FieldFilterable()
	info.AddField("Started At", "started_at", db.Timestamp).FieldSortable().FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField("Completed At", "completed_at", db.Timestamp).FieldSortable()
	info.SetTable("focus_sessions").SetTitle("Focus Sessions").SetDescription("Focus session history")

	formList := sessions.GetForm()
	formList.AddField("ID", "id", db.Varchar, form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("User ID", "user_id", db.Varchar, form.Text)
	formList.AddField("Started At", "started_at", db.Timestamp, form.Datetime)
	formList.AddField("Completed At", "completed_at", db.Timestamp, form.Datetime)
	formList.AddField("Focus Minutes", "focus_minutes", db.Int, form.Number).FieldDefault("30")
	formList.AddField("Rewarded", "rewarded", db.Bool, form.Radio).FieldOptions(types.FieldOptions{
		{Text: "No", Value: "false"},
		{Text: "Yes", Value: "true"},
	}).FieldDefault("false")
	formList.AddField("Created At", "created_at", db.Timestamp, form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("Updated At", "updated_at", db.Timestamp, form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.SetTable("focus_sessions").SetTitle("Focus Sessions").SetDescription("Focus session history")

	return sessions
}

func newTable(ctx *context.Context, keyType db.DatabaseType, keyName string) table.Table {
	return table.NewDefaultTable(ctx, table.Config{
		Driver:     db.DriverPostgresql,
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		Connection: table.DefaultConnectionName,
		PrimaryKey: table.PrimaryKey{
			Type: keyType,
			Name: keyName,
		},
	})
}

func rarityOptions() types.FieldOptions {
	return types.FieldOptions{
		{Text: "Common", Value: "common"},
		{Text: "Uncommon", Value: "uncommon"},
		{Text: "Rare", Value: "rare"},
		{Text: "Legendary", Value: "legendary"},
	}
}

func localeOptions() types.FieldOptions {
	return types.FieldOptions{
		{Text: "Tiếng Việt", Value: "vi"},
		{Text: "English", Value: "en"},
	}
}
