package models

import (
	"time"

	"net/textproto"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type FileUpload struct {
	orm.ModelBase       `bson:"-",json:"-"`
	Id                  bson.ObjectId ` bson:"_id" , json:"_id" `
	Filename            string
	Header              textproto.MIMEHeader
	UploadBy            string
	UploadTime          time.Time
	Cluster             string
	ClusterDesc         string
	ClusterConfident    float64
	SubCluster          string
	SubClusterDesc      string
	SubClusterConfident float64
	Filenamesource      string
	Filenametxt         string
	UploadId            string
	Type                string
	LibraryId           string
	ContractName        string //Default is filename but user can change it
	ContractManager     string
	ContractId          int
	Metadata            MetadataModel
	Questionnaire       QuestionnaireModel
	Status              string //NEW/EMPTY, "META" AND "QUIZ"
	Resolution          bool
	StatusSummary       string
	IsSaved             bool
}

type UploadParam struct {
	Id                  string
	Filename            string
	Header              textproto.MIMEHeader
	UploadBy            string
	UploadTime          time.Time
	Cluster             string
	ClusterDesc         string
	ClusterConfident    float64
	SubCluster          string
	SubClusterDesc      string
	SubClusterConfident float64
	Filenamesource      string
	Filenametxt         string
	UploadId            string
	LibraryId           string
	ContractName        string //Default is filename but user can change it
	ContractManager     string
	ContractId          int
}

type ClauseDetailModel struct {
	Clause              string
	CommonAncestorId    string `json:"common-ancestor-id"`
	Confidence          float64
	ContentFromUser     interface{}
	ContentFromUserDesc string
	ContentFromSystem   string `json:"contentfromsystem"`
	DataEaciitEnd       string `json:"data-eaciit-end"`
	DataEaciitId        string `json:"data-eaciit-id"`
	DataEaciitStart     string `json:"data-eaciit-start"`
	EndOffset           string `json:"end-offset"`
	StartOffset         string `json:"start-offset"`
	Identified          bool
	ConfirmOrTrained    string
}

type MetadataModel struct {
	DocumentTitle []ClauseDetailModel
	EffectiveDate []ClauseDetailModel
	GoverningLaw  []ClauseDetailModel
}
type QuestionnaireModel struct {
	ProtectionClause  []ClauseDetailModel
	AntiBriberyClause []ClauseDetailModel
}

type ConfirmClauseParam struct {
	Id     string
	Clause string
	Data   []ClauseDetailModel
	Status string
}

type SaveCloseParam struct {
	Id   string
	Data []struct {
		Clause   string
		Metadata []ClauseDetailModel
	}
}

func NewFileUpload() *FileUpload {
	m := new(FileUpload)

	return m
}
func (e *FileUpload) RecordID() interface{} {
	return e.Id
}

func (m *FileUpload) TableName() string {
	return "UploadInformation"
}
