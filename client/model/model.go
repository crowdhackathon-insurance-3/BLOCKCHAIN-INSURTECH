package model

// InsuranceProductType
type InsuranceProductType struct {// entry
	Id  		    	string `json:"id"`
	Name        		string `json:"name"`
	StartDate  			string `json:"startDate"`
	EndDate  			string `json:"endDate"`
	MinAmount 			int64  `json:"minAmount"`
	MaxAmount 			int64  `json:"maxAmount"`
	Terms				string `json:"terms"`
	DRate	 			int64  `json:"dRate"`
}

// InsuranceEntryType
type InsuranceEntryType struct {
	Id     		     string `json:"id"`
	Item	         string `json:"item"`
	StartDate 		 string `json:"startDate"`
	NumOfDays  		 int64  `json:"numOfDays"`
	Amount           int64  `json:"amount"`
	ProductId		 string	`json:"productId"` // foreign key
	PremiumAmount    int64  `json:"pAmount"`   // calculated
}
