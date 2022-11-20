package subscription

import "github.com/huoxue1/qinglong-go/models"

var (
	DISABLESTATUS = 1
	ENABLESTATUS  = 0
)

func AddSubscription(subscriptions *models.Subscriptions) (int, error) {
	subscriptions.Status = 1
	return models.AddSubscription(subscriptions)

}

func UpdateSubscription(subscriptions *models.Subscriptions) error {
	return models.UpdateSubscription(subscriptions)

}

func DeleteSubscription(ids []int) error {
	for _, id := range ids {
		err := models.DeleteSubscription(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func DisableSubscription(ids []int) error {
	for _, id := range ids {
		sub, err := models.GetSubscription(id)
		if err != nil {
			continue
		}
		sub.IsDisabled = 1
		err = models.UpdateSubscription(sub)
		if err != nil {
			return err
		}
	}
	return nil
}

func EnableSubscription(ids []int) error {
	for _, id := range ids {
		sub, err := models.GetSubscription(id)
		if err != nil {
			continue
		}
		sub.IsDisabled = 0
		err = models.UpdateSubscription(sub)
		if err != nil {
			return err
		}
	}
	return nil
}

func RunSubscription(ids []int) error {
	for _, id := range ids {
		sub, err := models.GetSubscription(id)
		if err != nil {
			continue
		}
		sub.IsDisabled = 0
		go downloadFiles(sub)
	}
	return nil
}

func StopSubscription(ids []int) error {
	for _, id := range ids {
		sub, err := models.GetSubscription(id)
		if err != nil {
			continue
		}
		stopSubscription(sub)
	}
	return nil
}
