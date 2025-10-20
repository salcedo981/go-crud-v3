package scpFeatureOne

var (
	//EXPENSES WITH USER OWNERSHIP
	AddExpense = "SELECT add_expense($1, $2, $3, $4, $5, $6)"
	GetExpenses = "SELECT get_expenses($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	GetExpense = "SELECT get_expense($1, $2)"
	UpdateExpense = "SELECT update_expense($1, $2, $3, $4, $5, $6, $7)"

	//EXPENSES CATEGORY
	AddExpenseCategory = `SELECT add_expense_category(?,?)`
	GetExpenseCategories = `SELECT get_expense_categories(?,?)`
)
