package db

import (
	"context"
	"errors"
	m "github.com/androzes/CovidCases/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type CovidStat struct {
	Code string `json:"code" bson:"_id"`
	Name string `json:"name" bson:"name"`
	NumCovidCases int `json:"num_covid_cases" bson:"num_covid_cases"`
	LastUpdated time.Time `json:"last_updated" bson:"last_updated"`
}

func GetStatsByStateCodes(ctx context.Context, stateCodes []string) ([]CovidStat, error) {
	ids := bson.A{}
	for _, code := range stateCodes {
		ids = append(ids, code)
	}

	filter := bson.M{
		"_id": bson.M {
			"$in": ids,
		},
	}

	cursor, err := collection().Find(ctx, filter)
	defer cursor.Close(ctx)
	if err != nil {
		return []CovidStat{}, err
	}

	stats := []CovidStat{}

	for cursor.Next(ctx) {
		var stat CovidStat
		err := cursor.Decode(&stat)
		if err != nil {
			return []CovidStat{}, err
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

func GetStatsByStateName(ctx context.Context, stateName string) (CovidStat, error) {
	filter := bson.M{
		"name": bson.M{
			"$regex" : primitive.Regex{
				Pattern: stateName,
				Options:"i",
			},
		},
	}

	var stat CovidStat

	err := collection().FindOne(ctx, filter).Decode(&stat)
	if err != nil {
		if err == mongo.ErrNoDocuments{
			return CovidStat{}, errors.New("Could not find state: "+ stateName)
		} else {
			return CovidStat{}, err
		}

	}

	return stat, nil
}

func GetTotalStatsForCountry(ctx context.Context) (CovidStat, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"name": bson.M{
					"$exists": true,
				},
				"num_covid_cases": bson.M{
					"$exists": true,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": 1,
				"num_covid_cases": bson.M{
					"$sum": "$num_covid_cases",
				},
				"last_updated": bson.M{
					"$max": "$last_updated",
				},
			},
		},
	}

	cursor, err := collection().Aggregate(ctx, pipeline)
	defer cursor.Close(ctx)
	if err != nil {
		return CovidStat{}, err
	}

	stat := CovidStat{
		Code: "IN",
		Name: "India",
		NumCovidCases: -1,
		LastUpdated: time.Now(),
	}

	for cursor.Next(ctx) {
		var result struct {
			NumCovidCases int `bson:"num_covid_cases"`
			LastUpdated time.Time `bson:"last_updated"`
		}
		err := cursor.Decode(&result)
		if err != nil {
			return CovidStat{}, err
		}

		stat.NumCovidCases = result.NumCovidCases
		stat.LastUpdated = result.LastUpdated
	}

	return stat, nil
}

func UpdateCovidStats(ctx context.Context, stats []CovidStat) error {
	operations := []mongo.WriteModel{}

	for _, stat := range stats {
		operation := mongo.NewUpdateOneModel()
		operation.SetFilter(bson.M{"_id":stat.Code})
		operation.SetUpdate(bson.M{
			"$set": bson.M{
				"num_covid_cases": stat.NumCovidCases,
				"last_updated":    stat.LastUpdated,
			},
		})

		operation.SetUpsert(true)
		operations = append(operations, operation)
	}

	bulkOptions := options.BulkWriteOptions{}
	bulkOptions.SetOrdered(false)

	_, err := collection().BulkWrite(ctx, operations, &bulkOptions)

	return err
}

func UpdateStates(ctx context.Context, stats []CovidStat) error {
	operations := []mongo.WriteModel{}

	for _, stat := range stats {
		operation := mongo.NewUpdateOneModel()
		operation.SetFilter(bson.M{"_id": stat.Code})
		operation.SetUpdate(bson.M{
			"$set": bson.M{
				"_id" : stat.Code,
				"name": stat.Name,
			},
		})

		operation.SetUpsert(true)
		operations = append(operations, operation)
	}

	bulkOptions := options.BulkWriteOptions{}
	bulkOptions.SetOrdered(false)

	_, err := collection().BulkWrite(ctx, operations, &bulkOptions)

	return err
}

func collection() *mongo.Collection {
	return m.GetDB().Collection("covid_numbers")
}