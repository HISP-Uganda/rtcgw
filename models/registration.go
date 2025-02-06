package models

type ECHISRequest struct {
	ECHISID             string `json:"echis_patient_id"`
	NIN                 string `json:"national_identification_number"`
	Name                string `json:"name"`
	Sex                 string `json:"patient_gender"`
	FacilityID          string `json:"facility_id"`
	PatientPhone        string `json:"patient_phone"`
	PatientCategory     string `json:"patient_category"`
	PatientAgeInYears   string `json:"patient_age_in_years"`
	PatientAgeInMonths  string `json:"patient_age_in_months,omitempty"`
	PatientAgeInDays    string `json:"patient_age_in_days,omitempty"`
	ClientCategory      string `json:"client_category,omitempty"`
	Cough               string `json:"cough,omitempty"`
	Fever               string `json:"fever,omitempty"`
	WeightLoss          string `json:"weight_loss,omitempty"`
	ExcessiveNightSweat string `json:"excessive_night_sweat,omitempty"`
	IsOnTBTreatment     string `json:"is_on_tb_treatment,omitempty"`
	PoorWeightGain      string `json:"poor_weight_gain,omitempty"`
}
