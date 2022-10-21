package stages

func WithJob(j Job) StageOpt {
	return func(s Stager) {
		s.SetJob(j)
	}
}

func WithNext(event, id int) StageOpt {
	return func(s Stager) {
		s.SetNext(event, id)
	}
}
