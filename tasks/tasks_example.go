package tasks

/*

// https://github.com/xh3b4sd/matic/blob/master/collector/collector.go
// https://github.com/xh3b4sd/matic/blob/master/collector/package_import.go

type Ctx struct {
	WorkingDir string

	// Variable name of the created middleware server, if any. We assume there is
	// only one created middleware server. Maybe that is not true for all cases.
	ServerName string

	Files []File
}

func (gcg *GoClientCollector) GenerateClient(wd string) error {
	// Create task context.
	ctx := &Ctx{
		WorkingDir: wd,
	}

	// Create a new queue.
	q := taskqPkg.NewQueue(ctx)

	// Run tasks.
	err := q.RunTasks(
		taskqPkg.InSeries(
			SourceCodeTask,
			PackageImportTask,
			ServerNameTask,
			ServeStmtTask,
			// find middlewares for each route
			// find possible responses for each route
		),
	)

	if err != nil {
		return Mask(err)
	}

	return nil
}

func PackageImportTask(ctx interface{}) error {
}
*/