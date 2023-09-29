package service

func (s *CLIService) Delete() error {
	if err := s.repo.Drop(); err != nil {
		return err
	}

	return nil
}
