# nws-ledger
Expense tracker with Needs-Wants-Savings support ( ⚠️ experimental - a learning project )

### Roadmap (phases)

1. **MVP:** 🟢 done
    - JSON storage
    - `add`, `list`, `summary` by NWS
2. **Phase 2:** 🟡 next todo
    - Domain categories support (DONE)
    - Filtering (`-last`, `-domain`, `-nws`) (need to add timestamps to expenses)
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
$ nws add --amount 250 --name Weekly --nws needs --domain groceries
$ nws add --amount 150 --name VTI --nws savings --domain stocks
$ nws add --amount 900 --name iPhone --nws wants --domain shopping
$ nws list
name,amount,nws,domain
Weekly,250,needs,groceries
VTI,150,savings,stocks
iPhone,900,wants,shopping
$ nws summary # expenses sums by NWS category
nws,amount
needs,250
wants,900
savings,150
total,1300
$ nws summary --domain # expenses sums by Domain category
domain,amount
shopping,900
groceries,250
stocks,150
total,1300
$ cat store.json # backed by plain JSON storage for now
[{"amount":250,"nws":"needs","domain":"groceries","name":"Weekly"},{"amount":150,"nws":"savings","domain":"stocks","name":"VTI"},{"amount":900,"nws":"wants","domain":"shopping","name":"iPhone"}]
```
