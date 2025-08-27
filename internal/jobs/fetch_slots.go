// package jobs - джобы
package jobs

import (
	"fmt"
	"kvparser/internal/domain"
	"kvparser/internal/logger"
	"kvparser/internal/services"

	"github.com/go-co-op/gocron"
)

func RegisterFetchSlotsJob(s *gocron.Scheduler, cronVal string, logger logger.Logger, cp services.ChromeParserService) error {

	_, err := s.Cron("*/5 * * * *").StartImmediately().Do(func() {
		page, err := cp.DoctorsSchedulePage()
		if err != nil {
			logger.Error("can't fetch page:", err)
			return
		}

		processed, err := domain.FindMatches(page, domain.DoctorOptions{})
		if err != nil {
			logger.Error("can't  process page:", err)
			return
		}

		for _, val := range processed.Matches {
			fmt.Printf("%s %s %s %s\n", val.Name, val.Speciality, val.Status, val.Subdivision)
		}
	})

	return err

}
