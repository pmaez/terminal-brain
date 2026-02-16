# **🧠 Terminal Brain — Technical Specification**

**Terminal Brain** is a modular, local-first TUI application written in Go that acts as a “cognitive operating system” for the user.

It reduces mental load by organizing tasks, thoughts, debugging reflection, time awareness, and command recall — all inside the terminal.

---

## **🎯 Core Product Vision**

Terminal Brain is:

- A **modular TUI platform**
- A **personal cognitive workspace**
- Extensible via new feature modules
- Fully local (no cloud required)

Each feature is implemented as a **Module** plugged into a shared application core.

---

## **🏗 Architecture Overview**

```
TUI Shell (Bubble Tea)
        ↓
Core Engine
  - App state
  - Module registry
  - Navigation
  - Persistence
        ↓
Modules (Feature Units)
```

---

## **🧩 Technology Stack**

|**Concern**|**Tool**|
|---|---|
|Language|Go|
|TUI Framework|Bubble Tea|
|Components|Bubbles|
|Styling|Lip Gloss|
|CLI Entry|Cobra|
|Database|SQLite|
|Storage Type|Local filesystem + SQLite|

---

## **🧠 Core Engine Responsibilities**

1. Global application state
2. Module registration
3. Active module switching
4. Event handling (future)
5. Shared database connection

### **AppState**

```go
type AppState struct {
    DB      *sql.DB
    Modules map[string]Module
    Active  string
}
```

### **Module Interface**

```go
type Module interface {
    Name() string
    Init(*AppState) error
    Update(tea.Msg) (Module, tea.Cmd)
    View() string
}
```

---

## **🖥 TUI Shell Requirements**

- Sidebar lists modules
- Main pane shows active module
- Global keys: `q` -> Quit; `up/down` -> Navigate modules; `Enter` -> Switch module; `Esc` -> Focus sidebar
- Footer shows global + module-specific hints (wraps to two lines if needed)
- Top title: "Terminal Brain"

---

## **🧩 Initial Modules (MVP)**

### **1️⃣ Stack (Mental Stack Manager)**

Top item = current focus.

**Features**

- Add task
- Remove task
- Move task up/down
- Reorder persists

```sql
CREATE TABLE stack_tasks (
  id INTEGER PRIMARY KEY,
  title TEXT,
  created_at DATETIME
);
```

---

### **2️⃣ Dump (Brain Dump)**

Stores raw unstructured text for later organization.

**Features**

- Add dump entries
- View recent entries (most recent first)

---

### **3️⃣ Duck (Rubber Duck Mode)**

Guided prompts:

1. What did you expect?
2. What happened?
3. What did you try?
4. What changed recently?

**Features**

- Save sessions to SQLite
- View recent sessions list
- View a full session
- Edit saved sessions
- Delete saved sessions

---

### **4️⃣ Command Memory**

Searchable saved commands.

**Features**

- Search by substring
- Add commands and descriptions
- Mark usage (updates `uses` and `last_used`)

```sql
CREATE TABLE commands (
  id INTEGER PRIMARY KEY,
  command TEXT,
  description TEXT,
  uses INTEGER,
  last_used DATETIME
);
```

---

### **5️⃣ Tiny Task Timer**

Track real time spent.

**Features**

- Start/stop a single active timer
- Delete selected logs
- Clear all logs (confirm on `D`)
- Elapsed time display (refreshes every 5 seconds)

```sql
CREATE TABLE time_logs (
  id INTEGER PRIMARY KEY,
  label TEXT,
  start DATETIME,
  end DATETIME
);
```

---

## **📁 Project Structure**

```
terminal-brain/
│
├── cmd/
├── internal/
│   ├── core/
│   ├── ui/
│   └── modules/
│       ├── stack/
│       ├── dump/
│       ├── duck/
│       ├── cmdmem/
│       └── timer/
└── main.go
```

---

## **🔌 Adding New Modules**

1. Create module folder
2. Implement Module interface
3. Register in AppState
4. Add to sidebar

---

## **🧠 Design Philosophy**

- Local-first
- Keyboard-driven
- Fast
- Low cognitive load
- Each module reduces mental friction

---

## **🚀 MVP Definition**

MVP complete when:

- Shell navigation works
- Stack module add/remove/reorder works
- Timer logs persist
- Command memory searchable
- Data persists between runs

**Status (as implemented)**

- Shell navigation works ✅
- Stack add/remove/reorder works ✅
- Timer logs persist ✅
- Command memory searchable ✅
- Data persists between runs ✅

MVP is currently met.

---

## **🔮 Future Extensions**

- Plugin system
- Cross-module analytics
- AI summaries
- Markdown export
- Optional sync
