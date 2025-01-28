package business

import (
	"log"

	"github.com/maSrifiyy/utils"
)

func SendThresholdNotification() error {
	thresholdedProducts, err := pgStore.GetLowStockProducts()
	if err != nil {
		log.Printf("Error fetching low stock products: %v", err)
		return err
	}

	if len(thresholdedProducts) == 0 {
		log.Println("No low-stock products found. No notifications sent.")
		return nil
	}

	for _, tped := range thresholdedProducts {
		log.Printf("Processing low stock alert: ProductID=%s, Name=%s, SellerEmail=%s", tped.ProductID, tped.Name, tped.Email)
		if err := utils.SendEmail(tped.Email, tped.FirstName, tped.Name, tped.ProductID); err != nil {
			log.Printf("Error sending email to %s for product %s: %v", tped.Email, tped.ProductID, err)
			continue
		}
		log.Printf("Email successfully sent to %s for product %s", tped.Email, tped.ProductID)
	}

	log.Println("All low stock notifications processed successfully.")
	return nil
}
