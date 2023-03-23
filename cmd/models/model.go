package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID               primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Name             string               `bson:"name,omitempty" json:"name"`
	Username         string               `bson:"username,omitempty" json:"username"`
	Email            string               `bson:"email,omitempty" json:"email"`
	Avatar           Avatar               `bson:"avatar,omitempty" json:"avatar"`
	DeliveryId       primitive.ObjectID   `bson:"deliveryId,omitempty" json:"deliveryId"`
	Posts            []primitive.ObjectID `bson:"posts,omitempty" json:"posts"`
	Partners         []primitive.ObjectID `bson:"partners,omitempty" json:"partners"`
	PartnerRequests  []PartnerRequest     `bson:"partnerrequests,omitempty" json:"partnerrequests"`
	PartnerRequested []PartnerRequest     `bson:"partnerrequested,omitempty" json:"partnerrequested"`
	AccessToken      string               `bson:"accessToken,omitempty" json:"accessToken,omitempty"`
	Role             string               `bson:"role,omitempty" json:"role,omitempty"`
	CreatedAt        string               `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt        string               `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	DeletedAt        string               `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

type PartnerRequest struct {
	ID                string `bson:"id" json:"id"`
	RequesterId       string `bson:"requesterId" json:"requesterId"`
	RequesterUsername string `bson:"requesterUsername" json:"requesterUsername"`
	RequesterName     string `bson:"requesterName" json:"requesterName"`
	TargetId          string `bson:"targetId" json:"targetId"`
	TargetUsername    string `bson:"targetUsername" json:"targetUsername"`
	TargetName        string `bson:"targetName" json:"targetName"`
	PublicKey         string `bson:"publicKey" json:"publicKey"`
	CreatedAt         string `bson:"createdAt" json:"createdAt"`
}

type Avatar struct {
	PublicId  string `json:"publicId" bson:"publicId"`
	SecureUrl string `json:"secureUrl" bson:"secureUrl"`
}

type FrequencyTable struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	Username  string             `bson:"username" json:"username"`
	Frequency int                `bson:"frequency" json:"frequency"`
}

type Log struct {
	TimeStamp   string `bson:"timeStamp" json:"timeStamp"`
	ServiceType string `bson:"serviceType" json:"serviceType"`
	Type        string `bson:"type" json:"type"`
	FileName    string `bson:"fileName" json:"fileName"`
	LineNumber  int    `bson:"lineNumber" json:"lineNumber"`
	Message     string `bson:"errorMessage" json:"errorMessage"`
}

type LogPacket struct {
	NodeName string `json:"name"`
	Type     string `json:"type"`
	Message  string `json:"message"`
}
