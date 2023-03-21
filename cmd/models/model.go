package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID               primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	Name             string               `bson:"name" json:"name,omitempty"`
	Username         string               `bson:"username" json:"username,omitempty"`
	Email            string               `bson:"email" json:"email,omitempty"`
	Avatar           Avatar               `bson:"avatar" json:"avatar,omitempty"`
	MessageID        primitive.ObjectID   `bson:"messageId" json:"messageId"`
	Posts            []primitive.ObjectID `bson:"posts" json:"posts"`
	Partners         []primitive.ObjectID `bson:"partners" json:"partners,omitempty"`
	PartnerRequests  []PartnerRequest     `bson:"partnerrequests" json:"partnerrequests,omitempty"`
	PartnerRequested []PartnerRequest     `bson:"partnerrequested" json:"partnerrequested,omitempty"`
	Role             string               `bson:"role" json:"role,omitempty"`
	Token            string               `bson:"token" json:"token,omitempty"`
	Logout           bool                 `bson:"logout" json:"logout,omitempty"`
	CreatedAt        string               `bson:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt        string               `bson:"updatedAt" json:"updatedAt,omitempty"`
	DeletedAt        string               `bson:"deletedAt" json:"deletedAt,omitempty"`
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
