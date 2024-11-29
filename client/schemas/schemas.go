package schemas

import "render-box/shared/db/repo"

type PageData struct {
	Tasks   []repo.Task
	Jobs    []repo.Job
	Workers []repo.Worker
}
