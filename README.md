# nws-ledger
Expense tracker with Needs-Wants-Savings support

### Roadmap (phases)

1. **MVP:**
    - JSON storage
    - `add`, `list`, `summary` by NWS
2. **Phase 2:**
    - Domain categories support
    - Filtering (`-last`, `-domain`, `-nws`)
3. **Phase 3:**
    - SQLite backend
    - Monthly budgets
    - Export/Import
4. **Phase 4:**
    - TUI
    - Sync (maybe via Git repo push/pull)


### Core commands:

- `add` -> add an expense
- `list` -> show recent expenses
- `summary` -> totals by NWS or by domain
- `categories` -> manage domain categories
- `import/export` -> backup or migrate data


```bash
nws add 250 "Groceries - supermarket" --nws need --domain groceries
nws list --last 7d
nws summary --nws
nws summary --domain
```
