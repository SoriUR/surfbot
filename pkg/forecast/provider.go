package forecast

import (
	// "context"
	"log"
	"strconv"

	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"

	// "u40apps.com/surfbot/pkg/setup"
)

// var forecastCollection *mongo.Collection

func GetForecast(spotName string, daysLimit int) (*string, error) {

	// setupCollection()

	log.Println("Requesting forecast for: ", spotName)
	response, err := getForecast(spotName)
	if err != nil {
		return nil, err
	}

	log.Println("Formatting forecast: ")
	formattedForecast := formatForecast(*response, daysLimit)

	var result = "Forecast at " + spotName + " for " + strconv.Itoa(daysLimit) + " days:" + formattedForecast

	return &result, nil
}

// MARK: - collection

// func setupCollection() {
// 	var err error

// 	if forecastCollection == nil {
// 		forecastCollection, err = setup.GetCollection("forecast")
// 		if err != nil {
// 			log.Printf("Error while setting up forecast collection. %v", err)
// 		} else {
// 			log.Printf("Successfully set up forecast collection")
// 		}
// 	}
// }

// MARK: - Forecast

func getForecast(spotName string) (*Forecast, error) {
	log.Printf("Getting forecast for %v", spotName)

	// if forecastCollection == nil {
	// 	return requestAndSaveForecast(spotName)
	// }

	// log.Printf("Fetching forecast from DB for %v", spotName)
	// var dbForecast Forecast
	// filter := bson.M{"spot": spotName}
	// err := forecastCollection.FindOne(context.TODO(), filter).Decode(&dbForecast)

	// if err == nil {
	// 	log.Printf("Found forecast in DB for %v", spotName)
	// 	return &dbForecast, nil
	// }

	return requestAndSaveForecast(spotName)
}

func requestAndSaveForecast(spotName string) (*Forecast, error) {
	forecast, err := requestForecast(spotName)

	// if err != nil && forecastCollection != nil {
	// 	log.Printf("Saving forecast in DB for %v", spotName)
	// 	newForecast := DBForecast{
	// 		Spot:     spotName,
	// 		Forecast: *forecast,
	// 	}
	// 	_, err := forecastCollection.InsertOne(context.TODO(), newForecast)
	// 	if err != nil {
	// 		log.Printf("Error while saving forecast in DB for %v. %v", spotName, err)
	// 	} else {
	// 		log.Printf("Successfully saved forecast in DB for %v", spotName)
	// 	}
	// }

	return forecast, err
}
