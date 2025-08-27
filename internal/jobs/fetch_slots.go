// package jobs - джобы
package jobs

import (
	"fmt"
	"kvparser/internal/domain"
	"kvparser/internal/infrastructure"
	"kvparser/internal/logger"
	"kvparser/internal/services"

	"github.com/go-co-op/gocron"
)

func RegisterFetchSlotsJob(s *gocron.Scheduler, rem services.Remeber[domain.DoctorMatch], cronVal string, logger logger.Logger, cp services.ChromeParserService, do domain.DoctorOptions, tg infrastructure.Bot) error {

	_, err := s.Cron("*/5 * * * *").StartImmediately().Do(func() {
		page, err := cp.DoctorsSchedulePage()
		if err != nil {
			logger.Error("can't fetch page:", err)
			return
		}

		processed, err := domain.FindMatches(page, do)
		if err != nil {
			logger.Error("can't  process page:", err)
			return
		}

		for _, val := range processed.Matches {
			if !rem.Remember(val) {
				tg.SendMessage(fmt.Sprintf("%s %s %s %s\n", val.Name, val.Speciality, val.Status, val.Subdivision))
				fmt.Printf("%s %s %s %s\n", val.Name, val.Speciality, val.Status, val.Subdivision)
			}
		}
	})

	return err

}
