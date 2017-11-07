package model

type MongoInsert struct {
	ID string `json:"_id" bson:"_id" `
	SearchKey string  `json:"SearchKey" bson:"SearchKey" `
	SaveUrl []string  `json:"SaveUrl" bson:"SaveUrl" `
}
type ListSearchKey struct {
	ID string `json:"_id" bson:"_id" `
	SearchKey string  `json:"SearchKey" bson:"SearchKey" `
}
