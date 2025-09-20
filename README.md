# nws-ledger
Expense tracker with Needs-Wants-Savings support

### Roadmap (phases)

1. **MVP:** 🟢 done
    - JSON storage
    - `add`, `list`, `summary` by NWS
2. **Phase 2:** 🟡 next todo
    - Domain categories support
    - Filtering (`-last`, `-domain`, `-nws`)
3. **Phase 3:** 🟡 todo
    - SQLite backend
    - Monthly budgets
    - Export/Import
4. **Phase 4:** 🟡 todo
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

### Current behavior (as of this commit)

```bash
$ nws add --amount 250 --name Groceries --nws needs
$ nws add --amount 150 --name Stocks --nws savings
$ nws add --amount 900 --name iPhone --nws wants
$ nws list
name,amount,nws
Groceries,250,needs
Stocks,150,savings
iPhone,900,wants
$ nws summary # expenses sums by NWS category
nws,amount
needs,250
wants,900
savings,150
total,1300
$ cat store.json # backed by plain JSON storage for now
[{"amount":250,"nws":"needs","domain":"","name":"Groceries"},{"amount":150,"nws":"savings","domain":"","name":"Stocks"},{"amount":900,"nws":"wants","domain":"","name":"iPhone"}]
```
