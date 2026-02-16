package core

import "context"

type Stores struct {
	Stack   StackStore
	Dump    DumpStore
	Duck    DuckStore
	Command CommandStore
	Time    TimeStore
}

type StackTask struct {
	ID        int64
	Title     string
	Position  int
	CreatedAt string
	UpdatedAt string
}

type DumpEntry struct {
	ID        int64
	Body      string
	CreatedAt string
}

type DuckSession struct {
	ID        int64
	Q1        string
	Q2        string
	Q3        string
	Q4        string
	CreatedAt string
}

type CommandEntry struct {
	ID          int64
	Command     string
	Description string
	Uses        int
	LastUsed    string
	CreatedAt   string
}

type TimeLog struct {
	ID        int64
	Label     string
	Start     string
	End       string
	CreatedAt string
}

type StackStore interface {
	List(ctx context.Context) ([]StackTask, error)
	Add(ctx context.Context, title string) (StackTask, error)
	Delete(ctx context.Context, id int64) error
	UpdateOrder(ctx context.Context, ids []int64) error
}

type DumpStore interface {
	List(ctx context.Context, limit int) ([]DumpEntry, error)
	Add(ctx context.Context, body string) (DumpEntry, error)
}

type DuckStore interface {
	List(ctx context.Context, limit int) ([]DuckSession, error)
	Add(ctx context.Context, q1, q2, q3, q4 string) (DuckSession, error)
	Update(ctx context.Context, id int64, q1, q2, q3, q4 string) error
	Delete(ctx context.Context, id int64) error
}

type CommandStore interface {
	Search(ctx context.Context, query string, limit int) ([]CommandEntry, error)
	Add(ctx context.Context, command, description string) (CommandEntry, error)
	Touch(ctx context.Context, id int64) error
}

type TimeStore interface {
	List(ctx context.Context, limit int) ([]TimeLog, error)
	Start(ctx context.Context, label, start string) (TimeLog, error)
	Stop(ctx context.Context, id int64, end string) error
	Delete(ctx context.Context, id int64) error
	Clear(ctx context.Context) error
}
