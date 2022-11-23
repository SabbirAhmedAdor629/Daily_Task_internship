package main()

type st struct{
	Action string		  `json:"action"`
	SubmittedAt string		  `json:"submitted_at"`
	Data Data				`json:"data"`
}

type Data struct{
	S3File string		`json:"s3_file"`
	CampaignId int		`json:"campaign_id"`
	LogChunkId int		`json:"log_chunk_id"`
}

